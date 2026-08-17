import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { ArrowLeft, Star, MapPin, Languages, Clock, CheckCircle, XCircle } from 'lucide-react'
import { panditsService } from '../services/pandits.service'
import { reviewsService } from '../services/reviews.service'
import { LoadingSpinner, EmptyState } from '../components/common'
import type { Pandit, Review } from '../types'

function RatingStars({ rating, size = 'md' }: { rating: number; size?: 'sm' | 'md' }) {
  const className = size === 'sm' ? 'w-3.5 h-3.5' : 'w-5 h-5'
  return (
    <div className="flex items-center gap-0.5">
      {[1, 2, 3, 4, 5].map((star) => (
        <Star
          key={star}
          className={`${className} ${star <= Math.round(rating) ? 'fill-yellow-400 text-yellow-400' : 'fill-gray-200 text-gray-200 dark:fill-gray-600 dark:text-gray-600'}`}
        />
      ))}
    </div>
  )
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('en-IN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export default function PanditDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [pandit, setPandit] = useState<Pandit | null>(null)
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)
  const [reviewsLoading, setReviewsLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    const fetchPandit = async () => {
      try {
        const data = await panditsService.getById(id)
        setPandit(data)
      } catch {
        toast.error('Failed to load pandit details')
      } finally {
        setLoading(false)
      }
    }
    fetchPandit()
  }, [id])

  useEffect(() => {
    if (!id) return
    const fetchReviews = async () => {
      try {
        const data = await reviewsService.getPanditReviews(id)
        setReviews(data.filter((r) => r.is_visible))
      } catch {
        // reviews are optional, don't show error
      } finally {
        setReviewsLoading(false)
      }
    }
    fetchReviews()
  }, [id])

  if (loading) return <LoadingSpinner size="lg" />
  if (!pandit) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <EmptyState
          title="Pandit not found"
          description="This pandit profile may have been removed or the link is invalid."
          action={<button onClick={() => navigate('/pandits')} className="btn-primary">Back to Pandits</button>}
        />
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <button
        onClick={() => navigate('/pandits')}
        className="flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white mb-6 transition-colors"
      >
        <ArrowLeft className="w-5 h-5" />
        Back to Pandits
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2">
          <div className="flex flex-wrap items-center gap-3 mb-4">
            <h1 className="text-3xl font-bold">{pandit.full_name}</h1>
            {pandit.is_available ? (
              <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-sacred-100 text-sacred-800 dark:bg-sacred-900 dark:text-sacred-200">
                <CheckCircle className="w-4 h-4" />
                Available
              </span>
            ) : (
              <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200">
                <XCircle className="w-4 h-4" />
                Unavailable
              </span>
            )}
          </div>

          <p className="text-lg font-medium text-primary-600 dark:text-primary-400 mb-4">
            {pandit.specialization}
          </p>

          <p className="text-gray-600 dark:text-gray-400 mb-8 leading-relaxed">
            {pandit.bio}
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-8">
            <div className="flex items-center gap-3">
              <Clock className="w-5 h-5 text-gray-400" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Experience</p>
                <p className="font-semibold">{pandit.experience_years} years</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <MapPin className="w-5 h-5 text-gray-400" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Service Area</p>
                <p className="font-semibold">{pandit.service_area}</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <Languages className="w-5 h-5 text-gray-400" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Languages</p>
                <div className="flex flex-wrap gap-1 mt-1">
                  {pandit.languages.map((lang) => (
                    <span key={lang} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 dark:bg-gray-700">
                      {lang}
                    </span>
                  ))}
                </div>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <Star className="w-5 h-5 text-gray-400" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Rating</p>
                <div className="flex items-center gap-2 mt-0.5">
                  <RatingStars rating={pandit.rating} />
                  <span className="font-semibold">{pandit.rating.toFixed(1)}</span>
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    ({pandit.total_reviews} reviews)
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Reviews</h2>
            {reviewsLoading ? (
              <LoadingSpinner size="sm" />
            ) : reviews.length === 0 ? (
              <EmptyState
                title="No reviews yet"
                description="Be the first to leave a review for this pandit."
              />
            ) : (
              <div className="space-y-4">
                {reviews.map((review) => (
                  <div key={review.id} className="card">
                    <div className="flex items-center justify-between mb-2">
                      <RatingStars rating={review.rating} size="sm" />
                      <span className="text-xs text-gray-500 dark:text-gray-400">
                        {formatDate(review.created_at)}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {review.comment}
                    </p>
                    {review.admin_reply && (
                      <div className="mt-3 pl-4 border-l-2 border-primary-300 dark:border-primary-700">
                        <p className="text-xs font-medium text-primary-600 dark:text-primary-400 mb-1">
                          Admin Response
                        </p>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {review.admin_reply}
                        </p>
                      </div>
                    )}
                    {review.is_verified && (
                      <span className="inline-flex items-center gap-1 mt-2 text-xs text-sacred-600 dark:text-sacred-400">
                        <CheckCircle className="w-3 h-3" />
                        Verified Booking
                      </span>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="lg:col-span-1">
          <div className="card space-y-6">
            <div className="text-center">
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Base Price</p>
              <p className="text-3xl font-bold text-primary-600 dark:text-primary-400">
                Rs. {pandit.base_price}
              </p>
            </div>

            <button
              onClick={() => navigate(`/book-pandit/${pandit.id}`)}
              className="btn-primary w-full"
            >
              Book This Pandit
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
