package mppuma

import (
	"encoding/json"
	"errors"
	"fmt"
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

// GCStats is convered from /gc-stats json
type GCStats struct {
	// Ruby2.0
	Count                int64 `json:"count"`
	HeapFinalNum         int64 `json:"heap_final_num"`
	HeapFreeNum          int64 `json:"heap_free_num"`
	HeapIncrement        int64 `json:"heap_increment"`
	HeapLength           int64 `json:"heap_length"`
	HeapLiveNum          int64 `json:"heap_live_num"`
	HeapUsed             int64 `json:"heap_used"`
	TotalAllocatedObject int64 `json:"total_allocated_object"`
	TotalFreedObject     int64 `json:"total_freed_object"`
	// Added since Ruby2.1
	HeapLiveSlot               int `json:"heap_live_slot"`
	HeapFreeSlot               int `json:"heap_free_slot"`
	HeapFinalSlot              int `json:"heap_final_slot"`
	HeapSweptSlot              int `json:"heap_swept_slot"`
	HeapEdenPageLength         int `json:"heap_eden_page_length"`
	HeapTombPageLength         int `json:"heap_tomb_page_length"`
	MallocIncrease             int `json:"malloc_increase"`
	MallocLimit                int `json:"malloc_limit"`
	MinorGcCount               int `json:"minor_gc_count"`
	MajorGcCount               int `json:"major_gc_count"`
	RememberedShadyObject      int `json:"remembered_shady_object"`
	RememberedShadyObjectLimit int `json:"remembered_shady_object_limit"`
	OldObject                  int `json:"old_object"`
	OldObjectLimit             int `json:"old_object_limit"`
	OldmallocIncrease          int `json:"oldmalloc_increase"`
	OldmallocLimit             int `json:"oldmalloc_limit"`
	// Added since Ruby2.2
	HeapAllocatedPages                  int `json:"heap_allocated_pages"`
	HeapSortedLength                    int `json:"heap_sorted_length"`
	HeapAllocatablePages                int `json:"heap_allocatable_pages"`
	HeapAvailableSlots                  int `json:"heap_available_slots"`
	HeapLiveSlots                       int `json:"heap_live_slots"`
	HeapFreeSlots                       int `json:"heap_free_slots"`
	HeapFinalSlots                      int `json:"heap_final_slots"`
	HeapMarkedSlots                     int `json:"heap_marked_slots"`
	HeapSweptSlots                      int `json:"heap_swept_slots"`
	HeapEdenPages                       int `json:"heap_eden_pages"`
	HeapTombPages                       int `json:"heap_tomb_pages"`
	TotalAllocatedPages                 int `json:"total_allocated_pages"`
	TotalFreedPages                     int `json:"total_freed_pages"`
	TotalAllocatedObjects               int `json:"total_allocated_objects"`
	TotalFreedObjects                   int `json:"total_freed_objects"`
	MallocIncreaseBytes                 int `json:"malloc_increase_bytes"`
	MallocIncreaseBytesLimit            int `json:"malloc_increase_bytes_limit"`
	RememberedWbUnprotectedObjects      int `json:"remembered_wb_unprotected_objects"`
	RememberedWbUnprotectedObjectsLimit int `json:"remembered_wb_unprotected_objects_limit"`
	OldObjects                          int `json:"old_objects"`
	OldObjectsLimit                     int `json:"old_objects_limit"`
	OldmallocIncreaseBytes              int `json:"oldmalloc_increase_bytes"`
	OldmallocIncreaseBytesLimit         int `json:"oldmalloc_increase_bytes_limit"`
	// Ruby2.3 is same as Ruby2.2
	// Ruby2.4 is almost same Ruby2.3 (deletes heap_swept_slots)
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
