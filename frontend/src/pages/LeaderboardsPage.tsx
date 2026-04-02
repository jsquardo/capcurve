import { useState } from 'react'
import type { LeaderboardCategory, LeaderboardEntry } from '@/types'
import LeaderboardHero from '@/components/leaderboards/LeaderboardHero'
import LeaderboardCategoryTabs from '@/components/leaderboards/LeaderboardCategoryTabs'
import LeaderboardTable from '@/components/leaderboards/LeaderboardTable'
import LeaderboardPagination from '@/components/leaderboards/LeaderboardPagination'

// ── Mock data ─────────────────────────────────────────────────────────────────
// Keyed by category. totalPages: 3 is intentionally larger than these 10 rows
// so pagination controls render and can be verified. Replace this block with a
// TanStack Query call to getLeaderboards() when wiring to the live API.

const MOCK_SEASON = 2025
const MOCK_TOTAL_PAGES = 3

const MOCK_LEADERS: Record<LeaderboardCategory, LeaderboardEntry[]> = {
  peak_arc: [
    { rank: 1,  player_id: 1,  player_name: 'Mike Trout',        position: 'CF', team: 'Los Angeles Angels',    value: 94.2, season: null },
    { rank: 2,  player_id: 2,  player_name: 'Barry Bonds',        position: 'LF', team: 'San Francisco Giants',  value: 93.7, season: null },
    { rank: 3,  player_id: 3,  player_name: 'Willie Mays',        position: 'CF', team: 'San Francisco Giants',  value: 92.1, season: null },
    { rank: 4,  player_id: 4,  player_name: 'Ken Griffey Jr.',    position: 'CF', team: 'Seattle Mariners',      value: 91.5, season: null },
    { rank: 5,  player_id: 5,  player_name: 'Hank Aaron',         position: 'RF', team: 'Milwaukee Braves',      value: 90.8, season: null },
    { rank: 6,  player_id: 6,  player_name: 'Mookie Betts',       position: 'RF', team: 'Los Angeles Dodgers',   value: 89.4, season: null },
    { rank: 7,  player_id: 7,  player_name: 'Juan Soto',          position: 'LF', team: 'New York Yankees',      value: 88.9, season: null },
    { rank: 8,  player_id: 8,  player_name: 'Frank Robinson',     position: 'RF', team: 'Cincinnati Reds',       value: 88.3, season: null },
    { rank: 9,  player_id: 9,  player_name: 'Roberto Clemente',   position: 'RF', team: 'Pittsburgh Pirates',    value: 87.7, season: null },
    { rank: 10, player_id: 10, player_name: 'Ted Williams',       position: 'LF', team: 'Boston Red Sox',        value: 87.1, season: null },
  ],
  hr: [
    { rank: 1,  player_id: 11, player_name: 'Aaron Judge',        position: 'RF', team: 'New York Yankees',      value: 58,   season: MOCK_SEASON },
    { rank: 2,  player_id: 12, player_name: 'Kyle Schwarber',     position: 'LF', team: 'Philadelphia Phillies', value: 47,   season: MOCK_SEASON },
    { rank: 3,  player_id: 13, player_name: 'Pete Alonso',        position: '1B', team: 'New York Mets',         value: 46,   season: MOCK_SEASON },
    { rank: 4,  player_id: 14, player_name: 'Yordan Alvarez',     position: 'DH', team: 'Houston Astros',        value: 44,   season: MOCK_SEASON },
    { rank: 5,  player_id: 15, player_name: 'Shohei Ohtani',      position: 'DH', team: 'Los Angeles Dodgers',   value: 44,   season: MOCK_SEASON },
    { rank: 6,  player_id: 16, player_name: 'Ronald Acuña Jr.',   position: 'RF', team: 'Atlanta Braves',        value: 41,   season: MOCK_SEASON },
    { rank: 7,  player_id: 17, player_name: 'Francisco Lindor',   position: 'SS', team: 'New York Mets',         value: 39,   season: MOCK_SEASON },
    { rank: 8,  player_id: 18, player_name: 'Bryce Harper',       position: '1B', team: 'Philadelphia Phillies', value: 38,   season: MOCK_SEASON },
    { rank: 9,  player_id: 19, player_name: 'Freddie Freeman',    position: '1B', team: 'Los Angeles Dodgers',   value: 36,   season: MOCK_SEASON },
    { rank: 10, player_id: 20, player_name: 'Matt Olson',         position: '1B', team: 'Atlanta Braves',        value: 35,   season: MOCK_SEASON },
  ],
  avg: [
    { rank: 1,  player_id: 21, player_name: 'Luis Arraez',        position: '2B', team: 'San Diego Padres',      value: 0.354, season: MOCK_SEASON },
    { rank: 2,  player_id: 22, player_name: 'Freddie Freeman',    position: '1B', team: 'Los Angeles Dodgers',   value: 0.341, season: MOCK_SEASON },
    { rank: 3,  player_id: 23, player_name: 'Corey Seager',       position: 'SS', team: 'Texas Rangers',         value: 0.327, season: MOCK_SEASON },
    { rank: 4,  player_id: 24, player_name: 'Steven Kwan',        position: 'LF', team: 'Cleveland Guardians',   value: 0.319, season: MOCK_SEASON },
    { rank: 5,  player_id: 25, player_name: 'Paul Goldschmidt',   position: '1B', team: 'St. Louis Cardinals',   value: 0.317, season: MOCK_SEASON },
    { rank: 6,  player_id: 26, player_name: 'Yordan Alvarez',     position: 'DH', team: 'Houston Astros',        value: 0.314, season: MOCK_SEASON },
    { rank: 7,  player_id: 27, player_name: 'Juan Soto',          position: 'LF', team: 'New York Yankees',      value: 0.312, season: MOCK_SEASON },
    { rank: 8,  player_id: 28, player_name: 'Trea Turner',        position: 'SS', team: 'Philadelphia Phillies', value: 0.308, season: MOCK_SEASON },
    { rank: 9,  player_id: 29, player_name: 'Shohei Ohtani',      position: 'DH', team: 'Los Angeles Dodgers',   value: 0.304, season: MOCK_SEASON },
    { rank: 10, player_id: 30, player_name: 'Nolan Arenado',      position: '3B', team: 'St. Louis Cardinals',   value: 0.299, season: MOCK_SEASON },
  ],
  era: [
    { rank: 1,  player_id: 31, player_name: 'Zack Wheeler',       position: 'SP', team: 'Philadelphia Phillies', value: 2.34, season: MOCK_SEASON },
    { rank: 2,  player_id: 32, player_name: 'Gerrit Cole',        position: 'SP', team: 'New York Yankees',      value: 2.51, season: MOCK_SEASON },
    { rank: 3,  player_id: 33, player_name: 'Spencer Strider',    position: 'SP', team: 'Atlanta Braves',        value: 2.67, season: MOCK_SEASON },
    { rank: 4,  player_id: 34, player_name: 'Blake Snell',        position: 'SP', team: 'Los Angeles Dodgers',   value: 2.79, season: MOCK_SEASON },
    { rank: 5,  player_id: 35, player_name: 'Kevin Gausman',      position: 'SP', team: 'Toronto Blue Jays',     value: 2.88, season: MOCK_SEASON },
    { rank: 6,  player_id: 36, player_name: 'Sandy Alcantara',    position: 'SP', team: 'Miami Marlins',         value: 2.94, season: MOCK_SEASON },
    { rank: 7,  player_id: 37, player_name: 'Logan Webb',         position: 'SP', team: 'San Francisco Giants',  value: 3.03, season: MOCK_SEASON },
    { rank: 8,  player_id: 38, player_name: 'Corbin Burnes',      position: 'SP', team: 'Baltimore Orioles',     value: 3.11, season: MOCK_SEASON },
    { rank: 9,  player_id: 39, player_name: 'Framber Valdez',     position: 'SP', team: 'Houston Astros',        value: 3.18, season: MOCK_SEASON },
    { rank: 10, player_id: 40, player_name: 'Pablo López',        position: 'SP', team: 'Minnesota Twins',       value: 3.24, season: MOCK_SEASON },
  ],
  k9: [
    { rank: 1,  player_id: 41, player_name: 'Spencer Strider',    position: 'SP', team: 'Atlanta Braves',        value: 13.8, season: MOCK_SEASON },
    { rank: 2,  player_id: 42, player_name: 'Gerrit Cole',        position: 'SP', team: 'New York Yankees',      value: 13.1, season: MOCK_SEASON },
    { rank: 3,  player_id: 43, player_name: 'Kevin Gausman',      position: 'SP', team: 'Toronto Blue Jays',     value: 12.6, season: MOCK_SEASON },
    { rank: 4,  player_id: 44, player_name: 'Dylan Cease',        position: 'SP', team: 'San Diego Padres',      value: 12.3, season: MOCK_SEASON },
    { rank: 5,  player_id: 45, player_name: 'Zack Wheeler',       position: 'SP', team: 'Philadelphia Phillies', value: 11.9, season: MOCK_SEASON },
    { rank: 6,  player_id: 46, player_name: 'Shane Bieber',       position: 'SP', team: 'Cleveland Guardians',   value: 11.6, season: MOCK_SEASON },
    { rank: 7,  player_id: 47, player_name: 'Chris Sale',         position: 'SP', team: 'Atlanta Braves',        value: 11.4, season: MOCK_SEASON },
    { rank: 8,  player_id: 48, player_name: 'Corbin Burnes',      position: 'SP', team: 'Baltimore Orioles',     value: 11.2, season: MOCK_SEASON },
    { rank: 9,  player_id: 49, player_name: 'Max Fried',          position: 'SP', team: 'New York Yankees',      value: 10.9, season: MOCK_SEASON },
    { rank: 10, player_id: 50, player_name: 'Logan Webb',         position: 'SP', team: 'San Francisco Giants',  value: 10.7, season: MOCK_SEASON },
  ],
}

// ── Page ──────────────────────────────────────────────────────────────────────

const MOCK_PAGE_SIZE = 4

export default function LeaderboardsPage() {
  const [activeCategory, setActiveCategory] = useState<LeaderboardCategory>('peak_arc')
  const [page, setPage] = useState(1)

  // Reset to page 1 whenever the category changes.
  function handleSelectCategory(category: LeaderboardCategory) {
    setActiveCategory(category)
    setPage(1)
  }

  const allLeaders = MOCK_LEADERS[activeCategory]
  // Slice to the current page so the pagination controls actually work against mock data.
  // Replace allLeaders + this slice with a TanStack Query call when wiring to the live API.
  const leaders = allLeaders.slice((page - 1) * MOCK_PAGE_SIZE, page * MOCK_PAGE_SIZE)

  return (
    <div className="shell-container space-y-6 py-8">
      <LeaderboardHero season={MOCK_SEASON} />
      <LeaderboardCategoryTabs
        activeCategory={activeCategory}
        onSelect={handleSelectCategory}
      />
      <LeaderboardTable leaders={leaders} category={activeCategory} />
      <LeaderboardPagination
        page={page}
        totalPages={MOCK_TOTAL_PAGES}
        onPrev={() => setPage(p => p - 1)}
        onNext={() => setPage(p => p + 1)}
      />
    </div>
  )
}
