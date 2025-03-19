package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0    | 1    | Buffer level                              | uint16                 |
// | 2    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 3-6  | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 7-10 | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 11-12| 2    | Altitude                                  | uint16, 1/10 meter     |
// | 13-16| 4    | Unix timestamp                            | uint32                 |
// | 17-18| 2    | Battery voltage                           | uint16, mV             |
// | 19   | 1    | TTF                                       | uint8                  |
// | 20   | 1    | PDOP  		                                 | uint8, 1/2 meter       |
// | 21   | 1    | Number of satellites                      | uint8, 		            |
// | 22-27| 6    | MAC1                                      | 6 x uint8              |
// | 28   | 1    | RSSI1                                     | int8                   |
// | …    |      |                                           |                        |
// |      | 6    | MACN                                      | 6 x uint8              |
// |      | 1    | RSSIN                                     | int8                   |
// +------+------+-------------------------------------------+------------------------+

// Timestamp for the Wi-Fi scanning is TSGNSS – TTF + 10 seconds.
type Port151Payload struct {
	Moving      bool          `json:"moving"`
	DutyCycle   bool          `json:"dutyCycle"`
	BufferLevel uint16        `json:"bufferLevel"`
	Latitude    float64       `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude   float64       `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude    float64       `json:"altitude"`
	Timestamp   time.Time     `json:"timestamp"`
	Battery     float64       `json:"battery" validate:"gte=1,lte=5"`
	TTF         time.Duration `json:"ttf"`
	PDOP        float64       `json:"pdop"`
	Satellites  uint8         `json:"satellites" validate:"gte=3,lte=27"`
	Mac1        string        `json:"mac1"`
	Rssi1       int8          `json:"rssi1"`
	Mac2        string        `json:"mac2"`
	Rssi2       int8          `json:"rssi2"`
	Mac3        string        `json:"mac3"`
	Rssi3       int8          `json:"rssi3"`
	Mac4        string        `json:"mac4"`
	Rssi4       int8          `json:"rssi4"`
	Mac5        string        `json:"mac5"`
	Rssi5       int8          `json:"rssi5"`
	Mac6        string        `json:"mac6"`
	Rssi6       int8          `json:"rssi6"`
	Mac7        string        `json:"mac7"`
	Rssi7       int8          `json:"rssi7"`
}

var _ decoder.UplinkFeatureBase = &Port151Payload{}
var _ decoder.UplinkFeatureGNSS = &Port151Payload{}
var _ decoder.UpLinkFeatureBattery = &Port151Payload{}
var _ decoder.UplinkFeatureWiFi = &Port151Payload{}
var _ decoder.UplinkFeatureBuffered = &Port151Payload{}
var _ decoder.UplinkFeatureMoving = &Port151Payload{}
var _ decoder.UplinkFeatureDutyCycle = &Port151Payload{}

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

func (p Port151Payload) GetAccessPoints() []decoder.AccessPoint {
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

	if p.Mac5 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	if p.Mac6 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac6,
			RSSI: p.Rssi6,
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
