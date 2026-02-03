import { Card } from '@/shared/ui'
import { formatNumber } from '@/shared/lib/formatters'
import { Skeleton } from '@/shared/ui'

interface StatCardProps {
  title: string
  value: number
  icon: string
  isLoading?: boolean
}

export function StatCard({ title, value, icon, isLoading }: StatCardProps) {
  return (
    <Card className="flex items-center gap-4">
      <div className="text-4xl">{icon}</div>
      <div>
        <p className="text-sm text-gray-500">{title}</p>
        {isLoading ? (
          <Skeleton className="h-8 w-20 mt-1" />
        ) : (
          <p className="text-2xl font-bold text-gray-900">
            {formatNumber(value)}
          </p>
        )}
      </div>
    </Card>
  )
}
