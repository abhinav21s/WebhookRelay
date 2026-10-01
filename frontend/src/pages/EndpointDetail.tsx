import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchEndpoint, createEvent, setMockMode, getMockMode } from '../services/api'
import Card from '../components/Card'
import Button from '../components/Button'

export default function EndpointDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [failureMode, setFailureMode] = useState('success')
  const [showTestModal, setShowTestModal] = useState(false)

  const { data: endpoint } = useQuery({
    queryKey: ['endpoint', id],
    queryFn: () => fetchEndpoint(id!),
    enabled: !!id,
  })

  const { data: mockModeData } = useQuery({
    queryKey: ['mockMode'],
    queryFn: getMockMode,
    refetchInterval: 3000,
  })

  const setModeMutation = useMutation({
    mutationFn: setMockMode,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mockMode'] })
    },
  })

  if (!endpoint) {
    return (
      <div className="text-center py-12">
        <div className="text-slate-400">Loading...</div>
      </div>
    )
  }

  const handleModeChange = (mode: string) => {
    setFailureMode(mode)
    setModeMutation.mutate(mode)
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">{endpoint.name}</h1>
          <p className="mt-1 text-slate-400">{endpoint.url}</p>
        </div>
        <Button variant="secondary" onClick={() => navigate('/endpoints')}>
          Back to Endpoints
        </Button>
      </div>

      {/* Endpoint Info */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">Endpoint Details</h2>
        <div className="space-y-4">
          <div>
            <div className="text-sm text-slate-400 mb-1">URL</div>
            <code className="text-sm text-slate-300 bg-slate-900 px-3 py-2 rounded block">
              {endpoint.url}
            </code>
          </div>
          <div>
            <div className="text-sm text-slate-400 mb-1">Secret</div>
            <code className="text-sm text-slate-300 bg-slate-900 px-3 py-2 rounded block">
              {endpoint.secret.substring(0, 20)}••••••••
            </code>
          </div>
          <div>
            <div className="text-sm text-slate-400 mb-1">Event Types</div>
            <div className="flex flex-wrap gap-2">
              {endpoint.event_types.map((type) => (
                <span
                  key={type}
                  className="text-xs bg-slate-700 text-slate-300 px-3 py-1 rounded"
                >
                  {type}
                </span>
              ))}
            </div>
          </div>
          <div>
            <div className="text-sm text-slate-400 mb-1">Status</div>
            <span
              className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium ${
                endpoint.is_active
                  ? 'bg-green-500 text-white'
                  : 'bg-slate-600 text-slate-200'
              }`}
            >
              {endpoint.is_active ? 'Active' : 'Disabled'}
            </span>
          </div>
        </div>
      </Card>

      {/* Test Event */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">Send Test Event</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">
              Failure Mode
            </label>
            <select
              value={mockModeData?.mode || failureMode}
              onChange={(e) => handleModeChange(e.target.value)}
              className="w-full md:w-64 px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="success">Success (200)</option>
              <option value="error500">500 Internal Server Error</option>
              <option value="timeout">Timeout</option>
              <option value="rate_limit">429 Rate Limited</option>
              <option value="bad_request">400 Bad Request</option>
            </select>
            <p className="mt-2 text-sm text-slate-400">
              Current mode: <span className="text-blue-400">{mockModeData?.mode || 'success'}</span>
            </p>
          </div>
          <Button onClick={() => setShowTestModal(true)}>
            Send Test Event
          </Button>
        </div>
      </Card>

      {/* Test Event Modal */}
      {showTestModal && (
        <TestEventModal
          endpoint={endpoint}
          onClose={() => setShowTestModal(false)}
          onSuccess={() => {
            setShowTestModal(false)
            queryClient.invalidateQueries({ queryKey: ['events'] })
          }}
        />
      )}
    </div>
  )
}

function TestEventModal({
  endpoint,
  onClose,
  onSuccess,
}: {
  endpoint: any
  onClose: () => void
  onSuccess: () => void
}) {
  const [eventType, setEventType] = useState(endpoint.event_types[0] || '')
  const [payload, setPayload] = useState(
    JSON.stringify({ test: true, timestamp: new Date().toISOString() }, null, 2)
  )

  const createMutation = useMutation({
    mutationFn: createEvent,
    onSuccess,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const parsedPayload = JSON.parse(payload)
      createMutation.mutate({
        event_type: eventType,
        payload: parsedPayload,
      })
    } catch (error) {
      alert('Invalid JSON payload')
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-slate-800 rounded-lg p-6 w-full max-w-2xl border border-slate-700">
        <h2 className="text-xl font-bold text-white mb-4">Send Test Event</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              Event Type
            </label>
            <select
              value={eventType}
              onChange={(e) => setEventType(e.target.value)}
              className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {endpoint.event_types.map((type: string) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              Payload (JSON)
            </label>
            <textarea
              value={payload}
              onChange={(e) => setPayload(e.target.value)}
              rows={10}
              className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div className="flex gap-3 mt-6">
            <Button type="submit" className="flex-1">
              Send Event
            </Button>
            <Button variant="secondary" onClick={onClose} className="flex-1">
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
