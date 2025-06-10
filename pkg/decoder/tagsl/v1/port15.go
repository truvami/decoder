package tagsl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+---------------------------------------------+------------+
// | Byte | Size | Description                                 | Format     |
// +------+------+---------------------------------------------+------------+
// | 0    | 1    | Duty cycle flag                             | uint1      |
// | 0    | 1    | Config id                                   | uint4      |
// | 0    | 1    | Config change flag                          | uint1      |
// | 0    | 1    | Reserved                                    | uint1      |
// | 0    | 1    | Low battery flag                            | uint1      |
// | 1    | 2    | Battery voltage                             | uint16, mV |
// +------+------+---------------------------------------------+------------+

type Port15Payload struct {
	DutyCycle    bool    `json:"dutyCycle"`
	ConfigId     uint8   `json:"configId" validate:"gte=0,lte=15"`
	ConfigChange bool    `json:"configChange"`
	LowBattery   bool    `json:"lowBattery"`
	Battery      float64 `json:"battery" validate:"gte=1,lte=5"`
}

func (p Port15Payload) MarshalJSON() ([]byte, error) {
	type Alias Port15Payload
	return json.Marshal(&struct {
		*Alias
		Battery string `json:"battery"`
	}{
		Alias:   (*Alias)(&p),
		Battery: fmt.Sprintf("%.3fv", p.Battery),
	})
}

var _ decoder.UplinkFeatureBase = &Port15Payload{}
var _ decoder.UplinkFeatureBattery = &Port15Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port15Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port15Payload{}

func (p Port15Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port15Payload) GetBatteryVoltage() float64 {
	return p.Battery
}

func (p Port15Payload) GetLowBattery() *bool {
	return &p.LowBattery
}

func (p Port15Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port15Payload) GetConfigId() *uint8 {
	return &p.ConfigId
}

func (p Port15Payload) GetConfigChange() bool {
	return p.ConfigChange
}
