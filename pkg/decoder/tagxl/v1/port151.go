package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port151Payload struct {
	Battery *float32 `json:"battery" validate:"gte=1,lte=5"`
}

var _ decoder.UplinkFeatureBase = &Port151Payload{}
var _ decoder.UplinkFeatureBattery = &Port151Payload{}

func (p Port151Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port151Payload) GetBatteryVoltage() float64 {
	return float64(*p.Battery)
}

func (p Port151Payload) GetLowBattery() *bool {
	return nil
}
