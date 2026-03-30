package handlers

type leaderboardsResponse struct {
	Data leaderboardData `json:"data"`
}

type leaderboardData struct {
	Category string            `json:"category"`
	Leaders  []leaderboardItem `json:"leaders"`
	Meta     leaderboardMeta   `json:"meta"`
}

type leaderboardMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type leaderboardItem struct {
	Rank       int     `json:"rank"`
	PlayerID   uint    `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Position   string  `json:"position"`
	Team       string  `json:"team"`
	Value      float64 `json:"value"`
	Season     *int    `json:"season,omitempty"`
}

func (r leaderboardRow) toResponse(rank int) leaderboardItem {
	return leaderboardItem{
		Rank:       rank,
		PlayerID:   r.PlayerID,
		PlayerName: r.PlayerName,
		Position:   r.Position,
		Team:       r.Team,
		Value:      r.Value,
		Season:     r.Season,
	}
}
