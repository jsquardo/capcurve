import axios from "axios";
import type {
  Player,
  PlayerListItem,
  PlayerListResponse,
  CareerArcResponse,
  Contract,
  AdminDashboard,
  LeaderboardCategory,
  LeaderboardsResponse,
} from "@/types";

const ensureVersionedApiBaseURL = (baseURL: string): string => {
  return baseURL.replace(/\/api\/?$/, "/api/v1");
};

const apiBaseURL = ensureVersionedApiBaseURL(
  import.meta.env.VITE_API_URL ?? "/api/v1",
);
const adminSecret = import.meta.env.VITE_ADMIN_SECRET;

const api = axios.create({
  baseURL: apiBaseURL,
  headers: { "Content-Type": "application/json" },
});

export const searchPlayers = async (
  query: string,
): Promise<PlayerListItem[]> => {
  const { data } = await api.get<PlayerListResponse>("/players", {
    params: { q: query },
  });
  return data.data;
};

export const getPlayer = async (id: number): Promise<Player> => {
  const { data } = await api.get(`/players/${id}`);
  return data;
};

export const listPlayers = async (params?: {
  active?: boolean;
  position?: string;
}): Promise<PlayerListItem[]> => {
  const { data } = await api.get<PlayerListResponse>("/players", { params });
  return data.data;
};

export const getCareerArc = async (
  playerId: number,
): Promise<CareerArcResponse> => {
  const { data } = await api.get(`/players/${playerId}/career-arc`);
  return data;
};

export const getPlayerContracts = async (
  playerId: number,
): Promise<Contract[]> => {
  const { data } = await api.get(`/players/${playerId}/contracts`);
  return data;
};

export const getMostOverpaid = async (): Promise<Contract[]> => {
  const { data } = await api.get("/leaderboards/most-overpaid");
  return data;
};

export const getBestValue = async (): Promise<Contract[]> => {
  const { data } = await api.get("/leaderboards/best-value");
  return data;
};

export const getLeaderboards = async (params: {
  category: LeaderboardCategory;
  season?: number;
  page?: number;
  page_size?: number;
}): Promise<LeaderboardsResponse> => {
  const { data } = await api.get<LeaderboardsResponse>("/leaderboards", {
    params,
  });
  return data;
};

export const getAdminDashboard = async (): Promise<AdminDashboard> => {
  const { data } = await api.get("/admin/dashboard", {
    headers: adminSecret
      ? {
          Authorization: `Bearer ${adminSecret}`,
        }
      : undefined,
  });
  return data;
};

export default api;
