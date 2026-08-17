import { Mail, MapPin, Phone, ShieldAlert } from 'lucide-react'

export default function Contact() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="max-w-3xl mb-10">
        <h1 className="text-4xl font-bold mb-4">Contact</h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          Reach out for booking support, pandit verification questions, payment issues, or security concerns.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-10">
        <div className="card">
          <Phone className="w-6 h-6 text-primary-600 mb-3" />
          <h2 className="font-semibold mb-1">Phone</h2>
          <p className="text-gray-600 dark:text-gray-400">+977 9800000000</p>
        </div>
        <div className="card">
          <Mail className="w-6 h-6 text-primary-600 mb-3" />
          <h2 className="font-semibold mb-1">Email</h2>
          <p className="text-gray-600 dark:text-gray-400">support@bishalpujasewa.com</p>
        </div>
        <div className="card">
          <MapPin className="w-6 h-6 text-primary-600 mb-3" />
          <h2 className="font-semibold mb-1">Service Area</h2>
          <p className="text-gray-600 dark:text-gray-400">Kathmandu Valley</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <form className="card space-y-4" onSubmit={(e) => e.preventDefault()}>
          <h2 className="text-xl font-semibold">Send a message</h2>
          <div>
            <label className="block text-sm font-medium mb-1">Name</label>
            <input className="input-field" type="text" placeholder="Your name" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Email</label>
            <input className="input-field" type="email" placeholder="you@example.com" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Message</label>
            <textarea className="input-field" rows={5} placeholder="How can we help?" />
          </div>
          <button className="btn-primary" type="submit">Submit</button>
        </form>

        <aside className="card h-fit">
          <div className="flex items-start gap-3">
            <ShieldAlert className="w-6 h-6 text-saffron-600 shrink-0 mt-1" />
            <div>
              <h2 className="text-xl font-semibold mb-2">Security or fraud issue?</h2>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Include the booking ID, payment transaction ID, and the account email involved. Admin audit logs
                and payment history can then be used to investigate the issue.
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                For urgent payment disputes, contact support before cancelling or refunding a booking.
              </p>
            </div>
          </div>
        </aside>
      </div>
    </div>
  )
}
