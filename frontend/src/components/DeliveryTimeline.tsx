import type { DeliveryAttempt } from '../types'
import clsx from 'clsx'

interface DeliveryTimelineProps {
  deliveries: DeliveryAttempt[]
}

export default function DeliveryTimeline({ deliveries }: DeliveryTimelineProps) {
  const getStatusColor = (statusCode?: number, errorMessage?: string) => {
    if (errorMessage && !statusCode) return 'bg-red-500'
    if (!statusCode) return 'bg-slate-500'
    if (statusCode >= 200 && statusCode < 300) return 'bg-green-500'
    if (statusCode >= 400 && statusCode < 500) return 'bg-amber-500'
    if (statusCode >= 500) return 'bg-red-500'
    return 'bg-slate-500'
  }

  const formatTime = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString()
  }

  const getStatusText = (statusCode?: number, errorMessage?: string) => {
    if (errorMessage && !statusCode) return errorMessage
    if (!statusCode) return 'No response'
    if (statusCode >= 200 && statusCode < 300) return `Succeeded (${statusCode})`
    return `Failed (${statusCode})`
  }

  if (deliveries.length === 0) {
    return (
      <div className="text-center text-slate-400 py-8">
        No delivery attempts yet
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {deliveries.map((delivery, index) => (
        <div key={delivery.id} className="flex gap-4">
          {/* Timeline indicator */}
          <div className="flex flex-col items-center">
            <div
              className={clsx(
                'w-3 h-3 rounded-full',
                getStatusColor(delivery.status_code, delivery.error_message)
              )}
            />
            {index < deliveries.length - 1 && (
              <div className="w-0.5 flex-1 bg-slate-700 min-h-[40px]" />
            )}
          </div>

          {/* Content */}
          <div className="flex-1 pb-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="font-medium text-white">
                  Attempt {delivery.attempt_number}
                </span>
                <span className="text-sm text-slate-400">
                  {formatTime(delivery.created_at)}
                </span>
              </div>
              <span className="text-sm text-slate-400">
                {delivery.latency_ms}ms
              </span>
            </div>
            <div className="mt-1 text-sm">
              <span
                className={clsx(
                  delivery.status_code && delivery.status_code >= 200 && delivery.status_code < 300
                    ? 'text-green-400'
                    : 'text-red-400'
                )}
              >
                {getStatusText(delivery.status_code, delivery.error_message)}
              </span>
            </div>
            {delivery.error_message && delivery.status_code && (
              <div className="mt-2 text-xs text-slate-400 bg-slate-900 rounded p-2 font-mono">
                {delivery.error_message}
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}
