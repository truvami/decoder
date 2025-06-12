package nomadxs

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
	"github.com/truvami/decoder/pkg/encoder"
)

type NomadXSv1Encoder struct{}

func NewNomadXSv1Encoder() encoder.Encoder {
	return &NomadXSv1Encoder{}
}

func (n NomadXSv1Encoder) Encode(data any, port uint8) (any, error) {
	config, err := n.getConfig(port)
	if err != nil {
		return nil, err
	}

	payload, err := common.Encode(data, config)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (n NomadXSv1Encoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
	case 1:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
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
			},
			TargetType: reflect.TypeOf(nomadxs.Port1Payload{}),
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
			TargetType: reflect.TypeOf(nomadxs.Port4Payload{}),
		}, nil
	case 15:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "LowBattery", Start: 0, Length: 1, Transform: lowBattery},
				{Name: "Battery", Start: 1, Length: 2, Transform: battery},
			},
			TargetType: reflect.TypeOf(nomadxs.Port15Payload{}),
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func moving(v any) any {
	return common.BoolToBytes(common.BytesToBool(v.([]byte)), 0)
}

func latitude(v any) any {
	return common.IntToBytes(int64(common.BytesToFloat64(v.([]byte))*1000000), 4)
}

func longitude(v any) any {
	return common.IntToBytes(int64(common.BytesToFloat64(v.([]byte))*1000000), 4)
}

func altitude(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*10), 2)
}

func ttf(v any) any {
	return common.UintToBytes(uint64(common.BytesToInt64(v.([]byte))/1000000000), 1)
}

func temperature(v any) any {
	return common.IntToBytes(int64((common.BytesToFloat32(v.([]byte)))*100), 2)
}

func pressure(v any) any {
	return common.UintToBytes(uint64((common.BytesToFloat32(v.([]byte)))*10), 2)
}

func lowBattery(v any) any {
	return common.BoolToBytes(common.BytesToBool(v.([]byte)), 0)
}

func battery(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat64(v.([]byte))*1000), 2)
}
