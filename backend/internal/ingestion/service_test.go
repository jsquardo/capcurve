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

	require.Contains(t, updatedColumns, "team_id")
	require.Contains(t, updatedColumns, "team_name")
	require.Contains(t, updatedColumns, "sweet_spot_pct")
}
