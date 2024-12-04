package nomadxl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
)

type NomadXLv1Decoder struct{}

func NewNomadXLv1Decoder() decoder.Decoder {
	return NomadXLv1Decoder{}
}

// https://docs.truvami.com/docs/payloads/nomad-XL
func (t NomadXLv1Decoder) getConfig(port int16) (decoder.PayloadConfig, error) {
	switch port {
	case 101:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "SystemTime", Start: 0, Length: 8},
				{Name: "UTCDate", Start: 8, Length: 4},
				{Name: "UTCTime", Start: 12, Length: 4},
				//{Name: "BufferLevelSTA", Start: 16, Length: 2},
				{Name: "BufferLevel", Start: 18, Length: 2}, //GPS
				//{Name: "BufferLevelACC", Start: 20, Length: 2},
				//{Name: "BufferLevelLOG", Start: 22, Length: 2},
				{Name: "Temperature", Start: 24, Length: 2, Transform: func(v interface{}) interface{} {
					return float32(v.(int)) / 10
				}},
				{Name: "Pressure", Start: 26, Length: 2, Transform: func(v interface{}) interface{} {
					return float32(v.(int))
				}},
				{Name: "AccelerometerXAxis", Start: 28, Length: 2},
				{Name: "AccelerometerYAxis", Start: 30, Length: 2},
				{Name: "AccelerometerZAxis", Start: 32, Length: 2},
				{Name: "Battery", Start: 34, Length: 2, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000
				}},
				{Name: "BatteryLorawan", Start: 36, Length: 1},
				{Name: "TimeToFix", Start: 37, Length: 1},
			},
			TargetType: reflect.TypeOf(Port101Payload{}),
		}, nil
	case 103:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{

				{Name: "UTCDate", Start: 0, Length: 4},
				{Name: "UTCTime", Start: 4, Length: 4},
				{},
				{Name: "Latitude", Start: 8, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 100000
				}},
				{Name: "Longitude", Start: 12, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 100000
				}},
				{Name: "Altitude", Start: 16, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 100
				}},
			},
			TargetType: reflect.TypeOf(Port103Payload{}),
		}, nil
	}

	return decoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}

func (t NomadXLv1Decoder) Decode(data string, port int16, devEui string, autoPadding bool) (interface{}, interface{}, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, nil, err
	}

	if autoPadding {
		data = helpers.HexNullPad(&data, &config)
	}

	decodedData, err := helpers.Parse(data, config)
	if err != nil {
		return nil, nil, err
	}

	return decodedData, nil, nil
}
