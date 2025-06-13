package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless/types"
	geojson "github.com/paulmach/go.geojson"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/truvami/decoder/pkg/decoder"
	"go.uber.org/zap"
)

var (
	awsPostionEstimatesTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_total",
		Help: "The total number of processed position estimate requests",
	})
	awsPostionEstimatesErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_errors_total",
		Help: "The total number of errors encountered while processing position estimate requests",
	})
	awsPostionEstimatesDurationHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "truvami_aws_position_estimates_duration_seconds",
		Help:    "The duration of position estimate requests in seconds",
		Buckets: []float64{0.1, 0.2, 0.3, 0.5, 1, 2, 5, 10, 30, 60},
	})
	awsPostionEstimatesSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_success_total",
		Help: "The total number of successful position estimate requests",
	})
	awsPostionEstimatesFailureCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_failure_total",
		Help: "The total number of failed position estimate requests",
	})
)

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
func Solve(logger *zap.Logger, payload string, captureTime time.Time) (*Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	awsPostionEstimatesTotalCounter.Inc()

	logger.Debug("Starting position estimate request",
		zap.String("payload", payload),
		zap.Time("captureTime", captureTime),
	)

	// Load AWS config with context (respects timeout)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		awsPostionEstimatesErrorsCounter.Inc()
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := iotwireless.NewFromConfig(cfg)
	input := &iotwireless.GetPositionEstimateInput{
		Gnss: &types.Gnss{
			Payload:     aws.String(payload),
			CaptureTime: aws.Float32(float32(captureTime.Unix())),

			// Optional improvements:
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
	output, err := client.GetPositionEstimate(ctx, input)
	if err != nil {
		awsPostionEstimatesFailureCounter.Inc()
		return nil, fmt.Errorf("failed to get position estimate: %w", err)
	}

	logger.Debug("Position estimate received",
		zap.String("payload", payload),
		zap.ByteString("geoJson", output.GeoJsonPayload),
		zap.Any("metadata", output.ResultMetadata),
	)

	position, err := geojson.UnmarshalFeature(output.GeoJsonPayload)
	if err != nil {
		awsPostionEstimatesErrorsCounter.Inc()
		return nil, fmt.Errorf("failed to unmarshal GeoJSON payload: %w", err)
	}
	if position == nil {
		awsPostionEstimatesErrorsCounter.Inc()
		return nil, fmt.Errorf("received nil position from AWS IoT Wireless")
	}
	if position.Type != string(geojson.GeometryPoint) {
		awsPostionEstimatesErrorsCounter.Inc()
		return nil, fmt.Errorf("expected GeoJSON Feature, got %s", position.Type)
	}
	if len(position.Geometry.Point) < 2 {
		awsPostionEstimatesErrorsCounter.Inc()
		return nil, fmt.Errorf("invalid GeoJSON point: %v", position.Geometry.Point)
	}

	// Log the position estimate success
	awsPostionEstimatesSuccessCounter.Inc()
	awsPostionEstimatesDurationHistogram.Observe(time.Since(start).Seconds())

	return &Position{
		Latitude:  position.Geometry.Point[1],
		Longitude: position.Geometry.Point[0],
		Timestamp: &captureTime,
		Altitude:  nil, // Altitude is optional, set to nil if not provided
	}, nil
}

type Position struct {
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Altitude  *float64   `json:"altitude"` // Optional altitude
	Timestamp *time.Time `json:"timestamp"`
}

var _ decoder.UplinkFeatureBase = &Position{}
var _ decoder.UplinkFeatureGNSS = &Position{}

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
	return nil
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
