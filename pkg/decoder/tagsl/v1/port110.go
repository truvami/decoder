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
// | 0    | 2    | Buffer level                              | uint16                 |
// | 2    | 1    | Duty cycle flag                           | uint1                  |
// | 2    | 1    | Config change id                          | uint4                  |
// | 2    | 1    | Config change success flag                | uint1                  |
// | 2    | 1    | Reserved                                  | uint1                  |
// | 2    | 1    | Moving flag                               | uint1                  |
// | 3    | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 7    | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 11   | 2    | Altitude                                  | uint16, 1/10 meter     |
// | 11   | 4    | Unix timestamp                            | uint32                 |
// | 17   | 2    | Battery voltage                           | uint16, mV             |
// | 19   | 1    | Time to fix                               | uint8, s               |
// | 20   | 1    | Position dilution of precision            | uint8, 0.5m            |
// | 21   | 1    | Number of satellites                      | uint8                  |
// +------+------+-------------------------------------------+------------------------+

type Port110Payload struct {
	BufferLevel         uint16         `json:"bufferLevel"`
	DutyCycle           bool           `json:"dutyCycle"`
	ConfigChangeId      uint8          `json:"configChangeId" validate:"gte=0,lte=15"`
	ConfigChangeSuccess bool           `json:"configChangeSuccess"`
	Moving              bool           `json:"moving"`
	Latitude            float64        `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude           float64        `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude            float64        `json:"altitude"`
	Timestamp           time.Time      `json:"timestamp"`
	Battery             float64        `json:"battery" validate:"gte=1,lte=5"`
	TTF                 *time.Duration `json:"ttf"`
	PDOP                *float64       `json:"pdop"`
	Satellites          *uint8         `json:"satellites" validate:"gte=3,lte=27"`
}

func (p Port110Payload) MarshalJSON() ([]byte, error) {
	type Alias Port110Payload
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

var _ decoder.UplinkFeatureBase = &Port110Payload{}
var _ decoder.UplinkFeatureGNSS = &Port110Payload{}
var _ decoder.UplinkFeatureBattery = &Port110Payload{}
var _ decoder.UplinkFeatureBuffered = &Port110Payload{}
var _ decoder.UplinkFeatureMoving = &Port110Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port110Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port110Payload{}

func (p Port110Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port110Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port110Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port110Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port110Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port110Payload) GetTTF() *time.Duration {
	return p.TTF
}

func (p Port110Payload) GetPDOP() *float64 {
	return p.PDOP
}

func (p Port110Payload) GetSatellites() *uint8 {
	return p.Satellites
}

func (p Port110Payload) GetBatteryVoltage() float64 {
	return p.Battery
}

func (p Port110Payload) GetLowBattery() *bool {
	return nil
}

func (p Port110Payload) GetBufferLevel() uint16 {
	return p.BufferLevel
}

func (p Port110Payload) IsMoving() bool {
	return p.Moving
}

func (p Port110Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port110Payload) GetConfigId() *uint8 {
	return &p.ConfigChangeId
}

func (p Port110Payload) GetConfigChange() bool {
	return p.ConfigChangeSuccess
}
