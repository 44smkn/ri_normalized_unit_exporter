package aws

import (
	"github.com/44smkn/aws_ri_exporter/pkg/aws/services"
)

type mockCloud struct {
	rds *services.MockRDS
}

func NewMockCloud(rds *services.MockRDS) *mockCloud {
	return &mockCloud{
		rds: rds,
	}
}

func (c *mockCloud) RDS() services.RDS {
	return c.rds
}

func (c *mockCloud) Region() string {
	return "ap-northeast-1"
}
