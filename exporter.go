package dg1670aexporter

import (
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &Exporter{}

type Exporter struct {
	mu     sync.Mutex
	client *client

	m *ModemCollector
}

func New(httpClient *http.Client, url string) (*Exporter, error) {
	cl := &client{client: httpClient, url: url}
	modemCollector := NewModemCollector(cl)
	return &Exporter{client: cl, m: modemCollector}, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if err := e.m.collect(); err != nil {
		log.Println("failed to collect modem metrics:", err)
	}

	for _, metric := range e.m.collectorList() {
		metric.Collect(ch)
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.m.describe(ch)
}
