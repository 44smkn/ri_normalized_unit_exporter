package aws

import (
	"context"

	"github.com/44smkn/ri_normalized_unit_exporter/pkg/aws/services"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

type Cloud interface {
	// RDS provides API to AWS RDS
	RDS() services.RDS
	Region() string
}

func NewCloud(ctx context.Context) (Cloud, error) {
	cfg, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &defaultCloud{
		rds:    services.NewRDS(cfg),
		region: cfg.Region,
	}, nil
}

var _ Cloud = &defaultCloud{}

type defaultCloud struct {
	rds    services.RDS
	region string
}

func (c *defaultCloud) RDS() services.RDS {
	return c.rds
}

func (c *defaultCloud) Region() string {
	return c.region
}
