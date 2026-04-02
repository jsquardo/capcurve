import { Routes, Route } from 'react-router-dom'
import Layout from '@/components/Layout'
import HomePage from '@/pages/HomePage'
import PlayersPage from '@/pages/PlayersPage'
import PlayerPage from '@/pages/PlayerPage'
import LeaderboardsPage from '@/pages/LeaderboardsPage'
import AdminDashboardPage from '@/pages/AdminDashboardPage'

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/players" element={<PlayersPage />} />
        <Route path="/players/:id" element={<PlayerPage />} />
        <Route path="/leaderboards" element={<LeaderboardsPage />} />
        <Route path="/admin" element={<AdminDashboardPage />} />
      </Route>
    </Routes>
  )
}
