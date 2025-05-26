package smartlabel

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | unknown tag probably bit field                | byte       |
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
// | 36   | 1    | rssi signal 6                                 | int8       |
// | 37   | 6    | mac address signal 6                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+

type Port197Payload struct {
	Tag   byte    `json:"tag"`
	Mac1  string  `json:"mac1"`
	Rssi1 int8    `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac2  *string `json:"mac2"`
	Rssi2 *int8   `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac3  *string `json:"mac3"`
	Rssi3 *int8   `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac4  *string `json:"mac4"`
	Rssi4 *int8   `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac5  *string `json:"mac5"`
	Rssi5 *int8   `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac6  *string `json:"mac6"`
	Rssi6 *int8   `json:"rssi6" validate:"gte=-120,lte=-20"`
}

var _ decoder.UplinkFeatureBase = &Port197Payload{}
var _ decoder.UplinkFeatureWiFi = &Port197Payload{}

func (p Port197Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port197Payload) GetAccessPoints() []decoder.AccessPoint {
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
