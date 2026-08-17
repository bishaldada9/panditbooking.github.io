import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { useAuth } from '../store/auth.context'
import { profileService } from '../services/profile.service'
import { panditsService } from '../services/pandits.service'
import api from '../services/api'
import { LoadingSpinner } from '../components/common'
import { User, Shield, Key, Smartphone, LogOut, Eye, EyeOff, Check, X, Calendar, MapPin, Star, Clock, Award, Globe, IndianRupee } from 'lucide-react'
import type { Pandit } from '../types'

interface Session {
  id: string
  device: string
  ip: string
  last_active: string
  is_current: boolean
}

export default function Profile() {
  const { user, setUser } = useAuth()
  const navigate = useNavigate()

  const [profile, setProfile] = useState(user)
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState(false)
  const [editName, setEditName] = useState('')
  const [editPhone, setEditPhone] = useState('')
  const [saving, setSaving] = useState(false)

  const [panditProfile, setPanditProfile] = useState<Pandit | null>(null)
  const [panditLoading, setPanditLoading] = useState(false)
  const [panditEditing, setPanditEditing] = useState(false)
  const [editPandit, setEditPandit] = useState<Partial<Pandit>>({})

  const [sessions, setSessions] = useState<Session[]>([])
  const [sessionsLoading, setSessionsLoading] = useState(false)
  const [loggingOutAll, setLoggingOutAll] = useState(false)

  const [showPasswordForm, setShowPasswordForm] = useState(false)
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [changingPassword, setChangingPassword] = useState(false)
  const [showOldPassword, setShowOldPassword] = useState(false)
  const [showNewPassword, setShowNewPassword] = useState(false)

  const [mfaLoading, setMfaLoading] = useState(false)
  const [recoveryCodes, setRecoveryCodes] = useState<string[] | null>(null)
  const [mfaCode, setMfaCode] = useState('')
  const [mfaStep, setMfaStep] = useState<'idle' | 'setup' | 'verify'>('idle')

  useEffect(() => {
    fetchProfile()
    fetchSessions()
  }, [])

  const fetchProfile = async () => {
    setLoading(true)
    try {
      const data = await profileService.get()
      setProfile(data)
      setUser(data)
      setEditName(data.full_name)
      setEditPhone(data.phone || '')
      if (data.role === 'pandit') {
        fetchPanditProfile()
      }
    } catch {
      toast.error('Failed to load profile')
    } finally {
      setLoading(false)
    }
  }

  const fetchPanditProfile = async () => {
    setPanditLoading(true)
    try {
      const data = await panditsService.getProfile()
      setPanditProfile(data)
      setEditPandit({
        bio: data.bio,
        specialization: data.specialization,
        experience_years: data.experience_years,
        languages: data.languages,
        base_price: data.base_price,
        service_area: data.service_area,
      })
    } catch {
      // pandit profile may not exist yet
    } finally {
      setPanditLoading(false)
    }
  }

  const fetchSessions = async () => {
    setSessionsLoading(true)
    try {
      const res = await api.get('/auth/sessions')
      setSessions(res.data.data || [])
    } catch {
      // silently fail
    } finally {
      setSessionsLoading(false)
    }
  }

  const handleSaveProfile = async () => {
    setSaving(true)
    try {
      const updated = await profileService.update({ full_name: editName, phone: editPhone })
      setProfile(updated)
      setUser(updated)
      setEditing(false)
      toast.success('Profile updated successfully')
    } catch {
      toast.error('Failed to update profile')
    } finally {
      setSaving(false)
    }
  }

  const handleCancelEdit = () => {
    setEditName(profile?.full_name || '')
    setEditPhone(profile?.phone || '')
    setEditing(false)
  }

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    if (newPassword !== confirmPassword) {
      toast.error('New passwords do not match')
      return
    }
    if (newPassword.length < 6) {
      toast.error('Password must be at least 6 characters')
      return
    }
    setChangingPassword(true)
    try {
      await api.post('/auth/change-password', { old_password: oldPassword, new_password: newPassword })
      toast.success('Password changed successfully')
      setShowPasswordForm(false)
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch {
      toast.error('Failed to change password')
    } finally {
      setChangingPassword(false)
    }
  }

  const handleSetupMfa = async () => {
    setMfaLoading(true)
    try {
      const res = await api.post('/auth/mfa/setup')
      setRecoveryCodes(res.data.data?.recovery_codes || [])
      setMfaStep('verify')
    } catch {
      toast.error('Failed to setup MFA')
    } finally {
      setMfaLoading(false)
    }
  }

  const handleVerifyMfa = async () => {
    if (!mfaCode.trim()) {
      toast.error('Please enter the verification code')
      return
    }
    setMfaLoading(true)
    try {
      await api.post('/auth/mfa/verify', { code: mfaCode })
      toast.success('MFA enabled successfully')
      setMfaStep('idle')
      setMfaCode('')
      setRecoveryCodes(null)
      if (profile) {
        setProfile({ ...profile, mfa_enabled: true })
        setUser({ ...profile, mfa_enabled: true })
      }
    } catch {
      toast.error('Invalid verification code')
    } finally {
      setMfaLoading(false)
    }
  }

  const handleDisableMfa = async () => {
    if (!window.confirm('Are you sure you want to disable MFA? This reduces your account security.')) return
    setMfaLoading(true)
    try {
      await api.post('/auth/mfa/disable')
      toast.success('MFA disabled successfully')
      if (profile) {
        setProfile({ ...profile, mfa_enabled: false })
        setUser({ ...profile, mfa_enabled: false })
      }
    } catch {
      toast.error('Failed to disable MFA')
    } finally {
      setMfaLoading(false)
    }
  }

  const handleLogoutAllSessions = async () => {
    if (!window.confirm('This will log you out of all other devices. Continue?')) return
    setLoggingOutAll(true)
    try {
      await api.post('/auth/logout-all')
      toast.success('Other sessions logged out')
      fetchSessions()
    } catch {
      toast.error('Failed to logout other sessions')
    } finally {
      setLoggingOutAll(false)
    }
  }

  const handleSavePanditProfile = async () => {
    setSaving(true)
    try {
      await api.put('/pandit/profile', editPandit)
      toast.success('Pandit profile updated')
      setPanditEditing(false)
      fetchPanditProfile()
    } catch {
      toast.error('Failed to update pandit profile')
    } finally {
      setSaving(false)
    }
  }

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return 'N/A'
    return new Date(dateStr).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
  }

  const getVerificationBadge = (status?: string) => {
    if (!status) return null
    const styles: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
      approved: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
      rejected: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
    }
    return (
      <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${styles[status] || ''}`}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </span>
    )
  }

  if (loading) return <LoadingSpinner size="lg" />

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold">My Profile</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">Manage your account settings</p>
        </div>
        {!editing ? (
          <button onClick={() => setEditing(true)} className="btn-primary">
            Edit Profile
          </button>
        ) : (
          <div className="flex gap-2">
            <button onClick={handleSaveProfile} disabled={saving} className="btn-primary">
              {saving ? 'Saving...' : 'Save'}
            </button>
            <button onClick={handleCancelEdit} className="btn-secondary">
              Cancel
            </button>
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className="card">
            <div className="flex items-start gap-4">
              <div className="w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center flex-shrink-0">
                <User className="w-8 h-8 text-primary-600" />
              </div>
              <div className="flex-1 min-w-0">
                {editing ? (
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium mb-1">Full Name</label>
                      <input
                        type="text"
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                        className="input-field"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1">Phone</label>
                      <input
                        type="text"
                        value={editPhone}
                        onChange={(e) => setEditPhone(e.target.value)}
                        className="input-field"
                      />
                    </div>
                  </div>
                ) : (
                  <>
                    <h2 className="text-2xl font-bold">{profile?.full_name}</h2>
                    <p className="text-gray-600 dark:text-gray-400">{profile?.email}</p>
                    {profile?.phone && <p className="text-gray-600 dark:text-gray-400">{profile?.phone}</p>}
                    <div className="flex items-center gap-2 mt-2">
                      <span className="text-xs bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200 px-2 py-0.5 rounded-full font-medium capitalize">
                        {profile?.role}
                      </span>
                      {profile?.is_email_verified && (
                        <span className="text-xs bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 px-2 py-0.5 rounded-full font-medium">
                          Verified
                        </span>
                      )}
                    </div>
                  </>
                )}
              </div>
            </div>
            <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
                <Calendar className="w-4 h-4" />
                <span>Member since {formatDate(profile?.created_at)}</span>
              </div>
            </div>
          </div>

          <div className="card">
            <h2 className="text-xl font-semibold mb-4">Account Security</h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                <div className="flex items-center gap-3">
                  <Key className="w-5 h-5 text-gray-500" />
                  <div>
                    <p className="font-medium">Password</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Change your account password</p>
                  </div>
                </div>
                <button onClick={() => setShowPasswordForm(!showPasswordForm)} className="btn-secondary text-sm">
                  {showPasswordForm ? 'Cancel' : 'Change'}
                </button>
              </div>
              {showPasswordForm && (
                <form onSubmit={handleChangePassword} className="space-y-3 pl-11">
                  <div>
                    <label className="block text-sm font-medium mb-1">Current Password</label>
                    <div className="relative">
                      <input
                        type={showOldPassword ? 'text' : 'password'}
                        value={oldPassword}
                        onChange={(e) => setOldPassword(e.target.value)}
                        className="input-field pr-10"
                        required
                      />
                      <button type="button" onClick={() => setShowOldPassword(!showOldPassword)} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
                        {showOldPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                      </button>
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-1">New Password</label>
                    <div className="relative">
                      <input
                        type={showNewPassword ? 'text' : 'password'}
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        className="input-field pr-10"
                        required
                      />
                      <button type="button" onClick={() => setShowNewPassword(!showNewPassword)} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
                        {showNewPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                      </button>
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-1">Confirm New Password</label>
                    <input
                      type="password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                      className="input-field"
                      required
                    />
                  </div>
                  <button type="submit" disabled={changingPassword} className="btn-primary">
                    {changingPassword ? 'Changing...' : 'Change Password'}
                  </button>
                </form>
              )}

              <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                <div className="flex items-center gap-3">
                  <Smartphone className="w-5 h-5 text-gray-500" />
                  <div>
                    <p className="font-medium">Two-Factor Authentication (MFA)</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {profile?.mfa_enabled ? 'MFA is enabled' : 'Add an extra layer of security'}
                    </p>
                  </div>
                </div>
                {profile?.mfa_enabled ? (
                  <button onClick={handleDisableMfa} disabled={mfaLoading} className="btn-danger text-sm">
                    {mfaLoading ? 'Disabling...' : 'Disable MFA'}
                  </button>
                ) : mfaStep === 'verify' ? (
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={mfaCode}
                      onChange={(e) => setMfaCode(e.target.value)}
                      placeholder="Enter code"
                      className="input-field w-32 text-sm"
                    />
                    <button onClick={handleVerifyMfa} disabled={mfaLoading} className="btn-primary text-sm">
                      {mfaLoading ? 'Verifying...' : 'Verify'}
                    </button>
                  </div>
                ) : (
                  <button onClick={handleSetupMfa} disabled={mfaLoading} className="btn-primary text-sm">
                    {mfaLoading ? 'Setting up...' : 'Enable MFA'}
                  </button>
                )}
              </div>
              {recoveryCodes && (
                <div className="pl-11">
                  <div className="bg-yellow-50 dark:bg-yellow-900/30 border border-yellow-200 dark:border-yellow-700 rounded-lg p-4">
                    <p className="text-sm font-medium text-yellow-800 dark:text-yellow-200 mb-2">Save these recovery codes securely!</p>
                    <div className="grid grid-cols-2 gap-1">
                      {recoveryCodes.map((code, i) => (
                        <code key={i} className="text-sm font-mono bg-white dark:bg-gray-800 px-2 py-1 rounded">{code}</code>
                      ))}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold">Active Sessions</h2>
              <button onClick={handleLogoutAllSessions} disabled={loggingOutAll} className="btn-danger text-sm flex items-center gap-1">
                <LogOut className="w-4 h-4" />
                {loggingOutAll ? 'Logging out...' : 'Logout All Other'}
              </button>
            </div>
            {sessionsLoading ? (
              <LoadingSpinner size="sm" />
            ) : sessions.length === 0 ? (
              <p className="text-gray-500 dark:text-gray-400 text-sm">No active sessions</p>
            ) : (
              <div className="space-y-3">
                {sessions.map((session) => (
                  <div key={session.id} className={`flex items-center justify-between p-3 rounded-lg ${session.is_current ? 'bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800' : 'bg-gray-50 dark:bg-gray-700/50'}`}>
                    <div>
                      <p className="font-medium text-sm">{session.device || 'Unknown device'}</p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">{session.ip} · Last active: {formatDate(session.last_active)}</p>
                    </div>
                    {session.is_current && (
                      <span className="text-xs bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200 px-2 py-0.5 rounded-full font-medium">
                        Current
                      </span>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="space-y-6">
          <div className="card">
            <h2 className="text-xl font-semibold mb-4">Quick Stats</h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600 dark:text-gray-400">Email Verified</span>
                {profile?.is_email_verified ? (
                  <span className="flex items-center gap-1 text-sm text-green-600"><Check className="w-4 h-4" /> Yes</span>
                ) : (
                  <span className="flex items-center gap-1 text-sm text-red-600"><X className="w-4 h-4" /> No</span>
                )}
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600 dark:text-gray-400">MFA Enabled</span>
                {profile?.mfa_enabled ? (
                  <span className="flex items-center gap-1 text-sm text-green-600"><Check className="w-4 h-4" /> Yes</span>
                ) : (
                  <span className="flex items-center gap-1 text-sm text-red-600"><X className="w-4 h-4" /> No</span>
                )}
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600 dark:text-gray-400">Account Status</span>
                <span className={`text-sm ${profile?.is_active ? 'text-green-600' : 'text-red-600'}`}>
                  {profile?.is_active ? 'Active' : 'Inactive'}
                </span>
              </div>
            </div>
          </div>

          {profile?.role === 'pandit' && (
            <div className="card">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold">Pandit Profile</h2>
                {!panditEditing ? (
                  <button onClick={() => setPanditEditing(true)} className="btn-secondary text-sm">
                    Edit
                  </button>
                ) : (
                  <div className="flex gap-2">
                    <button onClick={handleSavePanditProfile} disabled={saving} className="btn-primary text-sm">
                      {saving ? 'Saving...' : 'Save'}
                    </button>
                    <button onClick={() => { setPanditEditing(false); fetchPanditProfile() }} className="btn-secondary text-sm">
                      Cancel
                    </button>
                  </div>
                )}
              </div>
              {panditLoading ? (
                <LoadingSpinner size="sm" />
              ) : panditProfile ? (
                <div className="space-y-3">
                  {panditEditing ? (
                    <>
                      <div>
                        <label className="block text-sm font-medium mb-1">Bio</label>
                        <textarea
                          value={editPandit.bio || ''}
                          onChange={(e) => setEditPandit({ ...editPandit, bio: e.target.value })}
                          className="input-field"
                          rows={3}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium mb-1">Specialization</label>
                        <input
                          type="text"
                          value={editPandit.specialization || ''}
                          onChange={(e) => setEditPandit({ ...editPandit, specialization: e.target.value })}
                          className="input-field"
                        />
                      </div>
                      <div className="grid grid-cols-2 gap-3">
                        <div>
                          <label className="block text-sm font-medium mb-1">Experience (years)</label>
                          <input
                            type="number"
                            value={editPandit.experience_years || 0}
                            onChange={(e) => setEditPandit({ ...editPandit, experience_years: parseInt(e.target.value) || 0 })}
                            className="input-field"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium mb-1">Base Price (₹)</label>
                          <input
                            type="number"
                            value={editPandit.base_price || 0}
                            onChange={(e) => setEditPandit({ ...editPandit, base_price: parseInt(e.target.value) || 0 })}
                            className="input-field"
                          />
                        </div>
                      </div>
                      <div>
                        <label className="block text-sm font-medium mb-1">Languages (comma separated)</label>
                        <input
                          type="text"
                          value={(editPandit.languages || []).join(', ')}
                          onChange={(e) => setEditPandit({ ...editPandit, languages: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })}
                          className="input-field"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium mb-1">Service Area</label>
                        <input
                          type="text"
                          value={editPandit.service_area || ''}
                          onChange={(e) => setEditPandit({ ...editPandit, service_area: e.target.value })}
                          className="input-field"
                        />
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-gray-600 dark:text-gray-400">Verification</span>
                        {getVerificationBadge(panditProfile.verification_status)}
                      </div>
                      {panditProfile.bio && (
                        <p className="text-sm text-gray-600 dark:text-gray-400">{panditProfile.bio}</p>
                      )}
                      <div className="flex items-center gap-2 text-sm">
                        <Award className="w-4 h-4 text-primary-500" />
                        <span>{panditProfile.specialization || 'No specialization'}</span>
                      </div>
                      <div className="flex items-center gap-2 text-sm">
                        <Clock className="w-4 h-4 text-primary-500" />
                        <span>{panditProfile.experience_years} years experience</span>
                      </div>
                      <div className="flex items-center gap-2 text-sm">
                        <Globe className="w-4 h-4 text-primary-500" />
                        <span>{panditProfile.languages?.join(', ') || 'N/A'}</span>
                      </div>
                      <div className="flex items-center gap-2 text-sm">
                        <IndianRupee className="w-4 h-4 text-primary-500" />
                        <span>₹{panditProfile.base_price} base price</span>
                      </div>
                      <div className="flex items-center gap-2 text-sm">
                        <MapPin className="w-4 h-4 text-primary-500" />
                        <span>{panditProfile.service_area || 'N/A'}</span>
                      </div>
                      <div className="flex items-center gap-2 text-sm">
                        <Star className="w-4 h-4 text-yellow-500" />
                        <span>{panditProfile.rating.toFixed(1)} ({panditProfile.total_reviews} reviews)</span>
                      </div>
                    </>
                  )}
                </div>
              ) : (
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  No pandit profile found.{' '}
                  <button onClick={() => navigate('/pandit/register')} className="text-primary-600 hover:underline font-medium">
                    Create one
                  </button>
                </p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
