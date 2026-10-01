import clsx from 'clsx'
import type { Event } from '../types'

interface StatusBadgeProps {
  status: Event['status']
  attemptInfo?: string
}

export default function StatusBadge({ status, attemptInfo }: StatusBadgeProps) {
  const getStatusColor = () => {
    switch (status) {
      case 'pending':
        return 'bg-slate-600 text-slate-200'
      case 'delivering':
        return 'bg-blue-500 text-white animate-pulse'
      case 'retrying':
        return 'bg-amber-500 text-white'
      case 'delivered':
        return 'bg-green-500 text-white'
      case 'failed':
        return 'bg-red-500 text-white'
      case 'dead_letter':
        return 'bg-purple-600 text-white'
      default:
        return 'bg-slate-600 text-slate-200'
    }
  }

  const getStatusLabel = () => {
    switch (status) {
      case 'pending':
        return 'Pending'
      case 'delivering':
        return 'Delivering'
      case 'retrying':
        return attemptInfo ? `Retrying ${attemptInfo}` : 'Retrying'
      case 'delivered':
        return 'Delivered'
      case 'failed':
        return 'Failed'
      case 'dead_letter':
        return 'Dead Letter'
      default:
        return status
    }
  }

  return (
    <span
      className={clsx(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
        getStatusColor()
      )}
    >
      {getStatusLabel()}
    </span>
  )
}
