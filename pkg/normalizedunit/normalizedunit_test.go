package normalizedunit_test

import (
	"testing"

	"github.com/44smkn/ri_normalized_unit_exporter/pkg/normalizedunit"
)

func Test_defaultConverter_Convert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		instanceClass string
		instanceCount float64
		want          float64
	}{
		{
			name:          "running instance",
			instanceClass: "db.t3.micro",
			instanceCount: 1,
			want:          0.5,
		},
		{
			name:          "active reservation",
			instanceClass: "db.m5.xlarge",
			instanceCount: 3,
			want:          24,
		},
	}

	converter := normalizedunit.NewConverter()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := converter.Convert(tt.instanceClass, tt.instanceCount)
			if err != nil {
				t.Errorf("failed to convert: %v", err)
			}
			if tt.want != got {
				t.Errorf("want: %v got: %v", tt.want, got)
			}
		})
	}
}
