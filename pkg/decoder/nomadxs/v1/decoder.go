package nomadxs

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
)

type NomadXSv1Decoder struct{}

func NewNomadXSv1Decoder() decoder.Decoder {
	return NomadXSv1Decoder{}
}

// https://docs.truvami.com/docs/payloads/nomad-xs
func (t NomadXSv1Decoder) getConfig(port int16) (decoder.PayloadConfig, error) {
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
				{Name: "TimeToFix", Start: 17, Length: 1},
				{Name: "AmbientLight", Start: 18, Length: 2},
				{Name: "AccelerometerXAxis", Start: 20, Length: 2},
				{Name: "AccelerometerYAxis", Start: 22, Length: 2},
				{Name: "AccelerometerZAxis", Start: 24, Length: 2},
				{Name: "Temperature", Start: 26, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int)) / 100
				}},
				{Name: "Pressure", Start: 28, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int))
				}},
				{Name: "GyroscopeXAxis", Start: 30, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int)) / 10
				}},
				{Name: "GyroscopeYAxis", Start: 32, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int)) / 10
				}},
				{Name: "GyroscopeZAxis", Start: 34, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int)) / 10
				}},
				{Name: "MagnetometerXAxis", Start: 36, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int))
				}},
				{Name: "MagnetometerYAxis", Start: 38, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int))
				}},
				{Name: "MagnetometerZAxis", Start: 40, Length: 2, Optional: true, Transform: func(v interface{}) interface{} {
					return float32(v.(int))
				}},
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
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
				{Name: "FirmwareVersionMajor", Start: 18, Length: 1},
				{Name: "FirmwareVersionMinor", Start: 19, Length: 1},
				{Name: "FirmwareVersionPatch", Start: 20, Length: 1},
				{Name: "HardwareVersionType", Start: 21, Length: 1},
				{Name: "HardwareVersionRevision", Start: 22, Length: 1},
				{Name: "BatteryKeepAliveMessageInterval", Start: 23, Length: 4},
				{Name: "ReJoinInterval", Start: 27, Length: 4},
				{Name: "AccuracyEnhancement", Start: 31, Length: 1},
				{Name: "LightLowerThreshold", Start: 32, Length: 2},
				{Name: "LightUpperThreshold", Start: 34, Length: 2},
			},
			TargetType: reflect.TypeOf(Port4Payload{}),
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
	}

	return decoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}

func (t NomadXSv1Decoder) Decode(data string, port int16, devEui string) (interface{}, interface{}, error) {
	config, err := t.getConfig(port)
	if err != nil {
		return nil, nil, err
	}

	decodedData, err := helpers.Parse(data, config)
	if err != nil {
		return nil, nil, err
	}

	return decodedData, nil, nil
}
