import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import { Text } from '@react-three/drei'
import type { Group, Mesh } from 'three'

// Stadium Boards (Борта)
export function StadiumBoards() {
  const rinkLength = 12
  const rinkWidth = 6
  const boardHeight = 0.5
  const glassHeight = 0.8

  return (
    <group>
      {/* Main boards - 4 sides */}
      {/* Long sides */}
      {[-rinkWidth / 2, rinkWidth / 2].map((z, i) => (
        <group key={`long-${i}`} position={[0, 0, z]}>
          {/* Board base */}
          <mesh position={[0, boardHeight / 2, 0]}>
            <boxGeometry args={[rinkLength, boardHeight, 0.15]} />
            <meshStandardMaterial color="#1a1a2e" />
          </mesh>
          {/* Red stripe */}
          <mesh position={[0, boardHeight - 0.05, z > 0 ? -0.08 : 0.08]}>
            <boxGeometry args={[rinkLength, 0.1, 0.02]} />
            <meshStandardMaterial color="#ff3366" emissive="#ff3366" emissiveIntensity={0.5} />
          </mesh>
          {/* Glass */}
          <mesh position={[0, boardHeight + glassHeight / 2, 0]}>
            <boxGeometry args={[rinkLength, glassHeight, 0.05]} />
            <meshStandardMaterial color="#88ccff" transparent opacity={0.15} />
          </mesh>
        </group>
      ))}

      {/* Short sides (behind goals) */}
      {[-rinkLength / 2, rinkLength / 2].map((x, i) => (
        <group key={`short-${i}`} position={[x, 0, 0]} rotation={[0, Math.PI / 2, 0]}>
          <mesh position={[0, boardHeight / 2, 0]}>
            <boxGeometry args={[rinkWidth, boardHeight, 0.15]} />
            <meshStandardMaterial color="#1a1a2e" />
          </mesh>
          <mesh position={[0, boardHeight + glassHeight / 2, 0]}>
            <boxGeometry args={[rinkWidth, glassHeight, 0.05]} />
            <meshStandardMaterial color="#88ccff" transparent opacity={0.15} />
          </mesh>
        </group>
      ))}

      {/* Corner curves (simplified as angled boards) */}
      {[
        { pos: [-rinkLength / 2 + 0.5, 0, -rinkWidth / 2 + 0.5], rot: Math.PI / 4 },
        { pos: [rinkLength / 2 - 0.5, 0, -rinkWidth / 2 + 0.5], rot: -Math.PI / 4 },
        { pos: [-rinkLength / 2 + 0.5, 0, rinkWidth / 2 - 0.5], rot: -Math.PI / 4 },
        { pos: [rinkLength / 2 - 0.5, 0, rinkWidth / 2 - 0.5], rot: Math.PI / 4 },
      ].map((corner, i) => (
        <mesh key={i} position={[corner.pos[0], boardHeight / 2, corner.pos[2]]} rotation={[0, corner.rot, 0]}>
          <boxGeometry args={[1.2, boardHeight, 0.15]} />
          <meshStandardMaterial color="#1a1a2e" />
        </mesh>
      ))}

      {/* Advertising boards */}
      {[-3, 0, 3].map((x, i) => (
        <mesh key={i} position={[x, boardHeight / 2, -rinkWidth / 2 + 0.1]}>
          <planeGeometry args={[2, 0.4]} />
          <meshStandardMaterial
            color={i === 1 ? '#00d4ff' : '#8b5cf6'}
            emissive={i === 1 ? '#00d4ff' : '#8b5cf6'}
            emissiveIntensity={0.3}
          />
        </mesh>
      ))}
    </group>
  )
}

// Stadium Stands (Трибуны)
export function StadiumStands() {
  const rows = 8

  return (
    <group>
      {/* Main stands - both long sides */}
      {[-1, 1].map((side, sideIndex) => (
        <group key={sideIndex} position={[0, 0, side * 8]} rotation={[side * 0.3, 0, 0]}>
          {/* Stand structure */}
          {Array.from({ length: rows }).map((_, row) => (
            <group key={row} position={[0, row * 0.5, row * 0.3 * side]}>
              {/* Row platform */}
              <mesh position={[0, 0, 0]}>
                <boxGeometry args={[14, 0.1, 0.8]} />
                <meshStandardMaterial color="#1a1a2e" />
              </mesh>
              {/* Seats (simplified as blocks with glow) */}
              {Array.from({ length: 15 }).map((_, seat) => (
                <mesh key={seat} position={[-7 + seat * 1, 0.15, 0]}>
                  <boxGeometry args={[0.4, 0.3, 0.3]} />
                  <meshStandardMaterial
                    color={Math.random() > 0.7 ? '#00d4ff' : '#2a2a4e'}
                    emissive={Math.random() > 0.8 ? '#00d4ff' : '#000'}
                    emissiveIntensity={0.3}
                  />
                </mesh>
              ))}
            </group>
          ))}
        </group>
      ))}

      {/* End stands (behind goals) */}
      {[-1, 1].map((side, sideIndex) => (
        <group key={`end-${sideIndex}`} position={[side * 9, 0, 0]} rotation={[0, 0, -side * 0.3]}>
          {Array.from({ length: 5 }).map((_, row) => (
            <group key={row} position={[row * 0.4 * side, row * 0.5, 0]}>
              <mesh>
                <boxGeometry args={[0.8, 0.1, 8]} />
                <meshStandardMaterial color="#1a1a2e" />
              </mesh>
            </group>
          ))}
        </group>
      ))}
    </group>
  )
}

// Stadium Lighting
export function StadiumLighting() {
  const spotlightRef = useRef<Group>(null)

  useFrame((state) => {
    if (spotlightRef.current) {
      // Subtle light movement
      spotlightRef.current.rotation.y = Math.sin(state.clock.elapsedTime * 0.1) * 0.05
    }
  })

  return (
    <group ref={spotlightRef}>
      {/* Main arena spotlights */}
      {[
        [-4, 8, -4],
        [4, 8, -4],
        [-4, 8, 4],
        [4, 8, 4],
      ].map((pos, i) => (
        <group key={i} position={pos as [number, number, number]}>
          {/* Light fixture */}
          <mesh>
            <cylinderGeometry args={[0.3, 0.5, 0.3, 8]} />
            <meshStandardMaterial color="#333" metalness={0.8} />
          </mesh>
          {/* Light cone (visual) */}
          <mesh position={[0, -0.3, 0]} rotation={[Math.PI, 0, 0]}>
            <coneGeometry args={[0.8, 2, 8, 1, true]} />
            <meshStandardMaterial
              color="#ffeeaa"
              transparent
              opacity={0.1}
              side={2}
            />
          </mesh>
          <spotLight
            position={[0, 0, 0]}
            angle={0.6}
            penumbra={0.5}
            intensity={2}
            color="#fff8e0"
            castShadow
            target-position={[0, -10, 0]}
          />
        </group>
      ))}

      {/* Neon accent lights around rink */}
      <mesh position={[0, 0.02, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <ringGeometry args={[6.5, 6.6, 64]} />
        <meshStandardMaterial color="#00d4ff" emissive="#00d4ff" emissiveIntensity={1} />
      </mesh>
    </group>
  )
}

// Scoreboard (Табло)
export function Scoreboard({ homeScore = 2, awayScore = 1, time = '12:34', period = 2 }) {
  const boardRef = useRef<Mesh>(null)

  useFrame((state) => {
    if (boardRef.current) {
      // Gentle floating
      boardRef.current.position.y = 6 + Math.sin(state.clock.elapsedTime * 0.5) * 0.1
    }
  })

  return (
    <group position={[0, 6, 0]}>
      {/* Main board */}
      <mesh ref={boardRef}>
        <boxGeometry args={[4, 1.5, 0.3]} />
        <meshStandardMaterial color="#0a0e1a" />
      </mesh>

      {/* Neon border */}
      <mesh position={[0, 0, 0.16]}>
        <planeGeometry args={[4.1, 1.6]} />
        <meshStandardMaterial color="#00d4ff" emissive="#00d4ff" emissiveIntensity={0.5} transparent opacity={0.3} />
      </mesh>

      {/* Frame */}
      {[
        { pos: [0, 0.8, 0.16], size: [4.2, 0.05, 0.1] },
        { pos: [0, -0.8, 0.16], size: [4.2, 0.05, 0.1] },
        { pos: [-2.05, 0, 0.16], size: [0.05, 1.6, 0.1] },
        { pos: [2.05, 0, 0.16], size: [0.05, 1.6, 0.1] },
      ].map((frame, i) => (
        <mesh key={i} position={frame.pos as [number, number, number]}>
          <boxGeometry args={frame.size as [number, number, number]} />
          <meshStandardMaterial color="#00d4ff" emissive="#00d4ff" emissiveIntensity={1} />
        </mesh>
      ))}

      {/* Team names */}
      <Text position={[-1.2, 0.4, 0.2]} fontSize={0.2} color="#ffffff" anchorX="center">
        ДИНАМО
      </Text>
      <Text position={[1.2, 0.4, 0.2]} fontSize={0.2} color="#ffffff" anchorX="center">
        СКА
      </Text>

      {/* Scores */}
      <Text position={[-1.2, -0.1, 0.2]} fontSize={0.5} color="#00d4ff" anchorX="center">
        {homeScore}
      </Text>
      <Text position={[0, -0.1, 0.2]} fontSize={0.3} color="#ff3366" anchorX="center">
        :
      </Text>
      <Text position={[1.2, -0.1, 0.2]} fontSize={0.5} color="#00d4ff" anchorX="center">
        {awayScore}
      </Text>

      {/* Time and period */}
      <Text position={[0, 0.4, 0.2]} fontSize={0.25} color="#ffaa00" anchorX="center">
        {time}
      </Text>
      <Text position={[0, -0.55, 0.2]} fontSize={0.15} color="#888" anchorX="center">
        {`${period} ПЕРИОД`}
      </Text>
    </group>
  )
}

// Complete Stadium
export function Stadium() {
  return (
    <group>
      <StadiumBoards />
      <StadiumStands />
      <StadiumLighting />
      <Scoreboard homeScore={3} awayScore={2} time="15:42" period={2} />
    </group>
  )
}
