import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { panditsService } from '../services/pandits.service'

export default function PanditRegister() {
  const navigate = useNavigate()
  const [bio, setBio] = useState('')
  const [specialization, setSpecialization] = useState('')
  const [experienceYears, setExperienceYears] = useState(0)
  const [languages, setLanguages] = useState('Nepali, Sanskrit')
  const [basePrice, setBasePrice] = useState(1000)
  const [serviceArea, setServiceArea] = useState('Kathmandu Valley')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    setSubmitting(true)
    try {
      await panditsService.register({
        bio,
        specialization,
        experience_years: experienceYears,
        languages: languages.split(',').map((item) => item.trim()).filter(Boolean),
        base_price: basePrice,
        service_area: serviceArea,
      })
      toast.success('Pandit profile submitted for verification')
      navigate('/profile')
    } catch {
      toast.error('Failed to create pandit profile')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Pandit Profile</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Submit your service details for admin verification.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="card space-y-5">
        <div>
          <label className="block text-sm font-medium mb-1">Bio *</label>
          <textarea
            className="input-field"
            rows={4}
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Specialization *</label>
          <input
            className="input-field"
            value={specialization}
            onChange={(e) => setSpecialization(e.target.value)}
            placeholder="Griha Pravesh, Vivah, Shraddha"
            required
          />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Experience Years *</label>
            <input
              type="number"
              min={0}
              className="input-field"
              value={experienceYears}
              onChange={(e) => setExperienceYears(Number(e.target.value))}
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Base Price *</label>
            <input
              type="number"
              min={0}
              className="input-field"
              value={basePrice}
              onChange={(e) => setBasePrice(Number(e.target.value))}
              required
            />
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Languages *</label>
          <input
            className="input-field"
            value={languages}
            onChange={(e) => setLanguages(e.target.value)}
            placeholder="Nepali, Sanskrit, Hindi"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Service Area *</label>
          <input
            className="input-field"
            value={serviceArea}
            onChange={(e) => setServiceArea(e.target.value)}
            required
          />
        </div>
        <button className="btn-primary disabled:opacity-50" type="submit" disabled={submitting}>
          {submitting ? 'Submitting...' : 'Submit for Verification'}
        </button>
      </form>
    </div>
  )
}
