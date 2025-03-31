package tagxl

import (
	"fmt"
	"reflect"

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
func (t TagXLv1Decoder) getConfig(port int16) (common.PayloadConfig, error) {
	switch port {
	case 151:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Ble", Start: 0, Length: 1, Transform: func(v interface{}) interface{} {
					return ((v.(int) >> 7) & 0x1) != 0
				}},
				{Name: "Gnss", Start: 0, Length: 1, Transform: func(v interface{}) interface{} {
					return ((v.(int) >> 6) & 0x1) != 0
				}},
				{Name: "Wifi", Start: 0, Length: 1, Transform: func(v interface{}) interface{} {
					return ((v.(int) >> 5) & 0x1) != 0
				}},
				{Name: "Acceleration", Start: 0, Length: 1, Transform: func(v interface{}) interface{} {
					return ((v.(int) >> 4) & 0x1) != 0
				}},
				{Name: "Rfu", Start: 0, Length: 1, Transform: func(v interface{}) interface{} {
					return uint8(v.(int) & 0xf)
				}},
				{Name: "MovingInterval", Start: 1, Length: 2},
				{Name: "SteadyInterval", Start: 3, Length: 2},
				{Name: "AccelerationThreshold", Start: 5, Length: 2},
				{Name: "AccelerationDelay", Start: 7, Length: 2},
				{Name: "HeartbeatInterval", Start: 9, Length: 1},
				{Name: "FwuAdvertisementInterval", Start: 10, Length: 1},
				{Name: "BatteryVoltage", Start: 11, Length: 2, Transform: func(v interface{}) interface{} {
					return float32(v.(int)) / 1000
				}},
				{Name: "FirmwareHash", Start: 13, Length: 4, Transform: func(v interface{}) interface{} {
					return fmt.Sprintf("%8x", v.(int))
				}},
				{Name: "ResetCount", Start: 17, Length: 2},
				{Name: "ResetCause", Start: 19, Length: 4},
				{Name: "GnssScans", Start: 23, Length: 2},
				{Name: "WifiScans", Start: 25, Length: 2},
			},
			TargetType: reflect.TypeOf(Port151Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBattery, decoder.FeatureConfig},
		}, nil
	case 152:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "OldRotationState", Start: 2, Length: 1, Transform: func(v interface{}) interface{} {
					// get bit 0-3 and return
					return uint8((v.(int) >> 4) & 0xF)
				}},
				{Name: "NewRotationState", Start: 2, Length: 1, Transform: func(v interface{}) interface{} {
					// get bit 4-7 and return
					return uint8(v.(int) & 0xF)
				}},
				{Name: "Timestamp", Start: 3, Length: 4},
				{Name: "NumberOfRotations", Start: 7, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 10
				}},
				{Name: "ElapsedSeconds", Start: 9, Length: 4},
			},
			TargetType: reflect.TypeOf(Port152Payload{}),
		}, nil

	}
	return common.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}

func (t TagXLv1Decoder) Decode(data string, port int16, devEui string) (*decoder.DecodedUplink, error) {
	switch port {
	case 192, 197, 199:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    uint8(port),
			Payload: data,
			FCount:  t.fCount,
		})
		return decoder.NewDecodedUplink([]decoder.Feature{}, decodedData, nil), err
	default:
		config, err := t.getConfig(port)
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
		return decoder.NewDecodedUplink(config.Features, decodedData, nil), err
	}
}
