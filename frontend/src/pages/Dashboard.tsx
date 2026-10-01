import { useQuery } from '@tanstack/react-query'
import { fetchMetrics, fetchEvents, fetchEndpoints } from '../services/api'
import Card from '../components/Card'
import StatusBadge from '../components/StatusBadge'
import { Link } from 'react-router-dom'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

export default function Dashboard() {
  const { data: metrics } = useQuery({
    queryKey: ['metrics'],
    queryFn: fetchMetrics,
    refetchInterval: 3000,
  })

  const { data: events = [] } = useQuery({
    queryKey: ['events'],
    queryFn: fetchEvents,
    refetchInterval: 2000,
  })

  const { data: endpoints = [] } = useQuery({
    queryKey: ['endpoints'],
    queryFn: fetchEndpoints,
    refetchInterval: 5000,
  })

  // Mock chart data (in real app, fetch from API)
  const chartData = [
    { time: '00:00', rate: 95 },
    { time: '04:00', rate: 92 },
    { time: '08:00', rate: 97 },
    { time: '12:00', rate: 94 },
    { time: '16:00', rate: 96 },
    { time: '20:00', rate: metrics?.success_rate || 95 },
  ]

  const recentEvents = events.slice(0, 20)

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Dashboard</h1>
        <p className="mt-1 text-slate-400">Overview of webhook delivery system</p>
      </div>

      {/* Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <div className="text-sm font-medium text-slate-400">Total Events (24h)</div>
          <div className="mt-2 text-3xl font-bold text-white">
            {metrics?.total_events_24h || 0}
          </div>
        </Card>

        <Card>
          <div className="text-sm font-medium text-slate-400">Success Rate</div>
          <div className="mt-2 text-3xl font-bold text-green-500">
            {metrics?.success_rate?.toFixed(1) || 0}%
          </div>
        </Card>

        <Card>
          <div className="text-sm font-medium text-slate-400">Avg Latency</div>
          <div className="mt-2 text-3xl font-bold text-blue-500">
            {Math.round(metrics?.average_latency_ms || 0)}ms
          </div>
        </Card>

        <Card>
          <div className="text-sm font-medium text-slate-400">Active Endpoints</div>
          <div className="mt-2 text-3xl font-bold text-white">
            {metrics?.active_endpoints || 0}
          </div>
        </Card>
      </div>

      {/* Success Rate Chart */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">Success Rate (24h)</h2>
        <ResponsiveContainer width="100%" height={200}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis dataKey="time" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" domain={[0, 100]} />
            <Tooltip
              contentStyle={{
                backgroundColor: '#1e293b',
                border: '1px solid #334155',
                borderRadius: '8px',
                color: '#f8fafc',
              }}
            />
            <Line
              type="monotone"
              dataKey="rate"
              stroke="#22c55e"
              strokeWidth={2}
              dot={{ fill: '#22c55e' }}
            />
          </LineChart>
        </ResponsiveContainer>
      </Card>

      {/* Endpoint Health */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">Endpoint Health</h2>
        <div className="flex flex-wrap gap-3">
          {endpoints.map((endpoint) => (
            <Link
              key={endpoint.id}
              to={`/endpoints/${endpoint.id}`}
              className="flex items-center gap-2 px-3 py-2 bg-slate-700 hover:bg-slate-600 rounded-lg transition-colors"
            >
              <div
                className={`w-2 h-2 rounded-full ${
                  endpoint.is_active ? 'bg-green-500' : 'bg-slate-500'
                }`}
              />
              <span className="text-sm text-white">{endpoint.name}</span>
            </Link>
          ))}
          {endpoints.length === 0 && (
            <p className="text-slate-400 text-sm">No endpoints configured</p>
          )}
        </div>
      </Card>

      {/* Recent Events */}
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-white">Recent Events</h2>
          <Link to="/events" className="text-sm text-blue-500 hover:text-blue-400">
            View all →
          </Link>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left text-sm text-slate-400 border-b border-slate-700">
                <th className="pb-3 font-medium">Event ID</th>
                <th className="pb-3 font-medium">Type</th>
                <th className="pb-3 font-medium">Status</th>
                <th className="pb-3 font-medium">Created</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700">
              {recentEvents.map((event) => (
                <tr key={event.id} className="text-sm">
                  <td className="py-3">
                    <Link
                      to={`/events/${event.id}`}
                      className="text-blue-500 hover:text-blue-400 font-mono"
                    >
                      {event.id.substring(0, 8)}...
                    </Link>
                  </td>
                  <td className="py-3 text-slate-300">{event.event_type}</td>
                  <td className="py-3">
                    <StatusBadge status={event.status} />
                  </td>
                  <td className="py-3 text-slate-400">
                    {new Date(event.created_at).toLocaleString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {recentEvents.length === 0 && (
            <div className="text-center py-8 text-slate-400">
              No events yet. Publish your first event!
            </div>
          )}
        </div>
      </Card>
    </div>
  )
}
