package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless/types"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"go.uber.org/zap"
)

type PositionEstimateClient struct {
	client iotwirelessClient
	logger *zap.Logger
}

type iotwirelessClient interface {
	GetPositionEstimate(ctx context.Context, params *iotwireless.GetPositionEstimateInput, optFns ...func(*iotwireless.Options)) (*iotwireless.GetPositionEstimateOutput, error)
}

var _ solver.SolverV1 = &PositionEstimateClient{}

func NewAwsPositionEstimateClient(ctx context.Context, logger *zap.Logger) (*PositionEstimateClient, error) {
	// Load AWS config with context (respects timeout)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		awsPositionEstimatesErrorsCounter.Inc()
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &PositionEstimateClient{
		client: iotwireless.NewFromConfig(cfg),
		logger: logger,
	}, nil
}

// Solve sends a GNSS payload to AWS IoT Wireless to obtain a position estimate.
// It logs the request and response using the provided zap.Logger. The function
// creates a context with a 5-second timeout for AWS operations, loads the AWS
// configuration, and calls the GetPositionEstimate API. If successful, it logs
// the resulting GeoJSON payload and metadata. Returns an error if any step fails.
//
// Parameters:
//   - logger:     A zap.Logger for structured logging.
//   - payload:    The GNSS payload as a string.
//   - captureTime: The time the GNSS data was captured.
//
// Returns:
//   - error:      An error if the AWS config could not be loaded or the position
//     estimate request fails; otherwise, nil.
func (c PositionEstimateClient) Solve(ctx context.Context, payload string) (*decoder.DecodedUplink, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	start := time.Now()
	awsPositionEstimatesTotalCounter.Inc()

	c.logger.Debug("Starting position estimate request",
		zap.String("payload", payload),
	)

	// remove first 2 characters from the payload
	if len(payload) > 2 {
		payload = payload[2:]
	}

	input := &iotwireless.GetPositionEstimateInput{
		Gnss: &types.Gnss{
			Payload: aws.String(payload),
			// in seconds GPS time (GPST)
			// GPS Time (GPST) is a continuous time scale (no leap seconds) defined by the GPS Control segment on the basis of a set of atomic clocks at the Monitor Stations and onboard the satellites. It starts at 0h UTC (midnight) of January 5th to 6th 1980 (6.d0). At that epoch, the difference TAI−UTC was 19 seconds, thence GPS−UTC=n − 19s. GPS time is synchronised with the UTC(USNO) at 1 microsecond level (modulo one second), but actually is kept within 25 ns.[
			// CaptureTime: aws.Float32(getGPSTime(captureTime)),
			// AssistAltitude: aws.Float32(50.0),
			// AssistPosition: []float32{37.7749, -122.4194},
			// Use2DSolver:    aws.Bool(true),
		},
	}

	// The position information of the resource, displayed as a JSON payload. The
	// payload is of type blob and uses the [GeoJSON]format, which a format that's used to
	// encode geographic data structures. A sample payload contains the timestamp
	// information, the WGS84 coordinates of the location, and the accuracy and
	// confidence level. For more information and examples, see [Resolve device location (console)].
	//
	// [Resolve device location (console)]: https://docs.aws.amazon.com/iot/latest/developerguide/location-resolve-console.html
	// [GeoJSON]: https://geojson.org/
	output, err := c.client.GetPositionEstimate(ctx, input)
	if err != nil {
		awsPositionEstimatesFailureCounter.Inc()
		return nil, fmt.Errorf("failed to get position estimate: %w", err)
	}

	c.logger.Debug("Position estimate received",
		zap.String("payload", payload),
		zap.ByteString("geoJson", output.GeoJsonPayload),
		zap.Any("metadata", output.ResultMetadata),
	)

	var position *GeoJsonResponse
	err = json.Unmarshal(output.GeoJsonPayload, &position)
	if err != nil {
		awsPositionEstimatesErrorsCounter.Inc()
		return nil, ErrFailedToUnmarshalGeoJSON
	}
	if len(position.Coordinates) < 2 {
		awsPositionEstimatesErrorsCounter.Inc()
		return nil, ErrInvalidGeoJSONCoordinates
	}

	// Log the position estimate success
	awsPositionEstimatesSuccessCounter.Inc()
	awsPositionEstimatesDurationHistogram.Observe(time.Since(start).Seconds())

	var altitude *float64
	if len(position.Coordinates) > 2 {
		altitude = position.Coordinates[2]
	}

	buffered := false
	// add 10 seconds buffer to the timestamp
	if position.Properties.Timestamp != nil && position.Properties.Timestamp.Before(time.Now().Add(-1*time.Minute)) {
		buffered = true
	}

	return decoder.NewDecodedUplink([]decoder.Feature{
		decoder.FeatureGNSS,
		decoder.FeatureTimestamp,
		decoder.FeatureBuffered,
	}, Position{
		Latitude:  *position.Coordinates[1],
		Longitude: *position.Coordinates[0],
		Altitude:  altitude,
		Timestamp: position.Properties.Timestamp,
		Accuracy:  position.Properties.HorizontalAccuracy,
		Buffered:  buffered,
	}), nil
}

func getGPSTime(captureTime time.Time) float32 {
	// GPS time starts at 0h UTC on January 5th, 1980
	gpsEpoch := time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)
	return float32(captureTime.Sub(gpsEpoch).Seconds())
}

type GeoJsonResponse struct {
	Coordinates []*float64 `json:"coordinates"`
	Type        string     `json:"type"`
	Properties  struct {
		HorizontalAccuracy        *float64   `json:"horizontalAccuracy"`
		HorizontalConfidenceLevel *int       `json:"horizontalConfidenceLevel"`
		Timestamp                 *time.Time `json:"timestamp"`
	} `json:"properties"`
}

type Position struct {
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Altitude  *float64   `json:"altitude"` // Optional altitude
	Timestamp *time.Time `json:"timestamp"`
	Accuracy  *float64   `json:"accuracy"` // Optional accuracy
	Buffered  bool       `json:"buffered"` // Indicates if the position is buffered
}

var _ decoder.UplinkFeatureTimestamp = &Position{}
var _ decoder.UplinkFeatureGNSS = &Position{}
var _ decoder.UplinkFeatureBuffered = &Position{}

func (p Position) GetTimestamp() *time.Time {
	return p.Timestamp
}

func (p Position) GetLatitude() float64 {
	return p.Latitude
}

func (p Position) GetLongitude() float64 {
	return p.Longitude
}

func (p Position) GetAltitude() float64 {
	if p.Altitude == nil {
		return 0.0 // Return 0 if altitude is not set
	}
	return *p.Altitude
}

func (p Position) GetAccuracy() *float64 {
	return p.Accuracy
}

func (p Position) GetTTF() *time.Duration {
	return nil
}

func (p Position) GetPDOP() *float64 {
	return nil
}

func (p Position) GetSatellites() *uint8 {
	return nil
}

func (p Position) GetBufferLevel() *uint16 {
	return nil
}

func (p Position) IsBuffered() bool {
	return p.Buffered
}
