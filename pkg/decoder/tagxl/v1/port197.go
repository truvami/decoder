package tagxl

import (
	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | version                                       | byte       |
// | 1    | 1    | rssi signal 1                                 | int8       |
// | 2    | 6    | mac address signal 1                          | byte[6]    |
// | 8    | 1    | rssi signal 2                                 | int8       |
// | 9    | 6    | mac address signal 2                          | byte[6]    |
// | 15   | 1    | rssi signal 3                                 | int8       |
// | 16   | 6    | mac address signal 3                          | byte[6]    |
// | 22   | 1    | rssi signal 4                                 | int8       |
// | 23   | 6    | mac address signal 4                          | byte[6]    |
// | 29   | 1    | rssi signal 5                                 | int8       |
// | 30   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+

type Port197Payload struct {
	Rssi1 *int8   `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac1  string  `json:"mac1"`
	Rssi2 *int8   `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac2  *string `json:"mac2"`
	Rssi3 *int8   `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac3  *string `json:"mac3"`
	Rssi4 *int8   `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac4  *string `json:"mac4"`
	Rssi5 *int8   `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac5  *string `json:"mac5"`
}

var _ decoder.UplinkFeatureWiFi = &Port197Payload{}

func (p Port197Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	return accessPoints
}
