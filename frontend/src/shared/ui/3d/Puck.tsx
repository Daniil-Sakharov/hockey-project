import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import { MeshDistortMaterial } from '@react-three/drei'
import type { Mesh, Group } from 'three'

interface PuckProps {
  position?: [number, number, number]
  scale?: number
}

export function Puck({ position = [0, 0, 0], scale = 1 }: PuckProps) {
  const groupRef = useRef<Group>(null)
  const glowRef = useRef<Mesh>(null)

  useFrame((state) => {
    if (groupRef.current) {
      // Rotating animation
      groupRef.current.rotation.y += 0.01
      // Floating animation
      groupRef.current.position.y = position[1] + Math.sin(state.clock.elapsedTime) * 0.1
    }
    if (glowRef.current) {
      // Pulsing glow
      const pulse = 1 + Math.sin(state.clock.elapsedTime * 2) * 0.1
      glowRef.current.scale.setScalar(pulse * 1.2)
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main puck body */}
      <mesh castShadow receiveShadow>
        <cylinderGeometry args={[0.8, 0.8, 0.25, 64]} />
        <meshStandardMaterial
          color="#1a1a2e"
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>

      {/* Top edge ring - neon blue */}
      <mesh position={[0, 0.125, 0]}>
        <torusGeometry args={[0.75, 0.02, 16, 100]} />
        <meshStandardMaterial
          color="#00d4ff"
          emissive="#00d4ff"
          emissiveIntensity={2}
        />
      </mesh>

      {/* Bottom edge ring - neon cyan */}
      <mesh position={[0, -0.125, 0]}>
        <torusGeometry args={[0.75, 0.02, 16, 100]} />
        <meshStandardMaterial
          color="#00ffff"
          emissive="#00ffff"
          emissiveIntensity={2}
        />
      </mesh>

      {/* Center logo area */}
      <mesh position={[0, 0.13, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <circleGeometry args={[0.4, 64]} />
        <meshStandardMaterial
          color="#0d1224"
          metalness={0.9}
          roughness={0.1}
        />
      </mesh>

      {/* Outer glow effect */}
      <mesh ref={glowRef} position={[0, 0, 0]}>
        <cylinderGeometry args={[1, 1, 0.3, 64]} />
        <MeshDistortMaterial
          color="#00d4ff"
          transparent
          opacity={0.15}
          distort={0.2}
          speed={2}
        />
      </mesh>

      {/* Inner holographic ring */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0.14, 0]}>
        <ringGeometry args={[0.45, 0.55, 64]} />
        <meshStandardMaterial
          color="#8b5cf6"
          emissive="#8b5cf6"
          emissiveIntensity={1}
          transparent
          opacity={0.8}
        />
      </mesh>
    </group>
  )
}
