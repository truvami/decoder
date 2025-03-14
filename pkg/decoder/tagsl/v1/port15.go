package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+---------------------------------------------+------------+
// | Byte | Size | Description                                 | Format     |
// +------+------+---------------------------------------------+------------+
// | 0    | 1    | Status[6:2] + Low battery flag[0] (low = 1) | uint8      |
// | 1-2  | 2    | Battery voltage                             | uint16, mV |
// +------+------+---------------------------------------------+------------+

type Port15Payload struct {
	LowBattery bool    `json:"lowBattery"`
	Battery    float64 `json:"battery" validate:"gte=1,lte=5"`
}

var _ decoder.UplinkFeatureBase = &Port15Payload{}
var _ decoder.UpLinkFeatureBattery = &Port15Payload{}

func (p Port15Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port15Payload) GetBatteryVoltage() float64 {
	return p.Battery
}
