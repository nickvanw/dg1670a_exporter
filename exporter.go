package dg1670aexporter

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &Exporter{}

type Exporter struct {
	mu     sync.Mutex
	client *client
}

func New(httpClient *http.Client, url string) (*Exporter, error) {
	cl := &client{client: httpClient, url: url}
	return &Exporter{client: cl}, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {}
