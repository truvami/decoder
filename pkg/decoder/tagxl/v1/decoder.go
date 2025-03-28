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
		// return decoder.PayloadConfig{
		// 	Fields: []decoder.FieldConfig{
		// 		{Name: "DeviceFlags", Start: 2, Length: 1},
		// 		// {Name: "AssetTrackingIntervals", Start: 3, Length: 4, Transform: func(v interface{}) interface{} {
		// 		// 	return []uint16{uint16(v.(int) >> 16), uint16(v.(int) & 0xFFFF)}
		// 		// }},
		// 		// {Name: "AccelerationSensor", Start: 7, Length: 4, Transform: func(v interface{}) interface{} {
		// 		// 	return []uint16{uint16(v.(int) >> 16), uint16(v.(int) & 0xFFFF)}
		// 		// }},
		// 		{Name: "HeartbeatInterval", Start: 11, Length: 1, Optional: true},
		// 		{Name: "AdvertisementFwuInterval", Start: 12, Length: 1},
		// 		{Name: "Battery", Start: 13, Length: 2},
		// 		{Name: "FirmwareHash", Start: 15, Length: 4},
		// 	},
		// 	TargetType: reflect.TypeOf(Port151Payload{}),
		// }, nil
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
