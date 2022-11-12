package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RootProcessedLoops() {
	rootProcessedLoops.Inc()
}

func RootStartLoops() {
	rootStartLoops.Inc()
}

var (
	rootProcessedLoops = promauto.NewCounter(prometheus.CounterOpts{
		Name: "root_processed_loops_total",
		Help: "The total number of processed loops",
	})

	rootStartLoops = promauto.NewCounter(prometheus.CounterOpts{
		Name: "root_start_loops_total",
		Help: "The total number of loops started.",
	})
)

/*
kubectl port-forward pod/active-incident 2112:2112
http://localhost:2112/metrics
*/
func StartMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
