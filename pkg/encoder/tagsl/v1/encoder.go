package tagsl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/encoder"
	"github.com/truvami/decoder/pkg/encoder/helpers"
)

type Option func(*TagSLv1Encoder)

type TagSLv1Encoder struct {
	autoPadding    bool
	skipValidation bool
}

func NewTagSLv1Encoder(options ...Option) encoder.Encoder {
	tagSLv1Encoder := &TagSLv1Encoder{}

	for _, option := range options {
		option(tagSLv1Encoder)
	}

	return tagSLv1Encoder
}

func WithAutoPadding(autoPadding bool) Option {
	return func(t *TagSLv1Encoder) {
		t.autoPadding = autoPadding
	}
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *TagSLv1Encoder) {
		t.skipValidation = skipValidation
	}
}

// Encode encodes the provided data into a payload string
func (t TagSLv1Encoder) Encode(data interface{}, port int16, extra string) (interface{}, interface{}, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, nil, err
	}

	payload, err := helpers.Encode(data, config)
	if err != nil {
		return nil, nil, err
	}

	if !t.skipValidation {
		err := helpers.ValidateLength(&payload, &config)
		if err != nil {
			return nil, nil, err
		}
	}

	return payload, extra, nil
}

// https://docs.truvami.com/docs/payloads/tag-S
// https://docs.truvami.com/docs/payloads/tag-L
func (t TagSLv1Encoder) getConfig(port int16) (encoder.PayloadConfig, error) {
	switch port {
	case 128:
		return encoder.PayloadConfig{
			Fields: []encoder.FieldConfig{
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
	}

	return encoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}
