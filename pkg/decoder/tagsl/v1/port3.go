package tagsl

import (
	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+----------------------------------------+-----------+
// | Byte | Size | Description                            | Format    |
// +------+------+----------------------------------------+-----------+
// | 0    | 2    | Scan pointer                           | uint16    |
// | 2    | 1    | Total messages                         | uint8     |
// | 3    | 1    | Message                                | uint8     |
// | 4    | 6    | Mac 1                                  | uint8[6]  |
// | 10   | 1    | Rssi 1                                 | int8      |
// | 11   | 6    | Mac 2                                  | uint8[6]  |
// | 17   | 1    | Rssi 2                                 | int8      |
// | 18   | 6    | Mac 3                                  | uint8[6]  |
// | 24   | 1    | Rssi 3                                 | int8      |
// | 25   | 6    | Mac 4                                  | uint8[6]  |
// | 31   | 1    | Rssi 4                                 | int8      |
// | 32   | 6    | Mac 5                                  | uint8[6]  |
// | 38   | 1    | Rssi 5                                 | int8      |
// | 39   | 6    | Mac 6                                  | uint8[6]  |
// | 45   | 1    | Rssi 6                                 | int8      |
// +------+------+----------------------------------------+-----------+

type Port3Payload struct {
	ScanPointer    uint16  `json:"scanPointer"`
	TotalMessages  uint8   `json:"totalMessages"`
	CurrentMessage uint8   `json:"currentMessage"`
	Mac1           string  `json:"mac1"`
	Rssi1          int8    `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac2           *string `json:"mac2"`
	Rssi2          *int8   `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac3           *string `json:"mac3"`
	Rssi3          *int8   `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac4           *string `json:"mac4"`
	Rssi4          *int8   `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac5           *string `json:"mac5"`
	Rssi5          *int8   `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac6           *string `json:"mac6"`
	Rssi6          *int8   `json:"rssi6" validate:"gte=-120,lte=-20"`
}

var _ decoder.UplinkFeatureWiFi = &Port3Payload{}

func (p Port3Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" && p.Rssi1 != 0 {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: &p.Rssi1,
		})
	}

	if p.Mac2 != nil && p.Rssi2 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != nil && p.Rssi3 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != nil && p.Rssi4 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != nil && p.Rssi5 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	if p.Mac6 != nil && p.Rssi6 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac6,
			RSSI: p.Rssi6,
		})
	}

	return accessPoints
}
