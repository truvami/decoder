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

	// Legacy v1 solver for backward compatibility (kept for existing tests and ports)
	solver         solver.SolverV1
	fallbackSolver solver.SolverV1

	// Preferred v2 solver (used for GNSS NAV grouping ports 192/193/194/195/199/210/211 when available)
	v2Solver         solver.SolverV2
	fallbackV2Solver solver.SolverV2
}

func NewTagXLv1Decoder(ctx context.Context, solver solver.SolverV1, logger *zap.Logger, options ...Option) decoder.Decoder {
	if solver == nil {
		logger.Panic("solver cannot be nil", zap.String("decoder", "TagXLv1Decoder"))
	}

	tagXLv1Decoder := &TagXLv1Decoder{
		logger:         logger,
		solver:         solver,
		fallbackSolver: nil,
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

func WithFallbackSolver(fallbackSolver solver.SolverV1) Option {
	return func(t *TagXLv1Decoder) {
		t.fallbackSolver = fallbackSolver
	}
}

// WithSolverV2 sets the v2 solver which accepts explicit options (DevEUI, counter, port, optional timestamp/moving).
func WithSolverV2(v2 solver.SolverV2) Option {
	return func(t *TagXLv1Decoder) {
		t.v2Solver = v2
	}
}

// WithFallbackSolverV2 sets the fallback v2 solver used when the primary v2 solver fails.
func WithFallbackSolverV2(fallback solver.SolverV2) Option {
	return func(t *TagXLv1Decoder) {
		t.fallbackV2Solver = fallback
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
			Features:   []decoder.Feature{decoder.FeatureTimestamp},
		}, nil
	case 151:
		if len(payload) < 1 {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var payloadType = payload[0]
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
				{Name: "DataRate", Tag: 0x4e, Optional: true, Transform: func(v any) any {
					if b, ok := v.([]byte); ok && len(b) > 0 {
						return DataRateFromUint8(b[0])
					}
					return nil
				}},
			},
			TargetType: reflect.TypeOf(Port151Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDataRate},
		}, nil
	case 152:
		if len(payload) < 1 {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[0]
		switch version {
		case Port152Version1:
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
				Features:   []decoder.Feature{decoder.FeatureRotationState, decoder.FeatureTimestamp},
			}, nil
		case Port152Version2:
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
				Features:   []decoder.Feature{decoder.FeatureSequenceNumber, decoder.FeatureRotationState, decoder.FeatureTimestamp},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 197:
		if len(payload) < 1 {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[0]
		switch version {
		case Port197Version1:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "Moving", Start: 0, Length: 1, Transform: alwaysFalse},
					{Name: "Mac1", Start: 1, Length: 6, Hex: true},
					{Name: "Mac2", Start: 7, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 19, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 25, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port197Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving},
			}, nil
		case Port197Version2:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "Moving", Start: 0, Length: 1, Transform: alwaysFalse},
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
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 198:
		if len(payload) < 1 {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[0]
		switch version {
		case Port198Version1:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "Moving", Start: 0, Length: 1, Transform: alwaysTrue},
					{Name: "Mac1", Start: 1, Length: 6, Hex: true},
					{Name: "Mac2", Start: 7, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 19, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 25, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port198Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving},
			}, nil
		case Port198Version2:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "Moving", Start: 0, Length: 1, Transform: alwaysTrue},
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
				TargetType: reflect.TypeOf(Port198Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 200:
		if len(payload) < 5 {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[4]
		switch version {
		case Port200Version1:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysFalse},
					{Name: "Mac1", Start: 5, Length: 6, Hex: true},
					{Name: "Mac2", Start: 11, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 17, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 23, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 29, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port200Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp, decoder.FeatureBuffered},
			}, nil
		case Port200Version2:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysFalse},
					{Name: "Rssi1", Start: 5, Length: 1},
					{Name: "Mac1", Start: 6, Length: 6, Hex: true},
					{Name: "Rssi2", Start: 12, Length: 1, Optional: true},
					{Name: "Mac2", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi3", Start: 19, Length: 1, Optional: true},
					{Name: "Mac3", Start: 20, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi4", Start: 26, Length: 1, Optional: true},
					{Name: "Mac4", Start: 27, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi5", Start: 33, Length: 1, Optional: true},
					{Name: "Mac5", Start: 34, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port200Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp, decoder.FeatureBuffered},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 201:
		if len(payload) < 5 {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[4]
		switch version {
		case Port201Version1:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysTrue},
					{Name: "Mac1", Start: 5, Length: 6, Hex: true},
					{Name: "Mac2", Start: 11, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 17, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 23, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 29, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port201Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp, decoder.FeatureBuffered},
			}, nil
		case Port201Version2:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysTrue},
					{Name: "Rssi1", Start: 5, Length: 1},
					{Name: "Mac1", Start: 6, Length: 6, Hex: true},
					{Name: "Rssi2", Start: 12, Length: 1, Optional: true},
					{Name: "Mac2", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi3", Start: 19, Length: 1, Optional: true},
					{Name: "Mac3", Start: 20, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi4", Start: 26, Length: 1, Optional: true},
					{Name: "Mac4", Start: 27, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi5", Start: 33, Length: 1, Optional: true},
					{Name: "Mac5", Start: 34, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port201Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp, decoder.FeatureBuffered},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 212:
		if len(payload) < Port212HeaderLength {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[Port212VersionIndex]
		switch version {
		case Port212Version1:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysFalse},
					{Name: "Mac1", Start: 5, Length: 6, Hex: true},
					{Name: "Mac2", Start: 11, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 17, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 23, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 29, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port212Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp},
			}, nil
		case Port212Version2:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysFalse},
					{Name: "Rssi1", Start: 5, Length: 1},
					{Name: "Mac1", Start: 6, Length: 6, Hex: true},
					{Name: "Rssi2", Start: 12, Length: 1, Optional: true},
					{Name: "Mac2", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi3", Start: 19, Length: 1, Optional: true},
					{Name: "Mac3", Start: 20, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi4", Start: 26, Length: 1, Optional: true},
					{Name: "Mac4", Start: 27, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi5", Start: 33, Length: 1, Optional: true},
					{Name: "Mac5", Start: 34, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port212Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	case 213:
		if len(payload) < Port213HeaderLength {
			return common.PayloadConfig{}, common.ErrPayloadTooShort
		}
		var version = payload[Port213VersionIndex]
		switch version {
		case Port213Version1:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysTrue},
					{Name: "Mac1", Start: 5, Length: 6, Hex: true},
					{Name: "Mac2", Start: 11, Length: 6, Optional: true, Hex: true},
					{Name: "Mac3", Start: 17, Length: 6, Optional: true, Hex: true},
					{Name: "Mac4", Start: 23, Length: 6, Optional: true, Hex: true},
					{Name: "Mac5", Start: 29, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port213Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp},
			}, nil
		case Port213Version2:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
					{Name: "Version", Start: 4, Length: 1},
					{Name: "Moving", Start: 4, Length: 1, Transform: alwaysTrue},
					{Name: "Rssi1", Start: 5, Length: 1},
					{Name: "Mac1", Start: 6, Length: 6, Hex: true},
					{Name: "Rssi2", Start: 12, Length: 1, Optional: true},
					{Name: "Mac2", Start: 13, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi3", Start: 19, Length: 1, Optional: true},
					{Name: "Mac3", Start: 20, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi4", Start: 26, Length: 1, Optional: true},
					{Name: "Mac4", Start: 27, Length: 6, Optional: true, Hex: true},
					{Name: "Rssi5", Start: 33, Length: 1, Optional: true},
					{Name: "Mac5", Start: 34, Length: 6, Optional: true, Hex: true},
				},
				TargetType: reflect.TypeOf(Port213Payload{}),
				Features:   []decoder.Feature{decoder.FeatureWiFi, decoder.FeatureMoving, decoder.FeatureTimestamp},
			}, nil
		default:
			return common.PayloadConfig{}, fmt.Errorf("%w: version %v for port %d not supported", common.ErrPortNotSupported, version, port)
		}
	}
	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

/*
GNSS solver routing and semantics:
- Ports 192/193/194/195/199/210/211 are GNSS NAV grouping ports. When a v2 solver is configured, we prefer it.
- Movement semantics by port:
  - 192: steady (Moving=false)
  - 193: moving (Moving=true)
  - 194: steady (Moving=false), timestamped payload (first 4 bytes UNIX seconds) is stripped before solving
  - 195: moving (Moving=true), timestamped payload (first 4 bytes UNIX seconds) is stripped before solving
  - 199: unspecified; Moving and Timestamp left nil unless future protocol specifies otherwise
  - 210: steady (Moving=false), timestamped payload (first 4 bytes UNIX seconds), rotation-triggered
  - 211: moving (Moving=true), timestamped payload (first 4 bytes UNIX seconds), rotation-triggered

- When no v2 solver is provided:
  - Ports 194/195/210/211 are not supported (they require timestamp stripping and explicit options).
  - Ports 192/193/199 fall back to the legacy v1 solver for backward compatibility.
*/
func (t TagXLv1Decoder) Decode(ctx context.Context, data string, port uint8) (*decoder.DecodedUplink, error) {
	switch port {
	// GNSS NAV grouping ports now use the v2 solver when available.
	case 192, 193, 194, 195, 199, 210, 211:
		if t.v2Solver != nil {
			devEui, _ := ctx.Value(decoder.DEVEUI_CONTEXT_KEY).(string)
			fcnt, _ := ctx.Value(decoder.FCNT_CONTEXT_KEY).(int)
			var movingPtr *bool
			switch port {
			case 192, 194, 210:
				mv := false
				movingPtr = &mv
			case 193, 195, 211:
				mv := true
				movingPtr = &mv
			default:
				// leave nil unless explicitly known
				movingPtr = nil
			}

			var tsPtr *time.Time
			payloadForSolve := data

			// For timestamped GNSS ports (194, 195, 210, 211), strip the leading 4-byte timestamp (big-endian)
			if port == 194 || port == 195 || port == 210 || port == 211 {
				bytes, err := common.HexStringToBytes(data)
				if err != nil {
					return nil, err
				}
				if len(bytes) < 5 {
					return nil, common.ErrPayloadTooShort
				}
				secs := common.BytesToUint32(bytes[0:4])
				ts := time.Unix(int64(secs), 0).UTC()
				tsPtr = &ts

				// Remove first 4 bytes (8 hex chars) from payload passed to solver
				if len(data) < 8 {
					return nil, common.ErrPayloadTooShort
				}
				payloadForSolve = data[8:]
			}

			opts := solver.SolverV2Options{
				DevEui:        devEui,
				UplinkCounter: uint16(fcnt),
				Port:          192, // always 192 for GNSS NAV grouping
				Timestamp:     tsPtr,
				Moving:        movingPtr,
			}

			uplink, err := t.v2Solver.Solve(ctx, payloadForSolve, opts)
			if err != nil {
				if t.fallbackV2Solver == nil {
					tagXlDecoderSolverFailedCounter.Inc()
					return nil, common.WrapError(err, common.ErrSolverFailed)
				}
				uplink, err = t.fallbackV2Solver.Solve(ctx, payloadForSolve, opts)
				if err != nil {
					tagXlDecoderSolverFailedCounter.Inc()
					return nil, common.WrapError(err, common.ErrSolverFailed)
				}
				tagXlDecoderSuccessfullyUsedFallbackSolverCounter.Inc()
			}
			return uplink, nil
		}

		// Fallback to legacy v1 solver when v2 is not provided (keeps backward compatibility).
		// Note: legacy path does not support 194/195/210/211 since v1 solver expects header as first byte.
		if port == 194 || port == 195 || port == 210 || port == 211 {
			return nil, fmt.Errorf("%w: port %v not supported without v2 solver", common.ErrPortNotSupported, port)
		}
		uplink, err := t.solver.Solve(ctx, data)
		if err != nil {
			if t.fallbackSolver == nil {
				tagXlDecoderSolverFailedCounter.Inc()
				return nil, common.WrapError(err, common.ErrSolverFailed)
			}

			uplink, err = t.fallbackSolver.Solve(ctx, data)
			if err != nil {
				tagXlDecoderSolverFailedCounter.Inc()
				return nil, common.WrapError(err, common.ErrSolverFailed)
			}
			tagXlDecoderSuccessfullyUsedFallbackSolverCounter.Inc()
		}
		return uplink, nil

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

func alwaysTrue(v any) any {
	return true
}

func alwaysFalse(v any) any {
	return false
}
