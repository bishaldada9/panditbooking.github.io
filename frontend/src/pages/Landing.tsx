import { Link } from 'react-router-dom'
import { Sparkles, Shield, Users, CreditCard, Search, Star } from 'lucide-react'

export default function Landing() {
  return (
    <div>
      {/* Hero Section */}
      <section className="relative bg-gradient-to-br from-primary-50 via-white to-saffron-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6">
              Your Sacred Rituals,{' '}
              <span className="text-primary-600 dark:text-primary-400">Simplified</span>
            </h1>
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-3xl mx-auto">
              Find verified pandits, book religious ceremonies, and perform rituals with confidence.
              Secure, trusted, and deeply rooted in tradition.
            </p>
            <div className="flex flex-col sm:flex-row justify-center gap-4">
              <Link to="/register" className="btn-primary text-lg px-8 py-3">
                Get Started
              </Link>
              <Link to="/pandits" className="btn-secondary text-lg px-8 py-3">
                Find a Pandit
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 bg-white dark:bg-gray-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12">Why Choose Bishal Puja Sewa?</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <div key={index} className="card-hover text-center">
                <div className="w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center mx-auto mb-4">
                  {feature.icon}
                </div>
                <h3 className="text-xl font-semibold mb-3">{feature.title}</h3>
                <p className="text-gray-600 dark:text-gray-400">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How it works */}
      <section className="py-20 bg-gray-50 dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            {steps.map((step, index) => (
              <div key={index} className="text-center">
                <div className="w-12 h-12 bg-primary-500 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                  {index + 1}
                </div>
                <h3 className="font-semibold mb-2">{step.title}</h3>
                <p className="text-gray-600 dark:text-gray-400 text-sm">{step.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 bg-primary-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold text-white mb-6">Ready to Begin?</h2>
          <p className="text-primary-100 mb-8 text-lg">Join thousands of satisfied households performing rituals with verified pandits.</p>
          <Link to="/register" className="bg-white text-primary-600 font-semibold px-8 py-3 rounded-lg hover:bg-primary-50 transition-colors text-lg">
            Create Your Account
          </Link>
        </div>
      </section>
    </div>
  )
}

const features = [
  {
    icon: <Shield className="w-8 h-8 text-primary-600" />,
    title: 'Verified Pandits',
    description: 'Every pandit is thoroughly verified with documents and background checks.',
  },
  {
    icon: <Search className="w-8 h-8 text-primary-600" />,
    title: 'Easy Booking',
    description: 'Browse rituals, check availability, and book in just a few clicks.',
  },
  {
    icon: <CreditCard className="w-8 h-8 text-primary-600" />,
    title: 'Secure Payments',
    description: 'Multiple payment options with end-to-end encryption and secure processing.',
  },
  {
    icon: <Users className="w-8 h-8 text-primary-600" />,
    title: 'Trusted Community',
    description: 'Read reviews and ratings from real customers to make informed decisions.',
  },
  {
    icon: <Star className="w-8 h-8 text-primary-600" />,
    title: 'Quality Assured',
    description: 'All rituals follow traditional practices with experienced pandits.',
  },
  {
    icon: <Sparkles className="w-8 h-8 text-primary-600" />,
    title: '24/7 Support',
    description: 'Our team is always available to help with any questions or concerns.',
  },
]

const steps = [
  {
    title: 'Browse Rituals',
    description: 'Explore our comprehensive list of Hindu rituals and ceremonies.',
  },
  {
    title: 'Choose Pandit',
    description: 'Select from verified pandits based on ratings and specialization.',
  },
  {
    title: 'Book & Pay',
    description: 'Schedule your ritual and pay securely through our platform.',
  },
  {
    title: 'Perform Ritual',
    description: 'Your pandit arrives at the scheduled time to perform the ceremony.',
  },
]
