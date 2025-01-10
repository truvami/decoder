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
// | 17   | 1    | TTF (TimeToFix)                           | uint8, s               |
// | 18   | 1    | PDOP  		                                 | uint8, 1/2 meter       |
// | 19   | 1    | Number of satellites                      | uint8, 		            |
// +------+------+-------------------------------------------+------------------------+

type Port10Payload struct {
	Latitude   float64   `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude  float64   `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude   float64   `json:"altitude" validate:"gte=0,lte=20000"`
	Timestamp  time.Time `json:"timestamp"`
	Battery    float64   `json:"battery" validate:"gte=1,lte=5"`
	TTF        float64   `json:"ttf"`
	PDOP       float64   `json:"pdop"`
	Satellites uint8     `json:"satellites" validate:"gte=3,lte=27"`
}
