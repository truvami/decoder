package smartlabel

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 2    | battery output voltage                        | uint16, mV |
// | 2    | 2    | photovoltaic input voltage                    | uint16, mV |
// +------+------+-----------------------------------------------+------------+

type Port1Payload struct {
	BatteryVoltage      float32 `json:"batteryVoltage" validate:"gte=1,lte=5"`
	PhotovoltaicVoltage float32 `json:"photovoltaicVoltage" validate:"gte=0,lte=5"`
}

var _ decoder.UplinkFeatureBase = &Port1Payload{}
var _ decoder.UplinkFeatureBattery = &Port1Payload{}
var _ decoder.UplinkFeaturePhotovoltaic = &Port1Payload{}

func (p Port1Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port1Payload) GetBatteryVoltage() float64 {
	return float64(p.BatteryVoltage)
}

func (p Port1Payload) GetPhotovoltaicVoltage() float32 {
	return p.PhotovoltaicVoltage
}
