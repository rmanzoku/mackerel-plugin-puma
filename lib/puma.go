package mppuma

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"workers": {
		Label: "Puma workers",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "workers", Label: "Active workers", Diff: false},
		},
	},
	"phase": {
		Label: "Puma phase",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "phase", Label: "Active phase", Diff: false},
		},
	},
	"gc.count": {
		Label: "GC count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "gccount", Label: "GC count", Diff: false},
		},
	},
}

// PumaPlugin mackerel plugin for Puma
type PumaPlugin struct {
	Prefix string
	Host   string
	Port   string
	Token  string
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

// GCStats is convered from /gc-stats json
type GCStats struct {
	Count                               int `json:"count"`
	HeapAllocatedPages                  int `json:"heap_allocated_pages"`
	HeapSortedLength                    int `json:"heap_sorted_length"`
	HeapAllocatablePages                int `json:"heap_allocatable_pages"`
	HeapAvailableSlots                  int `json:"heap_available_slots"`
	HeapLiveSlots                       int `json:"heap_live_slots"`
	HeapFreeSlots                       int `json:"heap_free_slots"`
	HeapFinalSlots                      int `json:"heap_final_slots"`
	HeapMarkedSlots                     int `json:"heap_marked_slots"`
	HeapEdenPages                       int `json:"heap_eden_pages"`
	HeapTombPages                       int `json:"heap_tomb_pages"`
	TotalAllocatedPages                 int `json:"total_allocated_pages"`
	TotalFreedPages                     int `json:"total_freed_pages"`
	TotalAllocatedObjects               int `json:"total_allocated_objects"`
	TotalFreedObjects                   int `json:"total_freed_objects"`
	MallocIncreaseBytes                 int `json:"malloc_increase_bytes"`
	MallocIncreaseBytesLimit            int `json:"malloc_increase_bytes_limit"`
	MinorGcCount                        int `json:"minor_gc_count"`
	MajorGcCount                        int `json:"major_gc_count"`
	RememberedWbUnprotectedObjects      int `json:"remembered_wb_unprotected_objects"`
	RememberedWbUnprotectedObjectsLimit int `json:"remembered_wb_unprotected_objects_limit"`
	OldObjects                          int `json:"old_objects"`
	OldObjectsLimit                     int `json:"old_objects_limit"`
	OldmallocIncreaseBytes              int `json:"oldmalloc_increase_bytes"`
	OldmallocIncreaseBytesLimit         int `json:"oldmalloc_increase_bytes_limit"`
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

// Fetch /gc-stats
func (p PumaPlugin) fetchGCStats() (*GCStats, error) {

	var gcStats GCStats

	uri := fmt.Sprintf("http://%s:%s/%s?token=%s", p.Host, p.Port, "gc-stats", p.Token)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&gcStats); err != nil {
		return nil, err
	}

	return &gcStats, nil
}

// FetchMetrics interface for mackerelplugin
func (p PumaPlugin) FetchMetrics() (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	stats, err := p.fetchStats()
	if err != nil {
		return nil, err
	}
	gcStats, err := p.fetchGCStats()
	if err != nil {
		return nil, err
	}

	ret["workers"] = float64(stats.Workers)
	ret["phase"] = float64(stats.Phase)
	ret["gccount"] = float64(gcStats.Count)
	return ret, nil

}

// GraphDefinition interface for mackerelplugin
func (p PumaPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p PumaPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "puma"
	}
	return p.Prefix
}

// Do the plugin
func Do() {
	var (
		optPrefix   = flag.String("metric-key-prefix", "puma", "Metric key prefix")
		optHost     = flag.String("host", "127.0.0.1", "The bind url to use for the control server")
		optPort     = flag.String("port", "9293", "The bind port to use for the control server")
		optToken    = flag.String("token", "", "The token to use as authentication for the control server")
		optTempfile = flag.String("tempfile", "", "Temp file name")
	)
	flag.Parse()

	var puma PumaPlugin
	puma.Prefix = *optPrefix
	puma.Host = *optHost
	puma.Port = *optPort
	puma.Token = *optToken

	helper := mp.NewMackerelPlugin(puma)
	helper.Tempfile = *optTempfile
	helper.Run()
}
