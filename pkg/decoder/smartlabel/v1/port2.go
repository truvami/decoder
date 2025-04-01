package smartlabel

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+----------------+
// | Byte | Size | Description                               | Format         |
// +------+------+-------------------------------------------+----------------+
// | 0    | 2    | temperature                               | uint16, 0.01 C |
// | 2    | 1    | relative humidity                         | uint8, 0.5%    |
// +------+------+-------------------------------------------+----------------+

type Port2Payload struct {
	Temperature float32 `json:"temperature" validate:"gte=-20,lte=60"`
	Humidity    float32 `json:"humidity" validate:"gte=5,lte=95"`
}

var _ decoder.UplinkFeatureBase = &Port2Payload{}
var _ decoder.UplinkFeatureTemperature = &Port2Payload{}
var _ decoder.UplinkFeatureHumidity = &Port2Payload{}

func (p Port2Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port2Payload) GetTemperature() float32 {
	return p.Temperature
}

func (p Port2Payload) GetHumidity() float32 {
	return p.Humidity
}
