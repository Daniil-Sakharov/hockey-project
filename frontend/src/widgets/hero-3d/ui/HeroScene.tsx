import { Suspense } from 'react'
import { Canvas } from '@react-three/fiber'
import { PerspectiveCamera, Environment, OrbitControls } from '@react-three/drei'
import { EffectComposer, Bloom, ChromaticAberration } from '@react-three/postprocessing'
import { BlendFunction } from 'postprocessing'
import { IceRink, IceParticles, Goal, Stadium } from '@/shared/ui/3d'

function Lights() {
  return (
    <>
      <ambientLight intensity={0.4} />
      {/* Main arena lights */}
      <spotLight
        position={[0, 12, 0]}
        angle={0.8}
        penumbra={0.5}
        intensity={2}
        color="#ffffff"
        castShadow
      />
      {/* Colored accent lights */}
      <spotLight
        position={[8, 10, 8]}
        angle={0.4}
        penumbra={1}
        intensity={1}
        color="#00d4ff"
      />
      <spotLight
        position={[-8, 10, -8]}
        angle={0.4}
        penumbra={1}
        intensity={0.8}
        color="#8b5cf6"
      />
      <pointLight position={[0, 3, 0]} intensity={0.3} color="#00ffff" />
    </>
  )
}

function Scene() {
  return (
    <>
      <Lights />

      {/* Stadium elements */}
      <Stadium />

      {/* Goals */}
      <Goal position={[-5.5, 0, 0]} rotation={[0, Math.PI / 2, 0]} scale={1} />
      <Goal position={[5.5, 0, 0]} rotation={[0, -Math.PI / 2, 0]} scale={1} />

      {/* Ice surface */}
      <IceRink position={[0, -0.01, 0]} />

      {/* Atmospheric particles */}
      <IceParticles count={150} spread={8} size={0.02} />

      <Environment preset="night" />
    </>
  )
}

function PostProcessing() {
  return (
    <EffectComposer>
      <Bloom
        intensity={1}
        luminanceThreshold={0.3}
        luminanceSmoothing={0.9}
        mipmapBlur
      />
      <ChromaticAberration
        blendFunction={BlendFunction.NORMAL}
        offset={[0.0002, 0.0002]}
      />
    </EffectComposer>
  )
}

export function HeroScene() {
  return (
    <div className="absolute inset-0 z-0">
      <Canvas shadows dpr={[1, 2]} gl={{ antialias: true, alpha: true }}>
        <Suspense fallback={null}>
          <PerspectiveCamera makeDefault position={[0, 5, 12]} fov={50} />
          <Scene />
          <PostProcessing />
          <OrbitControls
            enableZoom={false}
            enablePan={false}
            maxPolarAngle={Math.PI / 2.2}
            minPolarAngle={Math.PI / 6}
            autoRotate
            autoRotateSpeed={0.3}
          />
        </Suspense>
      </Canvas>
    </div>
  )
}
