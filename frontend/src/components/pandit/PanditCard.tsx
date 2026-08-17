import { Star, MapPin, Clock } from 'lucide-react'
import type { Pandit } from '../../types'

interface PanditCardProps {
  pandit: Pandit
  onClick?: () => void
}

function Avatar({ name }: { name: string }) {
  const parts = name.trim().split(/\s+/)
  const initials = parts.length >= 2
    ? (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
    : parts[0][0].toUpperCase()
  return (
    <div className="w-20 h-20 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center flex-shrink-0">
      <span className="text-2xl font-bold text-primary-600 dark:text-primary-300">{initials}</span>
    </div>
  )
}

export default function PanditCard({ pandit, onClick }: PanditCardProps) {
  return (
    <div className="card-hover p-6" onClick={onClick}>
      <div className="flex items-start gap-5">
        <Avatar name={pandit.full_name} />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 mb-1">
            <h3 className="text-2xl font-bold truncate">{pandit.full_name}</h3>
            {pandit.is_available && (
              <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-sacred-100 text-sacred-700 dark:bg-sacred-900 dark:text-sacred-300 flex-shrink-0">
                Available
              </span>
            )}
          </div>
          <p className="text-base font-medium text-primary-600 dark:text-primary-400">
            {pandit.specialization}
          </p>
          <div className="flex items-center gap-4 mt-2 text-sm text-gray-500 dark:text-gray-400">
            <span className="flex items-center gap-1">
              <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
              {pandit.rating.toFixed(1)} ({pandit.total_reviews})
            </span>
            <span className="flex items-center gap-1">
              <Clock className="w-4 h-4" />
              {pandit.experience_years} yrs
            </span>
            <span className="flex items-center gap-1">
              <MapPin className="w-4 h-4" />
              {pandit.service_area}
            </span>
          </div>
        </div>
      </div>
      <div className="mt-5 pt-4 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between">
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {pandit.total_bookings} {pandit.total_bookings === 1 ? 'booking' : 'bookings'}
        </span>
        <span className="text-2xl font-bold text-primary-600 dark:text-primary-400">
          Rs. {pandit.base_price.toLocaleString('en-IN')}
        </span>
      </div>
    </div>
  )
}
