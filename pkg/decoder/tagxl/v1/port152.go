package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port152Payload struct {
	Version           uint8     `json:"version"`
	SequenceNumber    uint8     `json:"sequenceNumber"`
	NewRotationState  uint8     `json:"newRotationState"`
	OldRotationState  uint8     `json:"oldRotationState"`
	Timestamp         time.Time `json:"timestamp"`
	NumberOfRotations float64   `json:"numberOfRotations"`
	ElapsedSeconds    uint32    `json:"elapsedSeconds"`
}

var _ decoder.UplinkFeatureBase = &Port152Payload{}
var _ decoder.UplinkFeatureRotationState = &Port152Payload{}
var _ decoder.UplinkFeatureSequenceNumber = &Port152Payload{}

// GetTimestamp implements decoder.UplinkFeatureBase.
func (p Port152Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

// GetRotationState implements decoder.UplinkFeatureRotationState.
func (p Port152Payload) GetRotationState() decoder.RotationState {
	switch p.NewRotationState {
	case 1:
		return decoder.RotationStatePouring
	case 2:
		return decoder.RotationStateMixing
	case 3:
		return decoder.RotationStateError
	default:
		return decoder.RotationStateUndefined
	}
}

// GetSequenceNumber implements decoder.UplinkFeatureSequenceNumber.
func (p Port152Payload) GetSequenceNumber() int {
	return int(p.SequenceNumber)
}
