package tagsl

import (
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/common"
)

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 1-4  | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 5-8  | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 9-10 | 2    | Altitude                                  | uint16, 1/100 meter    |
// | 11   | 1    | Year                                      | uint8, year after 2000 |
// | 12   | 1    | Month                                     | uint8, [1..12]         |
// | 13   | 1    | Day                                       | uint8, [1..31]         |
// | 14   | 1    | Hour                                      | [0..23]                |
// | 15   | 1    | Minute                                    | [0..59]                |
// | 16   | 1    | Second                                    | [0..59]                |
// +------+------+-------------------------------------------+------------------------+

type Port1Payload struct {
	Latitude  float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude  float64 `json:"altitude"`
	Year      uint8   `json:"year" validate:"gte=0,lte=255"`
	Month     uint8   `json:"month" validate:"gte=1,lte=12"`
	Day       uint8   `json:"day" validate:"gte=1,lte=31"`
	Hour      uint8   `json:"hour" validate:"gte=0,lte=23"`
	Minute    uint8   `json:"minute" validate:"gte=0,lte=59"`
	Second    uint8   `json:"second" validate:"gte=0,lte=59"`
}

var _ common.Position = &Port1Payload{}

func (p Port1Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port1Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port1Payload) GetAltitude() *float64 {
	return &p.Altitude
}

func (p Port1Payload) GetSource() common.PositionSource {
	return common.PositionSource_GNSS
}

func (p Port1Payload) GetCapturedAt() *time.Time {
	capturedAt, err := time.Parse(time.RFC3339,
		fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
			int(p.Year)+2000,
			p.Month,
			p.Day,
			p.Hour,
			p.Minute,
			p.Second,
		))
	if err != nil {
		return nil
	}
	return &capturedAt
}

func (p Port1Payload) GetBuffered() bool {
	return false
}
