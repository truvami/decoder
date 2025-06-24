package smartlabel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"go.uber.org/zap"
)

type Option func(*SmartLabelv1Decoder)

type SmartLabelv1Decoder struct {
	skipValidation bool
	logger         *zap.Logger

	solver solver.SolverV1
}

func NewSmartLabelv1Decoder(ctx context.Context, solver solver.SolverV1, logger *zap.Logger, options ...Option) decoder.Decoder {
	if solver == nil {
		logger.Panic("solver cannot be nil", zap.String("decoder", "SmartLabelv1Decoder"))
	}

	smartLabelv1Decoder := &SmartLabelv1Decoder{
		logger: logger,
		solver: solver,
	}

	for _, option := range options {
		option(smartLabelv1Decoder)
	}

	return smartLabelv1Decoder
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *SmartLabelv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// https://docs.truvami.com/docs/payloads/smartlabel
func (t SmartLabelv1Decoder) getConfig(port uint8, data string) (common.PayloadConfig, error) {
	switch port {
	case 1:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BatteryVoltage", Start: 0, Length: 2, Transform: battery},
				{Name: "PhotovoltaicVoltage", Start: 2, Length: 2, Transform: photovoltaic},
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBattery, decoder.FeaturePhotovoltaic},
		}, nil
	case 2:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Temperature", Start: 0, Length: 2, Transform: temperature},
				{Name: "Humidity", Start: 2, Length: 1, Transform: humidity},
			},
			TargetType: reflect.TypeOf(Port2Payload{}),
			Features:   []decoder.Feature{decoder.FeatureTemperature, decoder.FeatureHumidity},
		}, nil
	case 4:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DataRate", Start: 0, Length: 1, Transform: func(v any) any {
					return v.([]byte)[0] & 0x7
				}},
				{Name: "Acceleration", Start: 0, Length: 1, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 3) & 0x1) != 0
				}},
				{Name: "Wifi", Start: 0, Length: 1, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 4) & 0x1) != 0
				}},
				{Name: "Gnss", Start: 0, Length: 1, Transform: func(v any) any {
					return ((v.([]byte)[0] >> 5) & 0x1) != 0
				}},
				{Name: "SteadyInterval", Start: 1, Length: 2},
				{Name: "MovingInterval", Start: 3, Length: 2},
				{Name: "HeartbeatInterval", Start: 5, Length: 1},
				{Name: "AccelerationThreshold", Start: 6, Length: 2},
				{Name: "AccelerationDelay", Start: 8, Length: 2},
				{Name: "TemperaturePollingInterval", Start: 10, Length: 2},
				{Name: "TemperatureUplinkInterval", Start: 12, Length: 2},
				{Name: "TemperatureLowerThreshold", Start: 14, Length: 1},
				{Name: "TemperatureUpperThreshold", Start: 15, Length: 1},
				{Name: "AccessPointsThreshold", Start: 16, Length: 1},
				{Name: "FirmwareVersionMajor", Start: 17, Length: 1, Optional: true}, // FIXME: after firmware changes field will be required
				{Name: "FirmwareVersionMinor", Start: 18, Length: 1, Optional: true}, // FIXME: after firmware changes field will be required
				{Name: "FirmwareVersionPatch", Start: 19, Length: 1, Optional: true}, // FIXME: after firmware changes field will be required
			},
			TargetType: reflect.TypeOf(Port4Payload{}),
			Features:   []decoder.Feature{decoder.FeatureConfig, decoder.FeatureFirmwareVersion},
		}, nil
	case 11:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BatteryVoltage", Start: 0, Length: 2, Transform: battery},
				{Name: "PhotovoltaicVoltage", Start: 2, Length: 2, Transform: photovoltaic},
				{Name: "Temperature", Start: 4, Length: 2, Transform: temperature},
				{Name: "Humidity", Start: 6, Length: 1, Transform: humidity},
			},
			TargetType: reflect.TypeOf(Port11Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBattery, decoder.FeaturePhotovoltaic, decoder.FeatureTemperature, decoder.FeatureHumidity},
		}, nil
	case 150:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Battery100Voltage", Start: 0, Length: 2, Transform: func(v any) any {
					return float32(common.BytesToUint16(v.([]byte))) / 1000
				}},
				{Name: "Battery80Voltage", Start: 2, Length: 2, Transform: func(v any) any {
					return float32(common.BytesToUint16(v.([]byte))) / 1000
				}},
				{Name: "Battery60Voltage", Start: 4, Length: 2, Transform: func(v any) any {
					return float32(common.BytesToUint16(v.([]byte))) / 1000
				}},
				{Name: "Battery40Voltage", Start: 6, Length: 2, Transform: func(v any) any {
					return float32(common.BytesToUint16(v.([]byte))) / 1000
				}},
				{Name: "Battery20Voltage", Start: 8, Length: 2, Transform: func(v any) any {
					return float32(common.BytesToUint16(v.([]byte))) / 1000
				}},
			},
			TargetType: reflect.TypeOf(Port150Payload{}),
			Features:   []decoder.Feature{},
		}, nil
	case 197:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Tag", Start: 0, Length: 1},
				{Name: "Rssi1", Start: 1, Length: 1},
				{Name: "Mac1", Start: 2, Length: 6, Hex: true},
				{Name: "Rssi2", Start: 8, Length: 1, Optional: true},
				{Name: "Mac2", Start: 9, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 15, Length: 1, Optional: true},
				{Name: "Mac3", Start: 16, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 22, Length: 1, Optional: true},
				{Name: "Mac4", Start: 23, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 29, Length: 1, Optional: true},
				{Name: "Mac5", Start: 30, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 36, Length: 1, Optional: true},
				{Name: "Mac6", Start: 37, Length: 6, Optional: true, Hex: true},
			},
			TargetType: reflect.TypeOf(Port197Payload{}),
			Features:   []decoder.Feature{decoder.FeatureWiFi},
		}, nil
	default:
		return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
	}
}

func (t SmartLabelv1Decoder) Decode(ctx context.Context, data string, port uint8) (*decoder.DecodedUplink, error) {
	switch port {
	case 192:
		return t.solver.Solve(ctx, data)
	default:
		config, err := t.getConfig(port, data)
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
}

func battery(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 1000
}

func photovoltaic(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 1000
}

func temperature(v any) any {
	return float32(common.BytesToUint16(v.([]byte))) / 100
}

func humidity(v any) any {
	return float32(common.BytesToUint8(v.([]byte))) / 2
}
