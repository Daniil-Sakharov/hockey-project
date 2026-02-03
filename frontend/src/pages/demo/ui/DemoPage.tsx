import { HeroCard } from '@/widgets/hero-card'
import { LandingHeader } from '@/widgets/landing-header'
import { LandingFeatures } from '@/widgets/landing-features'
import { LandingPricing } from '@/widgets/landing-pricing'
import { LandingAbout } from '@/widgets/landing-about'
import { LandingFooter } from '@/widgets/landing-footer'

export function DemoPage() {
  return (
    <div className="min-h-screen bg-[#050810]">
      <LandingHeader />
      <HeroCard />
      <LandingFeatures />
      <LandingPricing />
      <LandingAbout />
      <LandingFooter />
    </div>
  )
}
