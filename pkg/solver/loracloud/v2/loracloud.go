package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	v1 "github.com/truvami/decoder/pkg/solver/loracloud"
	"go.uber.org/zap"
)

const (
	SemtechLoRaCloudBaseUrl  = v1.SemtechLoRaCloudBaseUrl
	TraxmateLoRaCloudBaseUrl = v1.TraxmateLoRaCloudBaseUrl
)

// LoracloudClient implements solver.SolverV2 without reading values from context.
// Optional fields Moving and Timestamp control features and implemented interfaces.
type LoracloudClient struct {
	accessToken       string
	logger            *zap.Logger
	BaseUrl           string
	bufferedThreshold time.Duration
}

var _ solver.SolverV2 = &LoracloudClient{}

// Options for configuring the v2 client
type LoracloudClientOptions func(*LoracloudClient)

func WithBaseUrl(baseUrl string) LoracloudClientOptions {
	return func(c *LoracloudClient) {
		c.BaseUrl = baseUrl
	}
}

func WithBufferedThreshold(threshold time.Duration) LoracloudClientOptions {
	return func(c *LoracloudClient) {
		c.bufferedThreshold = threshold
	}
}

// NewLoracloudClient creates a new v2 client with sane defaults.
// Defaults: BaseUrl=TraxmateLoRaCloudBaseUrl, bufferedThreshold=1m
func NewLoracloudClient(ctx context.Context, accessToken string, logger *zap.Logger, options ...LoracloudClientOptions) (LoracloudClient, error) {
	client := LoracloudClient{
		accessToken:       accessToken,
		logger:            logger,
		BaseUrl:           TraxmateLoRaCloudBaseUrl,
		bufferedThreshold: time.Minute,
	}
	for _, opt := range options {
		opt(&client)
	}
	// Warn for Semtech LoRaCloud shutdown (defensive)
	if client.BaseUrl == SemtechLoRaCloudBaseUrl && time.Now().After(time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)) {
		logger.Warn("LoRa Cloud is Sunsetting on 31.07.2025", zap.String("url", "https://www.semtech.com/loracloud-shutdown"))
	}
	return client, nil
}

func (l LoracloudClient) Solve(ctx context.Context, payload string, options solver.SolverV2Options) (*decoder.DecodedUplink, error) {
	start := time.Now()
	baseURLLabel := l.BaseUrl
	defer func() {
		loracloudV2RequestDurationSeconds.WithLabelValues(baseURLLabel).Observe(time.Since(start).Seconds())
	}()

	// Validate options (do NOT read from context)
	if err := l.validateOptions(payload, options); err != nil {
		loracloudV2RequestsTotal.WithLabelValues(baseURLLabel, "error").Inc()
		loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "invalid_options").Inc()
		return nil, common.WrapError(ErrInvalidOptions, err)
	}

	// Build timestamp to send to LoRaCloud (seconds, UTC)
	var ts *float64
	if options.Timestamp != nil {
		sec := float64(options.Timestamp.UTC().Unix())
		ts = &sec
	} else {
		// If no timestamp provided, always use current time for better API compatibility
		sec := float64(time.Now().UTC().Unix())
		ts = &sec
	}

	// Reuse v1 client for actual HTTP and response shaping, to keep behavior aligned
	v1Client, err := v1.NewLoracloudClient(ctx, l.accessToken, l.logger, v1.WithBaseUrl(l.BaseUrl))
	if err != nil {
		loracloudV2RequestsTotal.WithLabelValues(baseURLLabel, "error").Inc()
		loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "build_request").Inc()
		return nil, common.WrapError(ErrBuildRequest, err)
	}

	uplink := v1.UplinkMsg{
		MsgType:   "updf",
		FCount:    uint32(options.UplinkCounter),
		Port:      options.Port,
		Payload:   payload,
		Timestamp: ts,
	}

	resp, err := v1Client.DeliverUplinkMessage(options.DevEui, uplink)
	if err != nil {
		loracloudV2RequestsTotal.WithLabelValues(baseURLLabel, "error").Inc()
		switch {
		case errors.Is(err, v1.ErrUnexpectedStatusCode):
			loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "unexpected_status").Inc()
			return nil, common.WrapError(ErrUnexpectedStatus, err)
		case errors.Is(err, v1.ErrDecodingResponse):
			loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "decode_failed").Inc()
			return nil, common.WrapError(ErrDecodeFailed, err)
		case errors.Is(err, v1.ErrPositionResolutionIsEmpty):
			loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "position_invalid").Inc()
			return nil, common.WrapError(ErrPositionInvalid, err)
		default:
			loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "request_failed").Inc()
			return nil, common.WrapError(ErrRequestFailed, err)
		}
	}

	// Defensive validation of response
	if resp == nil {
		loracloudV2RequestsTotal.WithLabelValues(baseURLLabel, "error").Inc()
		loracloudV2ResponseInvalidTotal.WithLabelValues(baseURLLabel).Inc()
		return nil, common.WrapError(ErrResponseInvalid, fmt.Errorf("nil response"))
	}

	devEui := resp.Result.Deveui // v1 client already normalized and removed hyphens

	// Visibility counters similar to v1 (best effort)
	if resp.GetTimestamp() == nil {
		loracloudV2PositionInvalidTotal.WithLabelValues(devEui).Inc()
	}
	if !resp.HasValidCoordinates() {
		loracloudV2PositionInvalidTotal.WithLabelValues(devEui).Inc()
	}

	validPosition := resp.HasValidPositionResolution()

	// Build features based on inputs and buffered logic
	features := []decoder.Feature{}
	if validPosition {
		features = append(features, decoder.FeatureGNSS)
	} else {
		loracloudV2ErrorsTotal.WithLabelValues(baseURLLabel, "position_invalid").Inc()
		l.logger.Debug("position resolution invalid (no GNSS feature set)", zap.Any("uplinkResponse", resp))
	}

	withTimestamp := options.Timestamp != nil
	withMoving := options.Moving != nil
	buffered := false

	if withTimestamp {
		features = append(features, decoder.FeatureTimestamp)

		thresholdAgo := time.Now().Add(-1 * l.bufferedThreshold)
		if options.Timestamp.Before(thresholdAgo) {
			buffered = true
			features = append(features, decoder.FeatureBuffered)
			loracloudV2BufferedDetectedTotal.WithLabelValues(devEui, l.bufferedThreshold.String()).Inc()
		}
	}

	if withMoving {
		features = append(features, decoder.FeatureMoving)
	}

	// Build Data that implements only the requested feature interfaces
	var data any
	switch {
	case withTimestamp && withMoving && buffered:
		data = &dataTSMovingBuffered{
			dataTSMoving: dataTSMoving{
				dataBase: dataBase{resp: resp},
				ts:       options.Timestamp,
				moving:   *options.Moving,
			},
		}
	case withTimestamp && withMoving && !buffered:
		data = &dataTSMoving{
			dataBase: dataBase{resp: resp},
			ts:       options.Timestamp,
			moving:   *options.Moving,
		}
	case withTimestamp && !withMoving && buffered:
		data = &dataTSBuffered{
			dataTS: dataTS{
				dataBase: dataBase{resp: resp},
				ts:       options.Timestamp,
			},
		}
	case withTimestamp && !withMoving && !buffered:
		data = &dataTS{
			dataBase: dataBase{resp: resp},
			ts:       options.Timestamp,
		}
	case !withTimestamp && withMoving:
		data = &dataMoving{
			dataBase: dataBase{resp: resp},
			moving:   *options.Moving,
		}
	default:
		// No optional interfaces
		data = &dataBase{resp: resp}
	}

	loracloudV2RequestsTotal.WithLabelValues(baseURLLabel, "success").Inc()
	return decoder.NewDecodedUplink(features, data), nil
}

// validateOptions validates DevEui, payload and basic constraints.
func (l LoracloudClient) validateOptions(payload string, options solver.SolverV2Options) error {
	if len(options.DevEui) != 16 {
		return ErrInvalidDevEui
	}
	if _, err := hex.DecodeString(options.DevEui); err != nil {
		return ErrInvalidDevEui
	}
	// Port is uint8 already; just sanity check payload
	if payload == "" {
		return fmt.Errorf("payload empty")
	}
	return nil
}

// ---------- Data wrapper types to implement optional feature interfaces ----------

type dataBase struct {
	resp *v1.UplinkMsgResponse
}

// GNSS delegates
var _ decoder.UplinkFeatureGNSS = &dataBase{}

func (d dataBase) GetLatitude() float64   { return d.resp.GetLatitude() }
func (d dataBase) GetLongitude() float64  { return d.resp.GetLongitude() }
func (d dataBase) GetAltitude() float64   { return d.resp.GetAltitude() }
func (d dataBase) GetAccuracy() *float64  { return d.resp.GetAccuracy() }
func (d dataBase) GetTTF() *time.Duration { return d.resp.GetTTF() }
func (d dataBase) GetPDOP() *float64      { return d.resp.GetPDOP() }
func (d dataBase) GetSatellites() *uint8  { return d.resp.GetSatellites() }

// Timestamp only when provided
type dataTS struct {
	dataBase
	ts *time.Time
}

var _ decoder.UplinkFeatureTimestamp = &dataTS{}

func (d dataTS) GetTimestamp() *time.Time { return d.ts }

// Moving only when provided
type dataMoving struct {
	dataBase
	moving bool
}

var _ decoder.UplinkFeatureMoving = &dataMoving{}

func (d dataMoving) IsMoving() bool { return d.moving }

// Timestamp + Moving
type dataTSMoving struct {
	dataBase
	ts     *time.Time
	moving bool
}

var (
	_ decoder.UplinkFeatureTimestamp = &dataTSMoving{}
	_ decoder.UplinkFeatureMoving    = &dataTSMoving{}
)

func (d dataTSMoving) GetTimestamp() *time.Time { return d.ts }
func (d dataTSMoving) IsMoving() bool           { return d.moving }

// Buffered variants (only when timestamp is past threshold)
type dataTSBuffered struct {
	dataTS
}

var _ decoder.UplinkFeatureBuffered = &dataTSBuffered{}

func (d dataTSBuffered) IsBuffered() bool        { return true }
func (d dataTSBuffered) GetBufferLevel() *uint16 { return nil }

type dataTSMovingBuffered struct {
	dataTSMoving
}

var _ decoder.UplinkFeatureBuffered = &dataTSMovingBuffered{}

func (d dataTSMovingBuffered) IsBuffered() bool        { return true }
func (d dataTSMovingBuffered) GetBufferLevel() *uint16 { return nil }
