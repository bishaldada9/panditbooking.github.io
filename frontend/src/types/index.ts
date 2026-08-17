export interface User {
  id: string
  email: string
  full_name: string
  phone: string
  role: 'customer' | 'pandit' | 'admin'
  is_email_verified: boolean
  mfa_enabled: boolean
  is_active: boolean
  is_suspended: boolean
  created_at?: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}

export interface Pandit {
  id: string
  user_id: string
  full_name: string
  email: string
  phone: string
  bio: string
  experience_years: number
  specialization: string
  languages: string[]
  base_price: number
  rating: number
  total_reviews: number
  total_bookings: number
  is_available: boolean
  verification_status: 'pending' | 'approved' | 'rejected'
  service_area: string
  created_at: string
}

export interface RitualCategory {
  id: string
  name: string
  slug: string
  description: string
  icon: string
  is_active: boolean
}

export interface Ritual {
  id: string
  category_id: string
  category_name: string
  name: string
  slug: string
  description: string
  duration: string
  base_price: number
  required_items: string
  procedure: string
  pandit_commission: number
}

export interface Booking {
  id: string
  customer_id: string
  pandit_id: string
  ritual_id: string
  status: 'pending' | 'confirmed' | 'completed' | 'cancelled' | 'rejected'
  scheduled_date: string
  start_time: string
  end_time: string
  address: string
  total_amount: number
  platform_fee: number
  special_notes: string
  created_at: string
}

export interface Payment {
  id: string
  booking_id: string
  amount: number
  gateway: string
  status: string
  transaction_id: string
  gateway_ref_id: string
  gateway_url?: string
  created_at: string
}

export interface Review {
  id: string
  booking_id: string
  customer_id: string
  pandit_id: string
  rating: number
  comment: string
  is_verified: boolean
  is_visible: boolean
  admin_reply?: string
  created_at: string
}

export interface Notification {
  id: string
  user_id: string
  type: string
  title: string
  message: string
  is_read: boolean
  reference_id?: string
  reference_type?: string
  created_at: string
}

export interface AuditLog {
  id: string
  user_id: string
  action: string
  resource: string
  resource_id: string
  detail: string
  ip: string
  user_agent: string
  status: string
  created_at: string
}

export interface DashboardMetrics {
  total_users: number
  total_pandits: number
  total_bookings: number
  total_revenue: number
  pending_verifications: number
  active_bookings: number
  failed_logins: number
  new_users_today: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
}
