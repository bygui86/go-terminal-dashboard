package main

import (
	"flag"

	"terminal-ui/application"
)

var debugUrl = flag.String("url", "", "url for debug/vars, e.g. http://foo.com/debug/vars")

/*
	✅ done
	❌ todo
	⚠️ in-progress / to be improved
	❓ TBD

	TODO improvements
		❌ flags/env-vars to make dashboard config more flexible
		❌ exit only on 'q' key event
		❌ improve logging
			❌❓ log to file as stdout is dedicated to the dashboard
		❌ implement multiple goroutines
		❌ add os signals listener
		❌❓ open multiple terminal windows: first for logs, second to show the dashboard
		❌❓ solve the 'GCCPUFraction *widgets.Gauge' issue causing dashboard crash
		❌❓ get terminal windows size
*/
func main() {
	flag.Parse()
	if len(*debugUrl) == 0 {
		flag.Usage()
		return
	}

	// TODO run in a separate goroutine
	application.Run(*debugUrl)

	// TODO add os signals listener
}
