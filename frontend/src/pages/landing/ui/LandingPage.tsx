import { VideoHero } from '@/widgets/video-hero'
import { Holographic3DSection } from '@/widgets/holographic-charts'
import { FeaturesSection } from './FeaturesSection'
import { CTASection } from './CTASection'

export function LandingPage() {
  return (
    <div className="min-h-screen overflow-x-hidden">
      {/* Hero Section with Video Background */}
      <VideoHero />

      {/* 3D Holographic Charts */}
      <Holographic3DSection />

      {/* Features Section */}
      <FeaturesSection />

      {/* Call to Action */}
      <CTASection />
    </div>
  )
}
