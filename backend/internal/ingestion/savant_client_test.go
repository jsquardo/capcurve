package ingestion

import "testing"

func TestParseSavantLeaderboard(t *testing.T) {
	t.Parallel()

	html := `<script>
	var leaderboard_data = [{"entity_id":"660271","est_ba":0.31,"est_slg":0.672,"est_woba":0.444,"xera":null,"barrels_per_bip":21.5,"hard_hit_percent":60.1,"exit_velocity_avg":"95.8","launch_angle_avg":"16.2","sweet_spot_percent":37.8}];
	</script>`

	rows, err := parseSavantLeaderboard(html)
	if err != nil {
		t.Fatalf("parseSavantLeaderboard returned error: %v", err)
	}

	enrichment, ok := rows[660271]
	if !ok {
		t.Fatal("missing enrichment for player 660271")
	}
	if enrichment.ExpectedBattingAvg == nil || *enrichment.ExpectedBattingAvg != 0.31 {
		t.Fatalf("ExpectedBattingAvg = %v, want 0.31", enrichment.ExpectedBattingAvg)
	}
	if enrichment.ExpectedERA != nil {
		t.Fatalf("ExpectedERA = %v, want nil", enrichment.ExpectedERA)
	}
	if enrichment.BarrelPct == nil || *enrichment.BarrelPct != 21.5 {
		t.Fatalf("BarrelPct = %v, want 21.5", enrichment.BarrelPct)
	}
}
