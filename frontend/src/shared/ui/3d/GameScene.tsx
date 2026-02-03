import { useRef, useState, useEffect } from 'react'
import { useFrame, useThree } from '@react-three/fiber'
import { Html } from '@react-three/drei'
import { Player } from './Player'
import { Goal } from './Goal'
import * as THREE from 'three'

interface PlayerData {
  id: number
  name: string
  number: number
  team: string
  color: string
  goals: number
  assists: number
  games: number
}

const players: PlayerData[] = [
  { id: 1, name: 'Иван Петров', number: 10, team: 'Динамо', color: '#00d4ff', goals: 15, assists: 23, games: 42 },
  { id: 2, name: 'Алексей Смирнов', number: 77, team: 'ЦСКА', color: '#ff3366', goals: 22, assists: 18, games: 45 },
  { id: 3, name: 'Дмитрий Козлов', number: 91, team: 'СКА', color: '#8b5cf6', goals: 18, assists: 31, games: 48 },
  { id: 4, name: 'Михаил Новиков', number: 23, team: 'Спартак', color: '#ec4899', goals: 12, assists: 15, games: 38 },
]

interface SkatingPlayerProps {
  data: PlayerData
  pathRadius: number
  speed: number
  offset: number
  isAnalyzing: boolean
}

function SkatingPlayer({ data, pathRadius, speed, offset, isAnalyzing }: SkatingPlayerProps) {
  const groupRef = useRef<THREE.Group>(null)
  const [position, setPosition] = useState<[number, number, number]>([0, 0, 0])
  const [rotation, setRotation] = useState<[number, number, number]>([0, 0, 0])
  const angleRef = useRef(offset)

  useFrame((_, delta) => {
    if (isAnalyzing) return

    angleRef.current += delta * speed
    const x = Math.sin(angleRef.current) * pathRadius
    const z = Math.cos(angleRef.current) * pathRadius

    // Calculate direction of movement for rotation
    const nextX = Math.sin(angleRef.current + 0.1) * pathRadius
    const nextZ = Math.cos(angleRef.current + 0.1) * pathRadius
    const angle = Math.atan2(nextX - x, nextZ - z)

    setPosition([x, 0, z])
    setRotation([0, angle, 0])
  })

  return (
    <group ref={groupRef}>
      <Player
        position={position}
        rotation={rotation}
        color={data.color}
        playerNumber={data.number}
        isAnalyzing={isAnalyzing}
      />
    </group>
  )
}

interface PlayerAnalyticsProps {
  player: PlayerData
  visible: boolean
}

function PlayerAnalytics({ player, visible }: PlayerAnalyticsProps) {
  if (!visible) return null

  return (
    <Html center distanceFactor={10} style={{ pointerEvents: 'none' }}>
      <div className="animate-fadeIn rounded-xl border border-[#00d4ff]/30 bg-[#0a0e1a]/90 p-4 backdrop-blur-xl"
        style={{ width: 220 }}>
        <div className="mb-3 flex items-center gap-3">
          <div
            className="flex h-10 w-10 items-center justify-center rounded-lg text-lg font-bold text-white"
            style={{ backgroundColor: player.color }}
          >
            {player.number}
          </div>
          <div>
            <div className="font-semibold text-white">{player.name}</div>
            <div className="text-xs text-gray-400">{player.team}</div>
          </div>
        </div>

        <div className="space-y-2">
          <StatBar label="Голы" value={player.goals} max={30} color="#00d4ff" />
          <StatBar label="Ассисты" value={player.assists} max={35} color="#8b5cf6" />
          <StatBar label="Игры" value={player.games} max={50} color="#ec4899" />
        </div>

        <div className="mt-3 text-center text-xs text-[#00d4ff]">
          Анализ игрока...
        </div>
      </div>
    </Html>
  )
}

function StatBar({ label, value, max, color }: { label: string; value: number; max: number; color: string }) {
  const percentage = (value / max) * 100
  return (
    <div>
      <div className="mb-1 flex justify-between text-xs">
        <span className="text-gray-400">{label}</span>
        <span className="font-medium text-white">{value}</span>
      </div>
      <div className="h-1.5 overflow-hidden rounded-full bg-gray-700">
        <div
          className="h-full rounded-full transition-all duration-500"
          style={{ width: `${percentage}%`, backgroundColor: color, boxShadow: `0 0 8px ${color}` }}
        />
      </div>
    </div>
  )
}

export function GameScene() {
  const { camera } = useThree()
  const [analyzingPlayer, setAnalyzingPlayer] = useState<number | null>(null)
  const [isTransitioning, setIsTransitioning] = useState(false)
  const targetPosition = useRef(new THREE.Vector3(0, 3, 8))
  const originalPosition = useRef(new THREE.Vector3(0, 3, 8))

  // Cycle through players for analysis
  useEffect(() => {
    const interval = setInterval(() => {
      if (isTransitioning) return

      setIsTransitioning(true)

      // Pick random player
      const playerIndex = Math.floor(Math.random() * players.length)
      setAnalyzingPlayer(playerIndex)

      // Reset after 4 seconds
      setTimeout(() => {
        setAnalyzingPlayer(null)
        setIsTransitioning(false)
      }, 4000)
    }, 8000)

    return () => clearInterval(interval)
  }, [isTransitioning])

  // Camera animation
  useFrame(() => {
    if (analyzingPlayer !== null) {
      // Zoom in effect
      targetPosition.current.set(0, 2, 4)
    } else {
      targetPosition.current.copy(originalPosition.current)
    }

    camera.position.lerp(targetPosition.current, 0.02)
  })

  return (
    <group>
      {/* Goals at both ends */}
      <Goal position={[0, 0, -5]} rotation={[0, 0, 0]} scale={1.2} />
      <Goal position={[0, 0, 5]} rotation={[0, Math.PI, 0]} scale={1.2} />

      {/* Skating players */}
      {players.map((player, index) => (
        <group key={player.id}>
          <SkatingPlayer
            data={player}
            pathRadius={2.5 + index * 0.5}
            speed={0.5 + index * 0.1}
            offset={(index * Math.PI) / 2}
            isAnalyzing={analyzingPlayer === index}
          />
          <PlayerAnalytics player={player} visible={analyzingPlayer === index} />
        </group>
      ))}
    </group>
  )
}
