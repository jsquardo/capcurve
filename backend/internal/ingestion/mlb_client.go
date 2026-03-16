package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultMLBBaseURL = "https://statsapi.mlb.com/api/v1"

type MLBClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewMLBClient builds the thin HTTP client used for official MLB Stats API fetches.
func NewMLBClient(httpClient *http.Client) *MLBClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &MLBClient{
		baseURL:    defaultMLBBaseURL,
		httpClient: httpClient,
	}
}

// FetchPlayer loads the canonical player bio used to seed the players table.
func (c *MLBClient) FetchPlayer(ctx context.Context, playerID int) (*MLBPlayer, error) {
	var response mlbPeopleResponse
	if err := c.get(ctx, fmt.Sprintf("/people/%d", playerID), nil, &response); err != nil {
		return nil, err
	}
	if len(response.People) == 0 {
		return nil, fmt.Errorf("mlb player %d not found", playerID)
	}

	return &response.People[0], nil
}

// FetchYearByYearStats loads one stat group at a time so hitting and pitching can
// be merged into a single season record for two-way players.
func (c *MLBClient) FetchYearByYearStats(ctx context.Context, playerID int, group string) ([]MLBSeasonSplit, error) {
	params := url.Values{}
	params.Set("stats", "yearByYear")
	params.Set("group", group)

	var response mlbStatsResponse
	if err := c.get(ctx, fmt.Sprintf("/people/%d/stats", playerID), params, &response); err != nil {
		return nil, err
	}

	if len(response.Stats) == 0 {
		return nil, nil
	}

	return response.Stats[0].Splits, nil
}

// FetchTeams exposes the team reference feed for later validation or seed work.
func (c *MLBClient) FetchTeams(ctx context.Context) ([]MLBTeam, error) {
	params := url.Values{}
	params.Set("sportId", "1")

	var response mlbTeamsResponse
	if err := c.get(ctx, "/teams", params, &response); err != nil {
		return nil, err
	}

	return response.Teams, nil
}

// get centralizes request construction and decoding for the MLB API endpoints.
func (c *MLBClient) get(ctx context.Context, path string, params url.Values, target any) error {
	requestURL := c.baseURL + path
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("mlb request %s returned %d", requestURL, res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(target)
}

type mlbPeopleResponse struct {
	People []MLBPlayer `json:"people"`
}

type MLBPlayer struct {
	ID              int               `json:"id"`
	FullName        string            `json:"fullName"`
	FirstName       string            `json:"firstName"`
	LastName        string            `json:"lastName"`
	UseName         string            `json:"useName"`
	UseLastName     string            `json:"useLastName"`
	BirthDate       string            `json:"birthDate"`
	Active          bool              `json:"active"`
	PrimaryPosition mlbPosition       `json:"primaryPosition"`
	BatSide         mlbHandDescriptor `json:"batSide"`
	PitchHand       mlbHandDescriptor `json:"pitchHand"`
}

type mlbPosition struct {
	Name string `json:"name"`
}

type mlbHandDescriptor struct {
	Code string `json:"code"`
}

type mlbStatsResponse struct {
	Stats []struct {
		Splits []MLBSeasonSplit `json:"splits"`
	} `json:"stats"`
}

type MLBSeasonSplit struct {
	Season string         `json:"season"`
	Stat   map[string]any `json:"stat"`
	Team   mlbStatTeam    `json:"team"`
	Player map[string]any `json:"player"`
	League map[string]any `json:"league"`
	Sport  map[string]any `json:"sport"`
}

type mlbStatTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type mlbTeamsResponse struct {
	Teams []MLBTeam `json:"teams"`
}

type MLBTeam struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TeamName     string `json:"teamName"`
	LocationName string `json:"locationName"`
	Abbreviation string `json:"abbreviation"`
	Active       bool   `json:"active"`
}
