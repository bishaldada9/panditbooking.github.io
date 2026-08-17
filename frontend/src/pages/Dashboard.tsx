import { useAuth } from '../store/auth.context'
import { Link } from 'react-router-dom'
import { Calendar, UserCircle, CreditCard, Star, Clock, Shield } from 'lucide-react'

export default function Dashboard() {
  const { user } = useAuth()

  const customerLinks = [
    { icon: <Calendar className="w-6 h-6" />, title: 'My Bookings', path: '/bookings', desc: 'View and manage bookings' },
    { icon: <UserCircle className="w-6 h-6" />, title: 'Find Pandit', path: '/pandits', desc: 'Search verified pandits' },
    { icon: <CreditCard className="w-6 h-6" />, title: 'Payments', path: '/payments', desc: 'Payment history and receipts' },
    { icon: <Star className="w-6 h-6" />, title: 'My Reviews', path: '/reviews', desc: 'Reviews you have written' },
  ]

  const panditLinks = [
    { icon: <Calendar className="w-6 h-6" />, title: 'My Bookings', path: '/pandit/bookings', desc: 'Manage booking requests' },
    { icon: <Clock className="w-6 h-6" />, title: 'Availability', path: '/pandit/availability', desc: 'Set your schedule' },
    { icon: <Shield className="w-6 h-6" />, title: 'Verification', path: '/pandit/verification', desc: 'Upload documents' },
    { icon: <UserCircle className="w-6 h-6" />, title: 'My Profile', path: '/profile', desc: 'Update your profile' },
  ]

  const links = user?.role === 'pandit' ? panditLinks : customerLinks

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Welcome, {user?.full_name}</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">Manage your rituals and bookings</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {links.map((link, index) => (
          <Link key={index} to={link.path} className="card-hover">
            <div className="w-12 h-12 bg-primary-100 dark:bg-primary-900 rounded-lg flex items-center justify-center mb-4 text-primary-600">
              {link.icon}
            </div>
            <h3 className="font-semibold mb-1">{link.title}</h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">{link.desc}</p>
          </Link>
        ))}
      </div>

      {user?.role === 'admin' && (
        <div className="mt-8">
          <Link to="/admin" className="btn-primary inline-block">Go to Admin Dashboard</Link>
        </div>
      )}
    </div>
  )
}
