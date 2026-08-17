import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { ritualsService } from '../services/rituals.service'
import { RitualCard } from '../components/ritual'
import { LoadingSpinner, EmptyState } from '../components/common'
import type { Ritual } from '../types'

export default function Rituals() {
  const navigate = useNavigate()
  const [rituals, setRituals] = useState<Ritual[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchRituals = async () => {
      try {
        const data = await ritualsService.getAll()
        setRituals(data)
      } catch {
        toast.error('Failed to load rituals')
      } finally {
        setLoading(false)
      }
    }
    fetchRituals()
  }, [])

  if (loading) return <LoadingSpinner size="lg" />

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Our Rituals</h1>
        <p className="text-gray-600 dark:text-gray-400">
          Browse our comprehensive list of Hindu rituals and ceremonies performed by experienced pandits.
        </p>
      </div>

      {rituals.length === 0 ? (
        <EmptyState
          title="No rituals found"
          description="Check back later for new rituals."
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {rituals.map((ritual) => (
            <div key={ritual.id} className="flex flex-col">
              <RitualCard
                ritual={ritual}
                onClick={() => navigate(`/rituals/${ritual.id}`)}
              />
              <button
                onClick={() => navigate(`/rituals/${ritual.id}`)}
                className="btn-primary mt-3"
              >
                Book Now
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
