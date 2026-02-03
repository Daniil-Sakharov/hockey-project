import { useState, useEffect } from 'react'
import type { RefObject } from 'react'

interface ScrollRevealState {
  scrollProgress: number
  phase: 'intro' | 'rotate' | 'reveal'
  revealProgress: number
}

export function useScrollReveal(containerRef: RefObject<HTMLElement | null>): ScrollRevealState {
  const [scrollProgress, setScrollProgress] = useState(0)
  const [phase, setPhase] = useState<'intro' | 'rotate' | 'reveal'>('intro')
  const [revealProgress, setRevealProgress] = useState(0)

  useEffect(() => {
    const handleScroll = () => {
      if (!containerRef.current) return

      const rect = containerRef.current.getBoundingClientRect()
      const containerHeight = containerRef.current.offsetHeight
      const viewportHeight = window.innerHeight

      // Calculate scroll progress within the container
      // Animation completes at 80% scroll, then holds
      const scrolled = viewportHeight - rect.top
      const totalScrollable = containerHeight * 0.75 // Animation completes faster
      const progress = Math.max(0, Math.min(1, scrolled / totalScrollable))

      setScrollProgress(progress)

      // Determine phase based on scroll progress
      if (progress < 0.1) {
        setPhase('intro')
        setRevealProgress(0)
      } else if (progress < 0.5) {
        setPhase('rotate')
        setRevealProgress(0)
      } else {
        setPhase('reveal')
        // Map 0.5-1.0 to 0-1
        const reveal = (progress - 0.5) / 0.5
        setRevealProgress(reveal)
      }
    }

    // Initial call
    handleScroll()

    window.addEventListener('scroll', handleScroll, { passive: true })
    return () => window.removeEventListener('scroll', handleScroll)
  }, [containerRef])

  return { scrollProgress, phase, revealProgress }
}
