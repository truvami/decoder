package tagsl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/common"
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
func (t TagSLv1Encoder) Encode(data interface{}, port int16, extra string) (interface{}, interface{}, error) {
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
func (t TagSLv1Encoder) getConfig(port int16) (common.PayloadConfig, error) {
	switch port {
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
				{Name: "AdvertisingName", Start: 5, Length: 10, Transform: func(v interface{}) interface{} {
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

	return common.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}
