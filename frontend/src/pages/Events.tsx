import { useQuery } from '@tanstack/react-query'
import { fetchEvents } from '../services/api'
import { Link } from 'react-router-dom'
import Card from '../components/Card'
import StatusBadge from '../components/StatusBadge'

export default function Events() {
  const { data: events = [] } = useQuery({
    queryKey: ['events'],
    queryFn: fetchEvents,
    refetchInterval: 2000,
  })

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Events</h1>
        <p className="mt-1 text-slate-400">View all published webhook events</p>
      </div>

      {/* Events List */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left text-sm text-slate-400 border-b border-slate-700">
                <th className="pb-3 font-medium">Event ID</th>
                <th className="pb-3 font-medium">Type</th>
                <th className="pb-3 font-medium">Status</th>
                <th className="pb-3 font-medium">Created</th>
                <th className="pb-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700">
              {events.map((event) => (
                <tr key={event.id} className="text-sm">
                  <td className="py-4">
                    <code className="text-xs text-slate-300 bg-slate-900 px-2 py-1 rounded font-mono">
                      {event.id.substring(0, 8)}...
                    </code>
                  </td>
                  <td className="py-4 text-slate-300">{event.event_type}</td>
                  <td className="py-4">
                    <StatusBadge status={event.status} />
                  </td>
                  <td className="py-4 text-slate-400">
                    {new Date(event.created_at).toLocaleString()}
                  </td>
                  <td className="py-4">
                    <Link
                      to={`/events/${event.id}`}
                      className="text-blue-500 hover:text-blue-400 text-sm font-medium"
                    >
                      View Details →
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {events.length === 0 && (
            <div className="text-center py-12 text-slate-400">
              <p className="text-lg mb-2">No events yet</p>
              <p className="text-sm">Publish your first event to see it here</p>
            </div>
          )}
        </div>
      </Card>
    </div>
  )
}
