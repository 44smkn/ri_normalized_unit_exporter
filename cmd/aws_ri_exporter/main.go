package main

import (
	"context"
	"net/http"
	"os"

	"github.com/44smkn/aws_ri_exporter/pkg/aws"
	"github.com/44smkn/aws_ri_exporter/pkg/collector"
	"github.com/44smkn/aws_ri_exporter/pkg/normalizedunit"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	exitCodeOK int = 0

	// Errors start at 10
	exitCodeInitializeAWSConfigError = 10 + iota
	exitCodeStartServerError
)

var (
	// common configuration
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9981").Envar("WEB_LISTEN_ADDRESS").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").Envar("WEB_TELEMETRY_PATH").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector("aws_ri_exporter"))
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("aws_ri_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting aws_ri_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	// setup aws client and register exporter
	http.Handle(*metricsPath, promhttp.Handler())

	ctx := context.TODO()
	cloud, err := aws.NewCloud(ctx)
	if err != nil {
		level.Error(logger).Log("failed to initialize aws config")
		return exitCodeInitializeAWSConfigError
	}
	converter := normalizedunit.NewConverter()
	c := collector.NewRINormalizedUnitCollector(cloud, converter, logger)
	prometheus.MustRegister(c)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`<html>
<head><title>AWS Reserved Instance Exporter</title></head>
<body>
<h1>AWS Reserved Instance Normalized Unit Per Hour Exporter</h1>
<p><a href='` + *metricsPath + `'>Metrics</a></p>
</body>
</html>`))
	})

	level.Info(logger).Log("msg", "Listening on address", "address", listenAddress)
	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		return exitCodeStartServerError
	}
	return exitCodeOK
}
