package mppuma

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdefStats = map[string]mp.Graphs{
	"workers": {
		Label: "Puma Workers",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "workers", Label: "Active workers", Diff: false},
			{Name: "spawn_workers", Label: "Spawn workers", Diff: true},
			{Name: "removed_workers", Label: "Removed workers", Diff: true},
		},
	},
	"backlog.#": {
		Label: "Puma Backlog",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "backlog", Label: "Backlog", Diff: false, Stacked: true},
		},
	},
	"running.#": {
		Label: "Puma Running Thread",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "running", Label: "running", Diff: false, Stacked: true},
		},
	},
	"phase": {
		Label: "Puma Phase",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "phase", Label: "Active phase", Diff: false},
		},
	},
}

// Stats is convered from /stats json
type Stats struct {
	Workers       int `json:"workers"`
	Phase         int `json:"phase"`
	BootedWorkers int `json:"booted_workers"`
	OldWorkers    int `json:"old_workers"`
	WorkerStatus  []struct {
		Pid         int       `json:"pid"`
		Index       int       `json:"index"`
		Phase       int       `json:"phase"`
		Booted      bool      `json:"booted"`
		LastCheckin time.Time `json:"last_checkin"`
		LastStatus  struct {
			Backlog int `json:"backlog"`
			Running int `json:"running"`
		} `json:"last_status"`
	} `json:"worker_status"`
}

// Fetch /stats
func (p PumaPlugin) fetchStats() (*Stats, error) {

	var stats Stats

	uri := fmt.Sprintf("http://%s:%s/%s?token=%s", p.Host, p.Port, "stats", p.Token)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
