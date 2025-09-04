package tagsl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/common"
	tagsl "github.com/truvami/decoder/pkg/decoder/tagsl/v1"
	"github.com/truvami/decoder/pkg/encoder"
)

type Option func(*TagSLv1Encoder)

type TagSLv1Encoder struct{}

func NewTagSLv1Encoder(options ...Option) encoder.Encoder {
	tagSLv1Encoder := &TagSLv1Encoder{}

	for _, option := range options {
		option(tagSLv1Encoder)
	}

	return tagSLv1Encoder
}

// Encode encodes the provided data into a payload string
func (t TagSLv1Encoder) Encode(data any, port uint8) (any, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, err
	}

	payload, err := common.Encode(data, config)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// https://docs.truvami.com/docs/payloads/tag-S
// https://docs.truvami.com/docs/payloads/tag-L
func (t TagSLv1Encoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
	case 1:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Year", Start: 11, Length: 1},
				{Name: "Month", Start: 12, Length: 1},
				{Name: "Day", Start: 13, Length: 1},
				{Name: "Hour", Start: 14, Length: 1},
				{Name: "Minute", Start: 15, Length: 1},
				{Name: "Second", Start: 16, Length: 1},
			},
			TargetType: reflect.TypeOf(tagsl.Port1Payload{}),
		}, nil
	case 4:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
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
				{Name: "BatchSize", Start: 28, Length: 2, Optional: true},
				{Name: "BufferSize", Start: 30, Length: 2, Optional: true},
			},
			TargetType: reflect.TypeOf(tagsl.Port4Payload{}),
		}, nil
	case 5:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Mac1", Start: 1, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 7, Length: 1},
				{Name: "Mac2", Start: 8, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi2", Start: 14, Length: 1, Optional: true},
				{Name: "Mac3", Start: 15, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi3", Start: 21, Length: 1, Optional: true},
				{Name: "Mac4", Start: 22, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi4", Start: 28, Length: 1, Optional: true},
				{Name: "Mac5", Start: 29, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi5", Start: 35, Length: 1, Optional: true},
				{Name: "Mac6", Start: 36, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi6", Start: 42, Length: 1, Optional: true},
				{Name: "Mac7", Start: 43, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi7", Start: 49, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(tagsl.Port5Payload{}),
		}, nil
	case 6:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "ButtonPressed", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(tagsl.Port6Payload{}),
		}, nil
	case 7:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
				{Name: "Moving", Start: 4, Length: 1, Transform: moving},
				{Name: "Mac1", Start: 5, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 11, Length: 1},
				{Name: "Mac2", Start: 12, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi2", Start: 18, Length: 1, Optional: true},
				{Name: "Mac3", Start: 19, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi3", Start: 25, Length: 1, Optional: true},
				{Name: "Mac4", Start: 26, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi4", Start: 32, Length: 1, Optional: true},
				{Name: "Mac5", Start: 33, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi5", Start: 39, Length: 1, Optional: true},
				{Name: "Mac6", Start: 40, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi6", Start: 46, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(tagsl.Port7Payload{}),
		}, nil
	case 10:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: battery},
				{Name: "TTF", Start: 17, Length: 1, Transform: ttf},
				{Name: "PDOP", Start: 18, Length: 1, Transform: pdop},
				{Name: "Satellites", Start: 19, Length: 1},
			},
			TargetType: reflect.TypeOf(tagsl.Port10Payload{}),
		}, nil
	case 15:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "LowBattery", Start: 0, Length: 1, Transform: lowBattery},
				{Name: "Battery", Start: 1, Length: 2, Transform: battery},
			},
			TargetType: reflect.TypeOf(tagsl.Port15Payload{}),
		}, nil
	case 50:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: battery},
				{Name: "TTF", Start: 17, Length: 1, Transform: ttf},
				{Name: "Mac1", Start: 18, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 24, Length: 1},
				{Name: "Mac2", Start: 25, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi2", Start: 31, Length: 1, Optional: true},
				{Name: "Mac3", Start: 32, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi3", Start: 38, Length: 1, Optional: true},
				{Name: "Mac4", Start: 39, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi4", Start: 45, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(tagsl.Port50Payload{}),
		}, nil
	case 51:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: battery},
				{Name: "TTF", Start: 17, Length: 1, Transform: ttf},
				{Name: "PDOP", Start: 18, Length: 1, Transform: pdop},
				{Name: "Satellites", Start: 19, Length: 1},
				{Name: "Mac1", Start: 20, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 26, Length: 1},
				{Name: "Mac2", Start: 27, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi2", Start: 33, Length: 1, Optional: true},
				{Name: "Mac3", Start: 34, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi3", Start: 40, Length: 1, Optional: true},
				{Name: "Mac4", Start: 41, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi4", Start: 47, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(tagsl.Port51Payload{}),
		}, nil
	case 105:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "Timestamp", Start: 2, Length: 4, Transform: timestamp},
				{Name: "Moving", Start: 6, Length: 1, Transform: moving},
				{Name: "Mac1", Start: 7, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 13, Length: 1},
				{Name: "Mac2", Start: 14, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi2", Start: 20, Length: 1, Optional: true},
				{Name: "Mac3", Start: 21, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi3", Start: 27, Length: 1, Optional: true},
				{Name: "Mac4", Start: 28, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi4", Start: 34, Length: 1, Optional: true},
				{Name: "Mac5", Start: 35, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi5", Start: 41, Length: 1, Optional: true},
				{Name: "Mac6", Start: 42, Length: 6, Hex: true, Optional: true},
				{Name: "Rssi6", Start: 48, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(tagsl.Port105Payload{}),
		}, nil
	case 110:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "Moving", Start: 2, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 3, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 7, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 11, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 13, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 17, Length: 2, Transform: battery},
				{Name: "TTF", Start: 19, Length: 1, Transform: ttf},
				{Name: "PDOP", Start: 20, Length: 1, Transform: pdop},
				{Name: "Satellites", Start: 21, Length: 1},
			},
			TargetType: reflect.TypeOf(tagsl.Port110Payload{}),
		}, nil
	case 128:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Ble", Start: 0, Length: 1},
				{Name: "Gnss", Start: 1, Length: 1},
				{Name: "Wifi", Start: 2, Length: 1},
				{Name: "MovingInterval", Start: 3, Length: 4},
				{Name: "SteadyInterval", Start: 7, Length: 4},
				{Name: "ConfigInterval", Start: 11, Length: 4},
				{Name: "GnssTimeout", Start: 15, Length: 2},
				{Name: "AccelerometerThreshold", Start: 17, Length: 2},
				{Name: "AccelerometerDelay", Start: 19, Length: 2},
				{Name: "BatteryInterval", Start: 21, Length: 4},
				{Name: "BatchSize", Start: 25, Length: 2, Optional: true},
				{Name: "BufferSize", Start: 27, Length: 2, Optional: true},
			},
			TargetType: reflect.TypeOf(Port128Payload{}),
		}, nil
	case 129:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "TimeToBuzz", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(Port129Payload{}),
		}, nil
	case 130:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "EraseFlash", Start: 0, Length: 1, Transform: func(v any) any {
					erase := common.BytesToBool(v.([]byte))
					if erase {
						return []byte{0xde}
					}
					return []byte{0x00}
				}},
			},
			TargetType: reflect.TypeOf(Port130Payload{}),
		}, nil
	case 131:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "AccuracyEnhancement", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(Port131Payload{}),
		}, nil
	case 132:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "EraseFlash", Start: 0, Length: 1, Transform: func(v any) any {
					return []byte{0x00}
				}},
			},
			TargetType: reflect.TypeOf(Port132Payload{}),
		}, nil
	case 134:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "ScanInterval", Start: 0, Length: 2},
				{Name: "ScanTime", Start: 2, Length: 1},
				{Name: "MaxBeacons", Start: 3, Length: 1},
				{Name: "MinRssi", Start: 4, Length: 1},
				{Name: "AdvertisingName", Start: 5, Length: 10, Transform: func(v any) any {
					if len(v.([]byte)) > 9 {
						v = v.([]byte)[:9]
					}
					return v
				}},
				{Name: "AccelerometerDelay", Start: 15, Length: 2},
				{Name: "AccelerometerThreshold", Start: 17, Length: 2},
				{Name: "ScanMode", Start: 19, Length: 1},
				{Name: "BleConfigUplinkInterval", Start: 20, Length: 2},
			},
			TargetType: reflect.TypeOf(Port134Payload{}),
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func moving(v any) any {
	return common.BoolToBytes(common.BytesToBool(v.([]byte)), 0)
}

func timestamp(v any) any {
	return common.IntToBytes(common.BytesToInt64(v.([]byte)), 4)
}

func latitude(v any) any {
	return common.IntToBytes(int64(common.BytesToFloat64(v.([]byte))*1000000), 4)
}

func longitude(v any) any {
	return common.IntToBytes(int64(common.BytesToFloat64(v.([]byte))*1000000), 4)
}

func altitude(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*10), 2)
}

func lowBattery(v any) any {
	return common.BoolToBytes(common.BytesToBool(v.([]byte)), 0)
}

func battery(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*1000), 2)
}

func ttf(v any) any {
	return common.UintToBytes(uint64(common.BytesToInt64(v.([]byte))/1000000000), 1)
}

func pdop(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*2), 1)
}
