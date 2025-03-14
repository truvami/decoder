package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port152Payload struct {
	NewRotationState  uint8   `json:"newRotationState"`
	OldRotationState  uint8   `json:"oldRotationState"`
	Timestamp         uint32  `json:"timestamp"`
	NumberOfRotations float64 `json:"numberOfRotations"`
	ElapsedSeconds    uint32  `json:"elapsedSeconds"`
}

var _ decoder.UplinkFeatureBase = &Port152Payload{}

// GetTimestamp implements decoder.UplinkFeatureBase.
func (p Port152Payload) GetTimestamp() *time.Time {
	return nil
}
