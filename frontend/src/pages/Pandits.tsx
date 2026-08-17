import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { panditsService } from '../services/pandits.service'
import { PanditCard } from '../components/pandit'
import { LoadingSpinner, EmptyState } from '../components/common'
import type { Pandit } from '../types'

export default function Pandits() {
  const navigate = useNavigate()
  const [pandits, setPandits] = useState<Pandit[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchPandits = async () => {
      try {
        const data = await panditsService.getAll()
        setPandits(data.filter((p) => p.verification_status === 'approved'))
      } catch {
        toast.error('Failed to load pandits')
      } finally {
        setLoading(false)
      }
    }
    fetchPandits()
  }, [])

  if (loading) return <LoadingSpinner size="lg" />

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Our Pandits</h1>
        <p className="text-gray-600 dark:text-gray-400">
          Browse our verified pandits and choose the one that best suits your needs.
        </p>
      </div>

      {pandits.length === 0 ? (
        <EmptyState
          title="No pandits available"
          description="Check back later for verified pandits in your area."
        />
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {pandits.map((pandit) => (
            <PanditCard
              key={pandit.id}
              pandit={pandit}
              onClick={() => navigate(`/pandits/${pandit.id}`)}
            />
          ))}
        </div>
      )}
    </div>
  )
}
