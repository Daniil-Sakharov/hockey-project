import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import { Text } from '@react-three/drei'
import type { Group, Mesh } from 'three'

interface BarData {
  value: number
  label: string
}

interface HoloChart3DProps {
  data: BarData[]
  position?: [number, number, number]
  rotation?: [number, number, number]
  title?: string
  color?: string
}

export function HoloBarChart3D({
  data,
  position = [0, 0, 0],
  rotation = [0, 0, 0],
  title = 'Chart',
  color = '#00d4ff',
}: HoloChart3DProps) {
  const groupRef = useRef<Group>(null)
  const barsRef = useRef<Mesh[]>([])

  const maxValue = Math.max(...data.map((d) => d.value))
  const barWidth = 0.3
  const gap = 0.15
  const totalWidth = data.length * (barWidth + gap) - gap

  useFrame((state) => {
    if (groupRef.current) {
      // Subtle floating animation
      groupRef.current.position.y = position[1] + Math.sin(state.clock.elapsedTime * 0.5) * 0.1
    }

    // Animate bars
    barsRef.current.forEach((bar, i) => {
      if (bar) {
        const targetHeight = (data[i].value / maxValue) * 2
        bar.scale.y += (targetHeight - bar.scale.y) * 0.05
        bar.position.y = bar.scale.y / 2
      }
    })
  })

  return (
    <group ref={groupRef} position={position} rotation={rotation}>
      {/* Base platform */}
      <mesh position={[0, -0.05, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[totalWidth + 1, 1.5]} />
        <meshStandardMaterial
          color={color}
          transparent
          opacity={0.1}
          emissive={color}
          emissiveIntensity={0.5}
        />
      </mesh>

      {/* Grid lines */}
      {[0, 0.5, 1, 1.5, 2].map((y, i) => (
        <mesh key={i} position={[0, y, 0]}>
          <boxGeometry args={[totalWidth + 0.5, 0.005, 0.005]} />
          <meshStandardMaterial color={color} transparent opacity={0.3} />
        </mesh>
      ))}

      {/* Bars */}
      {data.map((item, index) => {
        const x = index * (barWidth + gap) - totalWidth / 2 + barWidth / 2
        return (
          <group key={index} position={[x, 0, 0]}>
            {/* Main bar */}
            <mesh
              ref={(el) => { if (el) barsRef.current[index] = el }}
              position={[0, 0.5, 0]}
            >
              <boxGeometry args={[barWidth, 1, 0.2]} />
              <meshStandardMaterial
                color={color}
                emissive={color}
                emissiveIntensity={0.8}
                transparent
                opacity={0.9}
              />
            </mesh>

            {/* Bar glow */}
            <mesh position={[0, 0.5, 0]}>
              <boxGeometry args={[barWidth + 0.05, 1, 0.25]} />
              <meshStandardMaterial
                color={color}
                transparent
                opacity={0.2}
              />
            </mesh>

            {/* Value label */}
            <Text
              position={[0, 2.3, 0]}
              fontSize={0.15}
              color={color}
              anchorX="center"
            >
              {item.value}
            </Text>

            {/* Category label */}
            <Text
              position={[0, -0.2, 0]}
              fontSize={0.1}
              color="#666"
              anchorX="center"
            >
              {item.label}
            </Text>
          </group>
        )
      })}

      {/* Title */}
      <Text
        position={[0, 2.8, 0]}
        fontSize={0.2}
        color="#ffffff"
        anchorX="center"
      >
        {title}
      </Text>
    </group>
  )
}

export function HoloRingChart3D({
  position = [0, 0, 0],
  rotation = [0, 0, 0],
  value = 75,
  label = 'Progress',
  color = '#00d4ff',
}: {
  position?: [number, number, number]
  rotation?: [number, number, number]
  value?: number
  label?: string
  color?: string
}) {
  const groupRef = useRef<Group>(null)
  const ringRef = useRef<Mesh>(null)

  useFrame((state) => {
    if (groupRef.current) {
      groupRef.current.rotation.y += 0.005
      groupRef.current.position.y = position[1] + Math.sin(state.clock.elapsedTime * 0.7) * 0.08
    }
  })

  const segments = 64
  const filledSegments = Math.floor((value / 100) * segments)

  return (
    <group ref={groupRef} position={position} rotation={rotation}>
      {/* Background ring */}
      <mesh rotation={[-Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1, 0.08, 16, segments]} />
        <meshStandardMaterial color="#1a1a2e" transparent opacity={0.5} />
      </mesh>

      {/* Filled ring */}
      <mesh ref={ringRef} rotation={[-Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1, 0.1, 16, filledSegments, (value / 100) * Math.PI * 2]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={1.5}
        />
      </mesh>

      {/* Center value */}
      <Text
        position={[0, 0, 0]}
        fontSize={0.4}
        color="#ffffff"
        anchorX="center"
        anchorY="middle"
      >
        {value}%
      </Text>

      {/* Label */}
      <Text
        position={[0, -0.4, 0]}
        fontSize={0.15}
        color="#666"
        anchorX="center"
      >
        {label}
      </Text>
    </group>
  )
}
