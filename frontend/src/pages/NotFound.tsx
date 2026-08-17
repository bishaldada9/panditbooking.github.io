import { Link } from 'react-router-dom'
import { EmptyState } from '../components/common'

export default function NotFound() {
  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <EmptyState
        title="Page not found"
        description="This page does not exist or the link is no longer valid."
        action={<Link to="/dashboard" className="btn-primary">Go to Dashboard</Link>}
      />
    </div>
  )
}
