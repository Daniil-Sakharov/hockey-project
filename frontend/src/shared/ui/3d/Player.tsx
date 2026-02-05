import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group } from 'three'

interface PlayerProps {
  position?: [number, number, number]
  rotation?: [number, number, number]
  color?: string
  playerNumber?: number
  isAnalyzing?: boolean
}

export function Player({
  position = [0, 0, 0],
  rotation = [0, 0, 0],
  color = '#00d4ff',
  isAnalyzing = false,
}: PlayerProps) {
  const groupRef = useRef<Group>(null)
  const stickRef = useRef<Group>(null)

  useFrame((state) => {
    if (!groupRef.current) return

    if (isAnalyzing) {
      // Hover effect when analyzing
      groupRef.current.position.y = position[1] + Math.sin(state.clock.elapsedTime * 3) * 0.05
    }

    // Skating motion for stick
    if (stickRef.current && !isAnalyzing) {
      stickRef.current.rotation.x = Math.sin(state.clock.elapsedTime * 4) * 0.1 - 0.3
    }
  })

  const glowIntensity = isAnalyzing ? 2 : 0.5

  return (
    <group ref={groupRef} position={position} rotation={rotation}>
      {/* Body - torso */}
      <mesh position={[0, 0.6, 0]}>
        <capsuleGeometry args={[0.15, 0.3, 8, 16]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={glowIntensity * 0.3}
        />
      </mesh>

      {/* Jersey number */}
      <mesh position={[0, 0.65, 0.16]} rotation={[0, 0, 0]}>
        <planeGeometry args={[0.15, 0.15]} />
        <meshStandardMaterial
          color="#ffffff"
          emissive="#ffffff"
          emissiveIntensity={0.5}
          transparent
          opacity={0.9}
        />
      </mesh>

      {/* Head with helmet */}
      <mesh position={[0, 1, 0]}>
        <sphereGeometry args={[0.12, 16, 16]} />
        <meshStandardMaterial color="#1a1a2e" metalness={0.8} roughness={0.2} />
      </mesh>

      {/* Helmet visor */}
      <mesh position={[0, 0.98, 0.1]}>
        <boxGeometry args={[0.2, 0.06, 0.05]} />
        <meshStandardMaterial
          color="#00d4ff"
          emissive="#00d4ff"
          emissiveIntensity={glowIntensity}
          transparent
          opacity={0.7}
        />
      </mesh>

      {/* Left leg */}
      <mesh position={[-0.08, 0.2, 0]} rotation={[0.1, 0, 0]}>
        <capsuleGeometry args={[0.06, 0.25, 8, 16]} />
        <meshStandardMaterial color="#1a1a2e" />
      </mesh>

      {/* Right leg */}
      <mesh position={[0.08, 0.2, 0]} rotation={[-0.1, 0, 0]}>
        <capsuleGeometry args={[0.06, 0.25, 8, 16]} />
        <meshStandardMaterial color="#1a1a2e" />
      </mesh>

      {/* Skates */}
      {[-0.08, 0.08].map((x, i) => (
        <group key={i} position={[x, 0.02, 0]}>
          <mesh>
            <boxGeometry args={[0.08, 0.04, 0.2]} />
            <meshStandardMaterial color="#333" metalness={0.9} roughness={0.1} />
          </mesh>
          {/* Blade */}
          <mesh position={[0, -0.03, 0]}>
            <boxGeometry args={[0.01, 0.02, 0.22]} />
            <meshStandardMaterial
              color="#aaa"
              metalness={1}
              roughness={0}
              emissive="#00d4ff"
              emissiveIntensity={0.2}
            />
          </mesh>
        </group>
      ))}

      {/* Arms */}
      <mesh position={[-0.22, 0.65, 0]} rotation={[0, 0, 0.5]}>
        <capsuleGeometry args={[0.05, 0.2, 8, 16]} />
        <meshStandardMaterial color={color} />
      </mesh>
      <mesh position={[0.22, 0.65, 0]} rotation={[0, 0, -0.5]}>
        <capsuleGeometry args={[0.05, 0.2, 8, 16]} />
        <meshStandardMaterial color={color} />
      </mesh>

      {/* Hockey stick */}
      <group ref={stickRef} position={[0.3, 0.4, 0.1]} rotation={[-0.3, 0, -0.3]}>
        {/* Shaft */}
        <mesh rotation={[0, 0, 0]}>
          <cylinderGeometry args={[0.015, 0.015, 1, 8]} />
          <meshStandardMaterial color="#222" />
        </mesh>
        {/* Blade */}
        <mesh position={[0, -0.5, 0.1]} rotation={[Math.PI / 2, 0, 0]}>
          <boxGeometry args={[0.03, 0.25, 0.08]} />
          <meshStandardMaterial color="#111" />
        </mesh>
      </group>

      {/* Glow ring when analyzing */}
      {isAnalyzing && (
        <>
          <mesh position={[0, 0.01, 0]} rotation={[-Math.PI / 2, 0, 0]}>
            <ringGeometry args={[0.4, 0.5, 32]} />
            <meshStandardMaterial
              color="#00d4ff"
              emissive="#00d4ff"
              emissiveIntensity={2}
              transparent
              opacity={0.6}
            />
          </mesh>
          <pointLight position={[0, 1, 0]} color="#00d4ff" intensity={2} distance={3} />
        </>
      )}
    </group>
  )
}
