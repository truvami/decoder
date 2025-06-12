package smartlabel

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder/smartlabel/v1"
	"github.com/truvami/decoder/pkg/encoder"
)

type Smartlabelv1Encoder struct{}

func NewSmartlabelv1Encoder() encoder.Encoder {
	return &Smartlabelv1Encoder{}
}

func (s Smartlabelv1Encoder) Encode(data any, port uint8) (any, error) {
	config, err := s.getConfig(port)
	if err != nil {
		return nil, err
	}

	payload, err := common.Encode(data, config)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (t Smartlabelv1Encoder) getConfig(port uint8) (common.PayloadConfig, error) {
	switch port {
	case 1:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BatteryVoltage", Start: 0, Length: 2, Transform: battery},
				{Name: "PhotovoltaicVoltage", Start: 2, Length: 2, Transform: photovoltaic},
			},
			TargetType: reflect.TypeOf(smartlabel.Port1Payload{}),
		}, nil
	case 2:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "Temperature", Start: 0, Length: 2, Transform: temperature},
				{Name: "Humidity", Start: 2, Length: 1, Transform: humidity},
			},
			TargetType: reflect.TypeOf(smartlabel.Port2Payload{}),
		}, nil
	case 11:
		return common.PayloadConfig{
			Fields: []common.FieldConfig{
				{Name: "BatteryVoltage", Start: 0, Length: 2, Transform: battery},
				{Name: "PhotovoltaicVoltage", Start: 2, Length: 2, Transform: photovoltaic},
				{Name: "Temperature", Start: 4, Length: 2, Transform: temperature},
				{Name: "Humidity", Start: 6, Length: 1, Transform: humidity},
			},
			TargetType: reflect.TypeOf(smartlabel.Port11Payload{}),
		}, nil
	}

	return common.PayloadConfig{}, fmt.Errorf("%w: port %v not supported", common.ErrPortNotSupported, port)
}

func battery(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat32(v.([]byte))*1000), 2)
}

func photovoltaic(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat32(v.([]byte))*1000), 2)
}

func temperature(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat32(v.([]byte))*100), 2)
}

func humidity(v any) any {
	return common.UintToBytes(uint64(common.BytesToFloat32(v.([]byte))*2), 1)
}
