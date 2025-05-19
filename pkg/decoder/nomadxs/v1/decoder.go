package nomadxs

import (
	"fmt"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

type Option func(*NomadXSv1Decoder)

type NomadXSv1Decoder struct {
	autoPadding    bool
	skipValidation bool
}

func NewNomadXSv1Decoder(options ...Option) decoder.Decoder {
	nomadXSv1Decoder := &NomadXSv1Decoder{}

	for _, option := range options {
		option(nomadXSv1Decoder)
	}

	return nomadXSv1Decoder
}

func WithAutoPadding(autoPadding bool) Option {
	return func(t *NomadXSv1Decoder) {
		t.autoPadding = autoPadding
	}
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *NomadXSv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// https://docs.truvami.com/docs/payloads/nomad-xs
func (t NomadXSv1Decoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
	case 1:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "ConfigChangeId", Start: 0, Length: 1, Transform: configChangeId},
				{Name: "ConfigChangeSuccess", Start: 0, Length: 1, Transform: configChangeSuccess},
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: func(v any) any {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Longitude", Start: 5, Length: 4, Transform: func(v any) any {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Altitude", Start: 9, Length: 2, Transform: func(v any) any {
					return float64(v.(int)) / 10
				}},
				{Name: "Year", Start: 11, Length: 1},
				{Name: "Month", Start: 12, Length: 1},
				{Name: "Day", Start: 13, Length: 1},
				{Name: "Hour", Start: 14, Length: 1},
				{Name: "Minute", Start: 15, Length: 1},
				{Name: "Second", Start: 16, Length: 1},
				{Name: "TimeToFix", Start: 17, Length: 1, Transform: func(v any) any {
					return time.Duration(v.(int)) * time.Second
				}},
				{Name: "AmbientLight", Start: 18, Length: 2},
				{Name: "AccelerometerXAxis", Start: 20, Length: 2},
				{Name: "AccelerometerYAxis", Start: 22, Length: 2},
				{Name: "AccelerometerZAxis", Start: 24, Length: 2},
				{Name: "Temperature", Start: 26, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(v.(int)) / 100
				}},
				{Name: "Pressure", Start: 28, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(v.(int)) / 10
				}},
				{Name: "GyroscopeXAxis", Start: 30, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(int16(v.(int))) / 10
				}},
				{Name: "GyroscopeYAxis", Start: 32, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(int16(v.(int))) / 10
				}},
				{Name: "GyroscopeZAxis", Start: 34, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(int16(v.(int))) / 10
				}},
				{Name: "MagnetometerXAxis", Start: 36, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(int16(v.(int))) / 1000
				}},
				{Name: "MagnetometerYAxis", Start: 38, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(int16(v.(int))) / 1000
				}},
				{Name: "MagnetometerZAxis", Start: 40, Length: 2, Optional: true, Transform: func(v any) any {
					return float32(int16(v.(int))) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureTemperature, decoder.FeaturePressure},
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
			Features:   []decoder.Feature{decoder.FeatureConfig, decoder.FeatureFirmwareVersion, decoder.FeatureHardwareVersion},
		}, nil
	case 15:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "LowBattery", Start: 0, Length: 1, Transform: lowBattery},
				{Name: "Battery", Start: 1, Length: 2, Transform: func(v any) any {
					return float64(v.(int)) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port15Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureBattery},
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t NomadXSv1Decoder) Decode(data string, port uint8, devEui string) (*decoder.DecodedUplink, error) {
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
	return decoder.NewDecodedUplink(config.Features, decodedData), err
}

func dutyCycle(v any) any {
	i, ok := v.(int)
	if !ok {
		return nil
	}
	return (byte(i)>>7)&0x01 == 1
}

func configChangeId(v any) any {
	i, ok := v.(int)
	if !ok {
		return nil
	}
	return (byte(i) >> 3) & 0x0f
}

func configChangeSuccess(v any) any {
	i, ok := v.(int)
	if !ok {
		return nil
	}
	return (byte(i)>>2)&0x01 == 1
}

func moving(v any) any {
	i, ok := v.(int)
	if !ok {
		return nil
	}
	return byte(i)&0x01 == 1
}

func lowBattery(v any) any {
	i, ok := v.(int)
	if !ok {
		return nil
	}
	return byte(i)&0x01 == 1
}
