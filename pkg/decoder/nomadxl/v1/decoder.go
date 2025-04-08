package nomadxl

import (
	"fmt"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

type Option func(*NomadXLv1Decoder)

type NomadXLv1Decoder struct {
	autoPadding    bool
	skipValidation bool
}

func NewNomadXLv1Decoder(options ...Option) decoder.Decoder {
	nomadXLv1Decoder := &NomadXLv1Decoder{}

	for _, option := range options {
		option(nomadXLv1Decoder)
	}

	return nomadXLv1Decoder
}

func WithAutoPadding(autoPadding bool) Option {
	return func(t *NomadXLv1Decoder) {
		t.autoPadding = autoPadding
	}
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *NomadXLv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// https://docs.truvami.com/docs/payloads/nomad-XL
func (t NomadXLv1Decoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
	case 101:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "SystemTime", Start: 0, Length: 8},
				{Name: "UTCDate", Start: 8, Length: 4},
				{Name: "UTCTime", Start: 12, Length: 4},
				//{Name: "BufferLevelSTA", Start: 16, Length: 2},
				{Name: "BufferLevel", Start: 18, Length: 2}, //GPS
				//{Name: "BufferLevelACC", Start: 20, Length: 2},
				//{Name: "BufferLevelLOG", Start: 22, Length: 2},
				{Name: "Temperature", Start: 24, Length: 2, Transform: func(v any) any {
					return float32(v.(int)) / 10
				}},
				{Name: "Pressure", Start: 26, Length: 2, Transform: func(v any) any {
					return float32(v.(int)) / 10
				}},
				{Name: "AccelerometerXAxis", Start: 28, Length: 2},
				{Name: "AccelerometerYAxis", Start: 30, Length: 2},
				{Name: "AccelerometerZAxis", Start: 32, Length: 2},
				{Name: "Battery", Start: 34, Length: 2, Transform: func(v any) any {
					return float64(v.(int)) / 1000
				}},
				{Name: "BatteryLorawan", Start: 36, Length: 1},
				{Name: "TimeToFix", Start: 37, Length: 1, Transform: func(v any) any {
					return time.Duration(v.(int)) * time.Second
				}},
			},
			TargetType: reflect.TypeOf(Port101Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBattery, decoder.FeatureTemperature, decoder.FeatureBuffered},
		}, nil
	case 103:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{

				{Name: "UTCDate", Start: 0, Length: 4},
				{Name: "UTCTime", Start: 4, Length: 4},
				{},
				{Name: "Latitude", Start: 8, Length: 4, Transform: func(v any) any {
					return float64(v.(int)) / 100000
				}},
				{Name: "Longitude", Start: 12, Length: 4, Transform: func(v any) any {
					return float64(v.(int)) / 100000
				}},
				{Name: "Altitude", Start: 16, Length: 4, Transform: func(v any) any {
					return float64(v.(int)) / 100
				}},
			},
			TargetType: reflect.TypeOf(Port103Payload{}),
			Features:   []decoder.Feature{decoder.FeatureGNSS},
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t NomadXLv1Decoder) Decode(data string, port uint8, devEui string) (*decoder.DecodedUplink, error) {
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
