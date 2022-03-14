package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

//go:generate mockgen -destination=rds_mocks.go -package=services github.com/44smkn/ri_normalized_unit_exporter/pkg/aws/services RDS
type RDS interface {
	// Wrapper to DescribeDBInstances, which aggregates paged results into list.
	DescribeDBInstancesAsList(context.Context, *rds.DescribeDBInstancesInput) ([]rdstypes.DBInstance, error)

	// Wrapper to DescribeReservedDBInstances, which aggregates paged results into list.
	DescribeReservedDBInstancesAsList(context.Context, *rds.DescribeReservedDBInstancesInput) ([]rdstypes.ReservedDBInstance, error)
}

func NewRDS(cfg aws.Config, optFns ...func(*rds.Options)) RDS {
	return &defaultRDS{
		Client: rds.NewFromConfig(cfg, optFns...),
	}
}

type defaultRDS struct {
	*rds.Client
}

func (c *defaultRDS) DescribeDBInstancesAsList(ctx context.Context, params *rds.DescribeDBInstancesInput) ([]rdstypes.DBInstance, error) {
	var instances []rdstypes.DBInstance
	paginator := rds.NewDescribeDBInstancesPaginator(c.Client, params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("paging process of DescribeDBInstances API was failed: %w", err)
		}
		instances = append(instances, output.DBInstances...)
	}
	return instances, nil
}

func (c *defaultRDS) DescribeReservedDBInstancesAsList(ctx context.Context, params *rds.DescribeReservedDBInstancesInput) ([]rdstypes.ReservedDBInstance, error) {
	var instances []rdstypes.ReservedDBInstance
	paginator := rds.NewDescribeReservedDBInstancesPaginator(c.Client, params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("paging process of DescribeReservedDBInstances API was failed: %w", err)
		}
		instances = append(instances, output.ReservedDBInstances...)
	}
	return instances, nil
}
