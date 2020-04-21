package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/syepes/ping_exporter/monitor"
	"github.com/syepes/ping_exporter/pkg/ping"
)

var (
	// icmpLabelNames = []string{"alias", "target", "ip", "ip_version"}
	icmpLabelNames = []string{"alias", "target"}
	rttDesc        = prometheus.NewDesc("ping_rtt_seconds", "ICMP Round trip time in seconds", append(icmpLabelNames, "type"), nil)
	lossDesc       = prometheus.NewDesc("ping_loss_percent", "Packet loss in percent", icmpLabelNames, nil)
	icmpProgDesc   = prometheus.NewDesc("ping_up", "ping_exporter version", nil, prometheus.Labels{"version": "xzy"})
	icmpMutex      = &sync.Mutex{}
)

// PingCollector prom
type PingCollector struct {
	Monitor *monitor.MonitorPing
	metrics map[string]*ping.PingReturn
}

// Describe prom
func (p *PingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rttDesc
	ch <- lossDesc
	ch <- icmpProgDesc
}

// Collect prom
func (p *PingCollector) Collect(ch chan<- prometheus.Metric) {
	icmpMutex.Lock()
	defer icmpMutex.Unlock()

	if m := p.Monitor.Export(); len(m) > 0 {
		p.metrics = m
	}

	if len(p.metrics) > 0 {
		ch <- prometheus.MustNewConstMetric(icmpProgDesc, prometheus.GaugeValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(icmpProgDesc, prometheus.GaugeValue, 0)
	}

	for target, metrics := range p.metrics {
		// fmt.Printf("target: %v\n", target)
		// fmt.Printf("metrics: %v\n", metrics)
		// l := strings.SplitN(target, " ", 2)
		l := []string{target, metrics.DestAddr}
		// fmt.Printf("L: %v\n", l)

		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.BestTime/1000), append(l, "best")...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.WrstTime/1000), append(l, "worst")...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.AvgTime/1000), append(l, "mean")...)
		ch <- prometheus.MustNewConstMetric(lossDesc, prometheus.GaugeValue, metrics.DropRate, l...)
	}
}