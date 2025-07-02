package smartlabel

import (
	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+----------------+
// | Byte | Size | Description                               | Format         |
// +------+------+-------------------------------------------+----------------+
// | 0    | 2    | battery output voltage                    | uint16, mV     |
// | 2    | 2    | photovoltaic input voltage                | uint16, mV     |
// | 4    | 2    | temperature                               | uint16, 0.01 C |
// | 6    | 1    | relative humidity                         | uint8, 0.5%    |
// +------+------+-------------------------------------------+----------------+

type Port11Payload struct {
	BatteryVoltage      float32 `json:"batteryVoltage" validate:"gte=1,lte=5"`
	PhotovoltaicVoltage float32 `json:"photovoltaicVoltage" validate:"gte=0,lte=5"`
	Temperature         float32 `json:"temperature" validate:"gte=-20,lte=60"`
	Humidity            float32 `json:"humidity" validate:"gte=5,lte=95"`
}

var _ decoder.UplinkFeatureBattery = &Port11Payload{}
var _ decoder.UplinkFeaturePhotovoltaic = &Port11Payload{}
var _ decoder.UplinkFeatureTemperature = &Port11Payload{}
var _ decoder.UplinkFeatureHumidity = &Port11Payload{}

func (p Port11Payload) GetBatteryVoltage() float64 {
	return float64(p.BatteryVoltage)
}

func (p Port11Payload) GetLowBattery() *bool {
	return nil
}

func (p Port11Payload) GetPhotovoltaicVoltage() float32 {
	return p.PhotovoltaicVoltage
}

func (p Port11Payload) GetTemperature() float32 {
	return p.Temperature
}

func (p Port11Payload) GetHumidity() float32 {
	return p.Humidity
}
