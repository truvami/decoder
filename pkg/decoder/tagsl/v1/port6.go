package tagsl

import (
	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+----------------------------------------+--------+
// | Byte | Size | Description                            | Format |
// +------+------+----------------------------------------+--------+
// | 0    | 1    | In case of a button-press 0x01 is sent | uint8  |
// +------+------+----------------------------------------+--------+

type Port6Payload struct {
	ButtonPressed bool `json:"buttonPressed"`
}

var _ decoder.UplinkFeatureButton = &Port6Payload{}

func (p Port6Payload) GetPressed() bool {
	return p.ButtonPressed
}
