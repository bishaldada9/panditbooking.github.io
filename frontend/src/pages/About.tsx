import { ShieldCheck, Users, HandCoins, Landmark } from 'lucide-react'

const values = [
  {
    icon: <ShieldCheck className="w-6 h-6" />,
    title: 'Verified trust',
    text: 'Pandit profiles are reviewed before they appear publicly, helping families choose with confidence.',
  },
  {
    icon: <Landmark className="w-6 h-6" />,
    title: 'Cultural care',
    text: 'The platform is designed around Hindu household rituals, local service areas, and Kathmandu Valley needs.',
  },
  {
    icon: <HandCoins className="w-6 h-6" />,
    title: 'Transparent pricing',
    text: 'Ritual pricing and platform fees are shown before booking so there is less room for confusion.',
  },
  {
    icon: <Users className="w-6 h-6" />,
    title: 'Accountability',
    text: 'Bookings, reviews, admin actions, and payment events are tracked so disputes can be handled clearly.',
  },
]

export default function About() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="max-w-3xl mb-12">
        <h1 className="text-4xl font-bold mb-4">About Bishal Puja Sewa</h1>
        <p className="text-lg text-gray-600 dark:text-gray-400 leading-relaxed">
          Bishal Puja Sewa connects Hindu households with verified pandits for religious ceremonies,
          while bringing booking reliability, transparent pricing, and safer digital workflows into a
          traditionally offline service space.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
        {values.map((value) => (
          <div key={value.title} className="card">
            <div className="w-12 h-12 rounded-lg bg-primary-100 dark:bg-primary-900 text-primary-600 dark:text-primary-300 flex items-center justify-center mb-4">
              {value.icon}
            </div>
            <h2 className="text-xl font-semibold mb-2">{value.title}</h2>
            <p className="text-gray-600 dark:text-gray-400">{value.text}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <section className="lg:col-span-2">
          <h2 className="text-2xl font-semibold mb-3">Why this platform exists</h2>
          <p className="text-gray-600 dark:text-gray-400 leading-relaxed mb-4">
            Many families still arrange rituals through informal references, scattered phone calls, or
            social media messages. That can make it hard to verify identity, compare prices, manage
            availability, and resolve problems after payment.
          </p>
          <p className="text-gray-600 dark:text-gray-400 leading-relaxed">
            Version 2 focuses on turning those needs into working application flows: searchable rituals,
            verified pandit listings, accountable booking actions, review moderation, payment records,
            and audit logs that help admins understand what happened inside the system.
          </p>
        </section>

        <aside className="card">
          <h2 className="text-xl font-semibold mb-3">Project focus</h2>
          <ul className="space-y-3 text-sm text-gray-600 dark:text-gray-400">
            <li>Secure account and role management</li>
            <li>Verified pandit onboarding</li>
            <li>Ritual discovery and booking</li>
            <li>Admin monitoring and audit trails</li>
            <li>Payment status tracking</li>
          </ul>
        </aside>
      </div>
    </div>
  )
}
