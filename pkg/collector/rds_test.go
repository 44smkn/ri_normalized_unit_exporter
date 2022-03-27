package collector_test

import (
	"testing"
	"time"

	"github.com/44smkn/aws_ri_exporter/pkg/collector"
)

func TestGetRemainingDays(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		startTime           time.Time
		reservationDuration time.Duration
		now                 time.Time
		want                float64
	}{
		{
			name:                "a year",
			startTime:           parseTimeRFC3339(t, "2021-10-02T15:04:05Z"),
			reservationDuration: 31536000 * time.Second,
			now:                 parseTimeRFC3339(t, "2022-04-02T16:04:05Z"),
			want:                183,
		},
		{
			name:                "three year",
			startTime:           parseTimeRFC3339(t, "2022-01-07T15:04:05Z"),
			reservationDuration: 94608000 * time.Second,
			now:                 parseTimeRFC3339(t, "2023-07-02T14:04:05Z"),
			want:                555,
		},
		{
			name:                "just after reservation",
			startTime:           parseTimeRFC3339(t, "2022-04-02T15:04:05Z"),
			reservationDuration: 31536000 * time.Second,
			now:                 parseTimeRFC3339(t, "2022-04-03T15:06:05Z"),
			want:                364,
		},
		{
			name:                "just before expiration",
			startTime:           parseTimeRFC3339(t, "2022-04-02T15:04:05Z"),
			reservationDuration: 31536000 * time.Second,
			now:                 parseTimeRFC3339(t, "2023-04-02T14:56:05Z"),
			want:                1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := collector.GetRemainingDays(tt.startTime, tt.reservationDuration, tt.now)
			if got != tt.want {
				t.Errorf("want: %v got: %v", tt.want, got)
			}
		})
	}
}

func parseTimeRFC3339(t *testing.T, value string) time.Time {
	t.Helper()
	tm, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Errorf("failed to parse time: %v", err)
	}
	return tm
}
