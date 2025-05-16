package tagsl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
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
func (t TagSLv1Encoder) Encode(data any, port uint8, extra string) (any, any, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, nil, err
	}

	payload, err := common.Encode(data, config)
	if err != nil {
		return nil, nil, err
	}

	return payload, extra, nil
}

// https://docs.truvami.com/docs/payloads/tag-S
// https://docs.truvami.com/docs/payloads/tag-L
func (t TagSLv1Encoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
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
				{Name: "Latitude", Start: 1, Length: 4, Transform: func(v any) any {
					return common.IntToBytes(int64(common.BytesToFloat64(v.([]byte))*1000000), 4)
				}},
				{Name: "Longitude", Start: 5, Length: 4, Transform: func(v any) any {
					return common.IntToBytes(int64(common.BytesToFloat64(v.([]byte))*1000000), 4)
				}},
				{Name: "Altitude", Start: 9, Length: 2, Transform: func(v any) any {
					return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*10), 2)
				}},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: func(v any) any {
					return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*1000), 2)
				}},
				{Name: "TTF", Start: 17, Length: 1, Transform: func(v any) any {
					return common.UintToBytes(uint64(common.BytesToInt64(v.([]byte))/1000000000), 1)
				}},
				{Name: "PDOP", Start: 18, Length: 1, Transform: func(v any) any {
					return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*2), 1)
				}},
				{Name: "Satellites", Start: 19, Length: 1},
			},
			TargetType: reflect.TypeOf(tagsl.Port10Payload{}),
		}, nil
	case 15:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "LowBattery", Start: 0, Length: 1, Transform: func(v any) any {
					return common.BoolToBytes(common.BytesToBool(v.([]byte)), 0)
				}},
				{Name: "Battery", Start: 1, Length: 2, Transform: func(v any) any {
					return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*1000), 2)
				}},
			},
			TargetType: reflect.TypeOf(tagsl.Port15Payload{}),
		}, nil
	case 128:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BLE", Start: 0, Length: 1},
				{Name: "GPS", Start: 1, Length: 1},
				{Name: "WIFI", Start: 2, Length: 1},
				{Name: "LocalizationIntervalWhileMoving", Start: 3, Length: 4},
				{Name: "LocalizationIntervalWhileSteady", Start: 7, Length: 4},
				{Name: "HeartbeatInterval", Start: 11, Length: 4},
				{Name: "GPSTimeoutWhileWaitingForFix", Start: 15, Length: 2},
				{Name: "AccelerometerWakeupThreshold", Start: 17, Length: 2},
				{Name: "AccelerometerDelay", Start: 19, Length: 2},
				{Name: "BatteryKeepAliveMessageInterval", Start: 21, Length: 4},
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
	case 131:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "AccuracyEnhancement", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(Port131Payload{}),
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
				{Name: "AccelerometerTriggerHoldTimer", Start: 15, Length: 2},
				{Name: "AcceleratorThreshold", Start: 17, Length: 2},
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
