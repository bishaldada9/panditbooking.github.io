import { Calendar, Clock, MapPin, IndianRupee, ChevronRight } from 'lucide-react'
import type { Booking } from '../../types'

interface BookingCardProps {
  booking: Booking
  onClick?: () => void
  role: 'customer' | 'pandit'
}

const statusStyles: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  confirmed: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  completed: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  cancelled: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  rejected: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
}

function shortText(value: unknown) {
  const text = String(value ?? '')
  return text ? text.slice(0, 8) : '-'
}

function formatDate(value: unknown) {
  if (!value) return '-'
  const date = new Date(String(value))
  return Number.isNaN(date.getTime())
    ? '-'
    : date.toLocaleDateString('en-IN', { day: 'numeric', month: 'long', year: 'numeric' })
}

export default function BookingCard({ booking, onClick, role }: BookingCardProps) {
  const status = booking.status || 'pending'

  return (
    <div className="card-hover" onClick={onClick}>
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center gap-2">
          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${statusStyles[status] || statusStyles.pending}`}>
            {status}
          </span>
          <span className="text-xs text-gray-400 dark:text-gray-500 font-mono">
            #{shortText(booking.id)}
          </span>
        </div>
        <ChevronRight className="w-5 h-5 text-gray-400 shrink-0" />
      </div>

      <div className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4 shrink-0" />
          <span>{formatDate(booking.scheduled_date)}</span>
        </div>
        <div className="flex items-center gap-2">
          <Clock className="w-4 h-4 shrink-0" />
          <span>{booking.start_time || '-'} - {booking.end_time || '-'}</span>
        </div>
        <div className="flex items-start gap-2">
          <MapPin className="w-4 h-4 shrink-0 mt-0.5" />
          <span className="line-clamp-1">{booking.address || '-'}</span>
        </div>
        <div className="flex items-center gap-2 pt-1 border-t border-gray-100 dark:border-gray-700">
          <IndianRupee className="w-4 h-4 shrink-0" />
          <span className="font-semibold text-gray-900 dark:text-gray-100">₹{Number(booking.total_amount || 0).toLocaleString()}</span>
        </div>
      </div>
    </div>
  )
}
