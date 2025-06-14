package metric

import (
	"net/http"

	"github.com/FlyKarlik/orderService/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheus(cfg *config.Config) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(cfg.Infrastructure.Prometheus.Address, nil)
}
