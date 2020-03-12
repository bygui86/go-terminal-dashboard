
# Go terminal dashboard

## Run

1. server (first terminal window)
	```
	cd go-metrics
	go run main.go
	```

2. terminal-ui (second terminal window)
	```
	cd terminal-ui
	go run main.go -url http://localhost:6060/debug/vars
	```

---

## Links
- https://levelup.gitconnected.com/building-a-terminal-dashboard-in-golang-in-300-lines-of-code-3b9f83f363a8
- https://pkg.go.dev/expvar?tab=doc
- http://blog.ralch.com/tutorial/golang-metrics-with-expvar/
- https://sysdig.com/blog/golang-expvar-custom-metrics/
- https://www.freecodecamp.org/news/how-i-investigated-memory-leaks-in-go-using-pprof-on-a-large-codebase-4bec4325e192/
- https://scene-si.org/2018/08/06/basic-monitoring-of-go-apps-with-the-runtime-package/
- https://golang.org/pkg/runtime/pprof/
- https://github.com/google/gops
