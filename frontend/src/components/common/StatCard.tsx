interface StatCardProps {
  icon: React.ReactNode
  label: string
  value: string | number
  change?: string
  color?: 'primary' | 'saffron' | 'sacred' | 'red'
}

export default function StatCard({ icon, label, value, change, color = 'primary' }: StatCardProps) {
  const colors = {
    primary: 'bg-primary-100 dark:bg-primary-900 text-primary-600',
    saffron: 'bg-saffron-100 dark:bg-saffron-900 text-saffron-600',
    sacred: 'bg-sacred-100 dark:bg-sacred-900 text-sacred-600',
    red: 'bg-red-100 dark:bg-red-900 text-red-600',
  }
  return (
    <div className="card">
      <div className="flex items-center justify-between">
        <div className={`w-12 h-12 rounded-lg flex items-center justify-center ${colors[color]}`}>
          {icon}
        </div>
        {change && (
          <span className={`text-sm font-medium ${change.startsWith('+') ? 'text-sacred-600' : 'text-red-500'}`}>
            {change}
          </span>
        )}
      </div>
      <p className="text-2xl font-bold mt-4">{value}</p>
      <p className="text-sm text-gray-500 dark:text-gray-400">{label}</p>
    </div>
  )
}
