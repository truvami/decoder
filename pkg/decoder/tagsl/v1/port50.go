package tagsl

import "time"

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
// | 18-23| 6    | MAC1                                      | 6 x uint8              |
// | 24   | 1    | RSSI1                                     | int8                   |
// | …    |      |                                           |                        |
// |      | 6    | MACN                                      | 6 x uint8              |
// |      | 1    | RSSIN                                     | int8                   |
// +------+------+-------------------------------------------+------------------------+

// Timestamp for the Wi-Fi scanning is TSGNSS – TTF + 10 seconds.
type Port50Payload struct {
	Latitude  float64   `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude float64   `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude  float64   `json:"altitude"`
	Timestamp time.Time `json:"timestamp"`
	Battery   float64   `json:"battery" validate:"gte=1,lte=5"`
	TTF       uint16    `json:"ttf"`
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
	Mac7      string    `json:"mac7"`
	Rssi7     int8      `json:"rssi7"`
}
