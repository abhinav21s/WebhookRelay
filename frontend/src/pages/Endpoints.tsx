import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchEndpoints, createEndpoint, deleteEndpoint, updateEndpoint } from '../services/api'
import { Link } from 'react-router-dom'
import Card from '../components/Card'
import Button from '../components/Button'
import type { CreateEndpointRequest } from '../types'

export default function Endpoints() {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const queryClient = useQueryClient()

  const { data: endpoints = [] } = useQuery({
    queryKey: ['endpoints'],
    queryFn: fetchEndpoints,
    refetchInterval: 3000,
  })

  const deleteMutation = useMutation({
    mutationFn: deleteEndpoint,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['endpoints'] })
    },
  })

  const toggleMutation = useMutation({
    mutationFn: ({ id, isActive }: { id: string; isActive: boolean }) =>
      updateEndpoint(id, { is_active: isActive }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['endpoints'] })
    },
  })

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">Endpoints</h1>
          <p className="mt-1 text-slate-400">Manage webhook endpoints</p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          New Endpoint
        </Button>
      </div>

      {/* Endpoints List */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left text-sm text-slate-400 border-b border-slate-700">
                <th className="pb-3 font-medium">Name</th>
                <th className="pb-3 font-medium">URL</th>
                <th className="pb-3 font-medium">Event Types</th>
                <th className="pb-3 font-medium">Status</th>
                <th className="pb-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700">
              {endpoints.map((endpoint) => (
                <tr key={endpoint.id} className="text-sm">
                  <td className="py-4">
                    <Link
                      to={`/endpoints/${endpoint.id}`}
                      className="text-blue-500 hover:text-blue-400 font-medium"
                    >
                      {endpoint.name}
                    </Link>
                  </td>
                  <td className="py-4">
                    <code className="text-xs text-slate-300 bg-slate-900 px-2 py-1 rounded">
                      {endpoint.url}
                    </code>
                  </td>
                  <td className="py-4">
                    <div className="flex flex-wrap gap-1">
                      {endpoint.event_types.slice(0, 2).map((type) => (
                        <span
                          key={type}
                          className="text-xs bg-slate-700 text-slate-300 px-2 py-1 rounded"
                        >
                          {type}
                        </span>
                      ))}
                      {endpoint.event_types.length > 2 && (
                        <span className="text-xs text-slate-400">
                          +{endpoint.event_types.length - 2} more
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="py-4">
                    <button
                      onClick={() =>
                        toggleMutation.mutate({
                          id: endpoint.id,
                          isActive: !endpoint.is_active,
                        })
                      }
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        endpoint.is_active
                          ? 'bg-green-500 text-white'
                          : 'bg-slate-600 text-slate-200'
                      }`}
                    >
                      {endpoint.is_active ? 'Active' : 'Disabled'}
                    </button>
                  </td>
                  <td className="py-4">
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() => deleteMutation.mutate(endpoint.id)}
                    >
                      Delete
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {endpoints.length === 0 && (
            <div className="text-center py-12 text-slate-400">
              <p className="text-lg mb-2">No endpoints yet</p>
              <p className="text-sm">Create your first endpoint to get started</p>
            </div>
          )}
        </div>
      </Card>

      {/* Create Modal */}
      {showCreateModal && (
        <CreateEndpointModal
          onClose={() => setShowCreateModal(false)}
          onSuccess={() => {
            setShowCreateModal(false)
            queryClient.invalidateQueries({ queryKey: ['endpoints'] })
          }}
        />
      )}
    </div>
  )
}

function CreateEndpointModal({
  onClose,
  onSuccess,
}: {
  onClose: () => void
  onSuccess: () => void
}) {
  const [formData, setFormData] = useState({
    name: '',
    url: '',
    event_types: '',
    secret: '',
  })
  const [error, setError] = useState<string>('')

  const createMutation = useMutation({
    mutationFn: (data: CreateEndpointRequest) => createEndpoint(data),
    onSuccess: () => {
      console.log('✅ Endpoint created successfully')
      setError('')
      onSuccess()
    },
    onError: (error: any) => {
      console.error('❌ Create endpoint error:', error)
      setError(error.message || 'Failed to create endpoint. Check console for details.')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    
    const eventTypes = formData.event_types.split(',').map((t) => t.trim()).filter(Boolean)
    
    if (eventTypes.length === 0) {
      setError('Please enter at least one event type')
      return
    }
    
    console.log('Creating endpoint with data:', {
      name: formData.name,
      url: formData.url,
      event_types: eventTypes,
    })
    
    createMutation.mutate({
      name: formData.name,
      url: formData.url,
      event_types: eventTypes,
      secret: formData.secret || undefined,
    })
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-slate-800 rounded-lg p-6 w-full max-w-md border border-slate-700">
        <h2 className="text-xl font-bold text-white mb-4">Create Endpoint</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              Name
            </label>
            <input
              type="text"
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Payment Service"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              URL
            </label>
            <input
              type="url"
              required
              value={formData.url}
              onChange={(e) => setFormData({ ...formData, url: e.target.value })}
              className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="http://localhost:9090/webhook"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              Event Types (comma-separated)
            </label>
            <input
              type="text"
              required
              value={formData.event_types}
              onChange={(e) => setFormData({ ...formData, event_types: e.target.value })}
              className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="payment.completed, payment.failed"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              Secret (optional)
            </label>
            <input
              type="text"
              value={formData.secret}
              onChange={(e) => setFormData({ ...formData, secret: e.target.value })}
              className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Leave empty to auto-generate"
            />
          </div>
          
          {error && (
            <div className="bg-red-500/10 border border-red-500 rounded-lg p-3 text-red-400 text-sm">
              {error}
            </div>
          )}
          
          {createMutation.isPending && (
            <div className="text-blue-400 text-sm text-center">
              Creating endpoint...
            </div>
          )}
          
          <div className="flex gap-3 mt-6">
            <Button 
              type="submit" 
              className="flex-1"
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? 'Creating...' : 'Create'}
            </Button>
            <Button 
              variant="secondary" 
              onClick={onClose} 
              className="flex-1"
              disabled={createMutation.isPending}
            >
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
