package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+-----------+
// | Byte | Size | Description                               | Format    |
// +------+------+-------------------------------------------+-----------+
// | 0    | 1    | Duty cycle flag                           | uint1     |
// | 0    | 1    | Config id                                 | uint4     |
// | 0    | 1    | Config change flag                        | uint1     |
// | 0    | 1    | Reserved                                  | uint1     |
// | 0    | 1    | Moving flag                               | uint1     |
// | 1    | 6    | Mac 1                                     | uint8[6]  |
// | 7    | 1    | Rssi 1                                    | int8      |
// | 8    | 6    | Mac 2                                     | uint8[6]  |
// | 14   | 1    | Rssi 2                                    | int8      |
// | 15   | 6    | Mac 3                                     | uint8[6]  |
// | 21   | 1    | Rssi 3                                    | int8      |
// | 22   | 6    | Mac 4                                     | uint8[6]  |
// | 28   | 1    | Rssi 4                                    | int8      |
// | 29   | 6    | Mac 5                                     | uint8[6]  |
// | 35   | 1    | Rssi 5                                    | int8      |
// | 36   | 6    | Mac 6                                     | uint8[6]  |
// | 42   | 1    | Rssi 6                                    | int8      |
// | 43   | 6    | Mac 7                                     | uint8[6]  |
// | 49   | 1    | Rssi 7                                    | int8      |
// +------+------+-------------------------------------------+-----------+

type Port5Payload struct {
	DutyCycle    bool    `json:"dutyCycle"`
	ConfigId     uint8   `json:"configId" validate:"gte=0,lte=15"`
	ConfigChange bool    `json:"configChange"`
	Moving       bool    `json:"moving"`
	Mac1         string  `json:"mac1"`
	Rssi1        int8    `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac2         *string `json:"mac2"`
	Rssi2        *int8   `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac3         *string `json:"mac3"`
	Rssi3        *int8   `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac4         *string `json:"mac4"`
	Rssi4        *int8   `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac5         *string `json:"mac5"`
	Rssi5        *int8   `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac6         *string `json:"mac6"`
	Rssi6        *int8   `json:"rssi6" validate:"gte=-120,lte=-20"`
	Mac7         *string `json:"mac7"`
	Rssi7        *int8   `json:"rssi7" validate:"gte=-120,lte=-20"`
}

var _ decoder.UplinkFeatureBase = &Port5Payload{}
var _ decoder.UplinkFeatureWiFi = &Port5Payload{}
var _ decoder.UplinkFeatureMoving = &Port5Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port5Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port5Payload{}

func (p Port5Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port5Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" && p.Rssi1 != 0 {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != nil && p.Rssi2 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac2,
			RSSI: *p.Rssi2,
		})
	}

	if p.Mac3 != nil && p.Rssi3 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac3,
			RSSI: *p.Rssi3,
		})
	}

	if p.Mac4 != nil && p.Rssi4 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac4,
			RSSI: *p.Rssi4,
		})
	}

	if p.Mac5 != nil && p.Rssi5 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac5,
			RSSI: *p.Rssi5,
		})
	}

	if p.Mac6 != nil && p.Rssi6 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac6,
			RSSI: *p.Rssi6,
		})
	}

	if p.Mac7 != nil && p.Rssi7 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac7,
			RSSI: *p.Rssi7,
		})
	}

	return accessPoints
}

func (p Port5Payload) IsMoving() bool {
	return p.Moving
}

func (p Port5Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port5Payload) GetConfigId() *uint8 {
	return &p.ConfigId
}

func (p Port5Payload) GetConfigChange() bool {
	return p.ConfigChange
}
