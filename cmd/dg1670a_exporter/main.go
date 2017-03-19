package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/nickvanw/dg1670a_exporter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsPath = flag.String("metrics.path", "/metrics", "path to fetch metrics")
	metricsAddr = flag.String("metrics.addr", ":9191", "address to listen")

	modemAddr = flag.String("modem.host", "http://192.168.100.1/cgi-bin/status_cgi", "url to fetch modem")
)

func main() {
	flag.Parse()

	c, err := createExporter(*modemAddr)
	if err != nil {
		log.Fatalf("unable to create client: %s", err)
	}

	prometheus.MustRegister(c)

	mux := http.NewServeMux()
	mux.Handle(*metricsPath, prometheus.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Arris DG1670a Exporter</title></head>
			<body>
			<h1>Arris DG1670a Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})
	loggedMux := handlers.LoggingHandler(os.Stdout, mux)
	if err := http.ListenAndServe(*metricsAddr, loggedMux); err != nil {
		log.Fatalf("unable to start metrics server: %s", err)
	}
}

func createExporter(modem string) (*dg1670aexporter.Exporter, error) {
	client := http.Client{}
	return dg1670aexporter.New(&client, modem)
}
