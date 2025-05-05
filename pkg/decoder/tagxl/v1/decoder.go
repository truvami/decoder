package tagxl

import (
	"encoding/binary"
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
func (t TagXLv1Decoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
	case 152:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "OldRotationState", Start: 2, Length: 1, Transform: func(v any) any {
					// get bit 0-3 and return
					return uint8((v.(int) >> 4) & 0xF)
				}},
				{Name: "NewRotationState", Start: 2, Length: 1, Transform: func(v any) any {
					// get bit 4-7 and return
					return uint8(v.(int) & 0xF)
				}},
				{Name: "Timestamp", Start: 3, Length: 4},
				{Name: "NumberOfRotations", Start: 7, Length: 2, Transform: func(v any) any {
					return float64(v.(int)) / 10
				}},
				{Name: "ElapsedSeconds", Start: 9, Length: 4},
			},
			TargetType: reflect.TypeOf(Port152Payload{}),
		}, nil
	case 197:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Tag", Start: 0, Length: 1},
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
	}
	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t TagXLv1Decoder) Decode(data string, port uint8, devEui string) (*decoder.DecodedUplink, error) {
	switch port {
	case 151:
		payload, err := common.HexStringToBytes(data)
		if err != nil {
			return nil, err
		}

		var payloadLen uint8 = uint8(len(payload))
		if payloadLen < 3 {
			return nil, fmt.Errorf("%w: port %d minimum %d", common.ErrPayloadTooShort, port, 3)
		}

		var payloadType byte = payload[0]
		if payloadType != 0x4c {
			return nil, fmt.Errorf("%w: port %d tag %x", common.ErrPortNotSupported, port, payloadType)
		}

		var dataLen uint8 = payload[1]
		if payloadLen-2 != dataLen {
			return nil, fmt.Errorf("%w: port %d expected %d received %d", common.ErrInvalidPayloadLength, port, payloadLen-2, dataLen)
		}

		var decodedData = Port151Payload{}
		var features = []decoder.Feature{}
		var index uint8 = 3
		for index+1 < payloadLen {
			var tag byte = payload[index]
			var len uint8 = payload[index+1]
			switch tag {
			case 0x45:
				var value = payload[index+2 : index+2+len]
				var battery = float32(binary.BigEndian.Uint16(value)) / 1000
				decodedData.BatteryVoltage = battery
				features = append(features, decoder.FeatureBattery)
			}
			index += 2
			index += len
		}
		return decoder.NewDecodedUplink(features, decodedData, nil), nil
	case 192:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    port,
			Payload: data,
			FCount:  t.fCount,
		})
		return decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS}, decodedData, nil), err
	case 199:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    port,
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
