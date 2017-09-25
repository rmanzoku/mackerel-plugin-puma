package mppuma

import (
	"flag"
	"fmt"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"puma.workers": {
		Label: "Puma workers",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "workers", Label: "Active workers", Diff: false},
		},
	},
}

// PumaPlugin mackerel plugin for Puma
type PumaPlugin struct {
	URI string
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

// FetchMetrics interface for mackerelplugin
func (n PumaPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	stat["workers"] = 2.0
	return stat, nil

}

// GraphDefinition interface for mackerelplugin
func (n PumaPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	var (
		// optPrefix = flag.String("metric-key-prefix", "", "Metric key prefix")
		optHost     = flag.String("host", "127.0.0.1", "The bind url to use for the control server")
		optPort     = flag.String("port", "9293", "The bind port to use for the control server")
		optToken    = flag.String("token", "", "The token to use as authentication for the control server")
		optTempfile = flag.String("tempfile", "", "Temp file name")
	)
	flag.Parse()

	var puma PumaPlugin
	puma.URI = fmt.Sprintf("%s:%s?token=%s", *optHost, *optPort, *optToken)

	helper := mp.NewMackerelPlugin(puma)
	helper.Tempfile = *optTempfile
	helper.Run()
}
