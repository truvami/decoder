package tagxl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
	"github.com/truvami/decoder/pkg/loracloud"
)

type TagXLv1Decoder struct {
	loracloudMiddleware loracloud.LoracloudMiddleware
}

func NewTagXLv1Decoder(loracloudMiddleware loracloud.LoracloudMiddleware) decoder.Decoder {
	return TagXLv1Decoder{loracloudMiddleware}
}

// https://docs.truvami.com/docs/payloads/tag-xl
func (t TagXLv1Decoder) getConfig(port int16) (decoder.PayloadConfig, error) {
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
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
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
	return decoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}

func (t TagXLv1Decoder) Decode(data string, port int16, devEui string) (interface{}, interface{}, error) {
	switch port {
	case 192, 197, 199:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    uint8(port),
			Payload: data,
		})
		return decodedData, nil, err
	default:
		config, err := t.getConfig(port)
		if err != nil {
			return nil, nil, err
		}

		decodedData, err := helpers.Parse(data, config)
		if err != nil {
			return nil, nil, err
		}

		return decodedData, nil, nil
	}
}
