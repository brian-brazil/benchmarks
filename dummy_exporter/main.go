package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	numTimeseries = flag.Int("timeseries", 10000, "The number of timeseries to return.")
	listen        = flag.String("listen", ":1234", "Address to listen on.")
	metric        = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "example_timeseries_total",
			Help: "Timeseries being exposed",
		},
		[]string{"label"},
	)
)

type collector struct {
}

func (collector) Describe(ch chan<- *prometheus.Desc) {
	metric.Describe(ch)
}
func (collector) Collect(ch chan<- prometheus.Metric) {
	for i := 0; i < *numTimeseries; i++ {
		metric.WithLabelValues(fmt.Sprintf("labelabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij%d", i)).Add(math.Max(0, rand.NormFloat64()))
	}
	metric.Collect(ch)
}

func main() {
	prometheus.MustRegister(collector{})
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(*listen, nil)
}
