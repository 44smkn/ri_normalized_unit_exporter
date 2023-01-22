package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/44smkn/aws_ri_exporter/pkg/aws"
	"github.com/44smkn/aws_ri_exporter/pkg/aws/services"
	"github.com/44smkn/aws_ri_exporter/pkg/collector"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestHanler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRDS := services.NewMockRDS(ctrl)
	mockRDS.EXPECT().DescribeDBInstancesAsList(gomock.Any(), gomock.Any()).Return([]rdstypes.DBInstance{
		{
			DBInstanceClass:      awssdk.String("db.r6g.2xlarge"),
			Engine:               awssdk.String("mysql"),
			DBInstanceIdentifier: awssdk.String("test-1"),
		},
	}, nil)
	mockRDS.EXPECT().DescribeReservedDBInstancesAsList(gomock.Any(), gomock.Any()).Return([]rdstypes.ReservedDBInstance{
		{
			State:                awssdk.String("active"),
			StartTime:            awssdk.Time(time.Now().Add(-24 * time.Hour * 90)),
			DBInstanceClass:      awssdk.String("db.r6g.4xlarge"),
			ProductDescription:   awssdk.String("postgres"),
			ReservedDBInstanceId: awssdk.String("myreservedinstance"),
		},
	}, nil)

	// If not parsed, the default value set by kingpin will not be set.
	// If we don't do this, the collector will not be enable.
	_, err := kingpin.CommandLine.Parse([]string{})
	if err != nil {
		t.Errorf("failed to parse command line: %v", err)
		return
	}

	// Clear the collector set by init()
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry
	collector.EnableScrapeMetrics = false

	fc := aws.NewMockCloud(mockRDS)
	logger := promlog.New(&promlog.Config{})
	s := httptest.NewServer(initPromHandler(fc, false, logger))
	defer s.Close()

	resp, err := http.Get(s.URL + "/metrics")
	if err != nil {
		t.Errorf("http get err should be nil: %v", err)
		return
	}
	defer resp.Body.Close()

	golden := filepath.Join("../../testdata/fixtures", "all_metrics_exist_one_by_one"+".golden")
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Errorf("failed to read files: %v", err)
		return
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
		return
	}

	if string(got) != string(want) {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(got), string(want), false)
		t.Errorf("diff: \n%v", dmp.DiffPrettyText(diffs))
	}
}
