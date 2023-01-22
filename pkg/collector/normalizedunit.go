package collector

import (
	"fmt"
	"strings"
)

func convertToNormalizedUnits(instanceClass string, instanceCount float64) (float64, error) {
	for classSuffix, normalizedUnit := range classSuffixNormalizedUnitFactorTable {
		if strings.HasSuffix(instanceClass, classSuffix) {
			return normalizedUnit * instanceCount, nil
		}
	}
	return 0, fmt.Errorf("failed to retrive normalized unit, please check db instance class name")
}

var classSuffixNormalizedUnitFactorTable = map[string]float64{
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/apply_ri.html#ri-instance-size-flexibility
	".nano":      0.25,
	".micro":     0.5,
	".small":     1,
	".medium":    2,
	".large":     4,
	".xlarge":    8,
	".2xlarge":   16,
	".3xlarge":   24,
	".4xlarge":   32,
	".6xlarge":   48,
	".8xlarge":   64,
	".9xlarge":   72,
	".10xlarge":  80,
	".12xlarge":  96,
	".16xlarge":  128,
	".18xlarge":  144,
	".24xlarge":  192,
	".32xlarge":  256,
	".56xlarge":  448,
	".112xlarge": 896,
}
