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
	Sock   string
	Token  string
	Single bool
	WithGC bool
}

func merge(m1, m2 map[string]float64) map[string]float64 {
	ans := make(map[string]float64)

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}
	return (ans)
}

// FetchMetrics interface for mackerelplugin
func (p PumaPlugin) FetchMetrics() (map[string]float64, error) {
	ret := make(map[string]float64)

	stats, err := p.getStatsAPI()
	if err != nil {
		return nil, err
	}

	ret = p.fetchStatsMetrics(stats)

	if p.WithGC == false {
		return ret, nil
	}

	gcStats, err := p.getGCStatsAPI()
	if err != nil {
		return nil, err
	}

	gcStatsMetrics, _ := p.fetchGCStatsMetrics(gcStats)

	ret = merge(ret, gcStatsMetrics)

	return ret, nil

}

// GraphDefinition interface for mackerelplugin
func (p PumaPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := graphdefStats

	if p.Single == true {
		graphdef = graphdefStatsSingle
	}

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
		optSock     = flag.String("sock", "", "The bind socket to use for the control server")
		optToken    = flag.String("token", "", "The token to use as authentication for the control server")
		optSingle   = flag.Bool("single", false, "Puma in single mode")
		optWithGC   = flag.Bool("with-gc", false, "Output include GC stats for Puma 3.10.0~")
		optTempfile = flag.String("tempfile", "", "Temp file name")
	)
	flag.Parse()

	var puma PumaPlugin
	puma.Prefix = *optPrefix
	puma.Host = *optHost
	puma.Port = *optPort
	puma.Sock = *optSock
	puma.Token = *optToken
	puma.Single = *optSingle
	puma.WithGC = *optWithGC

	helper := mp.NewMackerelPlugin(puma)
	helper.Tempfile = *optTempfile
	helper.Run()
}
