import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import type { Group } from 'three'

interface GoalProps {
  position?: [number, number, number]
  rotation?: [number, number, number]
  scale?: number
}

export function Goal({ position = [0, 0, 0], rotation = [0, 0, 0], scale = 1 }: GoalProps) {
  const groupRef = useRef<Group>(null)
  const netGlowRef = useRef<Group>(null)

  useFrame((state) => {
    if (netGlowRef.current) {
      // Subtle pulsing glow on the net
      const pulse = 0.3 + Math.sin(state.clock.elapsedTime * 2) * 0.1
      netGlowRef.current.children.forEach((child) => {
        const mesh = child as THREE.Mesh
        if (mesh.material && 'opacity' in mesh.material) {
          (mesh.material as THREE.MeshStandardMaterial).opacity = pulse
        }
      })
    }
  })

  const pipeRadius = 0.05
  const goalWidth = 1.8
  const goalHeight = 1.2
  const goalDepth = 0.9

  return (
    <group ref={groupRef} position={position} rotation={rotation} scale={scale}>
      {/* Main frame - red pipes */}
      {/* Top crossbar */}
      <mesh position={[0, goalHeight, 0]} rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[pipeRadius, pipeRadius, goalWidth, 16]} />
        <meshStandardMaterial color="#ff3366" emissive="#ff3366" emissiveIntensity={0.3} />
      </mesh>

      {/* Left post */}
      <mesh position={[-goalWidth / 2, goalHeight / 2, 0]}>
        <cylinderGeometry args={[pipeRadius, pipeRadius, goalHeight, 16]} />
        <meshStandardMaterial color="#ff3366" emissive="#ff3366" emissiveIntensity={0.3} />
      </mesh>

      {/* Right post */}
      <mesh position={[goalWidth / 2, goalHeight / 2, 0]}>
        <cylinderGeometry args={[pipeRadius, pipeRadius, goalHeight, 16]} />
        <meshStandardMaterial color="#ff3366" emissive="#ff3366" emissiveIntensity={0.3} />
      </mesh>

      {/* Back frame */}
      <mesh position={[0, goalHeight, -goalDepth]} rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[pipeRadius * 0.7, pipeRadius * 0.7, goalWidth, 16]} />
        <meshStandardMaterial color="#666" metalness={0.8} roughness={0.2} />
      </mesh>

      {/* Back verticals */}
      {[-goalWidth / 2, goalWidth / 2].map((x, i) => (
        <mesh key={i} position={[x, goalHeight / 2, -goalDepth]}>
          <cylinderGeometry args={[pipeRadius * 0.7, pipeRadius * 0.7, goalHeight, 16]} />
          <meshStandardMaterial color="#666" metalness={0.8} roughness={0.2} />
        </mesh>
      ))}

      {/* Depth bars */}
      {[-goalWidth / 2, goalWidth / 2].map((x, i) => (
        <mesh key={i} position={[x, goalHeight, -goalDepth / 2]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[pipeRadius * 0.7, pipeRadius * 0.7, goalDepth, 16]} />
          <meshStandardMaterial color="#666" metalness={0.8} roughness={0.2} />
        </mesh>
      ))}

      {/* Net visualization with glow */}
      <group ref={netGlowRef}>
        {/* Back net */}
        <mesh position={[0, goalHeight / 2, -goalDepth]}>
          <planeGeometry args={[goalWidth - 0.1, goalHeight - 0.1, 8, 8]} />
          <meshStandardMaterial
            color="#00d4ff"
            wireframe
            transparent
            opacity={0.3}
          />
        </mesh>

        {/* Side nets */}
        {[-goalWidth / 2, goalWidth / 2].map((x, i) => (
          <mesh key={i} position={[x, goalHeight / 2, -goalDepth / 2]} rotation={[0, Math.PI / 2, 0]}>
            <planeGeometry args={[goalDepth, goalHeight - 0.1, 4, 8]} />
            <meshStandardMaterial
              color="#00d4ff"
              wireframe
              transparent
              opacity={0.2}
            />
          </mesh>
        ))}

        {/* Top net */}
        <mesh position={[0, goalHeight, -goalDepth / 2]} rotation={[Math.PI / 2, 0, 0]}>
          <planeGeometry args={[goalWidth - 0.1, goalDepth, 8, 4]} />
          <meshStandardMaterial
            color="#00d4ff"
            wireframe
            transparent
            opacity={0.2}
          />
        </mesh>
      </group>

      {/* Goal line glow */}
      <mesh position={[0, 0.01, 0.1]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[goalWidth + 0.5, 0.08]} />
        <meshStandardMaterial
          color="#ff3366"
          emissive="#ff3366"
          emissiveIntensity={1}
          transparent
          opacity={0.8}
        />
      </mesh>
    </group>
  )
}
