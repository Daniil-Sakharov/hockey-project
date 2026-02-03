import { Suspense } from 'react'
import { Canvas } from '@react-three/fiber'
import { OrbitControls, PerspectiveCamera, Environment, Float } from '@react-three/drei'
import { EffectComposer, Bloom } from '@react-three/postprocessing'
import { motion } from 'framer-motion'
import { HoloBarChart3D, HoloRingChart3D } from '@/shared/ui/3d/HoloChart3D'

const monthlyData = [
  { value: 42, label: 'Янв' },
  { value: 58, label: 'Фев' },
  { value: 35, label: 'Мар' },
  { value: 67, label: 'Апр' },
  { value: 89, label: 'Май' },
  { value: 54, label: 'Июн' },
]

const teamData = [
  { value: 78, label: 'СКА' },
  { value: 65, label: 'ЦСКА' },
  { value: 82, label: 'Динамо' },
  { value: 45, label: 'Спартак' },
]

function Scene() {
  return (
    <>
      <ambientLight intensity={0.2} />
      <spotLight position={[10, 10, 5]} intensity={1} color="#00d4ff" />
      <spotLight position={[-10, 10, -5]} intensity={0.8} color="#8b5cf6" />

      {/* Main bar chart - floating */}
      <Float speed={1} rotationIntensity={0.1} floatIntensity={0.3}>
        <HoloBarChart3D
          data={monthlyData}
          position={[-2.5, 0, 0]}
          title="Голы по месяцам"
          color="#00d4ff"
        />
      </Float>

      {/* Secondary bar chart */}
      <Float speed={1.2} rotationIntensity={0.1} floatIntensity={0.2}>
        <HoloBarChart3D
          data={teamData}
          position={[2.5, 0, 0]}
          rotation={[0, -0.3, 0]}
          title="Рейтинг команд"
          color="#8b5cf6"
        />
      </Float>

      {/* Ring charts */}
      <Float speed={0.8} floatIntensity={0.4}>
        <HoloRingChart3D
          position={[-4, 2.5, -1]}
          value={85}
          label="Победы"
          color="#00d4ff"
        />
      </Float>

      <Float speed={1.1} floatIntensity={0.3}>
        <HoloRingChart3D
          position={[4, 2.5, -1]}
          value={72}
          label="Реализация"
          color="#ec4899"
        />
      </Float>

      {/* Holographic table/platform */}
      <mesh position={[0, -1.5, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[12, 8]} />
        <meshStandardMaterial
          color="#00d4ff"
          transparent
          opacity={0.05}
          emissive="#00d4ff"
          emissiveIntensity={0.3}
        />
      </mesh>

      {/* Grid on platform */}
      <gridHelper
        args={[12, 24, '#00d4ff', '#0a2040']}
        position={[0, -1.49, 0]}
      />

      <Environment preset="night" />
    </>
  )
}

export function Holographic3DSection() {
  return (
    <section className="relative py-24">
      {/* Background gradient */}
      <div className="gradient-mesh absolute inset-0" />

      <div className="relative mx-auto max-w-7xl px-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-8 text-center"
        >
          <h2 className="mb-4 text-4xl font-bold">
            <span className="text-gradient">3D Аналитика</span>
            <span className="text-white"> в реальном времени</span>
          </h2>
          <p className="text-gray-400">
            Интерактивная визуализация данных
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          className="relative overflow-hidden rounded-2xl border border-[#00d4ff]/20"
          style={{ height: 500 }}
        >
          {/* Holographic frame corners */}
          <div className="absolute left-2 top-2 z-10 h-8 w-8 border-l-2 border-t-2 border-[#00d4ff]" />
          <div className="absolute right-2 top-2 z-10 h-8 w-8 border-r-2 border-t-2 border-[#00d4ff]" />
          <div className="absolute bottom-2 left-2 z-10 h-8 w-8 border-b-2 border-l-2 border-[#00d4ff]" />
          <div className="absolute bottom-2 right-2 z-10 h-8 w-8 border-b-2 border-r-2 border-[#00d4ff]" />

          {/* HUD elements */}
          <div className="absolute left-4 top-4 z-10 text-xs font-mono text-[#00d4ff]">
            <div>SYS: ONLINE</div>
            <div>DATA: LIVE</div>
          </div>
          <div className="absolute right-4 top-4 z-10 text-right text-xs font-mono text-[#00d4ff]">
            <div>STATS: 2024</div>
            <div>MODE: 3D</div>
          </div>

          <Canvas shadows dpr={[1, 2]}>
            <Suspense fallback={null}>
              <PerspectiveCamera makeDefault position={[0, 3, 8]} fov={45} />
              <Scene />
              <EffectComposer>
                <Bloom intensity={1} luminanceThreshold={0.3} mipmapBlur />
              </EffectComposer>
              <OrbitControls
                enableZoom={false}
                enablePan={false}
                maxPolarAngle={Math.PI / 2.2}
                minPolarAngle={Math.PI / 4}
                autoRotate
                autoRotateSpeed={0.3}
              />
            </Suspense>
          </Canvas>

          {/* Scan line effect */}
          <div
            className="pointer-events-none absolute inset-0 z-10"
            style={{
              background:
                'repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0, 212, 255, 0.03) 2px, rgba(0, 212, 255, 0.03) 4px)',
            }}
          />
        </motion.div>

        {/* Mini stats below */}
        <div className="mt-6 grid grid-cols-4 gap-4">
          {[
            { label: 'Всего игроков', value: '2,847' },
            { label: 'Команд', value: '156' },
            { label: 'Турниров', value: '42' },
            { label: 'Матчей', value: '12,459' },
          ].map((stat) => (
            <motion.div
              key={stat.label}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="glass-card rounded-xl p-4 text-center"
            >
              <div className="text-2xl font-bold text-[#00d4ff]">{stat.value}</div>
              <div className="text-xs text-gray-500">{stat.label}</div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}
