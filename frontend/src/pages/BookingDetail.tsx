import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { useAuth } from '../store/auth.context'
import { bookingsService } from '../services/bookings.service'
import { ritualsService } from '../services/rituals.service'
import { panditsService } from '../services/pandits.service'
import { LoadingSpinner, EmptyState } from '../components/common'
import { Calendar, Clock, MapPin, IndianRupee, CheckCircle, XCircle, ArrowLeft, MessageSquare, FileText } from 'lucide-react'
import type { Booking, Ritual, Pandit } from '../types'

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

function formatDate(value: unknown, withWeekday = false) {
  if (!value) return '-'
  const date = new Date(String(value))
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('en-IN', {
    ...(withWeekday ? { weekday: 'long' as const } : {}),
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

export default function BookingDetail() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const navigate = useNavigate()

  const [booking, setBooking] = useState<Booking | null>(null)
  const [ritual, setRitual] = useState<Ritual | null>(null)
  const [pandit, setPandit] = useState<Pandit | null>(null)
  const [loading, setLoading] = useState(true)
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const isPandit = user?.role === 'pandit'
  const isCustomer = user?.role === 'customer'

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await bookingsService.getById(id!)
        setBooking(data)

        const [rit, pan] = await Promise.allSettled([
          ritualsService.getById(data.ritual_id),
          panditsService.getById(data.pandit_id),
        ])
        if (rit.status === 'fulfilled') setRitual(rit.value)
        if (pan.status === 'fulfilled') setPandit(pan.value)
      } catch {
        toast.error('Failed to load booking details')
      } finally {
        setLoading(false)
      }
    }
    if (id) fetchData()
  }, [id])

  const nextStatus: Record<string, Booking['status']> = {
    cancel: 'cancelled',
    reject: 'rejected',
    confirm: 'confirmed',
    complete: 'completed',
  }

  const handleAction = async (action: string, handler: () => Promise<unknown>) => {
    setActionLoading(action)
    try {
      await handler()
      setBooking((current) => current ? { ...current, status: nextStatus[action] || current.status } : current)
      toast.success(`Booking ${action} successfully`)
    } catch {
      toast.error(`Failed to ${action} booking`)
    } finally {
      setActionLoading(null)
    }
  }

  if (loading) return <LoadingSpinner size="lg" />

  if (!booking) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <EmptyState title="Booking not found" description="This booking may have been removed or the link is invalid." />
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <button
        onClick={() => navigate('/bookings')}
        className="inline-flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 mb-6 transition-colors"
      >
        <ArrowLeft className="w-4 h-4" />
        Back to Bookings
      </button>

      <div className="card">
        <div className="flex items-start justify-between mb-6">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium capitalize ${statusStyles[booking.status] || statusStyles.pending}`}>
                {booking.status}
              </span>
              <span className="text-sm text-gray-400 dark:text-gray-500 font-mono">#{shortText(booking.id)}</span>
            </div>
            <h1 className="text-2xl font-bold">{ritual?.name || 'Loading...'}</h1>
            {pandit && (
              <p className="text-gray-600 dark:text-gray-400 mt-1">
                Pandit: {pandit.full_name}
              </p>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-4">
            <div className="flex items-center gap-3 text-sm">
              <Calendar className="w-5 h-5 text-primary-500 shrink-0" />
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400">Date</p>
                <p className="font-medium">{formatDate(booking.scheduled_date, true)}</p>
              </div>
            </div>

            <div className="flex items-center gap-3 text-sm">
              <Clock className="w-5 h-5 text-primary-500 shrink-0" />
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400">Time</p>
                <p className="font-medium">{booking.start_time || '-'} - {booking.end_time || '-'}</p>
              </div>
            </div>

            <div className="flex items-start gap-3 text-sm">
              <MapPin className="w-5 h-5 text-primary-500 shrink-0 mt-0.5" />
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400">Address</p>
                <p className="font-medium">{booking.address || '-'}</p>
              </div>
            </div>

            {ritual && (
              <div className="flex items-start gap-3 text-sm">
                <FileText className="w-5 h-5 text-primary-500 shrink-0 mt-0.5" />
                <div>
                  <p className="text-xs text-gray-500 dark:text-gray-400">Ritual Details</p>
                  <p className="font-medium">{ritual.name}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{ritual.category_name} &middot; {ritual.duration}</p>
                </div>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <div className="flex items-center gap-3 text-sm">
              <IndianRupee className="w-5 h-5 text-primary-500 shrink-0" />
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400">Amount Breakdown</p>
                <p className="font-medium">Total: ₹{Number(booking.total_amount || 0).toLocaleString()}</p>
                <p className="text-xs text-gray-500 dark:text-gray-400">Platform fee: ₹{Number(booking.platform_fee || 0).toLocaleString()}</p>
              </div>
            </div>

            {booking.special_notes && (
              <div className="flex items-start gap-3 text-sm">
                <MessageSquare className="w-5 h-5 text-primary-500 shrink-0 mt-0.5" />
                <div>
                  <p className="text-xs text-gray-500 dark:text-gray-400">Special Notes</p>
                  <p className="font-medium">{booking.special_notes}</p>
                </div>
              </div>
            )}
          </div>
        </div>

        {(isCustomer || isPandit) && (
          <div className="mt-8 pt-6 border-t border-gray-100 dark:border-gray-700">
            <div className="flex flex-wrap gap-3">
              {isCustomer && (booking.status === 'pending' || booking.status === 'confirmed') && (
                <button
                  onClick={() => handleAction('cancel', () => bookingsService.cancel(booking.id))}
                  disabled={actionLoading === 'cancel'}
                  className="btn-danger inline-flex items-center gap-2"
                >
                  <XCircle className="w-4 h-4" />
                  {actionLoading === 'cancel' ? 'Cancelling...' : 'Cancel Booking'}
                </button>
              )}

              {isPandit && booking.status === 'pending' && (
                <>
                  <button
                    onClick={() => handleAction('confirm', () => bookingsService.confirm(booking.id))}
                    disabled={actionLoading === 'confirm'}
                    className="btn-primary inline-flex items-center gap-2"
                  >
                    <CheckCircle className="w-4 h-4" />
                    {actionLoading === 'confirm' ? 'Confirming...' : 'Confirm Booking'}
                  </button>
                  <button
                    onClick={() => handleAction('reject', () => bookingsService.reject(booking.id))}
                    disabled={actionLoading === 'reject'}
                    className="btn-danger inline-flex items-center gap-2"
                  >
                    <XCircle className="w-4 h-4" />
                    {actionLoading === 'reject' ? 'Rejecting...' : 'Reject'}
                  </button>
                </>
              )}

              {isPandit && booking.status === 'confirmed' && (
                <button
                  onClick={() => handleAction('complete', () => bookingsService.complete(booking.id))}
                  disabled={actionLoading === 'complete'}
                  className="btn-primary inline-flex items-center gap-2"
                >
                  <CheckCircle className="w-4 h-4" />
                  {actionLoading === 'complete' ? 'Completing...' : 'Mark Completed'}
                </button>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
