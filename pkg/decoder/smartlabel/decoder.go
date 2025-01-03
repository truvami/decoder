package smartlabel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
	"github.com/truvami/decoder/pkg/loracloud"
)

type SmartLabelv1Decoder struct {
	loracloudMiddleware loracloud.LoracloudMiddleware
}

func NewSmartLabelv1Decoder(loracloudMiddleware loracloud.LoracloudMiddleware) decoder.Decoder {
	return SmartLabelv1Decoder{loracloudMiddleware}
}

func getPort11PayloadType(data string) (string, error) {
	if len(data) < 2 {
		return "", fmt.Errorf("data length is less than 2")
	}
	firstByte := strings.ToUpper(data[0:2])
	if firstByte == "0E" || firstByte == "11" { // 14 || 17
		return "configuration", nil
	}
	return "heartbeat", nil
}

// https://docs.truvami.com/docs/payloads/smartlabel
func (t SmartLabelv1Decoder) getConfig(port int16, data string) (decoder.PayloadConfig, error) {
	switch port {
	case 11:
		// Check first byte length to determine message type
		payloadType, err := getPort11PayloadType(data)
		if err != nil {
			return decoder.PayloadConfig{}, err
		}
		switch payloadType {
			case "configuration":
				return decoder.PayloadConfig{
					Fields: []decoder.FieldConfig{
							{Name: "Flags", Start: 2, Length: 1},
						{Name: "GNSSEnabled", Start: 2, Length: 1, Transform: func(v interface{}) interface{} {
							return uint8((v.(int) >> 1) & 0x1)
						}},
						{Name: "WiFiEnabled", Start: 2, Length: 1, Transform: func(v interface{}) interface{} {
							return uint8((v.(int) >> 2) & 0x1)
						}},
						{Name: "AccEnabled", Start: 2, Length: 1, Transform: func(v interface{}) interface{} {
							return uint8((v.(int) >> 3) & 0x1)
						}},
						{Name: "StaticSF", Start: 2, Length: 1, Transform: func(v interface{}) interface{} {
							return fmt.Sprintf("SF%d", 9) // TODO: Hardcoded to SF9 for now
						}},
						{Name: "SteadyIntervalS", Start: 3, Length: 2},
						{Name: "MovingIntervalS", Start: 5, Length: 2},
						{Name: "HeartbeatIntervalH", Start: 7, Length: 1},
						{Name: "LEDBlinkIntervalS", Start: 8, Length: 2},
						{Name: "AccThresholdMS", Start: 10, Length: 2},
						{Name: "AccDelayMS", Start: 12, Length: 2},
						{Name: "GitHash", Start: 14, Length: 4, Optional: true, Transform: func(v interface{}) interface{} {
							if v == nil {
								return nil
							}
							return fmt.Sprintf("%08x", v.(int))
						}},
					},
					TargetType: reflect.TypeOf(Port11ConfigurationPayload{}),
				}, nil
			case "heartbeat":
				return decoder.PayloadConfig{
					Fields: []decoder.FieldConfig{
						{Name: "Battery", Start: 2, Length: 2, Transform: func(v interface{}) interface{} {
							return float64(v.(int)) / 1000
						}},
						{Name: "Temperature", Start: 4, Length: 2, Transform: func(v interface{}) interface{} {
							return float64(v.(int)) / 100
						}},
						{Name: "RH", Start: 6, Length: 1, Transform: func(v interface{}) interface{} {
							return float64(v.(int)) / 2
						}},
						{Name: "GNSSScansCount", Start: 7, Length: 2},
						{Name: "WiFiScansCount", Start: 9, Length: 2},
					},
					TargetType: reflect.TypeOf(Port11HeartbeatPayload{}),
				}, nil
		}
		return decoder.PayloadConfig{}, fmt.Errorf("invalid data length")
	default:
		return decoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
	}
}

func (t SmartLabelv1Decoder) Decode(data string, port int16, devEui string) (interface{}, interface{}, error) {
	switch port {
	case 192, 197:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    uint8(port),
			Payload: data,
		})
		return decodedData, nil, err
	default:
		config, err := t.getConfig(port, data)
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
