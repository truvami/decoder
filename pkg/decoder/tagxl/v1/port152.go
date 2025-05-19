package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// Version: 1
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | version                                       | uint8      |
// | 1    | 1    | reserved                                      | uint8      |
// | 2    | 1    | new rotation state                            | uint4      |
// | 2    | 1    | old rotation state                            | uint4      |
// | 3    | 4    | timestamp in seconds since epoch              | uint32     |
// | 7    | 2    | number of rotations since last rotation       | uint32     |
// | 9    | 4    | elapsed seconds since last rotation           | uint32     |
// +------+------+-----------------------------------------------+------------+
//
// Version: 2
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | version                                       | uint8      |
// | 1    | 1    | reserved                                      | uint8      |
// | 2    | 1    | sequence number                               | uint8      |
// | 3    | 1    | new rotation state                            | uint4      |
// | 3    | 1    | old rotation state                            | uint4      |
// | 4    | 4    | timestamp in seconds since epoch              | uint32     |
// | 8    | 2    | number of rotations since last rotation       | uint32     |
// | 10   | 4    | elapsed seconds since last rotation           | uint32     |
// +------+------+-----------------------------------------------+------------+

type Port152Payload struct {
	Version           uint8     `json:"version" validate:"gte=1,lte=2"`
	SequenceNumber    uint8     `json:"sequenceNumber" validate:"lte=255"`
	NewRotationState  uint8     `json:"newRotationState" validate:"lte=3"`
	OldRotationState  uint8     `json:"oldRotationState" validate:"lte=3"`
	Timestamp         time.Time `json:"timestamp"`
	NumberOfRotations float64   `json:"numberOfRotations" validate:"gte=0"`
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
func (p Port152Payload) GetSequenceNumber() uint {
	return uint(p.SequenceNumber)
}
