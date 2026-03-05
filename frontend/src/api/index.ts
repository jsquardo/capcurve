import axios from 'axios'
import type { Player, CareerArcResponse, Contract } from '@/types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

export const searchPlayers = async (query: string): Promise<Player[]> => {
  const { data } = await api.get('/players/search', { params: { q: query } })
  return data
}

export const getPlayer = async (id: number): Promise<Player> => {
  const { data } = await api.get(`/players/${id}`)
  return data
}

export const listPlayers = async (params?: { active?: boolean; position?: string }): Promise<Player[]> => {
  const { data } = await api.get('/players', { params })
  return data
}

export const getCareerArc = async (playerId: number): Promise<CareerArcResponse> => {
  const { data } = await api.get(`/players/${playerId}/career-arc`)
  return data
}

export const getPlayerContracts = async (playerId: number): Promise<Contract[]> => {
  const { data } = await api.get(`/players/${playerId}/contracts`)
  return data
}

export const getMostOverpaid = async (): Promise<Contract[]> => {
  const { data } = await api.get('/leaderboards/most-overpaid')
  return data
}

export const getBestValue = async (): Promise<Contract[]> => {
  const { data } = await api.get('/leaderboards/best-value')
  return data
}

export default api
