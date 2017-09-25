package mppuma

import (
	"flag"
	"fmt"

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

// FetchMetrics interface for mackerelplugin
func (n PumaPlugin) FetchMetrics() (map[string]interface{}, error) {
	return nil, nil

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
