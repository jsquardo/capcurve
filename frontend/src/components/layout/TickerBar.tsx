import { motion } from 'framer-motion'

// Mock data — replace with API-driven arc delta feed when backend endpoint is ready
// TODO: wire to /api/v1/players/trending arc-delta endpoint (Phase 3 backend feature)
const tickerItems = [
  { name: 'Elly De La Cruz', val: 'Arc +9.1', dir: '▲' },
  { name: 'Paul Skenes', val: 'Arc +7.4', dir: '▲' },
  { name: 'Aaron Judge', val: '58 HR', dir: '' },
  { name: 'Bobby Witt Jr.', val: 'Arc +5.9', dir: '▲' },
  { name: 'Shohei Ohtani', val: 'OPS 1.038', dir: '' },
  { name: 'Jackson Chourio', val: 'Arc +6.8', dir: '▲' },
  { name: 'Spencer Strider', val: 'ERA 2.11', dir: '' },
  { name: 'Gunnar Henderson', val: 'Arc +4.2', dir: '▲' },
]

// The list is duplicated so the animation loops seamlessly:
// animate to -50% (exactly one copy width), then repeat from 0.
export default function TickerBar() {
  return (
    <div className="overflow-hidden bg-accent py-2">
      <motion.div
        initial={{ x: '0%' }}
        animate={{ x: '-50%' }}
        transition={{ duration: 32, repeat: Infinity, ease: 'linear' }}
        className="flex w-max whitespace-nowrap"
      >
        {[...tickerItems, ...tickerItems].map((item, index) => (
          <div
            key={index}
            className="inline-flex items-center gap-[10px] border-r border-black/15 px-8 text-[12px] font-semibold tracking-[0.3px] text-[#0a0d12]"
          >
            <span>{item.name}</span>
            <span className="font-mono text-[12px]">{item.val}</span>
            {item.dir ? <span className="text-[10px] font-bold opacity-70">{item.dir}</span> : null}
          </div>
        ))}
      </motion.div>
    </div>
  )
}
