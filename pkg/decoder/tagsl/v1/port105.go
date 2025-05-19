package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+-----------+
// | Byte | Size | Description                               | Format    |
// +------+------+-------------------------------------------+-----------+
// | 0    | 2    | Buffer level                              | uint16    |
// | 2    | 4    | Unix timestamp                            | uint32    |
// | 6    | 1    | Duty cycle flag                           | uint1     |
// | 6    | 1    | Config change id                          | uint4     |
// | 6    | 1    | Config change success flag                | uint1     |
// | 6    | 1    | Reserved                                  | uint1     |
// | 6    | 1    | Moving flag                               | uint1     |
// | 7    | 6    | Mac 1                                     | uint8[6]  |
// | 13   | 1    | Rssi 1                                    | int8      |
// | 14   | 6    | Mac 2                                     | uint8[6]  |
// | 20   | 1    | Rssi 2                                    | int8      |
// | 21   | 6    | Mac 3                                     | uint8[6]  |
// | 27   | 1    | Rssi 3                                    | int8      |
// | 28   | 6    | Mac 4                                     | uint8[6]  |
// | 34   | 1    | Rssi 4                                    | int8      |
// | 35   | 6    | Mac 5                                     | uint8[6]  |
// | 41   | 1    | Rssi 5                                    | int8      |
// | 42   | 6    | Mac 6                                     | uint8[6]  |
// | 48   | 1    | Rssi 6                                    | int8      |
// +-------+------+-------------------------------------------+-----------+

type Port105Payload struct {
	BufferLevel         uint16    `json:"bufferLevel"`
	Timestamp           time.Time `json:"timestamp"`
	DutyCycle           bool      `json:"dutyCycle"`
	ConfigChangeId      uint8     `json:"configChangeId" validate:"gte=0,lte=15"`
	ConfigChangeSuccess bool      `json:"configChangeSuccess"`
	Moving              bool      `json:"moving"`
	Mac1                string    `json:"mac1"`
	Rssi1               int8      `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac2                string    `json:"mac2"`
	Rssi2               int8      `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac3                string    `json:"mac3"`
	Rssi3               int8      `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac4                string    `json:"mac4"`
	Rssi4               int8      `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac5                string    `json:"mac5"`
	Rssi5               int8      `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac6                string    `json:"mac6"`
	Rssi6               int8      `json:"rssi6" validate:"gte=-120,lte=-20"`
}

var _ decoder.UplinkFeatureBase = &Port105Payload{}
var _ decoder.UplinkFeatureWiFi = &Port105Payload{}
var _ decoder.UplinkFeatureBuffered = &Port105Payload{}
var _ decoder.UplinkFeatureMoving = &Port105Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port105Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port105Payload{}

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

func (p Port105Payload) GetConfigId() *uint8 {
	return &p.ConfigChangeId
}

func (p Port105Payload) GetConfigChange() bool {
	return p.ConfigChangeSuccess
}
