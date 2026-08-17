import { Link } from 'react-router-dom'

export default function Footer() {
  return (
    <footer className="bg-gray-900 text-gray-300">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <span className="text-3xl">ॐ</span>
              <span className="text-xl font-bold text-primary-400">Bishal Puja Sewa</span>
            </div>
            <p className="text-gray-400">
              Connecting Hindu households with verified pandits for authentic religious ceremonies.
            </p>
          </div>
          <div>
            <h3 className="text-white font-semibold mb-4">Quick Links</h3>
            <ul className="space-y-2">
              <li><Link to="/rituals" className="hover:text-primary-400 transition-colors">Rituals</Link></li>
              <li><Link to="/pandits" className="hover:text-primary-400 transition-colors">Find Pandits</Link></li>
              <li><Link to="/about" className="hover:text-primary-400 transition-colors">About Us</Link></li>
              <li><Link to="/contact" className="hover:text-primary-400 transition-colors">Contact</Link></li>
            </ul>
          </div>
          <div>
            <h3 className="text-white font-semibold mb-4">Security</h3>
            <ul className="space-y-2">
              <li className="text-gray-400">✓ End-to-end encrypted</li>
              <li className="text-gray-400">✓ Verified pandits only</li>
              <li className="text-gray-400">✓ Secure payments</li>
              <li className="text-gray-400">✓ Data protection</li>
            </ul>
          </div>
        </div>
        <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-500">
          <p>&copy; {new Date().getFullYear()} Bishal Puja Sewa. All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}
