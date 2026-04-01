import HeroSection from '../components/home/HeroSection'
import TrendingSection from '../components/home/TrendingSection'
import StatLeadersSection from '../components/home/StatLeadersSection'
import FeedSection from '../components/home/FeedSection'

export default function HomePage() {
  return (
    <>
      <HeroSection />
      <TrendingSection />
      <section className="border-b border-border px-4 py-12 sm:px-6 lg:px-10">
        <div className="grid grid-cols-1 gap-12 lg:grid-cols-[1.1fr_1fr]">
          <StatLeadersSection />
          <FeedSection />
        </div>
      </section>
    </>
  )
}
