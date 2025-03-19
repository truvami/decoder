package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
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
	Moving    bool    `json:"moving"`
	DutyCycle bool    `json:"dutyCycle"`
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

// Enforce that Port1Payload implements interfaces
var _ decoder.UplinkFeatureBase = &Port1Payload{}
var _ decoder.UplinkFeatureGNSS = &Port1Payload{}
var _ decoder.UplinkFeatureMoving = &Port1Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port1Payload{}

func (p Port1Payload) GetTimestamp() *time.Time {
	timestamp := time.Date(
		int(p.Year)+2000,
		time.Month(p.Month),
		int(p.Day),
		int(p.Hour),
		int(p.Minute),
		int(p.Second),
		0,
		time.UTC,
	)
	return &timestamp
}

func (p Port1Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port1Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port1Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port1Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port1Payload) GetTTF() *time.Duration {
	return nil
}

func (p Port1Payload) GetPDOP() *float64 {
	return nil
}

func (p Port1Payload) GetSatellites() *uint8 {
	return nil
}

func (p Port1Payload) IsMoving() bool {
	return p.Moving
}

func (p Port1Payload) IsDutyCycle() bool {
	return p.DutyCycle
}
