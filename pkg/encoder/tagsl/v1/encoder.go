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
