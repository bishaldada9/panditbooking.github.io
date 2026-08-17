import { Link } from 'react-router-dom'
import { ShieldCheck } from 'lucide-react'
import { EmptyState, LoadingSpinner } from '../components/common'
import { panditsService } from '../services/pandits.service'
import { useEffect, useState } from 'react'
import type { Pandit } from '../types'

export default function PanditVerification() {
  const [profile, setProfile] = useState<Pandit | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await panditsService.getProfile()
        setProfile(data)
      } catch {
        setProfile(null)
      } finally {
        setLoading(false)
      }
    }
    fetchProfile()
  }, [])

  if (loading) return <LoadingSpinner size="lg" />

  if (!profile) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <EmptyState
          title="No pandit profile yet"
          description="Create your pandit profile before admin verification can begin."
          action={<Link to="/pandit/register" className="btn-primary">Create Pandit Profile</Link>}
        />
      </div>
    )
  }

  const messages: Record<string, string> = {
    pending: 'Your profile is waiting for admin review. You can keep your profile details updated while it is reviewed.',
    approved: 'Your profile has been approved. Customers can discover and book your services.',
    rejected: 'Your profile was not approved yet. Update your profile details and contact support or wait for admin review.',
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="card">
        <div className="flex items-start gap-4">
          <div className="w-12 h-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center text-primary-600">
            <ShieldCheck className="w-6 h-6" />
          </div>
          <div>
            <h1 className="text-3xl font-bold mb-2">Verification</h1>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              {messages[profile.verification_status] || 'Your verification status is being reviewed.'}
            </p>
            <span className="inline-flex px-3 py-1 rounded-full text-sm font-medium capitalize bg-gray-100 dark:bg-gray-800">
              {profile.verification_status}
            </span>
            <div className="mt-6">
              <Link to="/profile" className="btn-primary">Review Profile</Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
