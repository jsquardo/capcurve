package handlers

type playgroundQueryResponse struct {
	Data []playgroundQueryItem `json:"data"`
	Meta playerListMeta        `json:"meta"`
}

type playgroundQueryItem struct {
	Player   playgroundQueryPlayerItem `json:"player"`
	Season   playgroundQuerySeasonItem `json:"season"`
	Hitting  *playerHittingStats       `json:"hitting"`
	Pitching *playerPitchingStats      `json:"pitching"`
}

type playgroundQueryPlayerItem struct {
	ID        uint   `json:"id"`
	MLBID     int    `json:"mlb_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Position  string `json:"position"`
	Bats      string `json:"bats"`
	Throws    string `json:"throws"`
	Active    bool   `json:"active"`
	ImageURL  string `json:"image_url"`
}

type playgroundQuerySeasonItem struct {
	Year       int     `json:"year"`
	TeamID     int     `json:"team_id"`
	TeamName   string  `json:"team_name"`
	Age        int     `json:"age"`
	ValueScore float64 `json:"value_score"`
}
