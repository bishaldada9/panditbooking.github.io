import { Clock, IndianRupee } from 'lucide-react'
import type { Ritual } from '../../types'

interface RitualCardProps {
  ritual: Ritual
  onClick?: () => void
}

export default function RitualCard({ ritual, onClick }: RitualCardProps) {
  return (
    <div className="card-hover" onClick={onClick}>
      <div className="flex items-start justify-between mb-3">
        <h3 className="text-lg font-semibold">{ritual.name}</h3>
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200">
          {ritual.category_name}
        </span>
      </div>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
        {ritual.description}
      </p>
      <div className="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
        <span className="flex items-center gap-1">
          <Clock className="w-4 h-4" />
          {ritual.duration}
        </span>
        <span className="flex items-center gap-1">
          <IndianRupee className="w-4 h-4" />
          {ritual.base_price}
        </span>
      </div>
    </div>
  )
}
