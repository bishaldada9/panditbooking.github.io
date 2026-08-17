import { Routes, Route } from 'react-router-dom'
import { AuthProvider } from './store/auth.context'
import { ProtectedRoute } from './components/layout/ProtectedRoute'
import { AdminRoute } from './components/layout/AdminRoute'
import Layout from './components/layout/Layout'
import Landing from './pages/Landing'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Rituals from './pages/Rituals'
import RitualDetail from './pages/RitualDetail'
import Pandits from './pages/Pandits'
import PanditDetail from './pages/PanditDetail'
import Profile from './pages/Profile'
import Notifications from './pages/Notifications'
import Payments from './pages/Payments'
import Bookings from './pages/Bookings'
import BookingDetail from './pages/BookingDetail'
import Reviews from './pages/Reviews'
import AdminDashboard from './pages/AdminDashboard'
import About from './pages/About'
import Contact from './pages/Contact'
import BookRitual from './pages/BookRitual'
import PanditRegister from './pages/PanditRegister'
import PanditAvailability from './pages/PanditAvailability'
import PanditVerification from './pages/PanditVerification'
import NotFound from './pages/NotFound'

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Landing />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/rituals" element={<Rituals />} />
          <Route path="/rituals/:id" element={<RitualDetail />} />
          <Route path="/pandits" element={<Pandits />} />
          <Route path="/pandits/:id" element={<PanditDetail />} />
          <Route path="/about" element={<About />} />
          <Route path="/contact" element={<Contact />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          <Route
            path="/notifications"
            element={
              <ProtectedRoute>
                <Notifications />
              </ProtectedRoute>
            }
          />
          <Route
            path="/payments"
            element={
              <ProtectedRoute>
                <Payments />
              </ProtectedRoute>
            }
          />
          <Route
            path="/bookings"
            element={
              <ProtectedRoute>
                <Bookings />
              </ProtectedRoute>
            }
          />
          <Route
            path="/book-ritual/:ritualId"
            element={
              <ProtectedRoute>
                <BookRitual />
              </ProtectedRoute>
            }
          />
          <Route
            path="/book-pandit/:panditId"
            element={
              <ProtectedRoute>
                <BookRitual />
              </ProtectedRoute>
            }
          />
          <Route
            path="/pandit/bookings"
            element={
              <ProtectedRoute roles={['pandit']}>
                <Bookings />
              </ProtectedRoute>
            }
          />
          <Route
            path="/pandit/availability"
            element={
              <ProtectedRoute roles={['pandit']}>
                <PanditAvailability />
              </ProtectedRoute>
            }
          />
          <Route
            path="/pandit/verification"
            element={
              <ProtectedRoute roles={['pandit']}>
                <PanditVerification />
              </ProtectedRoute>
            }
          />
          <Route
            path="/pandit/register"
            element={
              <ProtectedRoute roles={['pandit']}>
                <PanditRegister />
              </ProtectedRoute>
            }
          />
          <Route
            path="/bookings/:id"
            element={
              <ProtectedRoute>
                <BookingDetail />
              </ProtectedRoute>
            }
          />
          <Route
            path="/reviews"
            element={
              <ProtectedRoute>
                <Reviews />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin"
            element={
              <AdminRoute>
                <AdminDashboard />
              </AdminRoute>
            }
          />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}
