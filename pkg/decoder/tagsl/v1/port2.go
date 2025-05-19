package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+--------+
// | Byte | Size | Description                               | Format |
// +------+------+-------------------------------------------+--------+
// | 0    | 1    | Duty cycle flag                           | uint1  |
// | 0    | 1    | Config id                                 | uint4  |
// | 0    | 1    | Config change flag                        | uint1  |
// | 0    | 1    | Reserved                                  | uint1  |
// | 0    | 1    | Moving flag                               | uint1  |
// +------+------+-------------------------------------------+--------+

type Port2Payload struct {
	DutyCycle    bool  `json:"dutyCycle"`
	ConfigId     uint8 `json:"configId" validate:"gte=0,lte=15"`
	ConfigChange bool  `json:"configChange"`
	Moving       bool  `json:"moving"`
}

var _ decoder.UplinkFeatureBase = &Port2Payload{}
var _ decoder.UplinkFeatureMoving = &Port2Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port2Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port2Payload{}

func (p Port2Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port2Payload) IsMoving() bool {
	return p.Moving
}

func (p Port2Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port2Payload) GetConfigId() *uint8 {
	return &p.ConfigId
}

func (p Port2Payload) GetConfigChange() bool {
	return p.ConfigChange
}
