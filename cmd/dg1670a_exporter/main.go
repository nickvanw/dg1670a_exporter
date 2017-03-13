package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/nickvanw/dg1670a_exporter"
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

	c.Collect(nil)

}

func createExporter(modem string) (*dg1670aexporter.Exporter, error) {
	client := http.Client{}
	return dg1670aexporter.New(&client, modem)
}
