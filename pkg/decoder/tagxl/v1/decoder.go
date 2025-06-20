package tagxl

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"go.uber.org/zap"
)

type Option func(*TagXLv1Decoder)

type TagXLv1Decoder struct {
	skipValidation bool
	logger         *zap.Logger

	solver solver.SolverV1

	// This will be deprecated in the future.
	fCount uint32
	devEui string
}

func NewTagXLv1Decoder(ctx context.Context, solver solver.SolverV1, logger *zap.Logger, options ...Option) decoder.Decoder {
	if solver == nil {
		logger.Fatal("solver cannot be nil", zap.String("decoder", "TagXLv1Decoder"))
	}

	tagXLv1Decoder := &TagXLv1Decoder{
		logger: logger,
		solver: solver,
	}

	for _, option := range options {
		option(tagXLv1Decoder)
	}

	return tagXLv1Decoder
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *TagXLv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// WithFCount sets the frame counter for the decoder.
// This is required for the loracloud middleware.
// Will be deprecated in the future.
func WithFCount(fCount uint32) Option {
	return func(t *TagXLv1Decoder) {
		t.fCount = fCount
	}
}

// WithDevEui sets the DevEUI for the decoder.
// This is required for the loracloud middleware.
// Will be deprecated in the future.
func WithDevEui(devEui string) Option {
	return func(t *TagXLv1Decoder) {
		t.devEui = devEui
	}
}

// https://docs.truvami.com/docs/payloads/tag-xl
func (t TagXLv1Decoder) getConfig(port uint8, payload []byte) (common.PayloadConfig, error) {
	switch port {
	case 150:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Timestamp", Start: 5, Length: 4, Transform: timestamp},
			},
			TargetType: reflect.TypeOf(Port150Payload{}),
			Features:   []decoder.Feature{},
		}, nil
	case 151:
		var payloadType byte = payload[0]
		if payloadType != 0x4c {
			return common.PayloadConfig{}, fmt.Errorf("%w: port %d tag %x", common.ErrPortNotSupported, port, payloadType)
		}
		return common.PayloadConfig{
			Tags: []common.TagConfig{
				{Name: "AccelerometerEnabled", Tag: 0x40, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 3) & 0x01) != 0
				}},
				{Name: "WifiEnabled", Tag: 0x40, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 2) & 0x01) != 0
				}},
				{Name: "GnssEnabled", Tag: 0x40, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 1) & 0x01) != 0
				}},
				{Name: "FirmwareUpgrade", Tag: 0x40, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return (v.([]byte)[0] & 0x01) != 0
				}},
				{Name: "LocalizationIntervalWhileMoving", Tag: 0x41, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return uint16((common.BytesToUint32(v.([]byte)) >> 16) & 0xffff)
				}},
				{Name: "LocalizationIntervalWhileSteady", Tag: 0x41, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return uint16(common.BytesToUint32(v.([]byte)) & 0xffff)
				}},
				{Name: "AccelerometerWakeupThreshold", Tag: 0x42, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return uint16((common.BytesToUint32(v.([]byte)) >> 16) & 0xffff)
				}},
				{Name: "AccelerometerDelay", Tag: 0x42, Optional: true, Feature: decoder.FeatureConfig, Transform: func(v any) any {
					return uint16(common.BytesToUint32(v.([]byte)) & 0xffff)
				}},
				{Name: "HeartbeatInterval", Tag: 0x43, Optional: true},
				{Name: "AdvertisementFirmwareUpgradeInterval", Tag: 0x44, Optional: true},
				{Name: "Battery", Tag: 0x45, Optional: true, Feature: decoder.FeatureBattery, Transform: func(v any) any {
					return float32(common.BytesToUint16(v.([]byte))) / 1000
				}},
				{Name: "FirmwareHash", Tag: 0x46, Optional: true, Feature: decoder.FeatureFirmwareVersion, Hex: true},
				{Name: "RotationInvert", Tag: 0x47, Optional: true, Transform: func(v any) any {
					return (v.([]byte)[0] & 0x01) != 0
				}},
				{Name: "RotationConfirmed", Tag: 0x47, Optional: true, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 1) & 0x01) != 0
				}},
				{Name: "ResetCount", Tag: 0x49, Optional: true},
				{Name: "ResetCause", Tag: 0x4a, Optional: true},
				{Name: "GnssScans", Tag: 0x4b, Optional: true, Transform: func(v any) any {
					return uint16((common.BytesToUint32(v.([]byte)) >> 16) & 0xffff)
				}},
				{Name: "WifiScans", Tag: 0x4b, Optional: true, Transform: func(v any) any {
					return uint16(common.BytesToUint32(v.([]byte)) & 0xffff)
				}},
			},
			TargetType: reflect.TypeOf(Port151Payload{}),
			Features:   []decoder.Feature{},
		}, nil
	case 152:
		var version uint8 = payload[0]
		switch version {
		case 0x01:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "OldRotationState", Start: 2, Length: 1, Transform: func(v any) any {
						return common.BytesToUint8(v.([]byte)) >> 4
					}},
					{Name: "NewRotationState", Start: 2, Length: 1, Transform: func(v any) any {
						return common.BytesToUint8(v.([]byte)) & 0x0f
					}},
					{Name: "Timestamp", Start: 3, Length: 4, Transform: timestamp},
					{Name: "NumberOfRotations", Start: 7, Length: 2, Transform: func(v any) any {
						return float64(common.BytesToUint16(v.([]byte))) / 10
					}},
					{Name: "ElapsedSeconds", Start: 9, Length: 4},
				},
				TargetType: reflect.TypeOf(Port152Payload{}),
				Features:   []decoder.Feature{decoder.FeatureRotationState},
			}, nil
		case 0x02:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "SequenceNumber", Start: 2, Length: 1},
					{Name: "OldRotationState", Start: 3, Length: 1, Transform: func(v any) any {
						return common.BytesToUint8(v.([]byte)) >> 4
					}},
					{Name: "NewRotationState", Start: 3, Length: 1, Transform: func(v any) any {
						return common.BytesToUint8(v.([]byte)) & 0x0f
					}},
					{Name: "Timestamp", Start: 4, Length: 4, Transform: timestamp},
					{Name: "NumberOfRotations", Start: 8, Length: 2, Transform: func(v any) any {
						return float64(common.BytesToUint16(v.([]byte))) / 10
					}},
					{Name: "ElapsedSeconds", Start: 10, Length: 4},
				},
				TargetType: reflect.TypeOf(Port152Payload{}),
				Features:   []decoder.Feature{decoder.FeatureRotationState, decoder.FeatureSequenceNumber},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 197:
		var version uint8 = payload[0]
		switch version {
		case 0x00:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Mac1", Start: 1, Length: 6, Hex: true},
					{Name: "Mac2", Start: 7, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 19, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 25, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port197Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi},
			}, nil
		case 0x01:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Rssi1", Start: 1, Length: 1},
					{Name: "Mac1", Start: 2, Length: 6, Hex: true},
					{Name: "Rssi2", Start: 8, Length: 1, Optional: true},
					{Name: "Mac2", Start: 9, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi3", Start: 15, Length: 1, Optional: true},
					{Name: "Mac3", Start: 16, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi4", Start: 22, Length: 1, Optional: true},
					{Name: "Mac4", Start: 23, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi5", Start: 29, Length: 1, Optional: true},
					{Name: "Mac5", Start: 30, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port197Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	}
	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t TagXLv1Decoder) Decode(data string, port uint8) (*decoder.DecodedUplink, error) {
	switch port {
	case 192, 199:
		return t.solver.Solve(data)
	default:
		bytes, err := common.HexStringToBytes(data)
		if err != nil {
			return nil, err
		}

		config, err := t.getConfig(port, bytes)
		if err != nil {
			return nil, err
		}

		if !t.skipValidation {
			err := common.ValidateLength(&data, &config)
			if err != nil {
				return nil, err
			}
		}

		decodedData, err := common.Decode(&data, &config)
		return decoder.NewDecodedUplink(config.Features, decodedData), err
	}
}

func timestamp(v any) any {
	return time.Unix(int64(common.BytesToUint32(v.([]byte))), 0).UTC()
}
