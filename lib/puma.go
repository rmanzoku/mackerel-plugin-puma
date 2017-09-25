package mppuma

import (
	"flag"
	"fmt"
)

// Do the plugin
func Do() {
	var (
		// optPrefix = flag.String("metric-key-prefix", "", "Metric key prefix")
		optHost  = flag.String("host", "127.0.0.1", "The bind url to use for the control server")
		optPort  = flag.String("port", "9293", "The bind port to use for the control server")
		optToken = flag.String("token", "", "The token to use as authentication for the control server")
	)
	flag.Parse()

	uri := *optHost + ":" + *optPort + "?token=" + *optToken

	fmt.Println(uri)
}
