package tagsl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0    | 1    | Duty cycle flag                           | uint1                  |
// | 0    | 1    | Config id                                 | uint4                  |
// | 0    | 1    | Config change flag                        | uint1                  |
// | 0    | 1    | Reserved                                  | uint1                  |
// | 0    | 1    | Moving flag                               | uint1                  |
// | 1    | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 5    | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 9    | 2    | Altitude                                  | uint16, 1/10 meter     |
// | 11   | 4    | Unix timestamp                            | uint32                 |
// | 15   | 2    | Battery voltage                           | uint16, mV             |
// | 17   | 1    | Time to fix                               | uint8, s               |
// | 18   | 1    | Position dilution of precision  		       | uint8, 1/2 meter       |
// | 19   | 1    | Number of satellites                      | uint8, 		            |
// +------+------+-------------------------------------------+------------------------+

type Port10Payload struct {
	DutyCycle    bool           `json:"dutyCycle"`
	ConfigId     uint8          `json:"configId" validate:"gte=0,lte=15"`
	ConfigChange bool           `json:"configChange"`
	Moving       bool           `json:"moving"`
	Latitude     float64        `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude    float64        `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude     float64        `json:"altitude"`
	Timestamp    time.Time      `json:"timestamp"`
	Battery      float64        `json:"battery" validate:"gte=1,lte=5"`
	TTF          *time.Duration `json:"ttf"`
	PDOP         *float64       `json:"pdop"`
	Satellites   *uint8         `json:"satellites" validate:"gte=3,lte=27"`
}

func (p Port10Payload) MarshalJSON() ([]byte, error) {
	type Alias Port10Payload
	var ttf *string = nil
	if p.TTF != nil {
		ttf = common.StringPtr(p.TTF.String())
	}
	var pdop *string = nil
	if p.PDOP != nil {
		pdop = common.StringPtr(fmt.Sprintf("%.1fm", *p.PDOP))
	}
	return json.Marshal(&struct {
		*Alias
		Altitude   string  `json:"altitude"`
		Timestamp  string  `json:"timestamp"`
		Battery    string  `json:"battery"`
		TTF        *string `json:"ttf"`
		PDOP       *string `json:"pdop"`
		Satellites *uint8  `json:"satellites"`
	}{
		Alias:      (*Alias)(&p),
		Altitude:   fmt.Sprintf("%.1fm", p.Altitude),
		Timestamp:  p.Timestamp.Format(time.RFC3339),
		Battery:    fmt.Sprintf("%.3fv", p.Battery),
		TTF:        ttf,
		PDOP:       pdop,
		Satellites: p.Satellites,
	})
}

var _ decoder.UplinkFeatureTimestamp = &Port10Payload{}
var _ decoder.UplinkFeatureGNSS = &Port10Payload{}
var _ decoder.UplinkFeatureBattery = &Port10Payload{}
var _ decoder.UplinkFeatureMoving = &Port10Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port10Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port10Payload{}

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
	return p.TTF
}

func (p Port10Payload) GetPDOP() *float64 {
	return p.PDOP
}

func (p Port10Payload) GetSatellites() *uint8 {
	return p.Satellites
}

func (p Port10Payload) GetBatteryVoltage() float64 {
	return p.Battery
}

func (p Port10Payload) GetLowBattery() *bool {
	return nil
}

func (p Port10Payload) IsMoving() bool {
	return p.Moving
}

func (p Port10Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port10Payload) GetConfigId() *uint8 {
	return &p.ConfigId
}

func (p Port10Payload) GetConfigChange() bool {
	return p.ConfigChange
}
