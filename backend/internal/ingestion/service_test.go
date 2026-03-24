package ingestion

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
)

func TestSeasonStatUpsertClauseTargetsActivePlayerYearRows(t *testing.T) {
	t.Parallel()

	conflict := seasonStatUpsertClause()

	require.Equal(t, []clause.Column{{Name: "player_id"}, {Name: "year"}}, conflict.Columns)
	require.Len(t, conflict.TargetWhere.Exprs, 1)

	targetDeletedAt, ok := conflict.TargetWhere.Exprs[0].(clause.Eq)
	require.True(t, ok)
	require.Equal(t, "deleted_at", targetDeletedAt.Column)
	require.Nil(t, targetDeletedAt.Value)

	updatedColumns := make([]string, 0, len(conflict.DoUpdates))
	for _, assignment := range conflict.DoUpdates {
		updatedColumns = append(updatedColumns, assignment.Column.Name)
	}

	require.Len(t, updatedColumns, len(seasonStatUpsertColumns))
	require.ElementsMatch(t, seasonStatUpsertColumns, updatedColumns)
}

func TestOrderedSeasonSplitsKeepsRealTeamsAheadOfAggregateRows(t *testing.T) {
	t.Parallel()

	ordered := orderedSeasonSplits([]MLBSeasonSplit{
		{Season: "2023", Team: mlbStatTeam{ID: 135, Name: "San Diego Padres"}},
		{Season: "2022", Team: mlbStatTeam{ID: 0, Name: "TOT"}},
		{Season: "2022", Team: mlbStatTeam{ID: 120, Name: "Washington Nationals"}},
		{Season: "2022", Team: mlbStatTeam{ID: 135, Name: "San Diego Padres"}},
	})

	require.Len(t, ordered, 4)
	require.Equal(t, 2022, parseStringInt(ordered[0].split.Season))
	require.Equal(t, 120, ordered[0].split.Team.ID)
	require.Equal(t, 2, ordered[0].sourceOrder)
	require.Equal(t, 135, ordered[1].split.Team.ID)
	require.Equal(t, 3, ordered[1].sourceOrder)
	require.Equal(t, 0, ordered[2].split.Team.ID)
	require.Equal(t, 2023, parseStringInt(ordered[3].split.Season))
}

func TestFilterSeasonSplitsKeepsOnlyRequestedYear(t *testing.T) {
	t.Parallel()

	filtered := filterSeasonSplits([]MLBSeasonSplit{
		{Season: "2024", Team: mlbStatTeam{ID: 147, Name: "New York Yankees"}},
		{Season: "2025", Team: mlbStatTeam{ID: 147, Name: "New York Yankees"}},
		{Season: "2025", Team: mlbStatTeam{ID: 0, Name: "TOT"}},
		{Season: "2026", Team: mlbStatTeam{ID: 111, Name: "Boston Red Sox"}},
	}, 2025)

	require.Len(t, filtered, 2)
	require.Equal(t, "2025", filtered[0].Season)
	require.Equal(t, 147, filtered[0].Team.ID)
	require.Equal(t, "2025", filtered[1].Season)
	require.Equal(t, 0, filtered[1].Team.ID)
}
