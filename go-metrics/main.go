package main

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/google/gops/agent"
	
	"go-metrics/logging"
	"go-metrics/monitoring"
	"go-metrics/rest"
)

func main() {
	// from "github.com/pkg/profile"
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.BlockProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.GoroutineProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MutexProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.ThreadcreationProfile, profile.ProfilePath(".")).Stop()

	logging.Log.Infoln("[MAIN] Starting echo-server...")

	monitorServer := startMonitor()
	defer monitorServer.Shutdown()

	restServer := startRest(monitorServer.CustomMetrics)
	defer restServer.Shutdown()

	startGopsAgent()
	
	startDebugServer()

	logging.Log.Infoln("[MAIN] echo-server ready!")

	startSysCallChannel()
}

func startMonitor() *monitoring.MonitorServer {
	server, err := monitoring.NewMonitorServer()
	if err != nil {
		logging.Log.Errorf("[MAIN] Monitoring server creation failed: %s", err.Error())
		os.Exit(404)
	}
	logging.Log.Debugln("[MAIN] Monitoring server successfully created")

	server.Start()
	logging.Log.Debugln("[MAIN] Monitoring successfully started")

	return server
}

func startRest(customMetrics monitoring.ICustomMetrics) *rest.RestServer {
	server, err := rest.NewRestServer(customMetrics)
	if err != nil {
		logging.Log.Errorf("[MAIN] Echo server creation failed: %s", err.Error())
		os.Exit(404)
	}
	logging.Log.Debugln("[MAIN] Echo server successfully created")

	server.Start()
	logging.Log.Debugln("[MAIN] Echo successfully started")

	return server
}

func startGopsAgent() {
	err := agent.Listen(agent.Options{})
	if err != nil {
		logging.Log.Errorf("[MAIN] gops agent start failed: %s", err.Error())
		os.Exit(404)
	}
}

func startDebugServer() {
	logging.Log.Info("Start debug REST server on port 6060")
	go http.ListenAndServe(":6060", nil)
}

func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
	logging.Log.Warnln("[MAIN] Termination signal received!")
}
