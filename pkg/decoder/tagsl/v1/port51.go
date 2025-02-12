package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/common"
)

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 1-4  | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 5-8  | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 9-10 | 2    | Altitude                                  | uint16, 1/10 meter     |
// | 11-14| 4    | Unix timestamp                            | uint32                 |
// | 15-16| 2    | Battery voltage                           | uint16, mV             |
// | 17   | 1    | TTF                                       | uint8                  |
// | 18   | 1    | PDOP  		                                 | uint8, 1/2 meter       |
// | 19   | 1    | Number of satellites                      | uint8, 		            |
// | 20-25| 6    | MAC1                                      | 6 x uint8              |
// | 26   | 1    | RSSI1                                     | int8                   |
// | …    |      |                                           |                        |
// |      | 6    | MACN                                      | 6 x uint8              |
// |      | 1    | RSSIN                                     | int8                   |
// +------+------+-------------------------------------------+------------------------+

// Timestamp for the Wi-Fi scanning is TSGNSS – TTF + 10 seconds.
type Port51Payload struct {
	Latitude   float64   `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude  float64   `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude   float64   `json:"altitude"`
	Timestamp  time.Time `json:"timestamp"`
	Battery    float64   `json:"battery" validate:"gte=1,lte=5"`
	TTF        uint16    `json:"ttf"`
	PDOP       float64   `json:"pdop"`
	Satellites uint8     `json:"satellites" validate:"gte=3,lte=27"`
	Mac1       string    `json:"mac1"`
	Rssi1      int8      `json:"rssi1"`
	Mac2       string    `json:"mac2"`
	Rssi2      int8      `json:"rssi2"`
	Mac3       string    `json:"mac3"`
	Rssi3      int8      `json:"rssi3"`
	Mac4       string    `json:"mac4"`
	Rssi4      int8      `json:"rssi4"`
	Mac5       string    `json:"mac5"`
	Rssi5      int8      `json:"rssi5"`
	Mac6       string    `json:"mac6"`
	Rssi6      int8      `json:"rssi6"`
	Mac7       string    `json:"mac7"`
	Rssi7      int8      `json:"rssi7"`
}

var _ common.Position = &Port51Payload{}

func (p Port51Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port51Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port51Payload) GetAltitude() *float64 {
	return &p.Altitude
}

func (p Port51Payload) GetSource() common.PositionSource {
	return common.PositionSource_GNSS
}

func (p Port51Payload) GetCapturedAt() *time.Time {
	return &p.Timestamp
}

func (p Port51Payload) GetBuffered() bool {
	return false
}

var _ common.WifiLocation = &Port51Payload{}

func (p Port51Payload) GetAccessPoints() []common.WifiAccessPoint {
	var accessPoints []common.WifiAccessPoint

	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac1, p.Rssi1)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac2, p.Rssi2)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac3, p.Rssi3)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac4, p.Rssi4)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac5, p.Rssi5)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac6, p.Rssi6)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac7, p.Rssi7)

	return accessPoints
}
