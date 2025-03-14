package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +-------+------+-------------------------------------------+-----------+
// | Byte  | Size | Description                               | Format    |
// +-------+------+-------------------------------------------+-----------+
// | 0     | 1    | Buffer level                              | uint16    |
// | 2     | 4    | Unix timestamp                            | uint32    |
// | 7     | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8     |
// | 7     | 6    | MAC1                                      | 6 x uint8 |
// | 13    | 1    | RSSI1                                     | int8      |
// | …     |      |                                           |           |
// |       | 6    | MACN                                      | 6 x uint8 |
// |       | 1    | RSSIN                                     | int8      |
// +-------+------+-------------------------------------------+-----------+

type Port105Payload struct {
	Moving      bool      `json:"moving"`
	DutyCycle   bool      `json:"dutyCycle"`
	BufferLevel uint16    `json:"bufferLevel"`
	Timestamp   time.Time `json:"timestamp"`
	Mac1        string    `json:"mac1"`
	Rssi1       int8      `json:"rssi1"`
	Mac2        string    `json:"mac2"`
	Rssi2       int8      `json:"rssi2"`
	Mac3        string    `json:"mac3"`
	Rssi3       int8      `json:"rssi3"`
	Mac4        string    `json:"mac4"`
	Rssi4       int8      `json:"rssi4"`
	Mac5        string    `json:"mac5"`
	Rssi5       int8      `json:"rssi5"`
	Mac6        string    `json:"mac6"`
	Rssi6       int8      `json:"rssi6"`
}

var _ decoder.UplinkFeatureBase = &Port105Payload{}
var _ decoder.UplinkFeatureWiFi = &Port105Payload{}
var _ decoder.UplinkFeatureBuffered = &Port105Payload{}
var _ decoder.UplinkFeatureMoving = &Port105Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port105Payload{}

func (p Port105Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port105Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	if p.Mac6 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac6,
			RSSI: p.Rssi6,
		})
	}

	return accessPoints
}

func (p Port105Payload) GetBufferLevel() uint16 {
	return p.BufferLevel
}

func (p Port105Payload) IsMoving() bool {
	return p.Moving
}

func (p Port105Payload) IsDutyCycle() bool {
	return p.DutyCycle
}
