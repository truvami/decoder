package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+--------+
// | Byte | Size | Description                               | Format |
// +------+------+-------------------------------------------+--------+
// | 0    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8  |
// +------+------+-------------------------------------------+--------+

type Port2Payload struct {
	Moving    bool `json:"moving"`
	DutyCycle bool `json:"dutyCycle"`
}

var _ decoder.UplinkFeatureBase = &Port2Payload{}

func (p Port2Payload) GetTimestamp() *time.Time {
	return nil
}
