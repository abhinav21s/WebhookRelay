import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchEvents, replayEvent } from '../services/api'
import { Link, useNavigate } from 'react-router-dom'
import Card from '../components/Card'
import Button from '../components/Button'

export default function DeadLetters() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: events = [] } = useQuery({
    queryKey: ['events'],
    queryFn: fetchEvents,
    refetchInterval: 3000,
  })

  const deadLetters = events.filter((e) => e.status === 'dead_letter' || e.status === 'failed')

  const replayMutation = useMutation({
    mutationFn: replayEvent,
    onSuccess: (newEvent) => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
      navigate(`/events/${newEvent.id}`)
    },
  })

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Dead Letters</h1>
        <p className="mt-1 text-slate-400">
          Events that failed after maximum retry attempts
        </p>
      </div>

      {/* Dead Letters List */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left text-sm text-slate-400 border-b border-slate-700">
                <th className="pb-3 font-medium">Event ID</th>
                <th className="pb-3 font-medium">Type</th>
                <th className="pb-3 font-medium">Failed At</th>
                <th className="pb-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700">
              {deadLetters.map((event) => (
                <tr key={event.id} className="text-sm">
                  <td className="py-4">
                    <Link
                      to={`/events/${event.id}`}
                      className="text-blue-500 hover:text-blue-400 font-mono text-xs"
                    >
                      {event.id.substring(0, 8)}...
                    </Link>
                  </td>
                  <td className="py-4 text-slate-300">{event.event_type}</td>
                  <td className="py-4 text-slate-400">
                    {new Date(event.created_at).toLocaleString()}
                  </td>
                  <td className="py-4">
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        onClick={() => replayMutation.mutate(event.id)}
                      >
                        Replay
                      </Button>
                      <Link to={`/events/${event.id}`}>
                        <Button size="sm" variant="secondary">
                          Inspect
                        </Button>
                      </Link>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {deadLetters.length === 0 && (
            <div className="text-center py-12 text-slate-400">
              <p className="text-lg mb-2">No dead letters</p>
              <p className="text-sm">All events are delivering successfully!</p>
            </div>
          )}
        </div>
      </Card>
    </div>
  )
}
