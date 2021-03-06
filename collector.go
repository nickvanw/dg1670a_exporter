package dg1670aexporter

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "dg1670"

var (
	labelsDownstream = []string{"downstream", "dcid"}
	labelsUpstream   = []string{"upstream", "ucid"}
)

type ModemCollector struct {
	c *client

	// Downstream metrics
	DownstreamFreq *prometheus.GaugeVec

	DownstreamPower *prometheus.GaugeVec

	DownstreamSNR *prometheus.GaugeVec

	DownstreamModulation *prometheus.GaugeVec

	DownstreamOctets *prometheus.GaugeVec

	DownstreamCorrecteds *prometheus.GaugeVec

	DownstreamUncorrectables *prometheus.GaugeVec

	// Upstrem Metrics

	UpstreamFreq *prometheus.GaugeVec

	UpstreamPower *prometheus.GaugeVec

	UpstreamSymbolRate *prometheus.GaugeVec

	UpstreamModulation *prometheus.GaugeVec
}

func NewModemCollector(c *client) *ModemCollector {
	return &ModemCollector{
		c: c,

		DownstreamFreq: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_freq_hertz",
				Help:      "Modem Downstream Frequency (Hz)",
			},
			labelsDownstream,
		),

		DownstreamPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_power_dbmv",
				Help:      "Modem Downstream Power (dBmV)",
			},
			labelsDownstream,
		),

		DownstreamSNR: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_snr_db",
				Help:      "Modem Downstream SNR (dB)",
			},
			labelsDownstream,
		),

		DownstreamModulation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_modulation_qam",
				Help:      "Modem Downstream Modulation (QAM)",
			},
			labelsDownstream,
		),

		DownstreamOctets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_octets_total",
				Help:      "Modem Downstream Octets",
			},
			labelsDownstream,
		),

		DownstreamCorrecteds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_correcteds_total",
				Help:      "Modem Downstream Correcteds",
			},
			labelsDownstream,
		),

		DownstreamUncorrectables: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_uncorrectables_total",
				Help:      "Modem Downstream Uncorrectables",
			},
			labelsDownstream,
		),

		UpstreamFreq: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_freq_hertz",
				Help:      "Modem Upstream Frequency (MHz)",
			},
			labelsUpstream,
		),

		UpstreamPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_power_dbmv",
				Help:      "Modem Upstream Power (dBmV)",
			},
			labelsUpstream,
		),

		UpstreamSymbolRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_symbol_rate",
				Help:      "Modem Upstream Symbol Rate (kSym/s)",
			},
			labelsUpstream,
		),

		UpstreamModulation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_modulation_qam",
				Help:      "Modem Upstream Modulation",
			},
			labelsUpstream,
		),
	}
}

func (m *ModemCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		m.DownstreamCorrecteds,
		m.DownstreamSNR,
		m.DownstreamOctets,
		m.DownstreamModulation,
		m.DownstreamUncorrectables,
		m.DownstreamFreq,
		m.DownstreamPower,
		m.UpstreamFreq,
		m.UpstreamModulation,
		m.UpstreamSymbolRate,
		m.UpstreamPower,
	}
}

func (m *ModemCollector) collect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := m.c.fetch(ctx)
	if err != nil {
		return err
	}
	for id, node := range data.ds {
		downstreamID := strconv.Itoa(id + 1)
		downstreamDCID := strconv.FormatInt(node.DCID, 10)
		m.DownstreamSNR.WithLabelValues(downstreamID, downstreamDCID).Set(node.SNR)
		m.DownstreamFreq.WithLabelValues(downstreamID, downstreamDCID).Set(node.Freq)
		m.DownstreamPower.WithLabelValues(downstreamID, downstreamDCID).Set(node.Power)
		m.DownstreamOctets.WithLabelValues(downstreamID, downstreamDCID).Set(float64(node.Octets))
		m.DownstreamModulation.WithLabelValues(downstreamID, downstreamDCID).Set(float64(node.Modulation))
		m.DownstreamCorrecteds.WithLabelValues(downstreamID, downstreamDCID).Set(float64(node.Correcteds))
		m.DownstreamUncorrectables.WithLabelValues(downstreamID, downstreamDCID).Set(float64(node.Uncorrectables))
	}

	for id, node := range data.us {
		upstreamID := strconv.Itoa(id + 1)
		upstreamUCID := strconv.FormatInt(node.UCID, 10)

		m.UpstreamFreq.WithLabelValues(upstreamID, upstreamUCID).Set(node.Freq)
		m.UpstreamPower.WithLabelValues(upstreamID, upstreamUCID).Set(node.Power)
		m.UpstreamSymbolRate.WithLabelValues(upstreamID, upstreamUCID).Set(float64(node.SymbolRate))
		m.UpstreamModulation.WithLabelValues(upstreamID, upstreamUCID).Set(float64(node.Modulation))
	}

	return nil
}

func (m *ModemCollector) describe(ch chan<- *prometheus.Desc) {
	for _, metric := range m.collectorList() {
		metric.Describe(ch)
	}
}
