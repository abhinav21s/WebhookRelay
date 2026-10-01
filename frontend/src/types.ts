export interface Endpoint {
  id: string
  name: string
  url: string
  secret: string
  event_types: string[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Event {
  id: string
  event_type: string
  payload: any
  status: 'pending' | 'delivering' | 'delivered' | 'retrying' | 'failed' | 'dead_letter'
  created_at: string
}

export interface DeliveryAttempt {
  id: string
  event_id: string
  endpoint_id: string
  attempt_number: number
  status_code?: number
  latency_ms: number
  error_message?: string
  created_at: string
}

export interface EventWithDeliveries {
  event: Event
  deliveries: DeliveryAttempt[]
}

export interface Metrics {
  total_events_24h: number
  success_rate: number
  average_latency_ms: number
  active_endpoints: number
}

export interface CreateEndpointRequest {
  name: string
  url: string
  event_types: string[]
  secret?: string
}

export interface CreateEventRequest {
  event_type: string
  payload: any
}
