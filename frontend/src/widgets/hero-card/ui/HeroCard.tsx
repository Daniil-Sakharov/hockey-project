import { Suspense, useRef, useState, useEffect, useMemo } from 'react'
import { Canvas } from '@react-three/fiber'
import { PerspectiveCamera, Environment } from '@react-three/drei'
import { EffectComposer, Bloom } from '@react-three/postprocessing'
import { motion } from 'framer-motion'
import { PlayerCard3D } from './PlayerCard3D'
import { useCardAnimation } from '../model/useCardAnimation'
import { useScrollReveal } from '../model/useScrollReveal'

// Player data with photo
const PLAYER_DATA = {
  name: 'АЛЕКСАНДР ОВЕЧКИН',
  number: 8,
  team: 'WASHINGTON CAPITALS',
  teamColor: '#00d4ff',
  position: 'LW',
  photo: 'https://cms.nhl.bamgrid.com/images/headshots/current/168x168/8471214.jpg',
  stats: {
    goals: 853,
    assists: 680,
    points: 1533,
    games: 1426,
  },
}

export function HeroCard() {
  const containerRef = useRef<HTMLDivElement>(null)
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 })

  const { animationProgress } = useCardAnimation(2000)
  const { scrollProgress } = useScrollReveal(containerRef)

  // Card position: starts on left, moves to right on scroll
  const cardPosition = useMemo(() => {
    return -1.1 + scrollProgress * 1.8  // -1.1 → +0.7
  }, [scrollProgress])

  // Animated stats based on scroll
  const animatedStats = useMemo(() => {
    const progress = Math.min(1, scrollProgress * 1.5)
    return {
      goals: Math.floor(PLAYER_DATA.stats.goals * progress),
      assists: Math.floor(PLAYER_DATA.stats.assists * progress),
      points: Math.floor(PLAYER_DATA.stats.points * progress),
      games: Math.floor(PLAYER_DATA.stats.games * progress),
    }
  }, [scrollProgress])

  // Track mouse position for card tilt
  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      const x = (e.clientX / window.innerWidth) * 2 - 1
      const y = -((e.clientY / window.innerHeight) * 2 - 1)
      setMousePosition({ x, y })
    }

    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [])

  return (
    <section
      ref={containerRef}
      className="relative bg-[#030812]"
      style={{ height: '380vh' }}
    >
      {/* Sticky container */}
      <div className="sticky top-0 h-screen w-full overflow-hidden">
        {/* Growing chart background */}
        <GrowingChartBackground scrollProgress={scrollProgress} />

        {/* 3D Canvas */}
        <div className="absolute inset-0">
          <Canvas>
            <Suspense fallback={null}>
              <PerspectiveCamera makeDefault position={[0, 0, 6]} fov={45} />

              {/* Lighting */}
              <ambientLight intensity={0.4} />
              <directionalLight position={[5, 5, 5]} intensity={1.2} color="#ffffff" />
              <pointLight position={[-5, 3, 5]} intensity={0.8} color="#00d4ff" />
              <pointLight position={[5, -3, 5]} intensity={0.6} color={PLAYER_DATA.teamColor} />

              <Environment preset="city" />

              <PlayerCard3D
                player={{ ...PLAYER_DATA, stats: animatedStats }}
                animationProgress={animationProgress}
                mousePosition={mousePosition}
                scrollProgress={scrollProgress}
                horizontalOffset={cardPosition}
              />

              <EffectComposer>
                <Bloom
                  intensity={0.5}
                  luminanceThreshold={0.6}
                  luminanceSmoothing={0.9}
                  mipmapBlur
                />
              </EffectComposer>
            </Suspense>
          </Canvas>
        </div>

        {/* Text content - Right side */}
        <motion.div
          className="absolute right-[6%] top-1/2 -translate-y-1/2 max-w-lg"
          initial={{ opacity: 0, x: 50 }}
          animate={{
            opacity: scrollProgress < 0.6 ? 1 : Math.max(0, 1 - (scrollProgress - 0.6) * 3),
            x: 0,
          }}
          transition={{ duration: 0.5, delay: 0.3 }}
        >
          <motion.p
            className="text-cyan-400 text-sm font-medium tracking-widest mb-3"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5 }}
          >
            ПЛАТФОРМА ДЛЯ ХОККЕИСТОВ
          </motion.p>

          <motion.h1
            className="text-5xl md:text-6xl font-bold text-white mb-4 leading-tight"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
          >
            ТВОЙ ПУТЬ
            <span className="block text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-500">
              К ВЕРШИНЕ
            </span>
          </motion.h1>

          <motion.p
            className="text-lg text-gray-300 mb-8 leading-relaxed"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.7 }}
          >
            Отслеживай свою статистику, сравнивай с другими игроками
            и стань заметным для скаутов.
          </motion.p>

          {/* Stats preview */}
          <motion.div
            className="flex gap-6"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.8 }}
          >
            <StatCounter
              label="Голов"
              value={animatedStats.goals}
              color="#00d4ff"
            />
            <StatCounter
              label="Передач"
              value={animatedStats.assists}
              color="#8b5cf6"
            />
            <StatCounter
              label="Очков"
              value={animatedStats.points}
              color="#10b981"
            />
          </motion.div>

          {/* Growth indicator */}
          <motion.div
            className="mt-8 flex items-center gap-3"
            initial={{ opacity: 0 }}
            animate={{ opacity: scrollProgress > 0.1 ? 1 : 0 }}
            transition={{ duration: 0.3 }}
          >
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-green-500/20 border border-green-500/30">
              <svg className="w-5 h-5 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
              <span className="text-green-400 font-medium">
                Рост: +{Math.floor(scrollProgress * 100)}%
              </span>
            </div>
          </motion.div>
        </motion.div>

        {/* Left side analytics panel - appears on scroll */}
        <motion.div
          className="absolute left-[6%] top-1/2 -translate-y-1/2 max-w-md"
          initial={{ opacity: 0, x: -50 }}
          animate={{
            opacity: scrollProgress > 0.4 ? Math.min(1, (scrollProgress - 0.4) * 2.5) : 0,
            x: scrollProgress > 0.4 ? 0 : -50,
          }}
          transition={{ duration: 0.5 }}
        >
          <motion.p
            className="text-purple-400 text-sm font-medium tracking-widest mb-3"
          >
            АНАЛИТИКА ИГРОКА
          </motion.p>

          <motion.h2
            className="text-3xl md:text-4xl font-bold text-white mb-4 leading-tight"
          >
            ТВОЙ
            <span className="block text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-500">
              ПРОГРЕСС
            </span>
          </motion.h2>

          <motion.p
            className="text-gray-300 mb-6 leading-relaxed"
          >
            Вся статистика в одном месте.
            Покажи свои результаты скаутам.
          </motion.p>

          {/* Player details */}
          <div className="space-y-3">
            <div className="flex items-center gap-3 px-4 py-2 rounded-lg bg-white/5 border border-white/10">
              <div className="w-2 h-2 rounded-full bg-cyan-400" />
              <span className="text-gray-400 text-sm">Возраст:</span>
              <span className="text-white font-medium ml-auto">39 лет</span>
            </div>
            <div className="flex items-center gap-3 px-4 py-2 rounded-lg bg-white/5 border border-white/10">
              <div className="w-2 h-2 rounded-full bg-green-400" />
              <span className="text-gray-400 text-sm">Рейтинг:</span>
              <span className="text-green-400 font-medium ml-auto">Топ 1%</span>
            </div>
            <div className="flex items-center gap-3 px-4 py-2 rounded-lg bg-white/5 border border-white/10">
              <div className="w-2 h-2 rounded-full bg-cyan-400" />
              <span className="text-gray-400 text-sm">Видимость:</span>
              <span className="text-cyan-400 font-medium ml-auto">Открыт скаутам ✓</span>
            </div>
          </div>
        </motion.div>

        {/* Scroll hint */}
        <motion.div
          className="absolute bottom-8 left-1/2 -translate-x-1/2"
          initial={{ opacity: 0 }}
          animate={{ opacity: animationProgress > 0.8 && scrollProgress < 0.1 ? 1 : 0 }}
          transition={{ duration: 0.5 }}
        >
          <div className="flex flex-col items-center text-cyan-300/60">
            <p className="mb-2 text-sm">Скролль для статистики</p>
            <div className="h-10 w-6 rounded-full border-2 border-cyan-400/40 p-1">
              <div className="h-2 w-1 animate-bounce rounded-full bg-cyan-400/60 mx-auto" />
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

// Stat counter component
function StatCounter({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div className="text-center">
      <p className="text-xs text-gray-500 uppercase tracking-wider mb-1">{label}</p>
      <p className="text-3xl font-bold" style={{ color }}>{value}</p>
    </div>
  )
}

// Growing chart background with hockey theme
function GrowingChartBackground({ scrollProgress }: { scrollProgress: number }) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const scrollProgressRef = useRef(scrollProgress)
  scrollProgressRef.current = scrollProgress

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    let animationId: number
    let time = 0

    const resize = () => {
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    }
    resize()
    window.addEventListener('resize', resize)

    // Chart data points - hockey stick shaped curve going up
    const generateChartPoints = (progress: number, width: number, height: number) => {
      const points: Array<{ x: number; y: number }> = []
      const numPoints = 100

      for (let i = 0; i <= numPoints; i++) {
        const t = i / numPoints
        const x = width * 0.1 + t * width * 0.85

        // Hockey stick curve - gradual then sharp rise (positioned very low)
        let y: number
        if (t < 0.55) {
          // Gradual flat section - very slight rise
          y = height * 0.92 - t * height * 0.03
        } else {
          // Sharp rise - steeper angle, goes higher
          const sharpT = (t - 0.55) / 0.45
          y = height * 0.90 - Math.pow(sharpT, 1.3) * height * 0.55
        }

        // Only show points up to scroll progress
        if (t <= progress) {
          points.push({ x, y })
        }
      }

      return points
    }

    const animate = () => {
      time += 0.02
      const width = canvas.width
      const height = canvas.height

      // Dark gradient background
      const bgGradient = ctx.createLinearGradient(0, 0, 0, height)
      bgGradient.addColorStop(0, '#020810')
      bgGradient.addColorStop(0.5, '#051020')
      bgGradient.addColorStop(1, '#030815')
      ctx.fillStyle = bgGradient
      ctx.fillRect(0, 0, width, height)

      // Grid pattern
      ctx.strokeStyle = 'rgba(0, 150, 200, 0.05)'
      ctx.lineWidth = 1

      // Vertical lines
      const gridSpacing = 60
      for (let x = 0; x < width; x += gridSpacing) {
        ctx.beginPath()
        ctx.moveTo(x, height * 0.3)
        ctx.lineTo(x, height)
        ctx.stroke()
      }

      // Horizontal lines
      for (let y = height * 0.3; y < height; y += gridSpacing) {
        ctx.beginPath()
        ctx.moveTo(0, y)
        ctx.lineTo(width, y)
        ctx.stroke()
      }

      // Perspective grid at bottom
      ctx.strokeStyle = 'rgba(0, 200, 255, 0.08)'
      const perspectiveLines = 30
      for (let i = 0; i < perspectiveLines; i++) {
        const x = (i / perspectiveLines) * width
        ctx.beginPath()
        ctx.moveTo(x, height)
        ctx.lineTo(width * 0.5, height * 0.5)
        ctx.stroke()
      }

      // === HOCKEY RINK MARKINGS ===
      const rinkCenterX = width * 0.5
      const rinkCenterY = height * 0.85
      const rinkScale = Math.min(width, height) * 0.003

      // Center ice circle (faceoff circle)
      ctx.strokeStyle = 'rgba(0, 150, 255, 0.12)'
      ctx.lineWidth = 2
      ctx.beginPath()
      ctx.ellipse(rinkCenterX, rinkCenterY, 80 * rinkScale, 40 * rinkScale, 0, 0, Math.PI * 2)
      ctx.stroke()

      // Center dot
      ctx.fillStyle = 'rgba(0, 150, 255, 0.2)'
      ctx.beginPath()
      ctx.arc(rinkCenterX, rinkCenterY, 5 * rinkScale, 0, Math.PI * 2)
      ctx.fill()

      // Blue lines (zone lines)
      ctx.strokeStyle = 'rgba(30, 100, 200, 0.1)'
      ctx.lineWidth = 3
      // Left blue line
      ctx.beginPath()
      ctx.moveTo(width * 0.3, height * 0.7)
      ctx.lineTo(width * 0.35, height)
      ctx.stroke()
      // Right blue line
      ctx.beginPath()
      ctx.moveTo(width * 0.7, height * 0.7)
      ctx.lineTo(width * 0.65, height)
      ctx.stroke()

      // Red center line
      ctx.strokeStyle = 'rgba(200, 50, 50, 0.08)'
      ctx.lineWidth = 2
      ctx.beginPath()
      ctx.moveTo(width * 0.5, height * 0.65)
      ctx.lineTo(width * 0.5, height)
      ctx.stroke()

      // === FLOATING PUCKS ===
      const pucks = [
        { x: width * 0.08, y: height * 0.2, size: 25, speed: 1.2 },
        { x: width * 0.92, y: height * 0.35, size: 20, speed: 0.8 },
        { x: width * 0.15, y: height * 0.7, size: 18, speed: 1.5 },
        { x: width * 0.88, y: height * 0.75, size: 22, speed: 1.0 },
      ]

      for (const puck of pucks) {
        const offsetY = Math.sin(time * puck.speed) * 15
        const offsetX = Math.cos(time * puck.speed * 0.7) * 8
        const rotation = time * puck.speed

        ctx.save()
        ctx.translate(puck.x + offsetX, puck.y + offsetY)
        ctx.rotate(rotation * 0.2)

        // Puck glow
        const puckGlow = ctx.createRadialGradient(0, 0, 0, 0, 0, puck.size * 1.5)
        puckGlow.addColorStop(0, 'rgba(0, 0, 0, 0.4)')
        puckGlow.addColorStop(0.5, 'rgba(0, 100, 150, 0.1)')
        puckGlow.addColorStop(1, 'transparent')
        ctx.fillStyle = puckGlow
        ctx.beginPath()
        ctx.arc(0, 0, puck.size * 1.5, 0, Math.PI * 2)
        ctx.fill()

        // Puck body (ellipse for 3D effect)
        ctx.fillStyle = 'rgba(20, 20, 30, 0.6)'
        ctx.beginPath()
        ctx.ellipse(0, 0, puck.size, puck.size * 0.4, 0, 0, Math.PI * 2)
        ctx.fill()

        // Puck edge highlight
        ctx.strokeStyle = 'rgba(100, 150, 200, 0.3)'
        ctx.lineWidth = 1
        ctx.beginPath()
        ctx.ellipse(0, 0, puck.size, puck.size * 0.4, 0, 0, Math.PI * 2)
        ctx.stroke()

        ctx.restore()
      }

      // === ICE SPARKLES ===
      const numSparkles = 40
      for (let i = 0; i < numSparkles; i++) {
        const sparkleX = (Math.sin(i * 1.7) * 0.5 + 0.5) * width
        const sparkleY = (Math.cos(i * 2.3) * 0.5 + 0.5) * height
        const sparkleOpacity = 0.1 + Math.sin(time * 4 + i * 0.8) * 0.1
        const sparkleSize = 1 + Math.sin(time * 3 + i) * 0.5

        if (sparkleOpacity > 0.05) {
          ctx.fillStyle = `rgba(200, 230, 255, ${sparkleOpacity})`
          ctx.beginPath()
          ctx.arc(sparkleX, sparkleY, sparkleSize, 0, Math.PI * 2)
          ctx.fill()
        }
      }

      // Generate chart points based on scroll
      const chartProgress = Math.min(1, scrollProgressRef.current * 1.2)
      const points = generateChartPoints(chartProgress, width, height)

      if (points.length > 1) {
        // Glow under the chart
        const glowGradient = ctx.createLinearGradient(0, height * 0.3, 0, height)
        glowGradient.addColorStop(0, 'rgba(0, 200, 180, 0.15)')
        glowGradient.addColorStop(0.5, 'rgba(0, 150, 200, 0.08)')
        glowGradient.addColorStop(1, 'rgba(0, 100, 150, 0.02)')

        ctx.beginPath()
        ctx.moveTo(points[0].x, height)
        for (const point of points) {
          ctx.lineTo(point.x, point.y)
        }
        ctx.lineTo(points[points.length - 1].x, height)
        ctx.closePath()
        ctx.fillStyle = glowGradient
        ctx.fill()

        // Main chart line - glowing
        ctx.beginPath()
        ctx.moveTo(points[0].x, points[0].y)
        for (let i = 1; i < points.length; i++) {
          ctx.lineTo(points[i].x, points[i].y)
        }

        // Outer glow
        ctx.strokeStyle = 'rgba(0, 220, 200, 0.3)'
        ctx.lineWidth = 12
        ctx.lineCap = 'round'
        ctx.lineJoin = 'round'
        ctx.stroke()

        // Middle glow
        ctx.strokeStyle = 'rgba(0, 255, 220, 0.5)'
        ctx.lineWidth = 6
        ctx.stroke()

        // Core line
        const lineGradient = ctx.createLinearGradient(points[0].x, 0, points[points.length - 1].x, 0)
        lineGradient.addColorStop(0, '#00d4aa')
        lineGradient.addColorStop(0.5, '#00e4cc')
        lineGradient.addColorStop(1, '#00ffee')
        ctx.strokeStyle = lineGradient
        ctx.lineWidth = 3
        ctx.stroke()

        // Arrow at the end
        if (chartProgress > 0.3) {
          const lastPoint = points[points.length - 1]
          const prevPoint = points[Math.max(0, points.length - 5)]
          const angle = Math.atan2(lastPoint.y - prevPoint.y, lastPoint.x - prevPoint.x)

          // Arrow head
          const arrowSize = 20
          ctx.beginPath()
          ctx.moveTo(lastPoint.x, lastPoint.y)
          ctx.lineTo(
            lastPoint.x - arrowSize * Math.cos(angle - 0.4),
            lastPoint.y - arrowSize * Math.sin(angle - 0.4)
          )
          ctx.moveTo(lastPoint.x, lastPoint.y)
          ctx.lineTo(
            lastPoint.x - arrowSize * Math.cos(angle + 0.4),
            lastPoint.y - arrowSize * Math.sin(angle + 0.4)
          )
          ctx.strokeStyle = '#00ffee'
          ctx.lineWidth = 3
          ctx.stroke()

          // Glowing dot at the end
          const dotGlow = ctx.createRadialGradient(
            lastPoint.x, lastPoint.y, 0,
            lastPoint.x, lastPoint.y, 30
          )
          dotGlow.addColorStop(0, 'rgba(0, 255, 230, 0.8)')
          dotGlow.addColorStop(0.3, 'rgba(0, 255, 230, 0.3)')
          dotGlow.addColorStop(1, 'transparent')
          ctx.fillStyle = dotGlow
          ctx.beginPath()
          ctx.arc(lastPoint.x, lastPoint.y, 30, 0, Math.PI * 2)
          ctx.fill()

          // Center dot
          ctx.fillStyle = '#ffffff'
          ctx.beginPath()
          ctx.arc(lastPoint.x, lastPoint.y, 5, 0, Math.PI * 2)
          ctx.fill()
        }

        // === FALLING PARTICLES FROM CHART LINE ===
        if (chartProgress > 0.1) {
          const numFallingParticles = Math.floor(chartProgress * 80)
          for (let i = 0; i < numFallingParticles; i++) {
            // Position along the chart line
            const linePos = (i / numFallingParticles) * (points.length - 1)
            const pointIndex = Math.floor(linePos)
            const basePoint = points[Math.min(pointIndex, points.length - 1)]
            if (!basePoint) continue

            // Each particle has its own fall animation
            const fallSpeed = 0.8 + (i % 5) * 0.2
            const fallTime = (time * fallSpeed + i * 0.3) % 2.5
            const fallDistance = fallTime * 80

            // Horizontal drift
            const driftX = Math.sin(time * 0.5 + i * 0.7) * 15

            const particleX = basePoint.x + driftX
            const particleY = basePoint.y + fallDistance

            // Only draw if particle is above the bottom
            if (particleY < height - 20) {
              // Opacity fades as particle falls
              const fadeOut = Math.max(0, 1 - fallTime / 2.5)
              const particleOpacity = fadeOut * 0.6

              // Size varies
              const particleSize = 1.5 + (i % 3) * 0.8

              if (particleOpacity > 0.05) {
                // Particle glow
                const glow = ctx.createRadialGradient(
                  particleX, particleY, 0,
                  particleX, particleY, particleSize * 4
                )
                glow.addColorStop(0, `rgba(0, 220, 200, ${particleOpacity * 0.8})`)
                glow.addColorStop(0.5, `rgba(0, 200, 180, ${particleOpacity * 0.3})`)
                glow.addColorStop(1, 'transparent')
                ctx.fillStyle = glow
                ctx.beginPath()
                ctx.arc(particleX, particleY, particleSize * 4, 0, Math.PI * 2)
                ctx.fill()

                // Particle core
                ctx.fillStyle = `rgba(150, 255, 230, ${particleOpacity})`
                ctx.beginPath()
                ctx.arc(particleX, particleY, particleSize, 0, Math.PI * 2)
                ctx.fill()
              }
            }
          }
        }

        // === VERTICAL LIGHT BARS FROM CHART (like data columns) ===
        if (chartProgress > 0.2) {
          for (let i = 0; i < points.length; i += 8) {
            const point = points[i]
            const barOpacity = 0.08 + Math.sin(time * 2 + i * 0.3) * 0.04

            // Gradient bar from point down to bottom
            const barGradient = ctx.createLinearGradient(point.x, point.y, point.x, height)
            barGradient.addColorStop(0, `rgba(0, 220, 200, ${barOpacity * 2})`)
            barGradient.addColorStop(0.3, `rgba(0, 200, 180, ${barOpacity})`)
            barGradient.addColorStop(1, 'transparent')

            ctx.fillStyle = barGradient
            ctx.fillRect(point.x - 1, point.y, 2, height - point.y)
          }
        }

        // Data points along the line
        if (chartProgress > 0.1) {
          for (let i = 0; i < points.length; i += 10) {
            const point = points[i]
            const dotOpacity = 0.3 + Math.sin(time * 3 + i * 0.5) * 0.2

            ctx.fillStyle = `rgba(0, 255, 220, ${dotOpacity})`
            ctx.beginPath()
            ctx.arc(point.x, point.y, 3, 0, Math.PI * 2)
            ctx.fill()
          }
        }

        // Floating particles around the chart
        const numParticles = Math.floor(chartProgress * 50)
        for (let i = 0; i < numParticles; i++) {
          const particleProgress = (i / numParticles)
          const basePoint = points[Math.floor(particleProgress * (points.length - 1))]
          if (!basePoint) continue

          const offsetX = Math.sin(time * 2 + i * 0.7) * 50
          const offsetY = Math.cos(time * 2.5 + i * 0.5) * 30 - 20
          const particleOpacity = 0.2 + Math.sin(time * 3 + i) * 0.15

          ctx.fillStyle = `rgba(0, 220, 200, ${particleOpacity})`
          ctx.beginPath()
          ctx.arc(basePoint.x + offsetX, basePoint.y + offsetY, 2, 0, Math.PI * 2)
          ctx.fill()
        }
      }

      // Vignette
      const vignetteGradient = ctx.createRadialGradient(
        width * 0.5, height * 0.5, width * 0.3,
        width * 0.5, height * 0.5, width * 0.9
      )
      vignetteGradient.addColorStop(0, 'transparent')
      vignetteGradient.addColorStop(1, 'rgba(0, 5, 15, 0.6)')
      ctx.fillStyle = vignetteGradient
      ctx.fillRect(0, 0, width, height)

      animationId = requestAnimationFrame(animate)
    }

    animate()

    return () => {
      window.removeEventListener('resize', resize)
      cancelAnimationFrame(animationId)
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div className="absolute inset-0">
      <canvas
        ref={canvasRef}
        className="absolute inset-0 w-full h-full"
      />
    </div>
  )
}
