package nomadxl

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

type Option func(*NomadXLv1Decoder)

type NomadXLv1Decoder struct {
	skipValidation bool
}

func NewNomadXLv1Decoder(options ...Option) decoder.Decoder {
	nomadXLv1Decoder := &NomadXLv1Decoder{}

	for _, option := range options {
		option(nomadXLv1Decoder)
	}

	return nomadXLv1Decoder
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
				{Name: "BufferLevelSTA", Start: 16, Length: 2},
				{Name: "BufferLevelGPS", Start: 18, Length: 2},
				{Name: "BufferLevelACC", Start: 20, Length: 2},
				{Name: "BufferLevelLOG", Start: 22, Length: 2},
				{Name: "Temperature", Start: 24, Length: 2, Transform: temperature},
				{Name: "Pressure", Start: 26, Length: 2, Transform: pressure},
				{Name: "AccelerometerXAxis", Start: 28, Length: 2},
				{Name: "AccelerometerYAxis", Start: 30, Length: 2},
				{Name: "AccelerometerZAxis", Start: 32, Length: 2},
				{Name: "Battery", Start: 34, Length: 2, Transform: battery},
				{Name: "BatteryLorawan", Start: 36, Length: 1},
				{Name: "TimeToFix", Start: 37, Length: 1, Transform: ttf},
			},
			TargetType: reflect.TypeOf(Port101Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBuffered, decoder.FeatureBattery, decoder.FeatureTemperature, decoder.FeaturePressure},
		}, nil
	case 103:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "UTCDate", Start: 0, Length: 4},
				{Name: "UTCTime", Start: 4, Length: 4},
				{Name: "Latitude", Start: 8, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 12, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 16, Length: 4, Transform: altitude},
			},
			TargetType: reflect.TypeOf(Port103Payload{}),
			Features:   []decoder.Feature{decoder.FeatureGNSS},
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t NomadXLv1Decoder) Decode(ctx context.Context, data string, port uint8) (*decoder.DecodedUplink, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, err
	}

	if !t.skipValidation {
		err := common.ValidateLength(&data, &config)
		if err != nil {
			return nil, err
		}
	}

	decodedData, err := common.Decode(&data, &config)
	return decoder.NewDecodedUplink(config.Features, decodedData), err
}

func temperature(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 10
}

func pressure(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 10
}

func battery(v any) any {
	return float64(common.BytesToUint16(v.([]byte))) / 1000
}

func ttf(v any) any {
	return time.Duration(int64(common.BytesToUint8(v.([]byte)))) * time.Second
}

func latitude(v any) any {
	return float64(common.BytesToInt32(v.([]byte))) / 100000
}

func longitude(v any) any {
	return float64(common.BytesToInt32(v.([]byte))) / 100000
}

func altitude(v any) any {
	return float64(common.BytesToUint16(v.([]byte))) / 100
}
