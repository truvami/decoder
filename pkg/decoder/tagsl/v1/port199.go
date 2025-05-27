package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port199Payload struct {
	Constant string `json:"constant"`
	Sequence uint32 `json:"sequence"`
	Number   uint32 `json:"number"`
	Id       uint32 `json:"id"`
}

var _ decoder.UplinkFeatureBase = &Port199Payload{}

func (p Port199Payload) GetTimestamp() *time.Time {
	return nil
}
