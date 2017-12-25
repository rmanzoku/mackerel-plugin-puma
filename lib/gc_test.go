package mppuma

import (
	"encoding/json"
	"testing"
)

func TestFetchGCStatsMetricsRuby20(t *testing.T) {

	gcStatJSON := `{
		"count": 4,
		"heap_used": 77,
		"heap_length": 77,
		"heap_increment": 0,
		"heap_live_num": 13071,
		"heap_free_num": 21512,
		"heap_final_num": 0,
		"total_allocated_object": 48497,
		"total_freed_object": 35426
	}`

	desired := map[string]float64{
		"total":       float64(4),
		"live_slots":  float64(13071),
		"free_slots":  float64(21512),
		"final_slots": float64(0),
	}

	var p PumaPlugin
	var gcStats GCStats
	json.Unmarshal([]byte(gcStatJSON), &gcStats)

	ret, _ := p.fetchGCStatsMetrics(&gcStats)

	if len(ret) != len(desired) {
		t.Errorf("fetchGCStatsMetrics: len(ret) = %d should be len(desired) = %d", len(ret), len(desired))
	}

	for k, v := range desired {
		if _, ok := ret[k]; !ok {
			t.Errorf("%s not xists", k)
		}

		if ret[k] != v {
			t.Errorf("%s should be %f, out %f", k, v, ret[k])
		}
	}
}
