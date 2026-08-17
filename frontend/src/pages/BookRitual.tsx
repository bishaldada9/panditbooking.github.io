import { FormEvent, useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import toast from 'react-hot-toast'
import { CalendarCheck, MapPin, UserCircle } from 'lucide-react'
import { bookingsService } from '../services/bookings.service'
import { panditsService } from '../services/pandits.service'
import { ritualsService } from '../services/rituals.service'
import { LoadingSpinner } from '../components/common'
import type { Pandit, Ritual } from '../types'

export default function BookRitual() {
  const { ritualId, panditId } = useParams<{ ritualId?: string; panditId?: string }>()
  const navigate = useNavigate()
  const [rituals, setRituals] = useState<Ritual[]>([])
  const [pandits, setPandits] = useState<Pandit[]>([])
  const [selectedRitual, setSelectedRitual] = useState(ritualId || '')
  const [selectedPandit, setSelectedPandit] = useState(panditId || '')
  const [scheduledDate, setScheduledDate] = useState('')
  const [startTime, setStartTime] = useState('08:00')
  const [endTime, setEndTime] = useState('10:00')
  const [address, setAddress] = useState('')
  const [specialNotes, setSpecialNotes] = useState('')
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    const load = async () => {
      try {
        const [ritualData, panditData] = await Promise.all([
          ritualsService.getAll(),
          panditsService.getAll(),
        ])
        setRituals(ritualData)
        setPandits(panditData.filter((p) => p.verification_status === 'approved'))
      } catch {
        toast.error('Failed to load booking options')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const ritual = useMemo(
    () => rituals.find((item) => item.id === selectedRitual),
    [rituals, selectedRitual],
  )

  const pandit = useMemo(
    () => pandits.find((item) => item.id === selectedPandit),
    [pandits, selectedPandit],
  )

  const total = ritual ? ritual.base_price + ritual.base_price * 0.1 : 0

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    if (!selectedRitual || !selectedPandit || !scheduledDate || !startTime || !endTime || !address.trim()) {
      toast.error('Please complete all required booking fields')
      return
    }

    setSubmitting(true)
    try {
      const booking = await bookingsService.create({
        ritual_id: selectedRitual,
        pandit_id: selectedPandit,
        scheduled_date: scheduledDate,
        start_time: startTime,
        end_time: endTime,
        address,
        special_notes: specialNotes,
      })
      toast.success('Booking created successfully')
      navigate(`/bookings/${booking.id}`)
    } catch {
      toast.error('Failed to create booking')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <LoadingSpinner size="lg" />

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Create Booking</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Choose a ritual, select a verified pandit, and schedule the ceremony.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 card space-y-5">
          <div>
            <label className="block text-sm font-medium mb-1">Ritual *</label>
            <select
              className="input-field"
              value={selectedRitual}
              onChange={(e) => setSelectedRitual(e.target.value)}
              required
            >
              <option value="">Select ritual</option>
              {rituals.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name} - Rs. {item.base_price}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Pandit *</label>
            <select
              className="input-field"
              value={selectedPandit}
              onChange={(e) => setSelectedPandit(e.target.value)}
              required
            >
              <option value="">Select pandit</option>
              {pandits.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.full_name} - {item.specialization}
                </option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Date *</label>
              <input
                type="date"
                className="input-field"
                value={scheduledDate}
                onChange={(e) => setScheduledDate(e.target.value)}
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Start *</label>
              <input
                type="time"
                className="input-field"
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">End *</label>
              <input
                type="time"
                className="input-field"
                value={endTime}
                onChange={(e) => setEndTime(e.target.value)}
                required
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Ceremony Address *</label>
            <textarea
              className="input-field"
              rows={3}
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="House number, street, ward, city"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Special Notes</label>
            <textarea
              className="input-field"
              rows={3}
              value={specialNotes}
              onChange={(e) => setSpecialNotes(e.target.value)}
              placeholder="Family gotra, ritual preferences, access instructions, or item notes"
            />
          </div>

          <button type="submit" disabled={submitting} className="btn-primary disabled:opacity-50">
            {submitting ? 'Creating...' : 'Create Booking'}
          </button>
        </div>

        <aside className="card h-fit space-y-5">
          <h2 className="text-xl font-semibold">Summary</h2>
          <div className="flex items-start gap-3">
            <CalendarCheck className="w-5 h-5 text-primary-500 shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-medium">{ritual?.name || 'No ritual selected'}</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">{ritual?.duration || 'Choose a ritual to see duration'}</p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <UserCircle className="w-5 h-5 text-primary-500 shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-medium">{pandit?.full_name || 'No pandit selected'}</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">{pandit?.service_area || 'Choose a pandit to see service area'}</p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <MapPin className="w-5 h-5 text-primary-500 shrink-0 mt-0.5" />
            <p className="text-sm text-gray-600 dark:text-gray-400">{address || 'Address not entered'}</p>
          </div>
          <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
            <div className="flex justify-between text-sm mb-1">
              <span>Ritual base</span>
              <span>Rs. {ritual?.base_price || 0}</span>
            </div>
            <div className="flex justify-between text-sm mb-3">
              <span>Platform fee</span>
              <span>Rs. {ritual ? ritual.base_price * 0.1 : 0}</span>
            </div>
            <div className="flex justify-between font-semibold text-lg">
              <span>Total</span>
              <span>Rs. {total}</span>
            </div>
          </div>
        </aside>
      </form>
    </div>
  )
}
