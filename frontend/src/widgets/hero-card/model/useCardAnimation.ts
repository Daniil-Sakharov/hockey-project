import { useState, useEffect, useRef } from 'react'

interface CardAnimationState {
  animationProgress: number
  isAnimating: boolean
}

export function useCardAnimation(duration: number = 2000): CardAnimationState {
  const [animationProgress, setAnimationProgress] = useState(0)
  const [isAnimating, setIsAnimating] = useState(true)
  const startTimeRef = useRef<number | null>(null)
  const frameRef = useRef<number | undefined>(undefined)

  useEffect(() => {
    const animate = (timestamp: number) => {
      if (!startTimeRef.current) {
        startTimeRef.current = timestamp
      }

      const elapsed = timestamp - startTimeRef.current
      const progress = Math.min(elapsed / duration, 1)

      setAnimationProgress(progress)

      if (progress < 1) {
        frameRef.current = requestAnimationFrame(animate)
      } else {
        setIsAnimating(false)
      }
    }

    // Small delay before starting animation
    const timeout = setTimeout(() => {
      frameRef.current = requestAnimationFrame(animate)
    }, 300)

    return () => {
      clearTimeout(timeout)
      if (frameRef.current) {
        cancelAnimationFrame(frameRef.current)
      }
    }
  }, [duration])

  return { animationProgress, isAnimating }
}
