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
// | 0    | 2    | Buffer level                              | uint16                 |
// | 2    | 1    | Duty cycle flag                           | uint1                  |
// | 2    | 1    | Config change id                          | uint4                  |
// | 2    | 1    | Config change success flag                | uint1                  |
// | 2    | 1    | Reserved                                  | uint1                  |
// | 2    | 1    | Moving flag                               | uint1                  |
// | 3    | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 7    | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 11   | 2    | Altitude                                  | uint16, 1/10 meter     |
// | 13   | 4    | Unix timestamp                            | uint32                 |
// | 17   | 2    | Battery voltage                           | uint16, mV             |
// | 19   | 1    | Time to fix                               | uint8                  |
// | 20   | 1    | Position dilution of precision            | uint8, 1/2 meter       |
// | 21   | 1    | Number of satellites                      | uint8, 		            |
// | 22   | 6    | Mac 1                                     | uint8[6]               |
// | 28   | 1    | Rssi 1                                    | int8                   |
// | 29   | 6    | Mac 2                                     | uint8[6]               |
// | 35   | 1    | Rssi 2                                    | int8                   |
// | 36   | 6    | Mac 3                                     | uint8[6]               |
// | 42   | 1    | Rssi 3                                    | int8                   |
// | 43   | 6    | Mac 4                                     | uint8[6]               |
// | 49   | 1    | Rssi 4                                    | int8                   |
// +------+------+-------------------------------------------+------------------------+

// Timestamp for the Wi-Fi scanning is TSGNSS – TTF + 10 seconds.
type Port151Payload struct {
	BufferLevel         uint16        `json:"bufferLevel"`
	DutyCycle           bool          `json:"dutyCycle"`
	ConfigChangeId      uint8         `json:"configChangeId" validate:"gte=0,lte=15"`
	ConfigChangeSuccess bool          `json:"configChangeSuccess"`
	Moving              bool          `json:"moving"`
	Latitude            float64       `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude           float64       `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude            float64       `json:"altitude"`
	Timestamp           time.Time     `json:"timestamp"`
	Battery             float64       `json:"battery" validate:"gte=1,lte=5"`
	TTF                 time.Duration `json:"ttf"`
	PDOP                float64       `json:"pdop"`
	Satellites          uint8         `json:"satellites" validate:"gte=3,lte=27"`
	Mac1                string        `json:"mac1"`
	Rssi1               int8          `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac2                *string       `json:"mac2"`
	Rssi2               *int8         `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac3                *string       `json:"mac3"`
	Rssi3               *int8         `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac4                *string       `json:"mac4"`
	Rssi4               *int8         `json:"rssi4" validate:"gte=-120,lte=-20"`
}

func (p Port151Payload) MarshalJSON() ([]byte, error) {
	type Alias Port151Payload
	return json.Marshal(&struct {
		*Alias
		Altitude   string `json:"altitude"`
		Timestamp  string `json:"timestamp"`
		Battery    string `json:"battery"`
		TTF        string `json:"ttf"`
		PDOP       string `json:"pdop"`
		Satellites uint8  `json:"satellites"`
	}{
		Alias:      (*Alias)(&p),
		Altitude:   fmt.Sprintf("%.1fm", p.Altitude),
		Timestamp:  p.Timestamp.Format(time.RFC3339),
		Battery:    fmt.Sprintf("%.3fv", p.Battery),
		TTF:        fmt.Sprintf("%.0fs", p.TTF.Seconds()),
		PDOP:       fmt.Sprintf("%.1fm", p.PDOP),
		Satellites: p.Satellites,
	})
}

var _ decoder.UplinkFeatureBase = &Port151Payload{}
var _ decoder.UplinkFeatureGNSS = &Port151Payload{}
var _ decoder.UplinkFeatureBattery = &Port151Payload{}
var _ decoder.UplinkFeatureWiFi = &Port151Payload{}
var _ decoder.UplinkFeatureBuffered = &Port151Payload{}
var _ decoder.UplinkFeatureMoving = &Port151Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port151Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port151Payload{}

func (p Port151Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port151Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port151Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port151Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port151Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port151Payload) GetTTF() *time.Duration {
	return nil
}

func (p Port151Payload) GetPDOP() *float64 {
	return nil
}

func (p Port151Payload) GetSatellites() *uint8 {
	return nil
}

func (p Port151Payload) GetBatteryVoltage() float64 {
	return p.Battery
}

func (p Port151Payload) GetLowBattery() *bool {
	return nil
}

func (p Port151Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" && p.Rssi1 != 0 {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != nil && p.Rssi2 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac2,
			RSSI: *p.Rssi2,
		})
	}

	if p.Mac3 != nil && p.Rssi3 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac3,
			RSSI: *p.Rssi3,
		})
	}

	if p.Mac4 != nil && p.Rssi4 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac4,
			RSSI: *p.Rssi4,
		})
	}

	return accessPoints
}

func (p Port151Payload) GetBufferLevel() uint16 {
	return p.BufferLevel
}

func (p Port151Payload) IsMoving() bool {
	return p.Moving
}

func (p Port151Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port151Payload) GetConfigId() *uint8 {
	return &p.ConfigChangeId
}

func (p Port151Payload) GetConfigChange() bool {
	return p.ConfigChangeSuccess
}
