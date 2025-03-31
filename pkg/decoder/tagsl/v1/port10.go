package tagsl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
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
	Moving     bool          `json:"moving"`
	DutyCycle  bool          `json:"dutyCycle"`
	Latitude   float64       `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude  float64       `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude   float64       `json:"altitude"`
	Timestamp  time.Time     `json:"timestamp"`
	Battery    float64       `json:"battery" validate:"gte=1,lte=5"`
	TTF        time.Duration `json:"ttf"`
	PDOP       float64       `json:"pdop"`
	Satellites uint8         `json:"satellites" validate:"gte=3,lte=27"`
}

func (p Port10Payload) MarshalJSON() ([]byte, error) {
	type Alias Port10Payload
	return json.Marshal(&struct {
		*Alias
		TTF string `json:"ttf"`
	}{
		Alias: (*Alias)(&p),
		TTF:   fmt.Sprintf("%.0fs", p.TTF.Seconds()),
	})
}

var _ decoder.UplinkFeatureBase = &Port10Payload{}
var _ decoder.UplinkFeatureGNSS = &Port10Payload{}
var _ decoder.UplinkFeatureBattery = &Port10Payload{}
var _ decoder.UplinkFeatureMoving = &Port10Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port10Payload{}

func (p Port10Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port10Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port10Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port10Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port10Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port10Payload) GetTTF() *time.Duration {
	return &p.TTF
}

func (p Port10Payload) GetPDOP() *float64 {
	return &p.PDOP
}

func (p Port10Payload) GetSatellites() *uint8 {
	return &p.Satellites
}

func (p Port10Payload) GetBatteryVoltage() float64 {
	return p.Battery
}

func (p Port10Payload) IsMoving() bool {
	return p.Moving
}

func (p Port10Payload) IsDutyCycle() bool {
	return p.DutyCycle
}
