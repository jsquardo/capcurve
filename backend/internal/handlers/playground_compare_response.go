package handlers

type playgroundCompareResponse struct {
	Data []playgroundCompareItem `json:"data"`
}

type playgroundCompareItem struct {
	Player  playgroundQueryPlayerItem     `json:"player"`
	Seasons []playgroundCompareSeasonItem `json:"seasons"`
}

type playgroundCompareSeasonItem struct {
	Season   playgroundQuerySeasonItem `json:"season"`
	Hitting  *playerHittingStats       `json:"hitting"`
	Pitching *playerPitchingStats      `json:"pitching"`
}
