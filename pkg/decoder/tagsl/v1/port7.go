package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+-----------+
// | Byte | Size | Description                               | Format    |
// +------+------+-------------------------------------------+-----------+
// | 0    | 4    | Timestamp                                 | uint32    |
// | 4    | 1    | Status Byte                               | uint8     |
// | 5    | 6    | MAC1                                      | uint8[6]  |
// | 11   | 1    | RSSI1                                     | int8      |
// | 12   | 6    | MAC2                                      | uint8[6]  |
// | 18   | 1    | RSSI2                                     | int8      |
// | 19   | 6    | MAC3                                      | uint8[6]  |
// | 25   | 1    | RSSI3                                     | int8      |
// | 26   | 6    | MAC4                                      | uint8[6]  |
// | 32   | 1    | RSSI4                                     | int8      |
// | 33   | 6    | MAC5                                      | uint8[6]  |
// | 39   | 1    | RSSI5                                     | int8      |
// | 40   | 6    | MAC6                                      | uint8[6]  |
// | 46   | 1    | RSSI6                                     | int8      |
// +------+------+-------------------------------------------+-----------+

type Port7Payload struct {
	Moving    bool      `json:"moving"`
	DutyCycle bool      `json:"dutyCycle"`
	Timestamp time.Time `json:"timestamp"`
	Mac1      string    `json:"mac1"`
	Rssi1     int8      `json:"rssi1"`
	Mac2      string    `json:"mac2"`
	Rssi2     int8      `json:"rssi2"`
	Mac3      string    `json:"mac3"`
	Rssi3     int8      `json:"rssi3"`
	Mac4      string    `json:"mac4"`
	Rssi4     int8      `json:"rssi4"`
	Mac5      string    `json:"mac5"`
	Rssi5     int8      `json:"rssi5"`
	Mac6      string    `json:"mac6"`
	Rssi6     int8      `json:"rssi6"`
}

var _ decoder.UplinkFeatureBase = &Port7Payload{}
var _ decoder.UplinkFeatureWiFi = &Port7Payload{}
var _ decoder.UplinkFeatureMoving = &Port7Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port7Payload{}

func (p Port7Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port7Payload) GetAccessPoints() []decoder.AccessPoint {
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

func (p Port7Payload) IsMoving() bool {
	return p.Moving
}

func (p Port7Payload) IsDutyCycle() bool {
	return p.DutyCycle
}
