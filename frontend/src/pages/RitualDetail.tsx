import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { ArrowLeft, Clock, IndianRupee, Package, ListOrdered } from 'lucide-react'
import { ritualsService } from '../services/rituals.service'
import { LoadingSpinner, EmptyState } from '../components/common'
import type { Ritual } from '../types'

export default function RitualDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [ritual, setRitual] = useState<Ritual | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    const fetchRitual = async () => {
      try {
        const data = await ritualsService.getById(id)
        setRitual(data)
      } catch {
        toast.error('Failed to load ritual details')
      } finally {
        setLoading(false)
      }
    }
    fetchRitual()
  }, [id])

  if (loading) return <LoadingSpinner size="lg" />
  if (!ritual) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <EmptyState
          title="Ritual not found"
          description="This ritual may have been removed or the link is invalid."
          action={<button onClick={() => navigate('/rituals')} className="btn-primary">Back to Rituals</button>}
        />
      </div>
    )
  }

  const procedureSteps = ritual.procedure
    ? ritual.procedure.split('\n').filter((s) => s.trim())
    : []
  const requiredItems = ritual.required_items
    ? ritual.required_items.split(',').map((s) => s.trim()).filter((s) => s)
    : []

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <button
        onClick={() => navigate('/rituals')}
        className="flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white mb-6 transition-colors"
      >
        <ArrowLeft className="w-5 h-5" />
        Back to Rituals
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2">
          <div className="flex items-center gap-3 mb-4">
            <h1 className="text-3xl font-bold">{ritual.name}</h1>
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200">
              {ritual.category_name}
            </span>
          </div>
          <p className="text-gray-600 dark:text-gray-400 mb-8 text-lg">
            {ritual.description}
          </p>

          {procedureSteps.length > 0 && (
            <div className="mb-8">
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <ListOrdered className="w-5 h-5" />
                Procedure
              </h2>
              <ol className="space-y-3">
                {procedureSteps.map((step, index) => (
                  <li key={index} className="flex gap-3">
                    <span className="flex-shrink-0 w-7 h-7 rounded-full bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 flex items-center justify-center text-sm font-semibold">
                      {index + 1}
                    </span>
                    <span className="text-gray-600 dark:text-gray-400 pt-0.5">{step}</span>
                  </li>
                ))}
              </ol>
            </div>
          )}
        </div>

        <div className="lg:col-span-1">
          <div className="card space-y-6">
            <div className="flex items-center gap-3">
              <Clock className="w-5 h-5 text-gray-400" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Duration</p>
                <p className="font-semibold">{ritual.duration}</p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <IndianRupee className="w-5 h-5 text-gray-400" />
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Price</p>
                <p className="font-semibold text-2xl text-primary-600 dark:text-primary-400">
                  Rs. {ritual.base_price}
                </p>
              </div>
            </div>

            {requiredItems.length > 0 && (
              <div className="pt-4 border-t border-gray-100 dark:border-gray-700">
                <h3 className="font-semibold mb-3 flex items-center gap-2">
                  <Package className="w-5 h-5" />
                  Required Items
                </h3>
                <ul className="space-y-2">
                  {requiredItems.map((item, index) => (
                    <li key={index} className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                      <span className="w-1.5 h-1.5 rounded-full bg-primary-500 flex-shrink-0" />
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            )}

            <button
              onClick={() => navigate(`/book-ritual/${ritual.id}`)}
              className="btn-primary w-full"
            >
              Book This Ritual
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
