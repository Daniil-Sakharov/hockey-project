import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Points } from 'three'
import * as THREE from 'three'

interface IceParticlesProps {
  count?: number
  size?: number
  spread?: number
}

export function IceParticles({
  count = 500,
  size = 0.05,
  spread = 15
}: IceParticlesProps) {
  const pointsRef = useRef<Points>(null)

  const { positions, velocities } = useMemo(() => {
    const positions = new Float32Array(count * 3)
    const velocities = new Float32Array(count * 3)

    for (let i = 0; i < count; i++) {
      const i3 = i * 3
      // Random positions in a box
      positions[i3] = (Math.random() - 0.5) * spread * 2
      positions[i3 + 1] = Math.random() * spread
      positions[i3 + 2] = (Math.random() - 0.5) * spread * 2

      // Random velocities (mostly downward with some drift)
      velocities[i3] = (Math.random() - 0.5) * 0.01
      velocities[i3 + 1] = -Math.random() * 0.02 - 0.005
      velocities[i3 + 2] = (Math.random() - 0.5) * 0.01
    }

    return { positions, velocities }
  }, [count, spread])

  useFrame(() => {
    if (!pointsRef.current) return

    const positionAttribute = pointsRef.current.geometry.attributes.position
    const posArray = positionAttribute.array as Float32Array

    for (let i = 0; i < count; i++) {
      const i3 = i * 3

      // Update positions
      posArray[i3] += velocities[i3]
      posArray[i3 + 1] += velocities[i3 + 1]
      posArray[i3 + 2] += velocities[i3 + 2]

      // Reset particle if it goes below ground
      if (posArray[i3 + 1] < -2) {
        posArray[i3] = (Math.random() - 0.5) * spread * 2
        posArray[i3 + 1] = spread
        posArray[i3 + 2] = (Math.random() - 0.5) * spread * 2
      }

      // Add slight horizontal drift
      velocities[i3] += (Math.random() - 0.5) * 0.0005
      velocities[i3 + 2] += (Math.random() - 0.5) * 0.0005
    }

    positionAttribute.needsUpdate = true
  })

  const particleTexture = useMemo(() => {
    const canvas = document.createElement('canvas')
    canvas.width = 32
    canvas.height = 32
    const ctx = canvas.getContext('2d')!

    // Create radial gradient for soft particle
    const gradient = ctx.createRadialGradient(16, 16, 0, 16, 16, 16)
    gradient.addColorStop(0, 'rgba(255, 255, 255, 1)')
    gradient.addColorStop(0.3, 'rgba(200, 240, 255, 0.8)')
    gradient.addColorStop(1, 'rgba(200, 240, 255, 0)')

    ctx.fillStyle = gradient
    ctx.fillRect(0, 0, 32, 32)

    return new THREE.CanvasTexture(canvas)
  }, [])

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
