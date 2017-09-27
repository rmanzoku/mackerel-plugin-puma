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

var graphdefStats = map[string]mp.Graphs{
	"workers": {
		Label: "Puma workers",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "workers", Label: "Active workers", Diff: false},
			{Name: "spawn_workers", Label: "Spawn workers", Diff: true},
			{Name: "removed_workers", Label: "Removed workers", Diff: true},
		},
	},
	"backlog": {
		Label: "Puma backlog",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "total_backlog", Label: "Total backlog", Diff: false},
		},
	},
	"backlog_stats": {
		Label: "Puma backlog stats",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "max_backlog", Label: "Max backlog", Diff: false},
			{Name: "ave_backlog", Label: "Average backlog", Diff: false},
			{Name: "min_backlog", Label: "Min backlog", Diff: false},
		},
	},
	"thread": {
		Label: "Puma threads",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "total_threads", Label: "Total threads", Diff: false},
		},
	},
	"thread_stats": {
		Label: "Puma thread stats",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "max_threads", Label: "Max running threads", Diff: false},
			{Name: "ave_threads", Label: "Average running threads", Diff: false},
			{Name: "min_threads", Label: "Min running threads", Diff: false},
		},
	},
	"phase": {
		Label: "Puma phase",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "phase", Label: "Active phase", Diff: false},
		},
	},
}

var graphdefGC = map[string]mp.Graphs{
	"gc.count": {
		Label: "Puma GC Count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "total", Label: "Total GC count", Stacked: false},
			{Name: "minor", Label: "Minor GC count", Stacked: true},
			{Name: "major", Label: "Major GC count", Stacked: true},
		},
	},
}

// PumaPlugin mackerel plugin for Puma
type PumaPlugin struct {
	Prefix string
	Host   string
	Port   string
	Token  string
	WithGC bool
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

func (s Stats) getBacklogMaxMinAveSum() (float64, float64, float64, float64) {
	var sum int
	var count int
	var max = s.WorkerStatus[0].LastStatus.Backlog
	var min = s.WorkerStatus[0].LastStatus.Backlog

	for _, v := range s.WorkerStatus {
		value := v.LastStatus.Backlog
		sum += value
		count++

		if max < value {
			max = value
		}

		if min > value {
			min = value
		}
	}

	ave := float64(sum) / float64(count)
	return float64(max), float64(min), ave, float64(sum)
}

func (s Stats) getRunningMaxMinAveSum() (float64, float64, float64, float64) {
	var sum int
	var count int
	var max = s.WorkerStatus[0].LastStatus.Running
	var min = s.WorkerStatus[0].LastStatus.Running

	for _, v := range s.WorkerStatus {
		value := v.LastStatus.Running
		sum += value
		count++

		if max < value {
			max = value
		}

		if min > value {
			min = value
		}
	}

	ave := float64(sum) / float64(count)
	return float64(max), float64(min), ave, float64(sum)
}

// FetchMetrics interface for mackerelplugin
func (p PumaPlugin) FetchMetrics() (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	stats, err := p.fetchStats()
	if err != nil {
		return nil, err
	}

	ret["workers"] = float64(stats.Workers)
	ret["spawn_workers"] = float64(stats.BootedWorkers)
	ret["removed_workers"] = float64(stats.OldWorkers)
	ret["phase"] = float64(stats.Phase)

	ret["max_backlog"], ret["min_backlog"], ret["ave_backlog"], ret["total_backlog"] = stats.getBacklogMaxMinAveSum()
	ret["max_threads"], ret["min_threads"], ret["ave_threads"], ret["total_threads"] = stats.getRunningMaxMinAveSum()

	if p.WithGC == false {
		return ret, nil
	}

	gcStats, err := p.fetchGCStats()
	if err != nil {
		return nil, err
	}

	ret["total"] = float64(gcStats.Count)
	ret["minor"] = float64(gcStats.MinorGcCount)
	ret["major"] = float64(gcStats.MajorGcCount)
	return ret, nil

}

// GraphDefinition interface for mackerelplugin
func (p PumaPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := graphdefStats

	if p.WithGC == false {
		return graphdef
	}

	for k, v := range graphdefGC {
		graphdef[k] = v
	}
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
		optWithGC   = flag.Bool("with-gc", false, "Output include GC stats for Puma 3.10.0~")
		optTempfile = flag.String("tempfile", "", "Temp file name")
	)
	flag.Parse()

	var puma PumaPlugin
	puma.Prefix = *optPrefix
	puma.Host = *optHost
	puma.Port = *optPort
	puma.Token = *optToken
	puma.WithGC = *optWithGC

	helper := mp.NewMackerelPlugin(puma)
	helper.Tempfile = *optTempfile
	helper.Run()
}
