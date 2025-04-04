package smartlabel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/loracloud"
)

type Option func(*SmartLabelv1Decoder)

type SmartLabelv1Decoder struct {
	loracloudMiddleware loracloud.LoracloudMiddleware
	autoPadding         bool
	skipValidation      bool
	fCount              uint32
}

func NewSmartLabelv1Decoder(loracloudMiddleware loracloud.LoracloudMiddleware, options ...Option) decoder.Decoder {
	smartLabelv1Decoder := &SmartLabelv1Decoder{
		loracloudMiddleware: loracloudMiddleware,
	}

	for _, option := range options {
		option(smartLabelv1Decoder)
	}

	return smartLabelv1Decoder
}

func WithAutoPadding(autoPadding bool) Option {
	return func(t *SmartLabelv1Decoder) {
		t.autoPadding = autoPadding
	}
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *SmartLabelv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// WithFCount sets the frame counter for the decoder.
// This is required for the loracloud middleware.
func WithFCount(fCount uint32) Option {
	return func(t *SmartLabelv1Decoder) {
		t.fCount = fCount
	}
}

func getPort11PayloadType(data string) (string, error) {
	if len(data) < 2 {
		return "", fmt.Errorf("data length is less than 2")
	}
	firstByte := strings.ToUpper(data[0:2])
	if firstByte == "0E" || firstByte == "11" { // 14 || 17
		return "configuration", nil
	} else if firstByte == "0A" { // 10
		return "heartbeat", nil
	}
	return "", fmt.Errorf("invalid payload for port 11")
}

// https://docs.truvami.com/docs/payloads/smartlabel
func (t SmartLabelv1Decoder) getConfig(port uint8, data string) (common.PayloadConfig, error) {
	switch port {
	case 1:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BatteryVoltage", Start: 0, Length: 2, Transform: func(v any) any {
					return float32(v.(int)) / 1000
				}},
				{Name: "PhotovoltaicVoltage", Start: 2, Length: 2, Transform: func(v any) any {
					return float32(v.(int)) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBattery, decoder.FeaturePhotovoltaic},
		}, nil
	case 2:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Temperature", Start: 0, Length: 2, Transform: func(v any) any {
					return float32(v.(int)) / 100
				}},
				{Name: "Humidity", Start: 2, Length: 1, Transform: func(v any) any {
					return float32(v.(int)) / 2
				}},
			},
			TargetType: reflect.TypeOf(Port2Payload{}),
			Features:   []decoder.Feature{decoder.FeatureTemperature, decoder.FeatureHumidity},
		}, nil
	case 4:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DataRate", Start: 0, Length: 1, Transform: func(v any) any {
					return uint8(v.(int) & 0x7)
				}},
				{Name: "Acceleration", Start: 0, Length: 1, Transform: func(v any) any {
					return ((v.(int) >> 3) & 0x1) != 0
				}},
				{Name: "Wifi", Start: 0, Length: 1, Transform: func(v any) any {
					return ((v.(int) >> 4) & 0x1) != 0
				}},
				{Name: "Gnss", Start: 0, Length: 1, Transform: func(v any) any {
					return ((v.(int) >> 5) & 0x1) != 0
				}},
				{Name: "SteadyInterval", Start: 1, Length: 2},
				{Name: "MovingInterval", Start: 3, Length: 2},
				{Name: "HeartbeatInterval", Start: 5, Length: 1},
				{Name: "AccelerationThreshold", Start: 6, Length: 2},
				{Name: "AccelerationDelay", Start: 8, Length: 2},
				{Name: "TemperaturePollingInterval", Start: 10, Length: 2},
				{Name: "TemperatureUplinkInterval", Start: 12, Length: 2},
				{Name: "TemperatureLowerThreshold", Start: 14, Length: 1},
				{Name: "TemperatureUpperThreshold", Start: 15, Length: 1},
				{Name: "AccessPointsThreshold", Start: 16, Length: 1},
				{Name: "FirmwareVersionMajor", Start: 17, Length: 1, Optional: true}, // FIXME: after firmware changes field will be required
				{Name: "FirmwareVersionMinor", Start: 18, Length: 1, Optional: true}, // FIXME: after firmware changes field will be required
				{Name: "FirmwareVersionPatch", Start: 19, Length: 1, Optional: true}, // FIXME: after firmware changes field will be required
			},
			TargetType: reflect.TypeOf(Port4Payload{}),
			Features:   []decoder.Feature{decoder.FeatureConfig, decoder.FeatureFirmwareVersion},
		}, nil
	case 11:
		// Check first byte length to determine message type
		payloadType, err := getPort11PayloadType(data)
		if err != nil {
			return common.PayloadConfig{}, err
		}
		switch payloadType {
		case "configuration":
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Flags", Start: 2, Length: 1},
					{Name: "GNSSEnabled", Start: 2, Length: 1, Transform: func(v any) any {
						return uint8((v.(int) >> 1) & 0x1)
					}},
					{Name: "WiFiEnabled", Start: 2, Length: 1, Transform: func(v any) any {
						return uint8((v.(int) >> 2) & 0x1)
					}},
					{Name: "AccEnabled", Start: 2, Length: 1, Transform: func(v any) any {
						return uint8((v.(int) >> 3) & 0x1)
					}},
					{Name: "StaticSF", Start: 2, Length: 1, Transform: func(v any) any {
						return fmt.Sprintf("SF%d", 9) // TODO: Hardcoded to SF9 for now
					}},
					{Name: "SteadyIntervalS", Start: 3, Length: 2},
					{Name: "MovingIntervalS", Start: 5, Length: 2},
					{Name: "HeartbeatIntervalH", Start: 7, Length: 1},
					{Name: "LEDBlinkIntervalS", Start: 8, Length: 2},
					{Name: "AccThresholdMS", Start: 10, Length: 2},
					{Name: "AccDelayMS", Start: 12, Length: 2},
					{Name: "GitHash", Start: 14, Length: 4, Optional: true, Transform: func(v any) any {
						return fmt.Sprintf("%08x", v.(int))
					}},
				},
				TargetType: reflect.TypeOf(Port11ConfigurationPayload{}),
			}, nil
		case "heartbeat":
			return common.PayloadConfig{
				Fields: []common.FieldConfig{
					{Name: "Battery", Start: 2, Length: 2, Transform: func(v any) any {
						return float64(v.(int)) / 1000
					}},
					{Name: "Temperature", Start: 4, Length: 2, Transform: func(v any) any {
						return float64(v.(int)) / 100
					}},
					{Name: "RH", Start: 6, Length: 1, Transform: func(v any) any {
						return float64(v.(int)) / 2
					}},
					{Name: "GNSSScansCount", Start: 7, Length: 2},
					{Name: "WiFiScansCount", Start: 9, Length: 2},
				},
				TargetType: reflect.TypeOf(Port11HeartbeatPayload{}),
			}, nil
		}
		return common.PayloadConfig{}, fmt.Errorf("invalid payload for port 11")
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
				{Name: "Rssi6", Start: 36, Length: 1, Optional: true},
				{Name: "Mac6", Start: 37, Length: 6, Optional: true, Hex: true},
			},
			TargetType: reflect.TypeOf(Port197Payload{}),
			Features:   []decoder.Feature{decoder.FeatureWiFi},
		}, nil
	default:
		return common.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
	}
}

func (t SmartLabelv1Decoder) Decode(data string, port uint8, devEui string) (*decoder.DecodedUplink, error) {
	switch port {
	case 192:
		decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
			MsgType: "updf",
			Port:    port,
			Payload: data,
			FCount:  t.fCount,
		})
		return decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS}, decodedData, nil), err
	// case 197:
	// 	decodedData, err := t.loracloudMiddleware.DeliverUplinkMessage(devEui, loracloud.UplinkMsg{
	// 		MsgType: "updf",
	// 		Port:    port,
	// 		Payload: data,
	// 		FCount:  t.fCount,
	// 	})
	// 	return decoder.NewDecodedUplink([]decoder.Feature{}, decodedData, nil), err
	default:
		config, err := t.getConfig(port, data)
		if err != nil {
			return nil, err
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
