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
// | 0    | 1    | Duty cycle flag                           | uint1                  |
// | 0    | 1    | Config change id                          | uint4                  |
// | 0    | 1    | Config change success flag                | uint1                  |
// | 0    | 1    | Reserved                                  | uint1                  |
// | 0    | 1    | Moving flag                               | uint1                  |
// | 1    | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 5    | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 9    | 2    | Altitude                                  | uint16, 1/10 meter     |
// | 11   | 4    | Unix timestamp                            | uint32                 |
// | 15   | 2    | Battery voltage                           | uint16, mV             |
// | 17   | 1    | Time to fix                               | uint8                  |
// | 18   | 6    | Mac 1                                     | uint8[6]               |
// | 24   | 1    | Rssi 1                                    | int8                   |
// | 25   | 6    | Mac 2                                     | uint8[6]               |
// | 31   | 1    | Rssi 2                                    | int8                   |
// | 32   | 6    | Mac 3                                     | uint8[6]               |
// | 38   | 1    | Rssi 3                                    | int8                   |
// | 39   | 6    | Mac 4                                     | uint8[6]               |
// | 45   | 1    | Rssi 4                                    | int8                   |
// +------+------+-------------------------------------------+------------------------+

// Timestamp for the Wi-Fi scanning is TSGNSS – TTF + 10 seconds.
type Port50Payload struct {
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
	Mac1                string        `json:"mac1"`
	Rssi1               int8          `json:"rssi1"`
	Mac2                string        `json:"mac2"`
	Rssi2               int8          `json:"rssi2"`
	Mac3                string        `json:"mac3"`
	Rssi3               int8          `json:"rssi3"`
	Mac4                string        `json:"mac4"`
	Rssi4               int8          `json:"rssi4"`
}

func (p Port50Payload) MarshalJSON() ([]byte, error) {
	type Alias Port50Payload
	return json.Marshal(&struct {
		*Alias
		TTF string `json:"ttf"`
	}{
		Alias: (*Alias)(&p),
		TTF:   fmt.Sprintf("%.0fs", p.TTF.Seconds()),
	})
}

var _ decoder.UplinkFeatureBase = &Port50Payload{}
var _ decoder.UplinkFeatureGNSS = &Port50Payload{}
var _ decoder.UplinkFeatureBattery = &Port50Payload{}
var _ decoder.UplinkFeatureWiFi = &Port50Payload{}
var _ decoder.UplinkFeatureMoving = &Port50Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port50Payload{}
var _ decoder.UplinkFeatureConfigChange = &Port50Payload{}

func (p Port50Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port50Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port50Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port50Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port50Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port50Payload) GetTTF() *time.Duration {
	return nil
}

func (p Port50Payload) GetPDOP() *float64 {
	return nil
}

func (p Port50Payload) GetSatellites() *uint8 {
	return nil
}

func (p Port50Payload) GetBatteryVoltage() float64 {
	return p.Battery
}

func (p Port50Payload) GetLowBattery() *bool {
	return nil
}

func (p Port50Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	return accessPoints
}

func (p Port50Payload) IsMoving() bool {
	return p.Moving
}

func (p Port50Payload) IsDutyCycle() bool {
	return p.DutyCycle
}

func (p Port50Payload) GetConfigId() *uint8 {
	return &p.ConfigChangeId
}

func (p Port50Payload) GetConfigChange() bool {
	return p.ConfigChangeSuccess
}
