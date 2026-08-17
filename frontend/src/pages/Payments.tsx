import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { paymentsService } from '../services/payments.service'
import { bookingsService } from '../services/bookings.service'
import { LoadingSpinner, EmptyState } from '../components/common'
import { CreditCard, Wallet, Banknote, IndianRupee, Calendar, CheckCircle, XCircle, Clock, ExternalLink } from 'lucide-react'
import type { Booking, Payment } from '../types'

interface GatewayOption {
  id: 'esewa' | 'khalti' | 'cash'
  label: string
  icon: React.ReactNode
  description: string
  enabled: boolean
}

const gatewayOptions: GatewayOption[] = [
  {
    id: 'esewa',
    label: 'eSewa',
    icon: <Wallet className="w-8 h-8" />,
    description: 'Pay using eSewa wallet',
    enabled: false,
  },
  {
    id: 'khalti',
    label: 'Khalti',
    icon: <Wallet className="w-8 h-8" />,
    description: 'Pay using Khalti wallet',
    enabled: false,
  },
  {
    id: 'cash',
    label: 'Cash on Site',
    icon: <Banknote className="w-8 h-8" />,
    description: 'Pay the pandit directly at the venue',
    enabled: true,
  },
]

const statusStyles: Record<string, string> = {
  completed: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  failed: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  refunded: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
}

const statusIcons: Record<string, React.ReactNode> = {
  completed: <CheckCircle className="w-4 h-4 text-green-500" />,
  pending: <Clock className="w-4 h-4 text-yellow-500" />,
  failed: <XCircle className="w-4 h-4 text-red-500" />,
  refunded: <ExternalLink className="w-4 h-4 text-purple-500" />,
}

export default function Payments() {
  const [bookings, setBookings] = useState<Booking[]>([])
  const [payments, setPayments] = useState<Payment[]>([])
  const [selectedBookingId, setSelectedBookingId] = useState('')
  const [loadingBookings, setLoadingBookings] = useState(true)
  const [loadingPayments, setLoadingPayments] = useState(true)
  const [paying, setPaying] = useState(false)

  useEffect(() => {
    fetchBookings()
    fetchPayments()
  }, [])

  const fetchBookings = async () => {
    setLoadingBookings(true)
    try {
      const data = await bookingsService.getMyBookings()
      setBookings(data)
    } catch {
      toast.error('Failed to load bookings')
      setBookings([])
    } finally {
      setLoadingBookings(false)
    }
  }

  const fetchPayments = async () => {
    setLoadingPayments(true)
    try {
      const allPayments: Payment[] = []
      const bookingsData = await bookingsService.getMyBookings()
      setBookings(bookingsData)
      for (const b of bookingsData) {
        try {
          const payment = await paymentsService.getByBooking(b.id)
          if (payment) {
            allPayments.push(payment)
          }
        } catch {
          // no payment for this booking
        }
      }
      setPayments(allPayments)
    } catch {
      setPayments([])
    } finally {
      setLoadingPayments(false)
    }
  }

  const handlePay = async (gateway: 'esewa' | 'khalti' | 'cash') => {
    if (!selectedBookingId) {
      toast.error('Please select a booking first')
      return
    }

    if (gateway === 'cash') {
      setPaying(true)
      try {
        const booking = bookings.find((b) => b.id === selectedBookingId)
        await paymentsService.initiate({
          booking_id: selectedBookingId,
          gateway: 'cash',
          amount: booking?.total_amount || 0,
        })
        toast.success('Booking marked as Cash on Site. Pay the pandit directly.')
        fetchPayments()
        setSelectedBookingId('')
      } catch {
        toast.error('Failed to initiate payment')
      } finally {
        setPaying(false)
      }
    } else {
      toast.error('Not available yet. Please use Cash on Site.')
    }
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    if (Number.isNaN(date.getTime())) return '-'
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  const formatTime = (value: string) => {
    if (!value) return '-'
    if (/^\d{2}:\d{2}$/.test(value)) return value
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return value
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const shortText = (value: unknown, length = 8) => {
    const text = String(value ?? '')
    return text ? (text.length > length ? `${text.slice(0, length)}...` : text) : '-'
  }

  const selectedBooking = bookings.find((b) => b.id === selectedBookingId)

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Payments</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">Make payments and view transaction history</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
        <div className="lg:col-span-3 space-y-6">
          <div className="card">
            <h2 className="text-xl font-semibold mb-4">Make a Payment</h2>

            {loadingBookings ? (
              <LoadingSpinner size="sm" />
            ) : bookings.length === 0 ? (
              <p className="text-sm text-gray-500 dark:text-gray-400">
                No bookings found.{' '}
                <a href="/bookings" className="text-primary-600 hover:underline font-medium">
                  Create a booking first
                </a>
              </p>
            ) : (
              <>
                <div className="mb-6">
                  <label className="block text-sm font-medium mb-1">Select Booking</label>
                  <select
                    value={selectedBookingId}
                    onChange={(e) => setSelectedBookingId(e.target.value)}
                    className="input-field"
                  >
                    <option value="">-- Choose a booking --</option>
                    {bookings
                      .filter((b) => b.status !== 'cancelled')
                      .map((b) => (
                        <option key={b.id} value={b.id}>
                          {shortText(b.id)} - ₹{Number(b.total_amount || 0).toLocaleString()} - {formatDate(b.scheduled_date)}
                        </option>
                      ))}
                  </select>
                </div>

                {selectedBooking && (
                  <div className="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium">Booking Details</span>
                      <span className="text-sm text-gray-500 dark:text-gray-400">#{shortText(selectedBooking.id)}</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                      <Calendar className="w-4 h-4" />
                      <span>{formatDate(selectedBooking.scheduled_date)} at {formatTime(selectedBooking.start_time)}</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                      <IndianRupee className="w-4 h-4" />
                      <span className="font-semibold text-lg text-gray-900 dark:text-gray-100">₹{Number(selectedBooking.total_amount || 0).toLocaleString()}</span>
                    </div>
                  </div>
                )}

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  {gatewayOptions.map((gateway) => (
                    <button
                      key={gateway.id}
                      onClick={() => handlePay(gateway.id)}
                      disabled={!gateway.enabled || !selectedBookingId || paying}
                      className={`card-hover text-center p-6 flex flex-col items-center gap-3 transition-all ${
                        !gateway.enabled ? 'opacity-50 cursor-not-allowed hover:shadow-md hover:-translate-y-0' : ''
                      } ${!selectedBookingId ? 'opacity-60' : ''}`}
                    >
                      <div className={`w-14 h-14 rounded-full flex items-center justify-center ${
                        gateway.id === 'cash'
                          ? 'bg-sacred-100 dark:bg-sacred-900 text-sacred-600'
                          : 'bg-gray-100 dark:bg-gray-700 text-gray-500'
                      }`}>
                        {gateway.icon}
                      </div>
                      <div>
                        <div className="font-semibold flex items-center justify-center gap-2">
                          {gateway.label}
                          {!gateway.enabled && (
                            <span className="text-xs bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200 px-2 py-0.5 rounded-full">
                              Coming Soon
                            </span>
                          )}
                        </div>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">{gateway.description}</p>
                      </div>
                    </button>
                  ))}
                </div>
              </>
            )}
          </div>
        </div>

        <div className="lg:col-span-2 space-y-6">
          <div className="card">
            <h2 className="text-xl font-semibold mb-4">Payment History</h2>
            {loadingPayments ? (
              <LoadingSpinner size="sm" />
            ) : payments.length === 0 ? (
              <EmptyState
                title="No payments yet"
                description="Your payment transactions will appear here."
              />
            ) : (
              <div className="space-y-3">
                {payments.map((payment) => (
                  <div key={payment.id} className="p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                    <div className="flex items-center justify-between mb-1">
                      <div className="flex items-center gap-2">
                        <CreditCard className="w-4 h-4 text-gray-500" />
                        <span className="text-sm font-medium capitalize">{payment.gateway}</span>
                      </div>
                      <span className={`text-xs px-2 py-0.5 rounded-full font-medium flex items-center gap-1 ${
                        statusStyles[payment.status] || 'bg-gray-100 text-gray-800'
                      }`}>
                        {statusIcons[payment.status]}
                        {payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between mt-2">
                      <span className="text-lg font-bold">₹{Number(payment.amount || 0).toLocaleString()}</span>
                      <span className="text-xs text-gray-500 dark:text-gray-400">
                        {formatDate(payment.created_at)}
                      </span>
                    </div>
                    {payment.transaction_id && (
                      <p className="text-xs text-gray-400 mt-1">
                        TXN: {payment.transaction_id}
                      </p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
