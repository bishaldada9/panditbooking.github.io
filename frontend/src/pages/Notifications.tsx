import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { notificationsService } from '../services/notifications.service'
import { LoadingSpinner, EmptyState } from '../components/common'
import { Bell, Calendar, Star, CreditCard, Shield, Clock, CheckCheck } from 'lucide-react'
import type { Notification } from '../types'

function getRelativeTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`
  if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`
  if (diffDays === 1) return 'Yesterday'
  if (diffDays < 7) return `${diffDays} days ago`

  const year = date.getFullYear()
  const month = date.toLocaleString('en-US', { month: 'short' })
  const day = date.getDate()
  const thisYear = now.getFullYear()
  return year === thisYear ? `${month} ${day}` : `${month} ${day}, ${year}`
}

const typeIcons: Record<string, React.ReactNode> = {
  booking: <Calendar className="w-5 h-5 text-blue-500" />,
  review: <Star className="w-5 h-5 text-yellow-500" />,
  payment: <CreditCard className="w-5 h-5 text-green-500" />,
  verification: <Shield className="w-5 h-5 text-purple-500" />,
}

function getIcon(type: string): React.ReactNode {
  return typeIcons[type.toLowerCase()] || <Bell className="w-5 h-5 text-gray-500" />
}

export default function Notifications() {
  const navigate = useNavigate()
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)
  const [markingAll, setMarkingAll] = useState(false)

  useEffect(() => {
    fetchNotifications()
  }, [])

  const fetchNotifications = async () => {
    setLoading(true)
    try {
      const data = await notificationsService.getAll()
      setNotifications(data)
    } catch {
      toast.error('Failed to load notifications')
      setNotifications([])
    } finally {
      setLoading(false)
    }
  }

  const handleMarkRead = async (id: string) => {
    try {
      await notificationsService.markRead(id)
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
      )
    } catch {
      toast.error('Failed to mark notification as read')
    }
  }

  const handleMarkAllRead = async () => {
    setMarkingAll(true)
    try {
      await notificationsService.markAllRead()
      setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })))
      toast.success('All notifications marked as read')
    } catch {
      toast.error('Failed to mark all as read')
    } finally {
      setMarkingAll(false)
    }
  }

  const handleClick = (notification: Notification) => {
    if (!notification.is_read) {
      handleMarkRead(notification.id)
    }
    if (notification.reference_id && notification.reference_type) {
      const type = notification.reference_type.toLowerCase()
      if (type === 'booking') {
        navigate(`/bookings/${notification.reference_id}`)
      } else if (type === 'payment') {
        navigate(`/payments`)
      } else if (type === 'review') {
        navigate(`/reviews`)
      } else {
        navigate(`/${type}/${notification.reference_id}`)
      }
    }
  }

  const unreadCount = notifications.filter((n) => !n.is_read).length

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <h1 className="text-3xl font-bold">Notifications</h1>
          {unreadCount > 0 && (
            <span className="bg-primary-500 text-white text-xs font-bold px-2.5 py-1 rounded-full">
              {unreadCount}
            </span>
          )}
        </div>
        {unreadCount > 0 && (
          <button
            onClick={handleMarkAllRead}
            disabled={markingAll}
            className="btn-secondary text-sm flex items-center gap-1.5"
          >
            <CheckCheck className="w-4 h-4" />
            {markingAll ? 'Marking...' : 'Mark All as Read'}
          </button>
        )}
      </div>

      {loading ? (
        <LoadingSpinner size="lg" />
      ) : notifications.length === 0 ? (
        <EmptyState
          title="No notifications yet"
          description="When you get notifications, they'll appear here."
        />
      ) : (
        <div className="space-y-2">
          {notifications.map((notification) => (
            <button
              key={notification.id}
              onClick={() => handleClick(notification)}
              className={`w-full text-left card transition-all duration-200 hover:shadow-md ${
                !notification.is_read
                  ? 'border-l-4 border-l-primary-500 bg-primary-50/50 dark:bg-primary-900/10'
                  : 'border-l-4 border-l-transparent'
              }`}
            >
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0 mt-1">
                  {getIcon(notification.type)}
                </div>
                <div className="flex-1 min-w-0">
                  <p
                    className={`text-sm ${
                      !notification.is_read ? 'font-semibold' : 'font-medium text-gray-700 dark:text-gray-300'
                    }`}
                  >
                    {notification.title}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">
                    {notification.message}
                  </p>
                  <div className="flex items-center gap-1.5 mt-1.5">
                    <Clock className="w-3.5 h-3.5 text-gray-400" />
                    <span className="text-xs text-gray-400 dark:text-gray-500">
                      {getRelativeTime(notification.created_at)}
                    </span>
                    {!notification.is_read && (
                      <span className="w-2 h-2 bg-primary-500 rounded-full" />
                    )}
                  </div>
                </div>
                {notification.reference_id && (
                  <ChevronRight className="w-4 h-4 text-gray-400 flex-shrink-0 mt-1" />
                )}
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}

function ChevronRight({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
    </svg>
  )
}
