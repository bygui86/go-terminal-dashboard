package rest

import "go-metrics/monitoring"

const (
	// Custom metrics
	// .. general
	metricsNamespace = "echoserver"
	metricsSubsystem = "rest"
	// .. opsProcessed
	opsProcessedKey  = "opsProcessed"
	opsProcessedName = "total_processed_ops"
	opsProcessedHelp = "Total number of processed operations"
)

// addCustomMetrics -
func addCustomMetrics(customMetrics monitoring.ICustomMetrics) {

	customMetrics.AddCounter(metricsNamespace, metricsSubsystem, opsProcessedName, opsProcessedHelp, opsProcessedKey)
}
