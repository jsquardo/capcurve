package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const defaultSavantBaseURL = "https://baseballsavant.mlb.com"

var savantDataPattern = regexp.MustCompile(`(?s)var (?:leaderboard_data|data) = (\[.*?\]);`)

type SavantClient struct {
	baseURL    string
	httpClient *http.Client
	cache      map[string]map[int]SavantEnrichment
	mu         sync.RWMutex
}

// NewSavantClient builds the client that scrapes embedded leaderboard data from
// Baseball Savant's season-level Statcast pages.
func NewSavantClient(httpClient *http.Client) *SavantClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &SavantClient{
		baseURL:    defaultSavantBaseURL,
		httpClient: httpClient,
		cache:      make(map[string]map[int]SavantEnrichment),
	}
}

// FetchSeasonEnrichment loads one season/type leaderboard once, caches it, and
// returns the row for the requested player if Savant coverage exists.
func (c *SavantClient) FetchSeasonEnrichment(ctx context.Context, season int, playerID int, statType SavantType) (*SavantEnrichment, error) {
	cacheKey := fmt.Sprintf("%s:%d", statType, season)

	c.mu.RLock()
	if entries, ok := c.cache[cacheKey]; ok {
		if enrichment, found := entries[playerID]; found {
			c.mu.RUnlock()
			return &enrichment, nil
		}
		c.mu.RUnlock()
		return nil, nil
	}
	c.mu.RUnlock()

	entries, err := c.fetchLeaderboard(ctx, season, statType)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[cacheKey] = entries
	c.mu.Unlock()

	enrichment, found := entries[playerID]
	if !found {
		return nil, nil
	}

	return &enrichment, nil
}

// fetchLeaderboard downloads the HTML page whose embedded JavaScript contains the
// Statcast leaderboard payload.
func (c *SavantClient) fetchLeaderboard(ctx context.Context, season int, statType SavantType) (map[int]SavantEnrichment, error) {
	requestURL := fmt.Sprintf("%s/leaderboard/statcast?type=%s&year=%d&position=&team=&min=q", c.baseURL, statType, season)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("savant request %s returned %d", requestURL, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return parseSavantLeaderboard(string(body))
}

// parseSavantLeaderboard extracts the embedded leaderboard array and converts it
// into a lookup keyed by MLB player ID.
func parseSavantLeaderboard(html string) (map[int]SavantEnrichment, error) {
	match := savantDataPattern.FindStringSubmatch(html)
	if len(match) != 2 {
		return nil, fmt.Errorf("embedded savant data not found")
	}

	var rows []savantLeaderboardRow
	if err := json.Unmarshal([]byte(match[1]), &rows); err != nil {
		return nil, fmt.Errorf("decode savant leaderboard data: %w", err)
	}

	enrichments := make(map[int]SavantEnrichment, len(rows))
	for _, row := range rows {
		playerID := parseStringInt(row.EntityID)
		if playerID == 0 {
			continue
		}

		enrichments[playerID] = SavantEnrichment{
			ExpectedBattingAvg: parseOptionalStringFloat(row.ExpectedBattingAvg),
			ExpectedSlugging:   parseOptionalStringFloat(row.ExpectedSlugging),
			ExpectedWOBA:       parseOptionalStringFloat(row.ExpectedWOBA),
			ExpectedERA:        parseOptionalStringFloat(row.ExpectedERA),
			BarrelPct:          parseOptionalStringFloat(row.BarrelsPerBIP),
			HardHitPct:         parseOptionalStringFloat(row.HardHitPercent),
			AvgExitVelocity:    parseOptionalStringFloat(row.ExitVelocityAvg),
			AvgLaunchAngle:     parseOptionalStringFloat(row.LaunchAngleAvg),
			SweetSpotPct:       parseOptionalStringFloat(row.SweetSpotPercent),
		}
	}

	return enrichments, nil
}

type savantLeaderboardRow struct {
	EntityID           string `json:"entity_id"`
	ExpectedBattingAvg string `json:"est_ba"`
	ExpectedSlugging   string `json:"est_slg"`
	ExpectedWOBA       string `json:"est_woba"`
	ExpectedERA        string `json:"xera"`
	BarrelsPerBIP      string `json:"barrels_per_bip"`
	HardHitPercent     string `json:"hard_hit_percent"`
	ExitVelocityAvg    string `json:"exit_velocity_avg"`
	LaunchAngleAvg     string `json:"launch_angle_avg"`
	SweetSpotPercent   string `json:"sweet_spot_percent"`
}

func (r *savantLeaderboardRow) UnmarshalJSON(data []byte) error {
	aux := map[string]any{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.EntityID = valueToString(aux["entity_id"])
	r.ExpectedBattingAvg = valueToString(aux["est_ba"])
	r.ExpectedSlugging = valueToString(aux["est_slg"])
	r.ExpectedWOBA = valueToString(aux["est_woba"])
	r.ExpectedERA = valueToString(aux["xera"])
	r.BarrelsPerBIP = valueToString(aux["barrels_per_bip"])
	r.HardHitPercent = valueToString(aux["hard_hit_percent"])
	r.ExitVelocityAvg = valueToString(aux["exit_velocity_avg"])
	r.LaunchAngleAvg = valueToString(aux["launch_angle_avg"])
	r.SweetSpotPercent = valueToString(aux["sweet_spot_percent"])
	return nil
}

// valueToString normalizes Savant's mixed string/number/null field shapes before parsing.
func valueToString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", v), "0"), ".")
	default:
		return ""
	}
}
