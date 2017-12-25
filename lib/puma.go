package mppuma

import (
	"flag"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// PumaPlugin mackerel plugin for Puma
type PumaPlugin struct {
	Prefix string
	Host   string
	Port   string
	Token  string
	WithGC bool
}

// FetchMetrics interface for mackerelplugin
func (p PumaPlugin) FetchMetrics() (map[string]float64, error) {
	ret := make(map[string]float64)

	stats, err := p.getStatsAPI()
	if err != nil {
		return nil, err
	}

	statsMetrics := p.fetchStatsMetrics(stats)

	for k, v := range statsMetrics {
		ret[k] = v
	}

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

	ret["available_slots"] = float64(gcStats.HeapAvailableSlots)
	ret["live_slots"] = float64(gcStats.HeapLiveSlots)
	ret["free_slots"] = float64(gcStats.HeapFreeSlots)
	ret["final_slots"] = float64(gcStats.HeapFinalSlots)
	ret["marked_slots"] = float64(gcStats.HeapMarkedSlots)

	ret["old_count"] = float64(gcStats.OldObjects)
	ret["old_limit"] = float64(gcStats.OldObjectsLimit)

	ret["old_malloc_bytes"] = float64(gcStats.OldmallocIncreaseBytes)
	ret["old_malloc_limit"] = float64(gcStats.OldmallocIncreaseBytesLimit)

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
