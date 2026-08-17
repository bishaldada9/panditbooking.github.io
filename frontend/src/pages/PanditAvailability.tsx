import { FormEvent, useState } from 'react'
import toast from 'react-hot-toast'
import { Calendar, Clock } from 'lucide-react'
import { EmptyState } from '../components/common'
import { panditsService } from '../services/pandits.service'

interface AvailabilitySlot {
  id: string
  date: string
  start_time: string
  end_time: string
}

export default function PanditAvailability() {
  const today = new Date().toISOString().slice(0, 10)
  const [date, setDate] = useState(today)
  const [startTime, setStartTime] = useState('09:00')
  const [endTime, setEndTime] = useState('11:00')
  const [slots, setSlots] = useState<AvailabilitySlot[]>([])
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    if (endTime <= startTime) {
      toast.error('End time must be after start time')
      return
    }

    setSaving(true)
    try {
      await panditsService.updateAvailability({
        date,
        start_time: startTime,
        end_time: endTime,
        is_booked: false,
      })
      setSlots((current) => [
        { id: `${date}-${startTime}-${endTime}`, date, start_time: startTime, end_time: endTime },
        ...current,
      ])
      toast.success('Availability slot saved')
    } catch {
      toast.error('Failed to save availability')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Availability</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Add time windows when customers can request your services.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="card grid grid-cols-1 sm:grid-cols-4 gap-4 items-end mb-8">
        <div className="sm:col-span-2">
          <label className="block text-sm font-medium mb-1">Date</label>
          <input
            type="date"
            min={today}
            value={date}
            onChange={(event) => setDate(event.target.value)}
            className="input-field"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Start</label>
          <input
            type="time"
            value={startTime}
            onChange={(event) => setStartTime(event.target.value)}
            className="input-field"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">End</label>
          <input
            type="time"
            value={endTime}
            onChange={(event) => setEndTime(event.target.value)}
            className="input-field"
            required
          />
        </div>
        <button type="submit" disabled={saving} className="btn-primary sm:col-span-4 disabled:opacity-50">
          {saving ? 'Saving...' : 'Save Availability'}
        </button>
      </form>

      {slots.length === 0 ? (
        <EmptyState
          title="No availability added yet"
          description="Add a date and time slot above. Saved slots from this session will appear here."
        />
      ) : (
        <div className="space-y-3">
          {slots.map((slot) => (
            <div key={slot.id} className="card flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
              <div className="flex items-center gap-3">
                <Calendar className="w-5 h-5 text-primary-500" />
                <span className="font-medium">{new Date(slot.date).toLocaleDateString()}</span>
              </div>
              <div className="flex items-center gap-3 text-gray-600 dark:text-gray-400">
                <Clock className="w-5 h-5 text-primary-500" />
                <span>{slot.start_time} - {slot.end_time}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
