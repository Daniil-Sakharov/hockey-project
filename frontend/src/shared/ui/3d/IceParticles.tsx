/* eslint-disable react-hooks/purity */
import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Points } from 'three'
import * as THREE from 'three'

interface IceParticlesProps {
  count?: number
  size?: number
  spread?: number
}

// Create particle texture outside component
const particleTexture = (() => {
  if (typeof document === 'undefined') return null
  const canvas = document.createElement('canvas')
  canvas.width = 32
  canvas.height = 32
  const ctx = canvas.getContext('2d')!

  const gradient = ctx.createRadialGradient(16, 16, 0, 16, 16, 16)
  gradient.addColorStop(0, 'rgba(255, 255, 255, 1)')
  gradient.addColorStop(0.3, 'rgba(200, 240, 255, 0.8)')
  gradient.addColorStop(1, 'rgba(200, 240, 255, 0)')

  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, 32, 32)

  return new THREE.CanvasTexture(canvas)
})()

export function IceParticles({
  count = 500,
  size = 0.05,
  spread = 15
}: IceParticlesProps) {
  const pointsRef = useRef<Points>(null)
  const dataRef = useRef<{ positions: Float32Array; velocities: Float32Array } | null>(null)

  // Initialize data once using ref (not during render)
  if (!dataRef.current) {
    const positions = new Float32Array(count * 3)
    const velocities = new Float32Array(count * 3)

    for (let i = 0; i < count; i++) {
      const i3 = i * 3
      positions[i3] = (Math.random() - 0.5) * spread * 2
      positions[i3 + 1] = Math.random() * spread
      positions[i3 + 2] = (Math.random() - 0.5) * spread * 2

      velocities[i3] = (Math.random() - 0.5) * 0.01
      velocities[i3 + 1] = -Math.random() * 0.02 - 0.005
      velocities[i3 + 2] = (Math.random() - 0.5) * 0.01
    }

    dataRef.current = { positions, velocities }
  }

  const positions = useMemo(() => dataRef.current!.positions, [])

  useFrame(() => {
    if (!pointsRef.current || !dataRef.current) return

    const { velocities } = dataRef.current
    const positionAttribute = pointsRef.current.geometry.attributes.position
    const posArray = positionAttribute.array as Float32Array

    for (let i = 0; i < count; i++) {
      const i3 = i * 3

      posArray[i3] += velocities[i3]
      posArray[i3 + 1] += velocities[i3 + 1]
      posArray[i3 + 2] += velocities[i3 + 2]

      if (posArray[i3 + 1] < -2) {
        posArray[i3] = (Math.random() - 0.5) * spread * 2
        posArray[i3 + 1] = spread
        posArray[i3 + 2] = (Math.random() - 0.5) * spread * 2
      }

      velocities[i3] += (Math.random() - 0.5) * 0.0005
      velocities[i3 + 2] += (Math.random() - 0.5) * 0.0005
    }

    positionAttribute.needsUpdate = true
  })

  if (!particleTexture) return null

  return (
    <points ref={pointsRef}>
      <bufferGeometry>
        <bufferAttribute
          attach="attributes-position"
          args={[positions, 3]}
        />
      </bufferGeometry>
      <pointsMaterial
        size={size}
        map={particleTexture}
        transparent
        opacity={0.6}
        depthWrite={false}
        blending={THREE.AdditiveBlending}
        color="#a0f0ff"
      />
    </points>
  )
}
