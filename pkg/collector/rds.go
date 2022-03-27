package collector

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/44smkn/aws_ri_exporter/pkg/aws"
	"github.com/44smkn/aws_ri_exporter/pkg/aws/services"
	nu "github.com/44smkn/aws_ri_exporter/pkg/normalizedunit"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	rdsCollectorSubsystem = "rds"
)

func init() {
	registerCollector("rds", true, NewRDSCollector)
}

type rdsCollector struct {
	runningInstance          *prometheus.Desc
	activeReservation        *prometheus.Desc
	reservationRemainingDays *prometheus.Desc

	logger      log.Logger
	rds         services.RDS
	nuConverter nu.Converter
	region      string
}

func NewRDSCollector(aws aws.Cloud, nuConverter nu.Converter, logger log.Logger) Collector {
	c := &rdsCollector{
		runningInstance: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, rdsCollectorSubsystem, "running_instance_normalized_unit"),
			"Normalized Units for each running RDS instance.",
			[]string{"region", "instance_class", "engine", "instance_id"}, nil,
		),
		activeReservation: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, rdsCollectorSubsystem, "active_reservation_normalized_unit"),
			"Normalized Units for each active reservation",
			[]string{"region", "instance_class", "engine", "reservation_id"}, nil,
		),
		reservationRemainingDays: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, rdsCollectorSubsystem, "reservation_remaining_days"),
			"Remaining days of each active reservation",
			[]string{"region", "reservation_id"}, nil,
		),
		logger:      logger,
		rds:         aws.RDS(),
		nuConverter: nuConverter,
		region:      aws.Region(),
	}
	return c
}

func (c *rdsCollector) Update(ctx context.Context, ch chan<- prometheus.Metric) error {
	if err := c.updateReservation(ctx, ch); err != nil {
		return err
	}
	return c.updateRunningInstance(ctx, ch)
}

func (c *rdsCollector) updateRunningInstance(context context.Context, ch chan<- prometheus.Metric) error {
	params := &rds.DescribeDBInstancesInput{}
	instances, err := c.rds.DescribeDBInstancesAsList(context, params)
	if err != nil {
		return fmt.Errorf("To execute DescribeDBInstancesAsList() was failed: %w", err)
	}
	for _, instance := range instances {
		value, err := c.nuConverter.Convert(*instance.DBInstanceClass, 1)
		if err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			c.runningInstance,
			prometheus.GaugeValue,
			value,
			c.region,
			*instance.DBInstanceClass,
			*instance.Engine,
			*instance.DBInstanceIdentifier,
		)
	}
	return nil
}

func (c *rdsCollector) updateReservation(context context.Context, ch chan<- prometheus.Metric) error {
	params := &rds.DescribeReservedDBInstancesInput{}
	reservations, err := c.rds.DescribeReservedDBInstancesAsList(context, params)
	if err != nil {
		return fmt.Errorf("To execute DescribeDBInstancesAsList() was failed: %w", err)
	}
	for _, reservation := range reservations {
		if *reservation.State != "active" {
			continue
		}

		// reservation_remaining_days
		remainingDays := getRemainingDays(*reservation.StartTime, time.Duration(reservation.Duration)*time.Second, time.Now())
		ch <- prometheus.MustNewConstMetric(
			c.reservationRemainingDays,
			prometheus.GaugeValue,
			remainingDays,
			c.region,
			*reservation.ReservedDBInstanceId,
		)

		// active_reservation_normalized_unit
		normalizedUnit, err := c.nuConverter.Convert(*reservation.DBInstanceClass, float64(reservation.DBInstanceCount))
		if err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			c.activeReservation,
			prometheus.GaugeValue,
			normalizedUnit,
			c.region,
			*reservation.DBInstanceClass,
			*reservation.ProductDescription,
			*reservation.ReservedDBInstanceId,
		)
	}
	return nil
}

func getRemainingDays(startTime time.Time, reservationDuration time.Duration, now time.Time) float64 {
	expiredTime := startTime.Add(reservationDuration)
	remainingHour := expiredTime.Sub(now).Hours()
	return math.Ceil(remainingHour / 24)
}
