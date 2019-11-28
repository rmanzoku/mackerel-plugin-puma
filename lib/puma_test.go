package mppuma

import (
	"encoding/json"
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	desired := 5

	var puma PumaPlugin

	graphdef := puma.GraphDefinition()

	if len(graphdef) != desired {
		t.Errorf("GraphDefinition: %d should be %d", len(graphdef), desired)
	}
}

func TestGraphDefinitionWithGC(t *testing.T) {
	desired := 9

	var puma PumaPlugin
	puma.WithGC = true

	graphdef := puma.GraphDefinition()

	if len(graphdef) != desired {
		t.Errorf("GraphDefinitionWithGC: %d should be %d", len(graphdef), desired)
	}
}

func TestGraphDefinitionCluster(t *testing.T) {

	statJSON := `{
	  "workers": 2,
	  "phase": 0,
	  "booted_workers": 2,
	  "old_workers": 0,
	  "worker_status": [
	    {
	      "pid": 1,
	      "index": 0,
	      "phase": 0,
	      "booted": true,
	      "last_checkin": "2018-04-17T01:24:16Z",
	      "last_status": {
	        "backlog": 1,
	        "running": 5,
	        "pool_capacity": 4
	      }
	    },
	    {
	      "pid": 2,
	      "index": 1,
	      "phase": 0,
	      "booted": true,
	      "last_checkin": "2018-04-17T01:24:16Z",
	      "last_status": {
	        "backlog": 1,
	        "running": 5,
	        "pool_capacity": 4
	      }
	    }
	  ]
	}`

	desired := map[string]float64{
		"workers":                       float64(2),
		"spawn_workers":                 float64(2),
		"removed_workers":               float64(0),
		"phase":                         float64(0),
		"backlog.worker0.backlog":       float64(1),
		"running.worker0.running":       float64(5),
		"running.worker0.pool_capacity": float64(4),
		"backlog.worker1.backlog":       float64(1),
		"running.worker1.running":       float64(5),
		"running.worker1.pool_capacity": float64(4),
	}

	var p PumaPlugin
	var stats Stats
	json.Unmarshal([]byte(statJSON), &stats)

	ret := p.fetchStatsMetrics(&stats)

	if len(ret) != len(desired) {
		t.Errorf("fetchStatsMetrics: len(ret) = %d should be len(desired) = %d", len(ret), len(desired))
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

func TestGraphDefinitionSingle(t *testing.T) {

	statJSON := `{
		"backlog": 1,
		"running": 5,
		"pool_capacity": 4
	}`

	desired := map[string]float64{
		"backlog":       float64(1),
		"running":       float64(5),
		"pool_capacity": float64(4),
	}

	var p PumaPlugin
	p.Single = true

	var stats Stats
	json.Unmarshal([]byte(statJSON), &stats)

	ret := p.fetchStatsMetrics(&stats)

	if len(ret) != len(desired) {
		t.Errorf("fetchStatsMetrics: len(ret) = %d should be len(desired) = %d", len(ret), len(desired))
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
