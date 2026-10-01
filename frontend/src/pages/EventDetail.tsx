import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchEvent, replayEvent } from '../services/api'
import Card from '../components/Card'
import Button from '../components/Button'
import StatusBadge from '../components/StatusBadge'
import DeliveryTimeline from '../components/DeliveryTimeline'

export default function EventDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['event', id],
    queryFn: () => fetchEvent(id!),
    refetchInterval: 2000,
    enabled: !!id,
  })

  const replayMutation = useMutation({
    mutationFn: () => replayEvent(id!),
    onSuccess: (newEvent) => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
      navigate(`/events/${newEvent.id}`)
    },
  })

  if (!data) {
    return (
      <div className="text-center py-12">
        <div className="text-slate-400">Loading...</div>
      </div>
    )
  }

  const { event, deliveries } = data

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">Event Detail</h1>
          <p className="mt-1 text-slate-400">
            <code className="font-mono text-sm">{event.id}</code>
          </p>
        </div>
        <div className="flex gap-3">
          {(event.status === 'failed' || event.status === 'dead_letter') && (
            <Button onClick={() => replayMutation.mutate()}>
              Replay Event
            </Button>
          )}
          <Button variant="secondary" onClick={() => navigate('/events')}>
            Back to Events
          </Button>
        </div>
      </div>

      {/* Event Info */}
      <Card>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <div className="text-sm text-slate-400 mb-1">Event Type</div>
            <div className="text-white font-medium">{event.event_type}</div>
          </div>
          <div>
            <div className="text-sm text-slate-400 mb-1">Status</div>
            <StatusBadge status={event.status} />
          </div>
          <div>
            <div className="text-sm text-slate-400 mb-1">Created At</div>
            <div className="text-white">
              {new Date(event.created_at).toLocaleString()}
            </div>
          </div>
        </div>
      </Card>

      {/* Payload */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">Payload</h2>
        <pre className="bg-slate-900 rounded-lg p-4 overflow-x-auto text-sm text-slate-300 font-mono">
          {JSON.stringify(event.payload, null, 2)}
        </pre>
      </Card>

      {/* Delivery Timeline */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-6">Delivery Timeline</h2>
        <DeliveryTimeline deliveries={deliveries} />
      </Card>
    </div>
  )
}
