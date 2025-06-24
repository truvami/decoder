package tagsl

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

type Option func(*TagSLv1Decoder)

type TagSLv1Decoder struct {
	skipValidation bool
}

func NewTagSLv1Decoder(options ...Option) decoder.Decoder {
	tagSLv1Decoder := &TagSLv1Decoder{}

	for _, option := range options {
		option(tagSLv1Decoder)
	}

	return tagSLv1Decoder
}

func WithSkipValidation(skipValidation bool) Option {
	return func(t *TagSLv1Decoder) {
		t.skipValidation = skipValidation
	}
}

// https://docs.truvami.com/docs/payloads/tag-S
// https://docs.truvami.com/docs/payloads/tag-L
func (t TagSLv1Decoder) getConfig(port uint8) (common.PayloadConfig, error) {
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
			},
			TargetType: reflect.TypeOf(Port1Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS},
		}, nil
	case 2:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
			},
			TargetType: reflect.TypeOf(Port2Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving},
		}, nil
	case 3:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "ScanPointer", Start: 0, Length: 2},
				{Name: "TotalMessages", Start: 2, Length: 1},
				{Name: "CurrentMessage", Start: 3, Length: 1},
				{Name: "Mac1", Start: 4, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 10, Length: 1},
				{Name: "Mac2", Start: 11, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 17, Length: 1, Optional: true},
				{Name: "Mac3", Start: 18, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 24, Length: 1, Optional: true},
				{Name: "Mac4", Start: 25, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 31, Length: 1, Optional: true},
				{Name: "Mac5", Start: 32, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 38, Length: 1, Optional: true},
				{Name: "Mac6", Start: 39, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 45, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port3Payload{}),
			Features:   []decoder.Feature{decoder.FeatureWiFi},
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
				{Name: "DeviceState", Start: 18, Length: 1},
				{Name: "FirmwareVersionMajor", Start: 19, Length: 1},
				{Name: "FirmwareVersionMinor", Start: 20, Length: 1},
				{Name: "FirmwareVersionPatch", Start: 21, Length: 1},
				{Name: "HardwareVersionType", Start: 22, Length: 1},
				{Name: "HardwareVersionRevision", Start: 23, Length: 1},
				{Name: "BatteryKeepAliveMessageInterval", Start: 24, Length: 4},
				{Name: "BatchSize", Start: 28, Length: 2, Optional: true},
				{Name: "BufferSize", Start: 30, Length: 2, Optional: true},
			},
			TargetType: reflect.TypeOf(Port4Payload{}),
			Features:   []decoder.Feature{decoder.FeatureMoving, decoder.FeatureConfig, decoder.FeatureHardwareVersion, decoder.FeatureFirmwareVersion},
		}, nil
	case 5:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Mac1", Start: 1, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 7, Length: 1},
				{Name: "Mac2", Start: 8, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 14, Length: 1, Optional: true},
				{Name: "Mac3", Start: 15, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 21, Length: 1, Optional: true},
				{Name: "Mac4", Start: 22, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 28, Length: 1, Optional: true},
				{Name: "Mac5", Start: 29, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 35, Length: 1, Optional: true},
				{Name: "Mac6", Start: 36, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 42, Length: 1, Optional: true},
				{Name: "Mac7", Start: 43, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi7", Start: 49, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port5Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureWiFi},
		}, nil
	case 6:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "ButtonPressed", Start: 0, Length: 1},
			},
			TargetType: reflect.TypeOf(Port6Payload{}),
			Features:   []decoder.Feature{decoder.FeatureButton},
		}, nil
	case 7:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Timestamp", Start: 0, Length: 4, Transform: timestamp},
				{Name: "DutyCycle", Start: 4, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 4, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 4, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 4, Length: 1, Transform: moving},
				{Name: "Mac1", Start: 5, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 11, Length: 1},
				{Name: "Mac2", Start: 12, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 18, Length: 1, Optional: true},
				{Name: "Mac3", Start: 19, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 25, Length: 1, Optional: true},
				{Name: "Mac4", Start: 26, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 32, Length: 1, Optional: true},
				{Name: "Mac5", Start: 33, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 39, Length: 1, Optional: true},
				{Name: "Mac6", Start: 40, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 46, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port7Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureWiFi},
		}, nil
	case 8:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "ScanInterval", Start: 0, Length: 2},
				{Name: "ScanTime", Start: 2, Length: 1},
				{Name: "MaxBeacons", Start: 3, Length: 1},
				{Name: "MinRssiValue", Start: 4, Length: 1},
				{Name: "AdvertisingFilter", Start: 5, Length: 10},
				{Name: "AccelerometerTriggerHoldTimer", Start: 15, Length: 2},
				{Name: "AccelerometerThreshold", Start: 17, Length: 2},
				{Name: "ScanMode", Start: 19, Length: 1},
				{Name: "BLECurrentConfigurationUplinkInterval", Start: 20, Length: 2},
			},
			TargetType: reflect.TypeOf(Port8Payload{}),
			Features:   []decoder.Feature{},
		}, nil
	case 10:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: battery},
				{Name: "TTF", Start: 17, Length: 1, Optional: true, Transform: ttf},
				{Name: "PDOP", Start: 18, Length: 1, Optional: true, Transform: pdop},
				{Name: "Satellites", Start: 19, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port10Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureBattery},
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
	case 50:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: battery},
				{Name: "TTF", Start: 17, Length: 1, Transform: ttf},
				{Name: "Mac1", Start: 18, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 24, Length: 1},
				{Name: "Mac2", Start: 25, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 31, Length: 1, Optional: true},
				{Name: "Mac3", Start: 32, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 38, Length: 1, Optional: true},
				{Name: "Mac4", Start: 39, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 45, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port50Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureBattery, decoder.FeatureWiFi},
		}, nil
	case 51:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "DutyCycle", Start: 0, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 0, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 0, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 0, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 1, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 5, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 9, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 11, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 15, Length: 2, Transform: battery},
				{Name: "TTF", Start: 17, Length: 1, Transform: ttf},
				{Name: "PDOP", Start: 18, Length: 1, Transform: pdop},
				{Name: "Satellites", Start: 19, Length: 1},
				{Name: "Mac1", Start: 20, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 26, Length: 1},
				{Name: "Mac2", Start: 27, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 33, Length: 1, Optional: true},
				{Name: "Mac3", Start: 34, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 40, Length: 1, Optional: true},
				{Name: "Mac4", Start: 41, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 47, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port51Payload{}),
			Features:   []decoder.Feature{decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureBattery, decoder.FeatureWiFi},
		}, nil
	case 105:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "Timestamp", Start: 2, Length: 4, Transform: timestamp},
				{Name: "DutyCycle", Start: 6, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 6, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 6, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 6, Length: 1, Transform: moving},
				{Name: "Mac1", Start: 7, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 13, Length: 1},
				{Name: "Mac2", Start: 14, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 20, Length: 1, Optional: true},
				{Name: "Mac3", Start: 21, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 27, Length: 1, Optional: true},
				{Name: "Mac4", Start: 28, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 34, Length: 1, Optional: true},
				{Name: "Mac5", Start: 35, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 41, Length: 1, Optional: true},
				{Name: "Mac6", Start: 42, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 48, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port105Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBuffered, decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureWiFi},
		}, nil
	case 110:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "DutyCycle", Start: 2, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 2, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 2, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 2, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 3, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 7, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 11, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 13, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 17, Length: 2, Transform: battery},
				{Name: "TTF", Start: 19, Length: 1, Optional: true, Transform: ttf},
				{Name: "PDOP", Start: 20, Length: 1, Optional: true, Transform: pdop},
				{Name: "Satellites", Start: 21, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port110Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBuffered, decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureBattery},
		}, nil
	case 150:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "DutyCycle", Start: 2, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 2, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 2, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 2, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 3, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 7, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 11, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 13, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 17, Length: 2, Transform: battery},
				{Name: "TTF", Start: 19, Length: 1, Transform: ttf},
				{Name: "Mac1", Start: 20, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 26, Length: 1},
				{Name: "Mac2", Start: 27, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 33, Length: 1, Optional: true},
				{Name: "Mac3", Start: 34, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 40, Length: 1, Optional: true},
				{Name: "Mac4", Start: 41, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 47, Length: 1, Optional: true},
				{Name: "Mac5", Start: 48, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi5", Start: 54, Length: 1, Optional: true},
				{Name: "Mac6", Start: 55, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi6", Start: 61, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port150Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBuffered, decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureBattery, decoder.FeatureWiFi},
		}, nil
	case 151:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BufferLevel", Start: 0, Length: 2},
				{Name: "DutyCycle", Start: 2, Length: 1, Transform: dutyCycle},
				{Name: "ConfigId", Start: 2, Length: 1, Transform: configId},
				{Name: "ConfigChange", Start: 2, Length: 1, Transform: configSuccess},
				{Name: "Moving", Start: 2, Length: 1, Transform: moving},
				{Name: "Latitude", Start: 3, Length: 4, Transform: latitude},
				{Name: "Longitude", Start: 7, Length: 4, Transform: longitude},
				{Name: "Altitude", Start: 11, Length: 2, Transform: altitude},
				{Name: "Timestamp", Start: 13, Length: 4, Transform: timestamp},
				{Name: "Battery", Start: 17, Length: 2, Transform: battery},
				{Name: "TTF", Start: 19, Length: 1, Transform: ttf},
				{Name: "PDOP", Start: 20, Length: 1, Transform: pdop},
				{Name: "Satellites", Start: 21, Length: 1},
				{Name: "Mac1", Start: 22, Length: 6, Hex: true},
				{Name: "Rssi1", Start: 28, Length: 1},
				{Name: "Mac2", Start: 29, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi2", Start: 35, Length: 1, Optional: true},
				{Name: "Mac3", Start: 36, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi3", Start: 42, Length: 1, Optional: true},
				{Name: "Mac4", Start: 43, Length: 6, Optional: true, Hex: true},
				{Name: "Rssi4", Start: 49, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(Port151Payload{}),
			Features:   []decoder.Feature{decoder.FeatureBuffered, decoder.FeatureDutyCycle, decoder.FeatureConfigChange, decoder.FeatureMoving, decoder.FeatureGNSS, decoder.FeatureBattery, decoder.FeatureWiFi},
		}, nil
	case 198:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Reason", Start: 0, Length: 1},
				{Name: "Line", Start: 1, Length: -1, Optional: true, Transform: func(v any) any {
					return stacktrace(v.([]byte), 0)
				}},
				{Name: "File", Start: 1, Length: -1, Optional: true, Transform: func(v any) any {
					return stacktrace(v.([]byte), 1)
				}},
				{Name: "Function", Start: 1, Length: -1, Optional: true, Transform: func(v any) any {
					return stacktrace(v.([]byte), 2)
				}},
			},
			TargetType: reflect.TypeOf(Port198Payload{}),
			Features:   []decoder.Feature{decoder.FeatureResetReason},
		}, nil
	case 199:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Constant", Start: 0, Length: 7, Hex: true},
				{Name: "Sequence", Start: 7, Length: 4},
				{Name: "Number", Start: 11, Length: 3},
				{Name: "Id", Start: 14, Length: 1},
			},
			TargetType: reflect.TypeOf(Port199Payload{}),
			Features:   []decoder.Feature{},
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func (t TagSLv1Decoder) Decode(ctx context.Context, data string, port uint8) (*decoder.DecodedUplink, error) {
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

func timestamp(v any) any {
	return time.Unix(int64(common.BytesToUint32(v.([]byte))), 0).UTC()
}

func battery(v any) any {
	return float64(common.BytesToUint16(v.([]byte))) / 1000
}

func ttf(v any) any {
	return time.Duration(int64(common.BytesToUint8(v.([]byte)))) * time.Second
}

func pdop(v any) any {
	return float64(common.BytesToUint8(v.([]byte))) / 2
}

func stacktrace(v []byte, i int) any {
	frags := strings.Split(string(v), ":")
	if len(frags) > i {
		return frags[i]
	}
	return nil
}
