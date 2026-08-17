import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../store/auth.context'
import { useState } from 'react'
import { Menu, X, User, LogOut, Moon, Sun, Bell, LayoutDashboard } from 'lucide-react'

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuth()
  const navigate = useNavigate()
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isDark, setIsDark] = useState(() => document.documentElement.classList.contains('dark'))

  const toggleDark = () => {
    document.documentElement.classList.toggle('dark')
    setIsDark(!isDark)
  }

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  return (
    <nav className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link to="/" className="flex items-center space-x-2">
              <span className="text-2xl">ॐ</span>
              <span className="text-xl font-bold text-primary-600 dark:text-primary-400">
                Bishal Puja Sewa
              </span>
            </Link>
          </div>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center space-x-4">
            <Link to="/rituals" className="text-gray-600 dark:text-gray-300 hover:text-primary-600">
              Rituals
            </Link>
            <Link to="/pandits" className="text-gray-600 dark:text-gray-300 hover:text-primary-600">
              Pandits
            </Link>

            {isAuthenticated ? (
              <>
                <Link to="/notifications" className="relative text-gray-600 dark:text-gray-300 hover:text-primary-600">
                  <Bell size={20} />
                </Link>
                <button onClick={toggleDark} className="text-gray-600 dark:text-gray-300 hover:text-primary-600">
                  {isDark ? <Sun size={20} /> : <Moon size={20} />}
                </button>
                <div className="relative group">
                  <button className="flex items-center space-x-1 text-gray-600 dark:text-gray-300 hover:text-primary-600">
                    <User size={20} />
                    <span>{user?.full_name}</span>
                  </button>
                  <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
                    <Link to="/dashboard" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700">
                      <LayoutDashboard size={16} className="inline mr-2" />Dashboard
                    </Link>
                    <Link to="/profile" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700">
                      <User size={16} className="inline mr-2" />Profile
                    </Link>
                    {user?.role === 'admin' && (
                      <Link to="/admin" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700">
                        Admin Panel
                      </Link>
                    )}
                    <button onClick={handleLogout} className="block w-full text-left px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700">
                      <LogOut size={16} className="inline mr-2" />Logout
                    </button>
                  </div>
                </div>
              </>
            ) : (
              <>
                <button onClick={toggleDark} className="text-gray-600 dark:text-gray-300 hover:text-primary-600">
                  {isDark ? <Sun size={20} /> : <Moon size={20} />}
                </button>
                <Link to="/login" className="btn-secondary text-sm">Login</Link>
                <Link to="/register" className="btn-primary text-sm">Register</Link>
              </>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button onClick={() => setIsMenuOpen(!isMenuOpen)} className="text-gray-600 dark:text-gray-300">
              {isMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {isMenuOpen && (
        <div className="md:hidden bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
          <div className="px-4 py-2 space-y-2">
            <Link to="/rituals" className="block py-2 text-gray-600 dark:text-gray-300">Rituals</Link>
            <Link to="/pandits" className="block py-2 text-gray-600 dark:text-gray-300">Pandits</Link>
            {isAuthenticated ? (
              <>
                <Link to="/dashboard" className="block py-2 text-gray-600 dark:text-gray-300">Dashboard</Link>
                <Link to="/profile" className="block py-2 text-gray-600 dark:text-gray-300">Profile</Link>
                <button onClick={handleLogout} className="block w-full text-left py-2 text-red-500">Logout</button>
              </>
            ) : (
              <>
                <Link to="/login" className="block py-2 text-gray-600 dark:text-gray-300">Login</Link>
                <Link to="/register" className="block py-2 text-gray-600 dark:text-gray-300">Register</Link>
              </>
            )}
          </div>
        </div>
      )}
    </nav>
  )
}
