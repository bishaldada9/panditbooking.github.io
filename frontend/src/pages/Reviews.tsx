import { FormEvent, useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { useAuth } from '../store/auth.context'
import { reviewsService } from '../services/reviews.service'
import { panditsService } from '../services/pandits.service'
import { bookingsService } from '../services/bookings.service'
import { ReviewCard } from '../components/review'
import { LoadingSpinner, EmptyState } from '../components/common'
import type { Booking, Review } from '../types'

export default function Reviews() {
  const { user } = useAuth()
  const [reviews, setReviews] = useState<Review[]>([])
  const [completedBookings, setCompletedBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedBookingId, setSelectedBookingId] = useState('')
  const [rating, setRating] = useState(5)
  const [comment, setComment] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const isPandit = user?.role === 'pandit'

  useEffect(() => {
    const fetchReviews = async () => {
      try {
        if (isPandit) {
          const data = await panditsService.getProfile()
          const reviewsData = await reviewsService.getPanditReviews(data.id)
          setReviews(reviewsData)
        } else {
          const bookings = await bookingsService.getMyBookings()
          setCompletedBookings(bookings.filter((booking) => booking.status === 'completed'))
        }
      } catch {
        toast.error('Failed to load reviews')
      } finally {
        setLoading(false)
      }
    }
    fetchReviews()
  }, [isPandit])

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    const booking = completedBookings.find((item) => item.id === selectedBookingId)
    if (!booking) {
      toast.error('Select a completed booking')
      return
    }
    if (!comment.trim()) {
      toast.error('Write a short review comment')
      return
    }

    setSubmitting(true)
    try {
      const review = await reviewsService.create({
        booking_id: booking.id,
        pandit_id: booking.pandit_id,
        rating,
        comment,
      })
      setReviews((prev) => [review, ...prev])
      setCompletedBookings((prev) => prev.filter((item) => item.id !== booking.id))
      setSelectedBookingId('')
      setRating(5)
      setComment('')
      toast.success('Review submitted')
    } catch {
      toast.error('Failed to submit review')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <LoadingSpinner size="lg" />

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">
          {isPandit ? 'Reviews Received' : 'My Reviews'}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          {isPandit ? 'Reviews left by customers for your services' : 'Reviews you have written for pandits'}
        </p>
      </div>

      {!isPandit ? (
        <div className="space-y-6">
          <form onSubmit={handleSubmit} className="card space-y-4">
            <h2 className="text-xl font-semibold">Write a Review</h2>
            {completedBookings.length === 0 ? (
              <p className="text-sm text-gray-500 dark:text-gray-400">
                You need a completed booking before you can review a pandit.
              </p>
            ) : (
              <>
                <div>
                  <label className="block text-sm font-medium mb-1">Completed Booking</label>
                  <select
                    className="input-field"
                    value={selectedBookingId}
                    onChange={(e) => setSelectedBookingId(e.target.value)}
                    required
                  >
                    <option value="">Select booking</option>
                    {completedBookings.map((booking) => (
                      <option key={booking.id} value={booking.id}>
                        {booking.id.slice(0, 8)} - {new Date(booking.scheduled_date).toLocaleDateString()}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Rating</label>
                  <select
                    className="input-field"
                    value={rating}
                    onChange={(e) => setRating(Number(e.target.value))}
                  >
                    {[5, 4, 3, 2, 1].map((value) => (
                      <option key={value} value={value}>{value} star{value > 1 ? 's' : ''}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Comment</label>
                  <textarea
                    className="input-field"
                    rows={4}
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                    required
                  />
                </div>
                <button type="submit" disabled={submitting} className="btn-primary disabled:opacity-50">
                  {submitting ? 'Submitting...' : 'Submit Review'}
                </button>
              </>
            )}
          </form>

          {reviews.length > 0 && (
            <div className="space-y-4">
              <h2 className="text-xl font-semibold">Submitted This Session</h2>
              {reviews.map((review) => (
                <ReviewCard key={review.id} review={review} />
              ))}
            </div>
          )}
        </div>
      ) : reviews.length === 0 ? (
        <EmptyState
          title="No reviews yet"
          description="Reviews from customers will appear here once they complete bookings with you."
        />
      ) : (
        <div className="space-y-4">
          {reviews.map((review) => (
            <ReviewCard key={review.id} review={review} showReply />
          ))}
        </div>
      )}
    </div>
  )
}
