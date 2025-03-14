package nomadxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// | Byte  | Size | Description | Format               |
// |-------|------|-------------|----------------------|
// | 0-3   | 4    | UTC Date    | uint32, DDMMYY       |
// | 4-7   | 4    | UTC Time    | uint32, HHMMSS       |
// | 8-11  | 4    | Latitude    | int32, 1/100'000 deg |
// | 12-15 | 4    | Longitude   | int32, 1/100'000 deg |
// | 16-19 | 4    | Altitude    | int32, 1/100 m       |

type Port103Payload struct {
	UTCDate   uint32  `json:"date"`
	UTCTime   uint32  `json:"time"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

var _ decoder.UplinkFeatureBase = &Port103Payload{}
var _ decoder.UplinkFeatureGNSS = &Port103Payload{}

// GetAccuracy implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetAccuracy() *float64 {
	return nil
}

// GetAltitude implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetAltitude() float64 {
	return p.Altitude
}

// GetLatitude implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetLatitude() float64 {
	return p.Latitude
}

// GetLongitude implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetLongitude() float64 {
	return p.Longitude
}

// GetPDOP implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetPDOP() *float64 {
	return nil
}

// GetSatellites implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetSatellites() *uint8 {
	return nil
}

// GetTTF implements decoder.UplinkFeatureGNSS.
func (p Port103Payload) GetTTF() *float64 {
	return nil
}

// GetTimestamp implements decoder.UplinkFeatureBase.
func (p Port103Payload) GetTimestamp() *time.Time {
	return nil
}
