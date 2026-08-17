import { useState, useEffect, useMemo, Fragment } from 'react'
import { useAuth } from '../store/auth.context'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import api from '../services/api'
import { adminService } from '../services/admin.service'
import { ritualsService } from '../services/rituals.service'
import { reviewsService } from '../services/reviews.service'
import { LoadingSpinner, EmptyState, StatCard, DataTable } from '../components/common'
import type { DashboardMetrics, User, Pandit, Booking, Payment, Review, AuditLog, RitualCategory, Ritual } from '../types'
import {
  Users, UserCheck, Calendar, IndianRupee, Shield, Clock, AlertTriangle, UserPlus,
  Search, X, CheckCircle, XCircle, Eye, EyeOff, RefreshCw, Plus, Activity,
  Star, Ban, UserCog, ArrowLeft, Filter, ChevronDown
} from 'lucide-react'

type TabId = 'overview' | 'users' | 'pandits' | 'bookings' | 'payments' | 'reviews' | 'audit' | 'rituals'

const tabs: { id: TabId; label: string; icon: React.ReactNode }[] = [
  { id: 'overview', label: 'Overview', icon: <Activity className="w-4 h-4" /> },
  { id: 'users', label: 'Users', icon: <Users className="w-4 h-4" /> },
  { id: 'pandits', label: 'Pandits', icon: <UserCheck className="w-4 h-4" /> },
  { id: 'bookings', label: 'Bookings', icon: <Calendar className="w-4 h-4" /> },
  { id: 'payments', label: 'Payments', icon: <IndianRupee className="w-4 h-4" /> },
  { id: 'reviews', label: 'Reviews', icon: <Star className="w-4 h-4" /> },
  { id: 'audit', label: 'Audit Logs', icon: <Activity className="w-4 h-4" /> },
  { id: 'rituals', label: 'Categories/Rituals', icon: <Shield className="w-4 h-4" /> },
]

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400',
  confirmed: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
  completed: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  cancelled: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
  rejected: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
  active: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  suspended: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
  approved: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  paid: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  failed: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
  refunded: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
  visible: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  hidden: 'bg-gray-100 text-gray-700 dark:bg-gray-900/40 dark:text-gray-400',
}

function Badge({ status }: { status: string }) {
  const color = statusColors[status.toLowerCase()] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${color}`}>
      {status}
    </span>
  )
}

function Stars({ rating }: { rating: number }) {
  const safeRating = Math.min(5, Math.max(0, Number.isFinite(rating) ? rating : 0))
  return (
    <span className="text-yellow-500">{'★'.repeat(Math.round(safeRating))}{'☆'.repeat(5 - Math.round(safeRating))}</span>
  )
}

function shortText(value: unknown, length = 8) {
  const text = String(value ?? '')
  if (!text) return '-'
  return text.length > length ? `${text.slice(0, length)}...` : text
}

function formatDate(value: unknown) {
  if (!value) return '-'
  const date = new Date(String(value))
  return Number.isNaN(date.getTime()) ? '-' : date.toLocaleDateString()
}

function formatDateTime(value: unknown) {
  if (!value) return '-'
  const date = new Date(String(value))
  return Number.isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

function formatMoney(value: unknown) {
  const amount = Number(value)
  return `₹${Number.isFinite(amount) ? amount.toLocaleString() : '0'}`
}

function SearchInput({ value, onChange, placeholder }: { value: string; onChange: (v: string) => void; placeholder?: string }) {
  return (
    <div className="relative w-full max-w-sm">
      <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder || 'Search...'}
        className="input-field pl-10 pr-10"
      />
      {value && (
        <button onClick={() => onChange('')} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
          <X className="w-4 h-4" />
        </button>
      )}
    </div>
  )
}

function TableActions({ children }: { children: React.ReactNode }) {
  return <div className="flex items-center gap-2">{children}</div>
}

export default function AdminDashboard() {
  const { user, isLoading: authLoading } = useAuth()
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<TabId>('overview')

  if (authLoading) return <LoadingSpinner />
  if (!user || user.role !== 'admin') {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 text-center">
        <Ban className="w-16 h-16 text-red-400 mx-auto mb-4" />
        <h1 className="text-2xl font-bold text-gray-800 dark:text-gray-200 mb-2">Access Denied</h1>
        <p className="text-gray-500 dark:text-gray-400 mb-6">You do not have admin privileges.</p>
        <button onClick={() => navigate('/')} className="btn-primary">Return Home</button>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Admin Dashboard</h1>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Manage your platform — users, pandits, bookings, and more</p>
          </div>
        </div>

        <div className="flex overflow-x-auto gap-1 mb-6 border-b border-gray-200 dark:border-gray-800">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex items-center gap-2 px-4 py-3 text-sm font-medium whitespace-nowrap border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
              }`}
            >
              {tab.icon}
              {tab.label}
            </button>
          ))}
        </div>

        {activeTab === 'overview' && <OverviewTab />}
        {activeTab === 'users' && <UsersTab />}
        {activeTab === 'pandits' && <PanditsTab />}
        {activeTab === 'bookings' && <BookingsTab />}
        {activeTab === 'payments' && <PaymentsTab />}
        {activeTab === 'reviews' && <ReviewsTab />}
        {activeTab === 'audit' && <AuditLogsTab />}
        {activeTab === 'rituals' && <RitualsTab />}
      </div>
    </div>
  )
}

function OverviewTab() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchMetrics = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminService.getDashboard()
      setMetrics(data)
    } catch (err: any) {
      setError(err?.message || 'Failed to load dashboard metrics')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchMetrics() }, [])

  if (loading) return <LoadingSpinner />
  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchMetrics} className="btn-secondary">Retry</button>
      </div>
    )
  }
  if (!metrics) return <EmptyState title="No data available" />

  const cards = [
    { icon: <Users className="w-6 h-6" />, label: 'Total Users', value: metrics.total_users, color: 'primary' as const },
    { icon: <UserCheck className="w-6 h-6" />, label: 'Total Pandits', value: metrics.total_pandits, color: 'saffron' as const },
    { icon: <Calendar className="w-6 h-6" />, label: 'Total Bookings', value: metrics.total_bookings, color: 'sacred' as const },
    { icon: <IndianRupee className="w-6 h-6" />, label: 'Total Revenue', value: `₹${metrics.total_revenue.toLocaleString()}`, color: 'primary' as const },
    { icon: <Shield className="w-6 h-6" />, label: 'Pending Verifications', value: metrics.pending_verifications, color: 'saffron' as const },
    { icon: <Clock className="w-6 h-6" />, label: 'Active Bookings', value: metrics.active_bookings, color: 'sacred' as const },
    { icon: <AlertTriangle className="w-6 h-6" />, label: 'Failed Logins', value: metrics.failed_logins, color: 'red' as const },
    { icon: <UserPlus className="w-6 h-6" />, label: 'New Users Today', value: metrics.new_users_today, color: 'sacred' as const },
  ]

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {cards.map((card, i) => (
        <StatCard key={i} {...card} />
      ))}
    </div>
  )
}

function UsersTab() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const fetchUsers = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminService.getUsers()
      setUsers(data)
    } catch (err: any) {
      setError(err?.message || 'Failed to load users')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchUsers() }, [])

  const filtered = useMemo(() => {
    if (!search.trim()) return users
    const term = search.toLowerCase()
    return users.filter((u) =>
      Object.values(u).some((val) => String(val ?? '').toLowerCase().includes(term))
    )
  }, [users, search])

  const handleSuspend = async (id: string, isSuspended: boolean) => {
    const reason = isSuspended ? '' : window.prompt('Enter suspension reason:', 'Policy or verification issue')?.trim()
    if (!isSuspended && !reason) return
    const confirmed = window.confirm(
      isSuspended ? 'Are you sure you want to activate this user?' : 'Are you sure you want to suspend this user?'
    )
    if (!confirmed) return
    setActionLoading(id)
    try {
      if (isSuspended) {
        await adminService.activateUser(id)
        toast.success('User activated successfully')
      } else {
        await adminService.suspendUser(id, reason || 'Suspended by admin')
        toast.success('User suspended successfully')
      }
      setUsers((prev) => prev.map((u) => (u.id === id ? { ...u, is_suspended: !isSuspended } : u)))
    } catch (err: any) {
      toast.error(err?.message || 'Failed to update user')
    } finally {
      setActionLoading(null)
    }
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchUsers} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-4">
        <SearchInput value={search} onChange={setSearch} placeholder="Search by name, email, role..." />
      </div>
      <DataTable
        isLoading={loading}
        columns={[
          { key: 'full_name', header: 'Name', render: (u: User) => (
            <div>
              <p className="font-medium">{u.full_name}</p>
              <p className="text-xs text-gray-500">{u.email}</p>
            </div>
          )},
          { key: 'phone', header: 'Phone' },
          { key: 'role', header: 'Role', render: (u: User) => <Badge status={u.role} /> },
          { key: 'is_active', header: 'Status', render: (u: User) => (
            u.is_suspended ? <Badge status="suspended" /> : <Badge status="active" />
          )},
          { key: 'created_at', header: 'Joined', render: (u: User) => u.created_at ? new Date(u.created_at).toLocaleDateString() : '-' },
          { key: 'actions', header: 'Actions', render: (u: User) => (
            <button
              onClick={() => handleSuspend(u.id, u.is_suspended)}
              disabled={actionLoading === u.id}
              className={`text-sm font-medium ${u.is_suspended ? 'text-green-600 hover:text-green-700' : 'text-red-600 hover:text-red-700'} disabled:opacity-50`}
            >
              {actionLoading === u.id ? '...' : u.is_suspended ? 'Activate' : 'Suspend'}
            </button>
          )},
        ]}
        data={filtered}
        emptyMessage={search ? 'No users match your search.' : 'No users found.'}
      />
    </div>
  )
}

function PanditsTab() {
  const [pandits, setPandits] = useState<Pandit[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const fetchPandits = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminService.getPandits()
      setPandits(data)
    } catch (err: any) {
      setError(err?.message || 'Failed to load pandits')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchPandits() }, [])

  const filtered = useMemo(() => {
    if (!search.trim()) return pandits
    const term = search.toLowerCase()
    return pandits.filter((p) =>
      Object.values(p).some((val) => {
        if (Array.isArray(val)) return val.some((v) => String(v).toLowerCase().includes(term))
        return String(val ?? '').toLowerCase().includes(term)
      })
    )
  }, [pandits, search])

  const handleVerificationStatus = async (id: string, status: 'pending' | 'approved' | 'rejected') => {
    const notes = status === 'approved'
      ? 'Approved by admin'
      : status === 'rejected'
        ? window.prompt('Enter rejection note:', 'Verification requirements not met') || 'Rejected by admin'
        : 'Moved back to pending by admin'
    setActionLoading(id)
    try {
      await adminService.verifyPandit(id, status, notes)
      toast.success(`Pandit marked as ${status}`)
      setPandits((prev) => prev.map((p) => (p.id === id ? { ...p, verification_status: status } : p)))
    } catch (err: any) {
      toast.error(err?.message || 'Failed to update pandit verification')
    } finally {
      setActionLoading(null)
    }
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchPandits} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-4">
        <SearchInput value={search} onChange={setSearch} placeholder="Search pandits..." />
      </div>
      <DataTable
        isLoading={loading}
        columns={[
          { key: 'full_name', header: 'Name', render: (p: Pandit) => (
            <div>
              <p className="font-medium">{p.full_name}</p>
              <p className="text-xs text-gray-500">{p.email}</p>
            </div>
          )},
          { key: 'specialization', header: 'Specialization' },
          { key: 'experience_years', header: 'Exp', render: (p: Pandit) => `${p.experience_years}y` },
          { key: 'rating', header: 'Rating', render: (p: Pandit) => (
            <div className="flex items-center gap-1">
              <Stars rating={p.rating} />
              <span className="text-xs text-gray-500">({p.total_reviews})</span>
            </div>
          )},
          { key: 'verification_status', header: 'Status', render: (p: Pandit) => <Badge status={p.verification_status} /> },
          { key: 'actions', header: 'Actions', render: (p: Pandit) => (
            <TableActions>
              {p.verification_status === 'pending' && (
                <>
                  <button
                    onClick={() => handleVerificationStatus(p.id, 'approved')}
                    disabled={actionLoading === p.id}
                    className="btn-primary text-xs px-3 py-1 disabled:opacity-50"
                  >
                    {actionLoading === p.id ? '...' : 'Approve'}
                  </button>
                  <button
                    onClick={() => handleVerificationStatus(p.id, 'rejected')}
                    disabled={actionLoading === p.id}
                    className="btn-danger text-xs px-3 py-1 disabled:opacity-50"
                  >
                    {actionLoading === p.id ? '...' : 'Reject'}
                  </button>
                </>
              )}
              {p.verification_status === 'approved' && (
                <button
                  onClick={() => handleVerificationStatus(p.id, 'pending')}
                  disabled={actionLoading === p.id}
                  className="btn-secondary text-xs px-3 py-1 disabled:opacity-50"
                >
                  {actionLoading === p.id ? '...' : 'Move to Pending'}
                </button>
              )}
              {p.verification_status === 'rejected' && (
                <button
                  onClick={() => handleVerificationStatus(p.id, 'pending')}
                  disabled={actionLoading === p.id}
                  className="btn-secondary text-xs px-3 py-1 disabled:opacity-50"
                >
                  {actionLoading === p.id ? '...' : 'Review Again'}
                </button>
              )}
            </TableActions>
          )},
        ]}
        data={filtered}
        emptyMessage={search ? 'No pandits match your search.' : 'No pandit profiles found.'}
      />
    </div>
  )
}

function BookingsTab() {
  const [bookings, setBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchBookings = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await api.get('/admin/bookings')
      setBookings(Array.isArray(res.data.data) ? res.data.data : [])
    } catch (err: any) {
      setError(err?.message || 'Failed to load bookings')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchBookings() }, [])

  const filtered = useMemo(() => {
    if (!search.trim()) return bookings
    const term = search.toLowerCase()
    return bookings.filter((b) =>
      Object.values(b).some((val) => String(val ?? '').toLowerCase().includes(term))
    )
  }, [bookings, search])

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchBookings} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-4">
        <SearchInput value={search} onChange={setSearch} placeholder="Search by ID or status..." />
      </div>
      <DataTable
        isLoading={loading}
        columns={[
          { key: 'id', header: 'Booking ID', render: (b: Booking) => (
            <code className="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">{shortText(b.id)}</code>
          )},
          { key: 'customer_id', header: 'Customer', render: (b: Booking) => (
            <span className="text-xs text-gray-500">{shortText(b.customer_id)}</span>
          )},
          { key: 'pandit_id', header: 'Pandit', render: (b: Booking) => (
            <span className="text-xs text-gray-500">{shortText(b.pandit_id)}</span>
          )},
          { key: 'ritual_id', header: 'Ritual', render: (b: Booking) => shortText(b.ritual_id) },
          { key: 'status', header: 'Status', render: (b: Booking) => <Badge status={b.status} /> },
          { key: 'scheduled_date', header: 'Date', render: (b: Booking) => formatDate(b.scheduled_date) },
          { key: 'total_amount', header: 'Amount', render: (b: Booking) => formatMoney(b.total_amount) },
        ]}
        data={filtered}
        emptyMessage={search ? 'No bookings match your search.' : 'No bookings found yet.'}
      />
    </div>
  )
}

function PaymentsTab() {
  const [payments, setPayments] = useState<Payment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [refundLoading, setRefundLoading] = useState<string | null>(null)

  const fetchPayments = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminService.getPayments()
      setPayments(data)
    } catch (err: any) {
      setError(err?.message || 'Failed to load payments')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchPayments() }, [])

  const filtered = useMemo(() => {
    if (!search.trim()) return payments
    const term = search.toLowerCase()
    return payments.filter((p) =>
      Object.values(p).some((val) => String(val ?? '').toLowerCase().includes(term))
    )
  }, [payments, search])

  const handleRefund = async (payment: Payment) => {
    const reason = window.prompt('Enter refund reason:')
    if (!reason) return
    if (!window.confirm(`Refund ${formatMoney(payment.amount)} for transaction ${shortText(payment.transaction_id, 12)}?`)) return
    setRefundLoading(payment.id)
    try {
      await adminService.refund({ payment_id: payment.id, amount: payment.amount, reason })
      toast.success('Refund processed successfully')
      setPayments((prev) => prev.map((p) => (p.id === payment.id ? { ...p, status: 'refunded' } : p)))
    } catch (err: any) {
      toast.error(err?.message || 'Refund failed')
    } finally {
      setRefundLoading(null)
    }
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchPayments} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-4">
        <SearchInput value={search} onChange={setSearch} placeholder="Search transactions..." />
      </div>
      <DataTable
        isLoading={loading}
        columns={[
          { key: 'transaction_id', header: 'Transaction ID', render: (p: Payment) => (
            <code className="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">{shortText(p.transaction_id, 12)}</code>
          )},
          { key: 'booking_id', header: 'Booking', render: (p: Payment) => (
            <code className="text-xs">{shortText(p.booking_id)}</code>
          )},
          { key: 'gateway', header: 'Gateway' },
          { key: 'amount', header: 'Amount', render: (p: Payment) => formatMoney(p.amount) },
          { key: 'status', header: 'Status', render: (p: Payment) => <Badge status={p.status} /> },
          { key: 'created_at', header: 'Date', render: (p: Payment) => formatDate(p.created_at) },
          { key: 'actions', header: 'Actions', render: (p: Payment) => (
            p.status === 'paid' || p.status === 'completed' ? (
              <button
                onClick={() => handleRefund(p)}
                disabled={refundLoading === p.id}
                className="btn-danger text-xs px-3 py-1 disabled:opacity-50"
              >
                {refundLoading === p.id ? '...' : 'Refund'}
              </button>
            ) : (
              <span className="text-xs text-gray-400">—</span>
            )
          )},
        ]}
        data={filtered}
        emptyMessage={search ? 'No payments match your search.' : 'No payments found yet.'}
      />
    </div>
  )
}

function ReviewsTab() {
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const fetchReviews = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await api.get('/admin/reviews')
      setReviews(Array.isArray(res.data.data) ? res.data.data : [])
    } catch (err: any) {
      setError(err?.message || 'Failed to load reviews')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchReviews() }, [])

  const filtered = useMemo(() => {
    if (!search.trim()) return reviews
    const term = search.toLowerCase()
    return reviews.filter((r) =>
      Object.values(r).some((val) => String(val ?? '').toLowerCase().includes(term))
    )
  }, [reviews, search])

  const handleModerate = async (id: string, action: 'approve' | 'reject') => {
    setActionLoading(id)
    try {
      await reviewsService.moderate(id, {
        is_visible: action === 'approve',
        admin_reply: action === 'approve' ? 'Approved by admin' : 'Hidden by admin',
      })
      toast.success(action === 'approve' ? 'Review approved' : 'Review hidden')
      setReviews((prev) =>
        prev.map((r) =>
          r.id === id ? { ...r, is_visible: action === 'approve', is_verified: action === 'approve' } : r
        )
      )
    } catch (err: any) {
      toast.error(err?.message || 'Failed to moderate review')
    } finally {
      setActionLoading(null)
    }
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchReviews} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-4">
        <SearchInput value={search} onChange={setSearch} placeholder="Search reviews..." />
      </div>
      <DataTable
        isLoading={loading}
        columns={[
          { key: 'id', header: 'Review ID', render: (r: Review) => (
            <code className="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">{shortText(r.id)}</code>
          )},
          { key: 'rating', header: 'Rating', render: (r: Review) => <Stars rating={r.rating} /> },
          { key: 'comment', header: 'Comment', render: (r: Review) => (
            <span className="text-sm truncate max-w-[200px] inline-block">{r.comment || '—'}</span>
          )},
          { key: 'is_visible', header: 'Status', render: (r: Review) => (
            <Badge status={r.is_visible ? 'visible' : 'hidden'} />
          )},
          { key: 'actions', header: 'Actions', render: (r: Review) => (
            <TableActions>
              {!r.is_visible && (
                <button
                  onClick={() => handleModerate(r.id, 'approve')}
                  disabled={actionLoading === r.id}
                  className="btn-primary text-xs px-3 py-1 disabled:opacity-50"
                >
                  {actionLoading === r.id ? '...' : <><CheckCircle className="w-3 h-3 inline mr-1" />Approve</>}
                </button>
              )}
              {r.is_visible && (
                <button
                  onClick={() => handleModerate(r.id, 'reject')}
                  disabled={actionLoading === r.id}
                  className="btn-danger text-xs px-3 py-1 disabled:opacity-50"
                >
                  {actionLoading === r.id ? '...' : <><EyeOff className="w-3 h-3 inline mr-1" />Hide</>}
                </button>
              )}
            </TableActions>
          )},
        ]}
        data={filtered}
        emptyMessage={search ? 'No reviews match your search.' : 'No reviews found yet.'}
      />
    </div>
  )
}

function AuditLogsTab() {
  const [logs, setLogs] = useState<AuditLog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchLogs = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminService.getAuditLogs()
      setLogs(data)
    } catch (err: any) {
      setError(err?.message || 'Failed to load audit logs')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchLogs() }, [])

  const filtered = useMemo(() => {
    if (!search.trim()) return logs
    const term = search.toLowerCase()
    return logs.filter((l) =>
      [l.action, l.resource, l.resource_id, l.user_id, l.detail, l.ip, l.status]
        .some((val) => String(val ?? '').toLowerCase().includes(term))
    )
  }, [logs, search])

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchLogs} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-4">
        <SearchInput value={search} onChange={setSearch} placeholder="Search by action, resource, user..." />
      </div>
      <DataTable
        isLoading={loading}
        columns={[
          { key: 'action', header: 'Action', render: (l: AuditLog) => (
            <code className="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">{l.action}</code>
          )},
          { key: 'resource', header: 'Resource', render: (l: AuditLog) => <span className="text-sm">{l.resource}</span> },
          { key: 'resource_id', header: 'Resource ID', render: (l: AuditLog) => (
            <code className="text-xs">{shortText(l.resource_id)}</code>
          )},
          { key: 'user_id', header: 'User ID', render: (l: AuditLog) => (
            <code className="text-xs">{shortText(l.user_id)}</code>
          )},
          { key: 'detail', header: 'Detail', render: (l: AuditLog) => (
            <span className="text-xs truncate max-w-[150px] inline-block">{l.detail || '—'}</span>
          )},
          { key: 'ip', header: 'IP', render: (l: AuditLog) => <code className="text-xs">{l.ip}</code> },
          { key: 'status', header: 'Status', render: (l: AuditLog) => <Badge status={l.status} /> },
          { key: 'created_at', header: 'Timestamp', render: (l: AuditLog) => (
            <span className="text-xs whitespace-nowrap">{formatDateTime(l.created_at)}</span>
          )},
        ]}
        data={filtered}
        emptyMessage={search ? 'No audit logs match your search.' : 'No audit logs have been recorded yet.'}
      />
    </div>
  )
}

function RitualsTab() {
  const [categories, setCategories] = useState<RitualCategory[]>([])
  const [rituals, setRituals] = useState<Ritual[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [showCategoryForm, setShowCategoryForm] = useState(false)
  const [showRitualForm, setShowRitualForm] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  const [catForm, setCatForm] = useState({ name: '', description: '', icon: '' })
  const [ritForm, setRitForm] = useState({ category_id: '', name: '', description: '', duration: '', base_price: 0, required_items: '', procedure: '' })

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [cats, rits] = await Promise.all([
        ritualsService.getCategories(),
        ritualsService.getAll(),
      ])
      setCategories(cats)
      setRituals(rits)
    } catch (err: any) {
      setError(err?.message || 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const filteredRituals = useMemo(() => {
    if (!search.trim()) return rituals
    const term = search.toLowerCase()
    return rituals.filter((r) =>
      Object.values(r).some((val) => String(val ?? '').toLowerCase().includes(term))
    )
  }, [rituals, search])

  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!catForm.name.trim()) { toast.error('Name is required'); return }
    setSubmitting(true)
    try {
      await adminService.createCategory(catForm)
      toast.success('Category created successfully')
      setShowCategoryForm(false)
      setCatForm({ name: '', description: '', icon: '' })
      const cats = await ritualsService.getCategories()
      setCategories(cats)
    } catch (err: any) {
      toast.error(err?.message || 'Failed to create category')
    } finally {
      setSubmitting(false)
    }
  }

  const handleCreateRitual = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!ritForm.name.trim() || !ritForm.category_id) { toast.error('Name and category are required'); return }
    setSubmitting(true)
    try {
      await adminService.createRitual(ritForm)
      toast.success('Ritual created successfully')
      setShowRitualForm(false)
      setRitForm({ category_id: '', name: '', description: '', duration: '', base_price: 0, required_items: '', procedure: '' })
      const rits = await ritualsService.getAll()
      setRituals(rits)
    } catch (err: any) {
      toast.error(err?.message || 'Failed to create ritual')
    } finally {
      setSubmitting(false)
    }
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-500 mb-4">{error}</p>
        <button onClick={fetchData} className="btn-secondary">Retry</button>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Categories Section */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200">Categories</h2>
          <button onClick={() => setShowCategoryForm(true)} className="btn-primary text-sm flex items-center gap-1">
            <Plus className="w-4 h-4" /> Add Category
          </button>
        </div>
        <DataTable
          isLoading={loading}
          columns={[
            { key: 'name', header: 'Name' },
            { key: 'slug', header: 'Slug', render: (c: RitualCategory) => (
              <code className="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">{c.slug}</code>
            )},
            { key: 'description', header: 'Description', render: (c: RitualCategory) => (
              <span className="text-sm truncate max-w-[250px] inline-block">{c.description || '—'}</span>
            )},
            { key: 'is_active', header: 'Status', render: (c: RitualCategory) => (
              <Badge status={c.is_active ? 'active' : 'hidden'} />
            )},
          ]}
          data={categories}
          emptyMessage="No ritual categories found."
        />
      </div>

      {/* Rituals Section */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200">Rituals</h2>
          <button onClick={() => setShowRitualForm(true)} className="btn-primary text-sm flex items-center gap-1">
            <Plus className="w-4 h-4" /> Add Ritual
          </button>
        </div>
        <div className="mb-4">
          <SearchInput value={search} onChange={setSearch} placeholder="Search rituals..." />
        </div>
        <DataTable
          isLoading={loading}
          columns={[
            { key: 'name', header: 'Name' },
            { key: 'category_name', header: 'Category' },
            { key: 'base_price', header: 'Price', render: (r: Ritual) => `₹${r.base_price}` },
            { key: 'duration', header: 'Duration', render: (r: Ritual) => r.duration },
          ]}
          data={filteredRituals}
          emptyMessage={search ? 'No rituals match your search.' : 'No rituals found.'}
        />
      </div>

      {/* Category Modal */}
      {showCategoryForm && (
        <Modal onClose={() => setShowCategoryForm(false)} title="Add Category">
          <form onSubmit={handleCreateCategory} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">Name *</label>
              <input
                type="text"
                value={catForm.name}
                onChange={(e) => setCatForm((f) => ({ ...f, name: e.target.value }))}
                className="input-field"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Description</label>
              <textarea
                value={catForm.description}
                onChange={(e) => setCatForm((f) => ({ ...f, description: e.target.value }))}
                className="input-field"
                rows={3}
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Icon</label>
              <input
                type="text"
                value={catForm.icon}
                onChange={(e) => setCatForm((f) => ({ ...f, icon: e.target.value }))}
                className="input-field"
                placeholder="e.g. om, fire, flower"
              />
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <button type="button" onClick={() => setShowCategoryForm(false)} className="btn-secondary">Cancel</button>
              <button type="submit" disabled={submitting} className="btn-primary disabled:opacity-50">
                {submitting ? 'Creating...' : 'Create'}
              </button>
            </div>
          </form>
        </Modal>
      )}

      {/* Ritual Modal */}
      {showRitualForm && (
        <Modal onClose={() => setShowRitualForm(false)} title="Add Ritual">
          <form onSubmit={handleCreateRitual} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">Category *</label>
              <select
                value={ritForm.category_id}
                onChange={(e) => setRitForm((f) => ({ ...f, category_id: e.target.value }))}
                className="input-field"
                required
              >
                <option value="">Select category</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Name *</label>
              <input
                type="text"
                value={ritForm.name}
                onChange={(e) => setRitForm((f) => ({ ...f, name: e.target.value }))}
                className="input-field"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Description</label>
              <textarea
                value={ritForm.description}
                onChange={(e) => setRitForm((f) => ({ ...f, description: e.target.value }))}
                className="input-field"
                rows={3}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Duration</label>
                <input
                  type="text"
                  value={ritForm.duration}
                  onChange={(e) => setRitForm((f) => ({ ...f, duration: e.target.value }))}
                  className="input-field"
                  placeholder="e.g. 2 hours"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Base Price *</label>
                <input
                  type="number"
                  value={ritForm.base_price || ''}
                  onChange={(e) => setRitForm((f) => ({ ...f, base_price: Number(e.target.value) }))}
                  className="input-field"
                  min={0}
                  required
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Required Items</label>
              <textarea
                value={ritForm.required_items}
                onChange={(e) => setRitForm((f) => ({ ...f, required_items: e.target.value }))}
                className="input-field"
                rows={2}
                placeholder="Comma separated or list"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Procedure</label>
              <textarea
                value={ritForm.procedure}
                onChange={(e) => setRitForm((f) => ({ ...f, procedure: e.target.value }))}
                className="input-field"
                rows={3}
              />
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <button type="button" onClick={() => setShowRitualForm(false)} className="btn-secondary">Cancel</button>
              <button type="submit" disabled={submitting} className="btn-primary disabled:opacity-50">
                {submitting ? 'Creating...' : 'Create'}
              </button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}

function Modal({ title, children, onClose }: { title: string; children: React.ReactNode; onClose: () => void }) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div
        className="bg-white dark:bg-gray-900 rounded-xl shadow-2xl w-full max-w-lg mx-4 max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between p-5 border-b border-gray-200 dark:border-gray-800">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">{title}</h3>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <X className="w-5 h-5" />
          </button>
        </div>
        <div className="p-5">{children}</div>
      </div>
    </div>
  )
}
