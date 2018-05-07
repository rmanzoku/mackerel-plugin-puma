package mppuma

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	// Single mode
	Backlog int `json:"backlog"`
	Running int `json:"running"`
}

// GET request to /stats
func (p PumaPlugin) getStatsAPI() (*Stats, error) {

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

// Fetch /stats
func (p PumaPlugin) fetchStatsMetrics(stats *Stats) map[string]float64 {
	ret := make(map[string]float64)

	if p.Single == true {
		ret["backlog.worker0.backlog"] = float64(stats.Backlog)
		ret["running.worker0.running"] = float64(stats.Running)
		return ret
	}

	ret["workers"] = float64(stats.Workers)
	ret["spawn_workers"] = float64(stats.BootedWorkers)
	ret["removed_workers"] = float64(stats.OldWorkers)
	ret["phase"] = float64(stats.Phase)

	for _, v := range stats.WorkerStatus {
		ret["backlog.worker"+strconv.Itoa(v.Index)+".backlog"] = float64(v.LastStatus.Backlog)
		ret["running.worker"+strconv.Itoa(v.Index)+".running"] = float64(v.LastStatus.Running)
	}

	return ret

}
