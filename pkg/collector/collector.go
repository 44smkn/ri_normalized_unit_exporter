package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/44smkn/aws_ri_exporter/pkg/aws"
	nu "github.com/44smkn/aws_ri_exporter/pkg/normalizedunit"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "aws_ri"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "duration_seconds"),
		"eaws_ri_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "success"),
		"eaws_ri_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)

	factories      = make(map[string]func(aws.Cloud, nu.Converter, log.Logger) Collector)
	collectorState = make(map[string]*bool)
)

func registerCollector(collector string, isDefaultEnabled bool,
	factory func(aws aws.Cloud, converter nu.Converter, logger log.Logger) Collector) {
	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %v).", collector, isDefaultEnabled)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	enabled := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	collectorState[collector] = enabled

	factories[collector] = factory
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(context.Context, chan<- prometheus.Metric) error
}

type riNormalizedUnitCollector struct {
	Collectors map[string]Collector
	logger     log.Logger
}

func NewRINormalizedUnitCollector(aws aws.Cloud, normalizedUnitConverter nu.Converter, logger log.Logger) *riNormalizedUnitCollector {
	collectors := make(map[string]Collector)
	for key, enabled := range collectorState {
		if !*enabled {
			continue
		}
		collector := factories[key](aws, normalizedUnitConverter, log.With(logger, "collector", key))
		collectors[key] = collector
	}
	c := &riNormalizedUnitCollector{
		Collectors: collectors,
		logger:     logger,
	}
	return c
}

// Describe implements the prometheus.Collector interface
func (c *riNormalizedUnitCollector) Describe(ch chan<- *prometheus.Desc) {}

// Collect implements the prometheus.Collector interface.
func (r *riNormalizedUnitCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	ctx := context.TODO()
	wg.Add(len(r.Collectors))
	for name, c := range r.Collectors {
		go func(name string, c Collector) {
			execute(ctx, name, c, ch, r.logger)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(ctx context.Context, name string, c Collector, ch chan<- prometheus.Metric, logger log.Logger) {
	begin := time.Now()
	err := c.Update(ctx, ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		level.Error(logger).Log("msg", "collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		success = 0
	} else {
		level.Debug(logger).Log("msg", "collector succeeded", "name", name, "duration_seconds", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}
