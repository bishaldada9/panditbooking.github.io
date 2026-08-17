import { Star, ShieldCheck } from 'lucide-react'
import type { Review } from '../../types'

interface ReviewCardProps {
  review: Review
  showReply?: boolean
}

export default function ReviewCard({ review, showReply }: ReviewCardProps) {
  return (
    <div className="card">
      <div className="flex items-center gap-2 mb-2">
        <div className="flex">
          {Array.from({ length: 5 }).map((_, i) => (
            <Star
              key={i}
              className={`w-4 h-4 ${i < review.rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300 dark:text-gray-600'}`}
            />
          ))}
        </div>
        {review.is_verified && (
          <span className="inline-flex items-center gap-1 text-xs text-sacred-600 dark:text-sacred-400">
            <ShieldCheck className="w-3.5 h-3.5" />
            Verified
          </span>
        )}
      </div>

      {review.comment && (
        <p className="text-sm text-gray-700 dark:text-gray-300 mb-3">{review.comment}</p>
      )}

      <p className="text-xs text-gray-400 dark:text-gray-500">
        {new Date(review.created_at).toLocaleDateString('en-IN', { day: 'numeric', month: 'long', year: 'numeric' })}
      </p>

      {showReply && review.admin_reply && (
        <div className="mt-3 pl-3 border-l-2 border-gray-200 dark:border-gray-600">
          <p className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Admin response:</p>
          <p className="text-sm text-gray-600 dark:text-gray-400">{review.admin_reply}</p>
        </div>
      )}
    </div>
  )
}
