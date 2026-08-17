import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { useAuth } from '../store/auth.context'
import { bookingsService } from '../services/bookings.service'
import { BookingCard } from '../components/booking'
import { LoadingSpinner, EmptyState } from '../components/common'
import { Calendar } from 'lucide-react'
import type { Booking } from '../types'

export default function Bookings() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [bookings, setBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)

  const isPandit = user?.role === 'pandit'

  useEffect(() => {
    const fetchBookings = async () => {
      try {
        const data = isPandit
          ? await bookingsService.getPanditBookings()
          : await bookingsService.getMyBookings()
        setBookings(data)
      } catch {
        toast.error('Failed to load bookings')
      } finally {
        setLoading(false)
      }
    }
    fetchBookings()
  }, [isPandit])

  if (loading) return <LoadingSpinner size="lg" />

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">
          {isPandit ? 'Booking Requests' : 'My Bookings'}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          {isPandit ? 'Manage your upcoming bookings and requests' : 'View and manage your ritual bookings'}
        </p>
      </div>

      {bookings.length === 0 ? (
        <EmptyState
          title={isPandit ? 'No booking requests yet' : 'No bookings yet'}
          description={isPandit ? 'When customers book your services, they will appear here.' : 'Book a ritual to get started.'}
          action={
            !isPandit ? (
              <button onClick={() => navigate('/rituals')} className="btn-primary inline-flex items-center gap-2">
                <Calendar className="w-4 h-4" />
                Browse Rituals
              </button>
            ) : undefined
          }
        />
      ) : (
        <div className="space-y-4">
          {bookings.map((booking) => (
            <BookingCard
              key={booking.id}
              booking={booking}
              role={isPandit ? 'pandit' : 'customer'}
              onClick={() => navigate(`/bookings/${booking.id}`)}
            />
          ))}
        </div>
      )}
    </div>
  )
}
