package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	"github.com/aws/aws-sdk-go-v2/service/iotwireless/types"
	"go.uber.org/zap"
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
func Solve(logger *zap.Logger, payload string, captureTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Debug("Starting position estimate request",
		zap.String("payload", payload),
		zap.Time("captureTime", captureTime),
	)

	// Load AWS config with context (respects timeout)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
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

	output, err := client.GetPositionEstimate(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to get position estimate: %w", err)
	}

	logger.Debug("Position estimate received",
		zap.String("payload", payload),
		zap.ByteString("geoJson", output.GeoJsonPayload),
		zap.Any("metadata", output.ResultMetadata),
	)
	return nil
}
