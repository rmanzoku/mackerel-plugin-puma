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

func TestFetchGCStatsMetricsRuby21(t *testing.T) {

	gcStatJSON := `{
                "count": 5,
                "heap_used": 75,
                "heap_length": 81,
                "heap_increment": 6,
                "heap_live_slot": 29819,
                "heap_free_slot": 752,
                "heap_final_slot": 0,
                "heap_swept_slot": 3861,
                "heap_eden_page_length": 75,
                "heap_tomb_page_length": 0,
                "total_allocated_object": 48629,
                "total_freed_object": 18810,
                "malloc_increase": 1076728,
                "malloc_limit": 16777216,
                "minor_gc_count": 3,
                "major_gc_count": 2,
                "remembered_shady_object": 151,
                "remembered_shady_object_limit": 300,
                "old_object": 5842,
                "old_object_limit": 11684,
                "oldmalloc_increase": 1077176,
                "oldmalloc_limit": 16777216
        }`

	desired := map[string]float64{
		"total":            float64(5),
		"minor":            float64(3),
		"major":            float64(2),
		"live_slots":       float64(29819),
		"free_slots":       float64(752),
		"final_slots":      float64(0),
		"old_count":        float64(5842),
		"old_limit":        float64(11684),
		"old_malloc_bytes": float64(1077176),
		"old_malloc_limit": float64(16777216),
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

func TestFetchGCStatsMetricsRuby22(t *testing.T) {

	gcStatJSON := `{
              "count": 5,
              "heap_allocated_pages": 74,
              "heap_sorted_length": 75,
              "heap_allocatable_pages": 0,
              "heap_available_slots": 30165,
              "heap_live_slots": 29204,
              "heap_free_slots": 961,
              "heap_final_slots": 0,
              "heap_marked_slots": 8805,
              "heap_swept_slots": 7197,
              "heap_eden_pages": 73,
              "heap_tomb_pages": 1,
              "total_allocated_pages": 74,
              "total_freed_pages": 0,
              "total_allocated_objects": 50243,
              "total_freed_objects": 21039,
              "malloc_increase_bytes": 152840,
              "malloc_increase_bytes_limit": 16777216,
              "minor_gc_count": 3,
              "major_gc_count": 2,
              "remembered_wb_unprotected_objects": 156,
              "remembered_wb_unprotected_objects_limit": 278,
              "old_objects": 7418,
              "old_objects_limit": 10932,
              "oldmalloc_increase_bytes": 153288,
              "oldmalloc_increase_bytes_limit": 16777216
        }`

	desired := map[string]float64{
		"total":            float64(5),
		"minor":            float64(3),
		"major":            float64(2),
		"available_slots":  float64(30165),
		"live_slots":       float64(29204),
		"free_slots":       float64(961),
		"final_slots":      float64(0),
		"marked_slots":     float64(8805),
		"old_count":        float64(7418),
		"old_limit":        float64(10932),
		"old_malloc_bytes": float64(153288),
		"old_malloc_limit": float64(16777216),
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

var TestFetchGCStatsMetricsRuby23 = TestFetchGCStatsMetricsRuby22

func TestFetchGCStatsMetricsRuby24(t *testing.T) {

	gcStatJSON := `{
              "count": 8,
              "heap_allocated_pages": 65,
              "heap_sorted_length": 65,
              "heap_allocatable_pages": 0,
              "heap_available_slots": 26494,
              "heap_live_slots": 26269,
              "heap_free_slots": 225,
              "heap_final_slots": 0,
              "heap_marked_slots": 11738,
              "heap_eden_pages": 65,
              "heap_tomb_pages": 0,
              "total_allocated_pages": 65,
              "total_freed_pages": 0,
              "total_allocated_objects": 66208,
              "total_freed_objects": 39939,
              "malloc_increase_bytes": 105872,
              "malloc_increase_bytes_limit": 16777216,
              "minor_gc_count": 7,
              "major_gc_count": 1,
              "remembered_wb_unprotected_objects": 165,
              "remembered_wb_unprotected_objects_limit": 286,
              "old_objects": 10929,
              "old_objects_limit": 14302,
              "oldmalloc_increase_bytes": 1351056,
              "oldmalloc_increase_bytes_limit": 16777216
        }`

	desired := map[string]float64{
		"total":            float64(8),
		"minor":            float64(7),
		"major":            float64(1),
		"available_slots":  float64(26494),
		"live_slots":       float64(26269),
		"free_slots":       float64(225),
		"final_slots":      float64(0),
		"marked_slots":     float64(11738),
		"old_count":        float64(10929),
		"old_limit":        float64(14302),
		"old_malloc_bytes": float64(1351056),
		"old_malloc_limit": float64(16777216),
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

var TestFetchGCStatsMetricsRuby25 = TestFetchGCStatsMetricsRuby24
var TestFetchGCStatsMetricsRuby26 = TestFetchGCStatsMetricsRuby25
