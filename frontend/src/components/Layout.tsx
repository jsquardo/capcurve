import { Outlet } from 'react-router-dom'
import Navbar from './layout/Navbar'
import TickerBar from './layout/TickerBar'
import Footer from './layout/Footer'

export default function Layout() {
  return (
    <div className="min-h-screen bg-transparent text-text">
      <Navbar />
      <TickerBar />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  )
}
