package aws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	"github.com/stretchr/testify/assert"
	"github.com/truvami/decoder/pkg/decoder"
	"go.uber.org/zap"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		Payload     string
		CaptureTime time.Time
		Expected    Position
	}{
		{
			Payload: "05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e",
			Expected: Position{
				Latitude:  47.35438919067383,
				Longitude: 8.55547046661377,
				Altitude:  aws.Float64(486.05999755859375),
				Timestamp: aws.Time(time.Date(2025, time.June, 19, 22, 19, 20, 652675294, time.UTC)),
				Accuracy:  aws.Float64(33.6),
				Buffered:  true,
			},
		},
		{
			Payload: "83ab812e9de68d0cc1006b9008acb2ab60b8c4dc4322d6f091c65dc11545946d8bd29879f0067bfbeee22fbef19db3cc0d",
			Expected: Position{
				Latitude:  47.350059509277344,
				Longitude: 8.561149597167969,
				Altitude:  aws.Float64(471),
				Timestamp: aws.Time(time.Date(2025, time.June, 19, 22, 31, 20, 652675294, time.UTC)),
				Accuracy:  aws.Float64(22.4),
				Buffered:  true,
			},
		},
	}

	logger := zap.NewExample()
	defer func() {
		_ = logger.Sync() // Flushes buffer, if any
	}()

	c, err := NewAwsPositionEstimateClient(context.TODO(), logger)
	assert.NoError(t, err, "expected no error during client creation")

	for _, test := range tests {
		t.Run(test.Payload, func(t *testing.T) {
			result, err := c.Solve(context.TODO(), test.Payload)
			assert.NoError(t, err, "expected no error during Solve")
			assert.NotNil(t, result, "expected result to be non-nil")

			// TODO: The timestamp in the expected result is not exact, so we cannot assert equality directly.
			// Instead, we can check if the timestamp is within a reasonable range.
			//
			// timeDiff := result.Timestamp.Sub(*test.Expected.Timestamp)
			// assert.LessOrEqual(t, timeDiff.Hours(), 1.0, "timestamp should be within 1h of expected")

			gnssFeature, ok := result.Data.(decoder.UplinkFeatureGNSS)
			assert.True(t, ok, "result should implement UplinkFeatureGNSS")

			bufferedFeature, ok := result.Data.(decoder.UplinkFeatureBuffered)
			assert.True(t, ok, "result should implement UplinkFeatureBuffered")

			// The assertions have been split to ensure each field is checked separately since the timestamp is not exact
			assert.Equal(t, test.Expected.Latitude, gnssFeature.GetLatitude(), "latitude does not match expected value")
			assert.Equal(t, test.Expected.Longitude, gnssFeature.GetLongitude(), "longitude does not match expected value")
			assert.Equal(t, *test.Expected.Altitude, gnssFeature.GetAltitude(), "altitude does not match expected value")
			assert.Equal(t, *test.Expected.Accuracy, *gnssFeature.GetAccuracy(), "accuracy does not match expected value")
			assert.Equal(t, test.Expected.Buffered, bufferedFeature.IsBuffered(), "buffered status does not match expected value")
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

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload         string
		allowNoFeatures bool
	}{
		{
			payload: "05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e",
		},
	}

	logger := zap.NewExample()
	defer func() {
		_ = logger.Sync() // Flushes buffer, if any
	}()

	s, err := NewAwsPositionEstimateClient(context.TODO(), logger)
	if err != nil {
		t.Fatalf("error creating AWS position estimate client: %s", err)
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPayload%v", test.payload), func(t *testing.T) {
			decodedPayload, err := s.Solve(context.TODO(), test.payload)
			if err != nil {
				t.Fatalf("error %s", err)
			}

			if len(decodedPayload.GetFeatures()) == 0 && !test.allowNoFeatures {
				t.Error("expected features, got none")
			}

			if decodedPayload.Is(decoder.FeatureTimestamp) {
				timestamp, ok := decodedPayload.Data.(decoder.UplinkFeatureTimestamp)
				if !ok {
					t.Fatalf("expected UplinkFeatureTimestamp, got %T", decodedPayload)
				}
				if timestamp.GetTimestamp() == nil {
					t.Fatalf("expected non nil timestamp")
				}
			}

			if decodedPayload.Is(decoder.FeatureGNSS) {
				gnss, ok := decodedPayload.Data.(decoder.UplinkFeatureGNSS)
				if !ok {
					t.Fatalf("expected UplinkFeatureGNSS, got %T", decodedPayload)
				}
				if gnss.GetLatitude() == 0 {
					t.Fatalf("expected non zero latitude")
				}
				if gnss.GetLongitude() == 0 {
					t.Fatalf("expected non zero longitude")
				}
				if gnss.GetAltitude() == 0 {
					t.Fatalf("expected non zero altitude")
				}
				// call function to check if it panics
				gnss.GetAltitude()
				gnss.GetPDOP()
				gnss.GetSatellites()
				gnss.GetTTF()
			}
			if decodedPayload.Is(decoder.FeatureBuffered) {
				buffered, ok := decodedPayload.Data.(decoder.UplinkFeatureBuffered)
				if !ok {
					t.Fatalf("expected UplinkFeatureBuffered, got %T", decodedPayload)
				}
				// call function to check if it panics
				buffered.GetBufferLevel()
				buffered.IsBuffered()
			}
		})
	}
}

type mockAwsPositionEstimateClient struct {
	GeoJsonResponse []byte
}

func (m *mockAwsPositionEstimateClient) GetPositionEstimate(ctx context.Context, params *iotwireless.GetPositionEstimateInput, optFns ...func(*iotwireless.Options)) (*iotwireless.GetPositionEstimateOutput, error) {
	return &iotwireless.GetPositionEstimateOutput{
		GeoJsonPayload: m.GeoJsonResponse,
	}, nil
}

func TestSolveWithMock(t *testing.T) {
	logger := zap.NewExample()
	defer func() {
		_ = logger.Sync() // Flushes buffer, if any
	}()

	tests := []struct {
		payload       string
		geoJson       []byte
		expectedError error
	}{
		{
			payload: "05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e",
			geoJson: []byte(`
				{
					"coordinates": [
						8.55547046661377,
						47.35438919067383,
						486.05999755859375
					],
					"type": "Point",
					"properties": {
						"horizontalAccuracy": 33.6,
						"horizontalConfidenceLevel": -1,
						"timestamp": "2025-06-20T21:31:38.146674492Z"
					}
				}
			`),
		},
		{
			payload: "05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e",
			geoJson: []byte(`
				{
					"coordinates": [
						8.55547046661377
					],
					"type": "Point",
					"properties": {
						"horizontalAccuracy": 33.6,
						"horizontalConfidenceLevel": -1,
						"timestamp": "2025-06-20T21:31:38.146674492Z"
					}
				}
			`),
			expectedError: ErrInvalidGeoJSONCoordinates,
		},
		{
			payload:       "05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e",
			geoJson:       []byte(``),
			expectedError: ErrFailedToUnmarshalGeoJSON,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSolveWithMockPayload%v", test.payload), func(t *testing.T) {
			// Create a mock client with the provided GeoJSON response
			mockClient := &mockAwsPositionEstimateClient{
				GeoJsonResponse: test.geoJson,
			}
			c := &PositionEstimateClient{
				client: mockClient,
				logger: logger,
			}

			result, err := c.Solve(context.TODO(), test.payload)

			if test.expectedError != nil {
				assert.Error(t, err, "expected an error during Solve with mock client")
				assert.Equal(t, test.expectedError, err, "error does not match expected error")
				return
			}
			assert.NoError(t, err, "expected no error during Solve with mock client")
			assert.NotNil(t, result, "expected result to be non-nil")

			gnssFeature, ok := result.Data.(decoder.UplinkFeatureGNSS)
			assert.True(t, ok, "result should implement UplinkFeatureGNSS")
			assert.Equal(t, 47.35438919067383, gnssFeature.GetLatitude(), "latitude does not match expected value")
			assert.Equal(t, 8.55547046661377, gnssFeature.GetLongitude(), "longitude does not match expected value")
			assert.Equal(t, 486.05999755859375, gnssFeature.GetAltitude(), "altitude does not match expected value")
			assert.Equal(t, 33.6, *gnssFeature.GetAccuracy(), "accuracy does not match expected value")
		})
	}
}

func TestPositionGetAltitude(t *testing.T) {
	tests := []struct {
		name     string
		altitude *float64
		want     float64
	}{
		{
			name:     "Altitude is set",
			altitude: aws.Float64(123.45),
			want:     123.45,
		},
		{
			name:     "Altitude is nil",
			altitude: nil,
			want:     0.0,
		},
		{
			name:     "Altitude is zero",
			altitude: aws.Float64(0.0),
			want:     0.0,
		},
		{
			name:     "Altitude is negative",
			altitude: aws.Float64(-50.5),
			want:     -50.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Position{
				Altitude: tt.altitude,
			}
			got := p.GetAltitude()
			assert.Equal(t, tt.want, got, "GetAltitude() should return the correct value")
		})
	}
}
