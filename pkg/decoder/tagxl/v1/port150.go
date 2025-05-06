package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port150Payload struct {
	Timestamp time.Time `json:"timestamp"`
}

var _ decoder.UplinkFeatureBase = &Port150Payload{}

func (p Port150Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}
