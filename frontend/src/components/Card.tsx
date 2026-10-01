import { ReactNode } from 'react'
import clsx from 'clsx'

interface CardProps {
  children: ReactNode
  className?: string
}

export default function Card({ children, className }: CardProps) {
  return (
    <div
      className={clsx(
        'bg-slate-800 rounded-lg border border-slate-700 p-6',
        className
      )}
    >
      {children}
    </div>
  )
}
