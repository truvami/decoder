package aws

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		Payload     string
		CaptureTime time.Time
		Expected    Position
	}{
		{
			Payload:     "05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e",
			CaptureTime: time.Date(2025, time.June, 18, 14, 40, 00, 0, time.UTC),
			Expected: Position{
				Latitude:  47.35438919067383,
				Longitude: 8.55547046661377,
				Altitude:  aws.Float64(486.05999755859375),
				Timestamp: aws.Time(time.Date(2025, time.June, 19, 22, 19, 20, 652675294, time.UTC)),
				Accuracy:  aws.Float64(33.6),
				Buffered:  false,
			},
		},
		{
			Payload:     "83ab812e9de68d0cc1006b9008acb2ab60b8c4dc4322d6f091c65dc11545946d8bd29879f0067bfbeee22fbef19db3cc0d",
			CaptureTime: time.Date(2025, time.June, 18, 14, 30, 00, 0, time.UTC),
			Expected: Position{
				Latitude:  47.350059509277344,
				Longitude: 8.561149597167969,
				Altitude:  aws.Float64(471),
				Timestamp: aws.Time(time.Date(2025, time.June, 19, 22, 31, 20, 652675294, time.UTC)),
				Accuracy:  aws.Float64(22.4),
				Buffered:  false,
			},
		},
	}

	logger := zap.NewExample()
	defer logger.Sync() // flushes buffer, if any

	for _, test := range tests {
		t.Run(test.Payload, func(t *testing.T) {
			result, err := Solve(logger, test.Payload, test.CaptureTime)
			assert.NoError(t, err, "expected no error during Solve")

			timeDiff := result.Timestamp.Sub(*test.Expected.Timestamp)
			assert.LessOrEqual(t, timeDiff.Hours(), 1.0, "timestamp should be within 1h of expected")

			// The assertions have been split to ensure each field is checked separately since the timestamp is not exact
			assert.Equal(t, test.Expected.Latitude, result.Latitude, "latitude does not match expected value")
			assert.Equal(t, test.Expected.Longitude, result.Longitude, "longitude does not match expected value")
			assert.Equal(t, *test.Expected.Altitude, *result.Altitude, "altitude does not match expected value")
			assert.Equal(t, *test.Expected.Accuracy, *result.Accuracy, "accuracy does not match expected value")
			assert.Equal(t, test.Expected.Buffered, result.Buffered, "buffered status does not match expected value")
		})
	}
}

func TestGetGPSTime(t *testing.T) {
	tests := []struct {
		name        string
		captureTime time.Time
		want        float32
	}{
		{
			name:        "At GPS epoch",
			captureTime: time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC),
			want:        0,
		},
		{
			name:        "One second after GPS epoch",
			captureTime: time.Date(1980, time.January, 6, 0, 0, 1, 0, time.UTC),
			want:        1,
		},
		{
			name:        "One day after GPS epoch",
			captureTime: time.Date(1980, time.January, 7, 0, 0, 0, 0, time.UTC),
			want:        86400,
		},
		{
			name:        "Forty years after GPS epoch",
			captureTime: time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC),
			want:        float32((40*365 + 10) * 86400), // 10 leap days between 1980 and 2020
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getGPSTime(tt.captureTime)
			assert.InDelta(t, tt.want, got, 1, "getGPSTime(%v) = %v, want %v", tt.captureTime, got, tt.want)
		})
	}
}
