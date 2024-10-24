package tagsl

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
)

type TagSLv1Decoder struct{}

func NewTagSLv1Decoder() decoder.Decoder {
	return TagSLv1Decoder{}
}

// https://docs.truvami.com/docs/payloads/tag-S
// https://docs.truvami.com/docs/payloads/tag-L
func (t TagSLv1Decoder) getConfig(port int16) (decoder.PayloadConfig, error) {
	switch port {
	case 1:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
				{Name: "Latitude", Start: 1, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Longitude", Start: 5, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Altitude", Start: 9, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 10
				}},
				{Name: "Year", Start: 11, Length: 1},
				{Name: "Month", Start: 12, Length: 1},
				{Name: "Day", Start: 13, Length: 1},
				{Name: "Hour", Start: 14, Length: 1},
				{Name: "Minute", Start: 15, Length: 1},
				{Name: "Second", Start: 16, Length: 1},
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
		}, nil
	case 2:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(Port2Payload{}),
		}, nil
	case 3:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "ScanPointer", Start: 0, Length: 2},
				{Name: "TotalMessages", Start: 2, Length: 1},
				{Name: "CurrentMessage", Start: 3, Length: 1},
				{Name: "Mac1", Start: 4, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi1", Start: 10, Length: 1, Optional: true},
				{Name: "Mac2", Start: 11, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 17, Length: 1, Optional: true},
				{Name: "Mac3", Start: 18, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 24, Length: 1, Optional: true},
				{Name: "Mac4", Start: 25, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 31, Length: 1, Optional: true},
				{Name: "Mac5", Start: 32, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 38, Length: 1, Optional: true},
				{Name: "Mac6", Start: 39, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 45, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port3Payload{}),
		}, nil
	case 4:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "LocalizationIntervalWhileMoving", Start: 0, Length: 4},
				{Name: "LocalizationIntervalWhileSteady", Start: 4, Length: 4},
				{Name: "HeartbeatInterval", Start: 8, Length: 4},
				{Name: "GPSTimeoutWhileWaitingForFix", Start: 12, Length: 2},
				{Name: "AccelerometerWakeupThreshold", Start: 14, Length: 2},
				{Name: "AccelerometerDelay", Start: 16, Length: 2},
				{Name: "DeviceState", Start: 18, Length: 1},
				{Name: "FirmwareVersionMajor", Start: 19, Length: 1},
				{Name: "FirmwareVersionMinor", Start: 20, Length: 1},
				{Name: "FirmwareVersionPatch", Start: 21, Length: 1},
				{Name: "HardwareVersionType", Start: 22, Length: 1},
				{Name: "HardwareVersionRevision", Start: 23, Length: 1},
				{Name: "BatteryKeepAliveMessageInterval", Start: 24, Length: 4},
			},
			TargetType: reflect.TypeOf(Port4Payload{}),
		}, nil
	case 5:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
				{Name: "Mac1", Start: 1, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi1", Start: 7, Length: 1, Optional: true},
				{Name: "Mac2", Start: 8, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 14, Length: 1, Optional: true},
				{Name: "Mac3", Start: 15, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 21, Length: 1, Optional: true},
				{Name: "Mac4", Start: 22, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 28, Length: 1, Optional: true},
				{Name: "Mac5", Start: 29, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 35, Length: 1, Optional: true},
				{Name: "Mac6", Start: 36, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 42, Length: 1, Optional: true},
				{Name: "Mac7", Start: 43, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi7", Start: 49, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port5Payload{}),
		}, nil
	case 6:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "ButtonPressed", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(Port6Payload{}),
		}, nil
	case 7:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Timestamp", Start: 0, Length: 4},
				{Name: "Moving", Start: 4, Length: 1},
				{Name: "Mac1", Start: 5, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi1", Start: 11, Length: 1, Optional: true},
				{Name: "Mac2", Start: 12, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 18, Length: 1, Optional: true},
				{Name: "Mac3", Start: 19, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 25, Length: 1, Optional: true},
				{Name: "Mac4", Start: 26, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 32, Length: 1, Optional: true},
				{Name: "Mac5", Start: 33, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 39, Length: 1, Optional: true},
				{Name: "Mac6", Start: 40, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 46, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port7Payload{}),
		}, nil
	case 10:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
				{Name: "Latitude", Start: 1, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Longitude", Start: 5, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Altitude", Start: 9, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 10
				}},
				{Name: "Timestamp", Start: 11, Length: 4},
				{Name: "Battery", Start: 15, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port10Payload{}),
		}, nil
	case 15:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "LowBattery", Start: 0, Length: 1},
				{Name: "Battery", Start: 1, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port15Payload{}),
		}, nil
	case 50:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
				{Name: "Latitude", Start: 1, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Longitude", Start: 5, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Altitude", Start: 9, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 10
				}},
				{Name: "Timestamp", Start: 11, Length: 4},
				{Name: "Battery", Start: 15, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000
				}},
				{Name: "TTF", Start: 17, Length: 1},
				{Name: "Mac1", Start: 18, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi1", Start: 24, Length: 1, Optional: true},
				{Name: "Mac2", Start: 25, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 31, Length: 1, Optional: true},
				{Name: "Mac3", Start: 32, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 38, Length: 1, Optional: true},
				{Name: "Mac4", Start: 39, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 45, Length: 1, Optional: true},
				{Name: "Mac5", Start: 46, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 52, Length: 1, Optional: true},
				{Name: "Mac6", Start: 53, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 59, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port50Payload{}),
		}, nil
	case 105:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "Timestamp", Start: 2, Length: 4},
				{Name: "Moving", Start: 7, Length: 1},
				{Name: "Mac1", Start: 7, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi1", Start: 13, Length: 1, Optional: true},
				{Name: "Mac2", Start: 14, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 20, Length: 1, Optional: true},
				{Name: "Mac3", Start: 21, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 27, Length: 1, Optional: true},
				{Name: "Mac4", Start: 28, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 34, Length: 1, Optional: true},
				{Name: "Mac5", Start: 35, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 41, Length: 1, Optional: true},
				{Name: "Mac6", Start: 42, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 48, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port105Payload{}),
		}, nil
	case 110:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "Moving", Start: 2, Length: 1},
				{Name: "Latitude", Start: 3, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Longitude", Start: 7, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Altitude", Start: 11, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 10
				}},
				{Name: "Timestamp", Start: 13, Length: 4},
				{Name: "Battery", Start: 17, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port110Payload{}),
		}, nil
	case 150:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "Moving", Start: 2, Length: 1},
				{Name: "Latitude", Start: 3, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Longitude", Start: 7, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Altitude", Start: 11, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 10
				}},
				{Name: "Timestamp", Start: 13, Length: 4},
				{Name: "Battery", Start: 17, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000
				}},
				{Name: "TTF", Start: 19, Length: 1},
				{Name: "Mac1", Start: 20, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi1", Start: 26, Length: 1, Optional: true},
				{Name: "Mac2", Start: 27, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 33, Length: 1, Optional: true},
				{Name: "Mac3", Start: 34, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 40, Length: 1, Optional: true},
				{Name: "Mac4", Start: 41, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 47, Length: 1, Optional: true},
				{Name: "Mac5", Start: 48, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 54, Length: 1, Optional: true},
				{Name: "Mac6", Start: 55, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 61, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port150Payload{}),
		}, nil
	}

	return decoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}

func (t TagSLv1Decoder) Decode(data string, port int16, devEui string) (interface{}, interface{}, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, nil, err
	}

	decodedData, err := helpers.Parse(data, config)
	if err != nil {
		return nil, nil, err
	}

	statusData, err := parseStatusByte(data)
	if err != nil {
		return nil, nil, err
	}

	return decodedData, statusData, nil
}

type Status struct {
	DutyCycle           bool `json:"dutyCycle"`
	ConfigChangeId      int  `json:"configChangeId"`
	ConfigChangeSuccess bool `json:"configChangeSuccess"`
	Moving              bool `json:"moving"`
}

func parseStatusByte(hexInput string) (Status, error) {
	// Decode the hex string to bytes
	data, err := hex.DecodeString(hexInput)
	if err != nil {
		return Status{}, err
	}

	if len(data) == 0 {
		return Status{}, fmt.Errorf("invalid input, no data found")
	}

	statusByte := data[0] // Get the first byte

	// Extract bits as per the requirements
	dcFlag := (statusByte >> 7) & 0x01       // Bit 7
	confChangeID := (statusByte >> 3) & 0x0F // Bits 6:3 (4-bit)
	confSuccess := (statusByte >> 2) & 0x01  // Bit 2
	movingFlag := statusByte & 0x01          // Bit 0

	return Status{
		DutyCycle:           dcFlag == 1,
		ConfigChangeId:      int(confChangeID),
		ConfigChangeSuccess: confSuccess == 1,
		Moving:              movingFlag == 1,
	}, nil
}
