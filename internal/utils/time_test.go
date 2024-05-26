package utils_test

import (
	"kmtt/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func parseTime(t *testing.T, timeStr string) time.Time {
	t.Helper()
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	require.NoError(t, err)
	return parsedTime
}

func TestRoundToNearestInterval(t *testing.T) {
	tests := []struct {
		time     string
		interval time.Duration
		expected string
	}{
		{
			time:     "2024-05-23 10:23:45",
			interval: 5 * time.Minute,
			expected: "2024-05-23 10:25:00",
		},
		{
			time:     "2024-05-23 10:24:30",
			interval: 10 * time.Minute,
			expected: "2024-05-23 10:20:00",
		},
		{
			time:     "2024-05-23 10:26:00",
			interval: 15 * time.Minute,
			expected: "2024-05-23 10:30:00",
		},
	}

	for _, test := range tests {
		inputTime := parseTime(t, test.time)
		expectedTime := parseTime(t, test.expected)
		result := utils.RoundToNearestInterval(inputTime, test.interval)
		require.Truef(t, result.Equal(expectedTime), "got %v; want %v", result, expectedTime)
	}
}

func TestGetRangeBetween(t *testing.T) {
	tests := []struct {
		from     string
		to       string
		interval time.Duration
		expected []string
	}{
		{
			from:     "2024-05-23 10:00:00",
			to:       "2024-05-23 10:30:00",
			interval: 10 * time.Minute,
			expected: []string{
				"2024-05-23 10:00:00",
				"2024-05-23 10:10:00",
				"2024-05-23 10:20:00",
			},
		},
		{
			from:     "2024-05-23 10:00:00",
			to:       "2024-05-23 10:30:00",
			interval: 15 * time.Minute,
			expected: []string{
				"2024-05-23 10:00:00",
				"2024-05-23 10:15:00",
			},
		},
	}

	for _, test := range tests {
		fromTime := parseTime(t, test.from)
		toTime := parseTime(t, test.to)
		expectedTimes := make([]time.Time, len(test.expected))
		for i, expectedStr := range test.expected {
			expectedTimes[i] = parseTime(t, expectedStr)
		}

		result := utils.GetRangeBetween(fromTime, toTime, test.interval)
		require.Equal(t, len(result), len(expectedTimes))
		for i := range result {
			require.Truef(t, result[i].Equal(expectedTimes[i]), "got %v; want %v", result[i], expectedTimes[i])
		}
	}
}
