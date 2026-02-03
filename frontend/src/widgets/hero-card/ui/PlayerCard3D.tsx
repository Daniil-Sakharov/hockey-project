import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import { RoundedBox, Text } from '@react-three/drei'
import * as THREE from 'three'

interface PlayerCardData {
  name: string
  number: number
  team: string
  teamColor: string
  position: string
  photo?: string
  stats: {
    goals: number
    assists: number
    points: number
    games: number
  }
}

interface PlayerCard3DProps {
  player: PlayerCardData
  animationProgress: number
  mousePosition: { x: number; y: number }
  scrollProgress: number
  horizontalOffset: number
}

export function PlayerCard3D({
  player,
  animationProgress,
  mousePosition,
  scrollProgress,
  horizontalOffset,
}: PlayerCard3DProps) {
  const groupRef = useRef<THREE.Group>(null)
  const cardRef = useRef<THREE.Mesh>(null)
  const materialRef = useRef<THREE.MeshPhysicalMaterial>(null)

  // Smooth values
  const rotationRef = useRef({ x: 0, y: 0 })
  const positionRef = useRef({ x: 0 })
  const scaleRef = useRef(0)

  // Card dimensions
  const cardWidth = 2.4
  const cardHeight = 3.4
  const cardDepth = 0.08

  useFrame((state, delta) => {
    if (!groupRef.current) return

    const time = state.clock.elapsedTime

    // Smooth mouse following
    const targetRotX = mousePosition.y * 0.15
    const targetRotY = mousePosition.x * 0.15

    rotationRef.current.x += (targetRotX - rotationRef.current.x) * delta * 3
    rotationRef.current.y += (targetRotY - rotationRef.current.y) * delta * 3

    // Smooth horizontal position
    const targetX = horizontalOffset * 3
    positionRef.current.x += (targetX - positionRef.current.x) * delta * 2.5

    // Scroll-based rotation for 3D demonstration
    const scrollRotation = Math.sin(scrollProgress * Math.PI) * 0.2

    groupRef.current.rotation.x = rotationRef.current.x
    groupRef.current.rotation.y = rotationRef.current.y + scrollRotation

    // Subtle idle animation
    groupRef.current.rotation.z = Math.sin(time * 0.5) * 0.01

    groupRef.current.position.x = positionRef.current.x
    groupRef.current.scale.setScalar(1)
    groupRef.current.position.z = 0

    // Update iridescence
    if (materialRef.current) {
      const angle = Math.abs(rotationRef.current.x) + Math.abs(rotationRef.current.y)
      materialRef.current.iridescence = 0.8 + Math.sin(time * 2 + angle * 5) * 0.2
    }
  })

  return (
    <group ref={groupRef} scale={1}>
      {/* Main card body */}
      <RoundedBox
        ref={cardRef}
        args={[cardWidth, cardHeight, cardDepth]}
        radius={0.12}
        smoothness={4}
      >
        <meshPhysicalMaterial
          ref={materialRef}
          color="#0a0a1a"
          metalness={0.95}
          roughness={0.05}
          iridescence={1}
          iridescenceIOR={1.3}
          iridescenceThicknessRange={[100, 800]}
          clearcoat={1}
          clearcoatRoughness={0}
          reflectivity={1}
          envMapIntensity={1}
        />
      </RoundedBox>

      {/* Card front content */}
      <group position={[0, 0, cardDepth / 2 + 0.001]}>
        {/* Dark background */}
        <mesh position={[0, 0, 0]}>
          <planeGeometry args={[cardWidth - 0.15, cardHeight - 0.15]} />
          <meshBasicMaterial color="#0a0f18" />
        </mesh>

        {/* Team color accent at top */}
        <mesh position={[0, cardHeight / 2 - 0.22, 0.001]}>
          <planeGeometry args={[cardWidth - 0.15, 0.3]} />
          <meshBasicMaterial color={player.teamColor} />
        </mesh>

        {/* Player avatar area */}
        <mesh position={[0, 0.55, 0.001]}>
          <planeGeometry args={[1.5, 1.5]} />
          <meshBasicMaterial color="#151c28" />
        </mesh>

        {/* Player silhouette circle */}
        <mesh position={[0, 0.6, 0.002]}>
          <circleGeometry args={[0.5, 32]} />
          <meshBasicMaterial color="#1a2535" />
        </mesh>

        {/* Player number large - as avatar */}
        <Text
          position={[0, 0.55, 0.003]}
          fontSize={0.7}
          color={player.teamColor}
          anchorX="center"
          anchorY="middle"
          outlineWidth={0.015}
          outlineColor="#000000"
        >
          {player.number}
        </Text>

        {/* Player number - top right */}
        <Text
          position={[cardWidth / 2 - 0.35, cardHeight / 2 - 0.22, 0.01]}
          fontSize={0.25}
          color="#ffffff"
          anchorX="center"
          anchorY="middle"
          fontWeight="bold"
        >
          #{player.number}
        </Text>

        {/* Position badge - top left */}
        <group position={[-cardWidth / 2 + 0.3, cardHeight / 2 - 0.22, 0.01]}>
          <mesh>
            <circleGeometry args={[0.12, 32]} />
            <meshBasicMaterial color="#ffffff" />
          </mesh>
          <Text
            position={[0, 0, 0.01]}
            fontSize={0.09}
            color="#000000"
            anchorX="center"
            anchorY="middle"
          >
            {player.position}
          </Text>
        </group>

        {/* Player name */}
        <Text
          position={[0, -0.35, 0.01]}
          fontSize={0.18}
          color="#ffffff"
          anchorX="center"
          anchorY="middle"
          maxWidth={cardWidth - 0.4}
        >
          {player.name}
        </Text>

        {/* Team name */}
        <Text
          position={[0, -0.6, 0.01]}
          fontSize={0.1}
          color="#888888"
          anchorX="center"
          anchorY="middle"
        >
          {player.team}
        </Text>

        {/* Divider line */}
        <mesh position={[0, -0.8, 0.01]}>
          <planeGeometry args={[cardWidth - 0.4, 0.006]} />
          <meshBasicMaterial color={player.teamColor} transparent opacity={0.6} />
        </mesh>

        {/* Stats row */}
        <group position={[0, -1.1, 0.01]}>
          <StatItem label="ГОЛЫ" value={player.stats.goals} position={[-0.7, 0]} color="#00d4ff" />
          <StatItem label="ПАС" value={player.stats.assists} position={[0, 0]} color="#ffffff" />
          <StatItem label="ОЧКИ" value={player.stats.points} position={[0.7, 0]} color="#10b981" />
        </group>

        {/* Games played - bottom */}
        <Text
          position={[0, -1.45, 0.01]}
          fontSize={0.08}
          color="#555555"
          anchorX="center"
          anchorY="middle"
        >
          {player.stats.games} игр сыграно
        </Text>
      </group>

      {/* Card edge glow */}
      <mesh position={[0, 0, -cardDepth / 2 - 0.01]}>
        <planeGeometry args={[cardWidth + 0.1, cardHeight + 0.1]} />
        <meshBasicMaterial
          color={player.teamColor}
          transparent
          opacity={0.3}
          blending={THREE.AdditiveBlending}
        />
      </mesh>

      {/* Card lighting */}
      <pointLight
        position={[0, 0, 2]}
        intensity={0.5}
        color="#ffffff"
        distance={5}
      />
    </group>
  )
}

// Stat item component
function StatItem({
  label,
  value,
  position,
  color = '#ffffff',
}: {
  label: string
  value: number
  position: [number, number]
  color?: string
}) {
  return (
    <group position={[position[0], position[1], 0]}>
      <Text
        position={[0, 0.12, 0]}
        fontSize={0.06}
        color="#555555"
        anchorX="center"
        anchorY="middle"
      >
        {label}
      </Text>
      <Text
        position={[0, -0.05, 0]}
        fontSize={0.22}
        color={color}
        anchorX="center"
        anchorY="middle"
      >
        {value}
      </Text>
    </group>
  )
}

// Easing function
function easeOutExpo(t: number): number {
  return t === 1 ? 1 : 1 - Math.pow(2, -10 * t)
}
