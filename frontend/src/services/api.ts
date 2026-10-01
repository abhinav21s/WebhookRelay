import type { 
  Endpoint, 
  Event, 
  EventWithDeliveries,
  DeliveryAttempt,
  Metrics,
  CreateEndpointRequest,
  CreateEventRequest
} from '../types'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Endpoints
export async function fetchEndpoints(): Promise<Endpoint[]> {
  const response = await fetch(`${API_URL}/api/webhooks`)
  if (!response.ok) throw new Error('Failed to fetch endpoints')
  return response.json()
}

export async function fetchEndpoint(id: string): Promise<Endpoint> {
  const response = await fetch(`${API_URL}/api/webhooks/${id}`)
  if (!response.ok) throw new Error('Failed to fetch endpoint')
  return response.json()
}

export async function createEndpoint(data: CreateEndpointRequest): Promise<Endpoint> {
  const response = await fetch(`${API_URL}/api/webhooks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!response.ok) throw new Error('Failed to create endpoint')
  return response.json()
}

export async function updateEndpoint(id: string, data: Partial<Endpoint>): Promise<void> {
  const response = await fetch(`${API_URL}/api/webhooks/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!response.ok) throw new Error('Failed to update endpoint')
}

export async function deleteEndpoint(id: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/webhooks/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) throw new Error('Failed to delete endpoint')
}

// Events
export async function fetchEvents(): Promise<Event[]> {
  const response = await fetch(`${API_URL}/api/events`)
  if (!response.ok) throw new Error('Failed to fetch events')
  return response.json()
}

export async function fetchEvent(id: string): Promise<EventWithDeliveries> {
  const response = await fetch(`${API_URL}/api/events/${id}`)
  if (!response.ok) throw new Error('Failed to fetch event')
  return response.json()
}

export async function createEvent(data: CreateEventRequest): Promise<Event> {
  const response = await fetch(`${API_URL}/api/events`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!response.ok) throw new Error('Failed to create event')
  return response.json()
}

export async function replayEvent(id: string): Promise<Event> {
  const response = await fetch(`${API_URL}/api/events/${id}/replay`, {
    method: 'POST',
  })
  if (!response.ok) throw new Error('Failed to replay event')
  return response.json()
}

// Deliveries
export async function fetchDeliveries(eventId?: string): Promise<DeliveryAttempt[]> {
  const url = eventId 
    ? `${API_URL}/api/deliveries?event_id=${eventId}`
    : `${API_URL}/api/deliveries`
  const response = await fetch(url)
  if (!response.ok) throw new Error('Failed to fetch deliveries')
  return response.json()
}

// Metrics
export async function fetchMetrics(): Promise<Metrics> {
  const response = await fetch(`${API_URL}/api/metrics`)
  if (!response.ok) throw new Error('Failed to fetch metrics')
  return response.json()
}

// Mock receiver control
export async function setMockMode(mode: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/mock/control`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ mode }),
  })
  if (!response.ok) throw new Error('Failed to set mock mode')
}

export async function getMockMode(): Promise<{ mode: string }> {
  const response = await fetch(`${API_URL}/api/mock/mode`)
  if (!response.ok) throw new Error('Failed to get mock mode')
  return response.json()
}
