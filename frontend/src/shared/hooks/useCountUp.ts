import { useEffect, useRef, useState } from 'react'

export function useCountUp(end: number, duration = 800): number {
  const [value, setValue] = useState(0)
  const prevEnd = useRef(0)

  useEffect(() => {
    if (end === prevEnd.current) return
    const start = prevEnd.current
    prevEnd.current = end
    const startTime = performance.now()

    function tick(now: number) {
      const elapsed = now - startTime
      const progress = Math.min(elapsed / duration, 1)
      const eased = 1 - Math.pow(1 - progress, 3) // easeOutCubic
      setValue(Math.round(start + (end - start) * eased))
      if (progress < 1) requestAnimationFrame(tick)
    }

    requestAnimationFrame(tick)
  }, [end, duration])

  return value
}
