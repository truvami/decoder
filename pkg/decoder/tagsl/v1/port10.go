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
// | 17   | 1    | TTF (TimeToFix)                           | uint8, s               |
// | 18   | 1    | PDOP  		                                 | uint8, 1/2 meter       |
// | 19   | 1    | Number of satellites                      | uint8, 		            |
// +------+------+-------------------------------------------+------------------------+

type Port10Payload struct {
	Latitude   float64   `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude  float64   `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude   float64   `json:"altitude"`
	Timestamp  time.Time `json:"timestamp"`
	Battery    float64   `json:"battery" validate:"gte=1,lte=5"`
	TTF        float64   `json:"ttf"`
	PDOP       float64   `json:"pdop"`
	Satellites uint8     `json:"satellites" validate:"gte=3,lte=27"`
}

var _ common.Position = &Port10Payload{}

func (p Port10Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port10Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port10Payload) GetAltitude() *float64 {
	return &p.Altitude
}

func (p Port10Payload) GetSource() common.PositionSource {
	return common.PositionSource_GNSS
}

func (p Port10Payload) GetCapturedAt() *time.Time {
	return &p.Timestamp
}

func (p Port10Payload) GetBuffered() bool {
	return false
}
