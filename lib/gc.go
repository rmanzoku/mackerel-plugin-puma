package mppuma

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

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
	"gc.heap_slot": {
		Label: "Puma GC Heap Slot",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "available_slots", Label: "Heap available slots", Stacked: false},
			{Name: "live_slots", Label: "Heap live slots", Stacked: true},
			{Name: "free_slots", Label: "Heap free slots", Stacked: true},
			{Name: "final_slots", Label: "Heap final slots", Stacked: false},
			{Name: "marked_slots", Label: "Heap marked slots", Stacked: false},
		},
	},
	"gc.old_objects": {
		Label: "Puma GC Old Objects",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "old_count", Label: "Old Objects", Stacked: false},
			{Name: "old_limit", Label: "Old Objects Limit", Stacked: true},
		},
	},
	"gc.old_malloc": {
		Label: "Puma GC Old Malloc Increase",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "old_malloc_bytes", Label: "Old Malloc Increase Bytes", Stacked: false},
			{Name: "old_malloc_limit", Label: "Old Malloc Increase Bytes Limit", Stacked: true},
		},
	},
}

type response map[string]float64

// GCStats is convered from /gc-stats json
type GCStats struct {
	// Ruby2.0
	Count                json.Number `json:"count"`
	HeapFinalNum         json.Number `json:"heap_final_num"`
	HeapFreeNum          json.Number `json:"heap_free_num"`
	HeapIncrement        json.Number `json:"heap_increment"`
	HeapLength           json.Number `json:"heap_length"`
	HeapLiveNum          json.Number `json:"heap_live_num"`
	HeapUsed             json.Number `json:"heap_used"`
	TotalAllocatedObject json.Number `json:"total_allocated_object"`
	TotalFreedObject     json.Number `json:"total_freed_object"`
	// Added since Ruby2.1
	HeapLiveSlot               json.Number `json:"heap_live_slot"`
	HeapFreeSlot               json.Number `json:"heap_free_slot"`
	HeapFinalSlot              json.Number `json:"heap_final_slot"`
	HeapSweptSlot              json.Number `json:"heap_swept_slot"`
	HeapEdenPageLength         json.Number `json:"heap_eden_page_length"`
	HeapTombPageLength         json.Number `json:"heap_tomb_page_length"`
	MallocIncrease             json.Number `json:"malloc_increase"`
	MallocLimit                json.Number `json:"malloc_limit"`
	MinorGcCount               json.Number `json:"minor_gc_count"`
	MajorGcCount               json.Number `json:"major_gc_count"`
	RememberedShadyObject      json.Number `json:"remembered_shady_object"`
	RememberedShadyObjectLimit json.Number `json:"remembered_shady_object_limit"`
	OldObject                  json.Number `json:"old_object"`
	OldObjectLimit             json.Number `json:"old_object_limit"`
	OldmallocIncrease          json.Number `json:"oldmalloc_increase"`
	OldmallocLimit             json.Number `json:"oldmalloc_limit"`
	// Added since Ruby2.2
	HeapAllocatedPages                  json.Number `json:"heap_allocated_pages"`
	HeapSortedLength                    json.Number `json:"heap_sorted_length"`
	HeapAllocatablePages                json.Number `json:"heap_allocatable_pages"`
	HeapAvailableSlots                  json.Number `json:"heap_available_slots"`
	HeapLiveSlots                       json.Number `json:"heap_live_slots"`
	HeapFreeSlots                       json.Number `json:"heap_free_slots"`
	HeapFinalSlots                      json.Number `json:"heap_final_slots"`
	HeapMarkedSlots                     json.Number `json:"heap_marked_slots"`
	HeapSweptSlots                      json.Number `json:"heap_swept_slots"`
	HeapEdenPages                       json.Number `json:"heap_eden_pages"`
	HeapTombPages                       json.Number `json:"heap_tomb_pages"`
	TotalAllocatedPages                 json.Number `json:"total_allocated_pages"`
	TotalFreedPages                     json.Number `json:"total_freed_pages"`
	TotalAllocatedObjects               json.Number `json:"total_allocated_objects"`
	TotalFreedObjects                   json.Number `json:"total_freed_objects"`
	MallocIncreaseBytes                 json.Number `json:"malloc_increase_bytes"`
	MallocIncreaseBytesLimit            json.Number `json:"malloc_increase_bytes_limit"`
	RememberedWbUnprotectedObjects      json.Number `json:"remembered_wb_unprotected_objects"`
	RememberedWbUnprotectedObjectsLimit json.Number `json:"remembered_wb_unprotected_objects_limit"`
	OldObjects                          json.Number `json:"old_objects"`
	OldObjectsLimit                     json.Number `json:"old_objects_limit"`
	OldmallocIncreaseBytes              json.Number `json:"oldmalloc_increase_bytes"`
	OldmallocIncreaseBytesLimit         json.Number `json:"oldmalloc_increase_bytes_limit"`
	// Ruby2.3 is same as Ruby2.2
	// Ruby2.4 is almost same Ruby2.3 (deletes heap_swept_slots)
}

// Fetch /gc-stats
func (p PumaPlugin) getGCStatsAPI() (*GCStats, error) {

	var gcStats GCStats

	var client http.Client

	if p.Sock != "" {
		client = http.Client{Transport: &http.Transport{
			Dial: func(proto, addr string) (conn net.Conn, err error) {
				return net.Dial("unix", p.Sock)
			},
		}}
	} else {
		client = http.Client{}
	}

	uri := fmt.Sprintf("http://%s:%s/%s?token=%s", p.Host, p.Port, "gc-stats", p.Token)
	resp, err := client.Get(uri)
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

func (p PumaPlugin) fetchGCStatsMetrics(gcStats *GCStats) (map[string]float64, error) {
	ret := make(map[string]float64)

	// gc.count
	ret["total"], _ = gcStats.Count.Float64()

	if gcStats.MinorGcCount.String() != "" {
		ret["minor"], _ = gcStats.MinorGcCount.Float64()
	}
	if gcStats.MajorGcCount.String() != "" {
		ret["major"], _ = gcStats.MajorGcCount.Float64()
	}

	// gc.heap_slot
	if gcStats.HeapAvailableSlots.String() != "" {
		ret["available_slots"], _ = gcStats.HeapAvailableSlots.Float64()
	}

	for _, v := range []json.Number{gcStats.HeapLiveNum, gcStats.HeapLiveSlot, gcStats.HeapLiveSlots} {
		if v.String() != "" {
			ret["live_slots"], _ = v.Float64()
		}
	}

	for _, v := range []json.Number{gcStats.HeapFreeNum, gcStats.HeapFreeSlot, gcStats.HeapFreeSlots} {
		if v.String() != "" {
			ret["free_slots"], _ = v.Float64()
		}
	}

	for _, v := range []json.Number{gcStats.HeapFinalNum, gcStats.HeapFinalSlot, gcStats.HeapFinalSlots} {
		if v.String() != "" {
			ret["final_slots"], _ = v.Float64()
		}
	}

	if gcStats.HeapMarkedSlots.String() != "" {
		ret["marked_slots"], _ = gcStats.HeapMarkedSlots.Float64()
	}

	// old
	for _, v := range []json.Number{gcStats.OldObject, gcStats.OldObjects} {
		if v.String() != "" {
			ret["old_count"], _ = v.Float64()
		}
	}
	for _, v := range []json.Number{gcStats.OldObjectLimit, gcStats.OldObjectsLimit} {
		if v.String() != "" {
			ret["old_limit"], _ = v.Float64()
		}
	}
	for _, v := range []json.Number{gcStats.OldmallocIncrease, gcStats.OldmallocIncreaseBytes} {
		if v.String() != "" {
			ret["old_malloc_bytes"], _ = v.Float64()
		}
	}
	for _, v := range []json.Number{gcStats.OldmallocLimit, gcStats.OldmallocIncreaseBytesLimit} {
		if v.String() != "" {
			ret["old_malloc_limit"], _ = v.Float64()
		}
	}

	return ret, nil

}
