package tagxl

import (
	"fmt"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/loracloud"
)

type Option func(*TagXLv1Decoder)

type TagXLv1Decoder struct {
	loracloudMiddleware loracloud.LoracloudMiddleware
	autoPadding         bool
	skipValidation      bool
	fCount              uint32
}

func NewTagXLv1Decoder(loracloudMiddleware loracloud.LoracloudMiddleware, options ...Option) decoder.Decoder {
	tagXLv1Decoder := &TagXLv1Decoder{
		loracloudMiddleware: loracloudMiddleware,
	}

	for _, option := range options {
		option(tagXLv1Decoder)
	}

	return tagXLv1Decoder
}

func WithAutoPadding(autoPadding bool) Option {
	return func(t *TagXLv1Decoder) {
		t.autoPadding = autoPadding
	}
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *TagXLv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// WithFCount sets the frame counter for the decoder.
// This is required for the loracloud middleware.
func WithFCount(fCount uint32) Option {
	return func(t *TagXLv1Decoder) {
		t.fCount = fCount
	}
}

// https://docs.truvami.com/docs/payloads/tag-xl
func (t TagXLv1Decoder) getConfig(port uint8, payload []byte) (common.PayloadConfig, error) {
	switch port {
	case 150:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Timestamp", Start: 5, Length: 4, Optional: false, Transform: func(v any) any {
					return time.Unix(int64(v.(int)), 0).UTC()
				}},
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
				{Name: "GnssEnabled", Tag: 0x40, Optional: true, Transform: func(v any) any {
					// bit 1: GNSS_ENABLE
					return (v.(int) & 0x02) != 0
				}},
				{Name: "WiFiEnabled", Tag: 0x40, Optional: true, Transform: func(v any) any {
					// bit 2: WIFI_ENABLE
					return (v.(int) & 0x04) != 0
				}},
				{Name: "AccelerometerEnabled", Tag: 0x40, Optional: true, Transform: func(v any) any {
					// bit 3: ACCELERATION_ENABLE
					return (v.(int) & 0x08) != 0
				}},
				{Name: "LocalizationIntervalWhileMoving", Tag: 0x41, Optional: true, Transform: func(v any) any {
					// data 0: MOVING_INTERVAL
					return uint16((v.(int) >> 16) & 0xffff)
				}},
				{Name: "LocalizationIntervalWhileSteady", Tag: 0x41, Optional: true, Transform: func(v any) any {
					// data 1: STEADY_INTERVAL
					return uint16(v.(int) & 0xffff)
				}},
				{Name: "AccelerometerWakeupThreshold", Tag: 0x42, Optional: true, Transform: func(v any) any {
					// data 0: WAKEUP_THRESHOLD
					return uint16((v.(int) >> 16) & 0xffff)
				}},
				{Name: "AccelerometerDelay", Tag: 0x42, Optional: true, Transform: func(v any) any {
					// data 1: WAKEUP_DELAY
					return uint16(v.(int) & 0xffff)
				}},
				{Name: "HeartbeatInterval", Tag: 0x43, Optional: true},
				{Name: "AdvertisementFirmwareUpgradeInterval", Tag: 0x44, Optional: true},
				{Name: "Battery", Tag: 0x45, Optional: true, Feature: []decoder.Feature{decoder.FeatureBattery}, Transform: func(v any) any {
					return float32(v.(int)) / 1000
				}},
				{Name: "FirmwareHash", Tag: 0x46, Optional: true, Hex: true},
				{Name: "ResetCount", Tag: 0x49, Optional: true},
				{Name: "ResetCause", Tag: 0x4a, Optional: true},
				{Name: "GnssScans", Tag: 0x4b, Optional: true, Transform: func(v any) any {
					return uint16((v.(int) >> 16) & 0xffff)
				}},
				{Name: "WifiScans", Tag: 0x4b, Optional: true, Transform: func(v any) any {
					return uint16(v.(int) & 0xffff)
				}},
			},
			Features:   []decoder.Feature{decoder.FeatureConfig},
			TargetType: reflect.TypeOf(Port151Payload{}),
		}, nil
	case 152:
		var version uint8 = payload[0]
		switch version {
		case 0x01:
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Version", Start: 0, Length: 1},
					{Name: "OldRotationState", Start: 2, Length: 1, Transform: func(v any) any {
						return uint8((v.(int) >> 4))
					}},
					{Name: "NewRotationState", Start: 2, Length: 1, Transform: func(v any) any {
						return uint8(v.(int) & 0x0f)
					}},
					{Name: "Timestamp", Start: 3, Length: 4},
					{Name: "NumberOfRotations", Start: 7, Length: 2, Transform: func(v any) any {
						return float64(v.(int)) / 10
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
						return uint8((v.(int) >> 4))
					}},
					{Name: "NewRotationState", Start: 3, Length: 1, Transform: func(v any) any {
						return uint8(v.(int) & 0x0f)
					}},
					{Name: "Timestamp", Start: 4, Length: 4},
					{Name: "NumberOfRotations", Start: 8, Length: 2, Transform: func(v any) any {
						return float64(v.(int)) / 10
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

func (t TagXLv1Decoder) Decode(data string, port uint8, devEui string) (*decoder.DecodedUplink, error) {
	switch port {
	case 192:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    port,
			Payload: data,
			FCount:  t.fCount,
		})
		return decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS}, decodedData), err
	case 199:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    port,
			Payload: data,
			FCount:  t.fCount,
		})
		return decoder.NewDecodedUplink([]decoder.Feature{}, decodedData), err
	default:
		bytes, err := common.HexStringToBytes(data)
		if err != nil {
			return nil, err
		}

		config, err := t.getConfig(port, bytes)
		if err != nil {
			return nil, err
		}

		if t.autoPadding {
			data = common.HexNullPad(&data, &config)
		}

		if !t.skipValidation {
			err := common.ValidateLength(&data, &config)
			if err != nil {
				return nil, err
			}
		}

		decodedData, err := common.Parse(data, &config)
		return decoder.NewDecodedUplink(config.Features, decodedData), err
	}
}
