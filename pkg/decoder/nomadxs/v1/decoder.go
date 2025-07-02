package nomadxs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

type Option func(*NomadXSv1Decoder)

type NomadXSv1Decoder struct {
	skipValidation bool
}

func NewNomadXSv1Decoder(options ...Option) decoder.Decoder {
	nomadXSv1Decoder := &NomadXSv1Decoder{}

	for _, option := range options {
		option(nomadXSv1Decoder)
	}

	return nomadXSv1Decoder
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
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
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
				{Name: "TimeToFix", Start: 17, Length: 1, Transform: ttf},
				{Name: "AmbientLight", Start: 18, Length: 2},
				{Name: "AccelerometerXAxis", Start: 20, Length: 2},
				{Name: "AccelerometerYAxis", Start: 22, Length: 2},
				{Name: "AccelerometerZAxis", Start: 24, Length: 2},
				{Name: "Temperature", Start: 26, Length: 2, Optional: true, Transform: temperature},
				{Name: "Pressure", Start: 28, Length: 2, Optional: true, Transform: pressure},
				{Name: "GyroscopeXAxis", Start: 30, Length: 2, Optional: true, Transform: gyroscope},
				{Name: "GyroscopeYAxis", Start: 32, Length: 2, Optional: true, Transform: gyroscope},
				{Name: "GyroscopeZAxis", Start: 34, Length: 2, Optional: true, Transform: gyroscope},
				{Name: "MagnetometerXAxis", Start: 36, Length: 2, Optional: true, Transform: magnetometer},
				{Name: "MagnetometerYAxis", Start: 38, Length: 2, Optional: true, Transform: magnetometer},
				{Name: "MagnetometerZAxis", Start: 40, Length: 2, Optional: true, Transform: magnetometer},
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureTimestamp, decoder.FeatureTemperature, decoder.FeaturePressure},
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
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
				{Name: "LowBattery", Start: 0, Length: 1, Transform: lowBattery},
				{Name: "Battery", Start: 1, Length: 2, Transform: battery},
			},
			TargetType: reflect.TypeOf(Port15Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureBattery},
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t NomadXSv1Decoder) Decode(ctx context.Context, data string, port uint8) (*decoder.DecodedUplink, error) {
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

func dutyCycle(v any) any {
	return ((v.([]byte))[0]>>7)&0x01 == 1
}

func configId(v any) any {
	return ((v.([]byte))[0] >> 3) & 0x0f
}

func configSuccess(v any) any {
	return ((v.([]byte))[0]>>2)&0x01 == 1
}

func moving(v any) any {
	return (v.([]byte))[0]&0x01 == 1
}

func lowBattery(v any) any {
	return (v.([]byte))[0]&0x01 == 1
}

func latitude(v any) any {
	return float64(common.BytesToInt32(v.([]byte))) / 1000000
}

func longitude(v any) any {
	return float64(common.BytesToInt32(v.([]byte))) / 1000000
}

func altitude(v any) any {
	return float64(common.BytesToUint16(v.([]byte))) / 10
}

func ttf(v any) any {
	return time.Duration(int64(common.BytesToUint8(v.([]byte)))) * time.Second
}

func temperature(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 100
}

func pressure(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 10
}

func gyroscope(v any) any {
	return float32(common.BytesToInt16(v.([]byte))) / 10
}

func magnetometer(v any) any {
	return float32(common.BytesToInt16(v.([]byte))) / 1000
}

func battery(v any) any {
	return float64(common.BytesToUint16(v.([]byte))) / 1000
}
