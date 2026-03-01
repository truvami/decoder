package smartlabel

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0 	  | 1    | EOG (End of group)                        | uint1                  |
// |      |      | RFU (Reserved for future use)             | uint2                  |
// |      |      | GRP_TOKEN (Group token)                   | uint5                  |
// | 1 	  | n    | U-GNSSLOC-NAV message(s)                  |                        |
// +------+------+-------------------------------------------+------------------------+

type Port180Payload struct {
	EndOfGroup bool    `json:"endOfGroup"`
	GroupToken uint8   `json:"groupToken" validate:"gte=2, lte=31"`
	NavMessage []byte  `json:"navMessage"`
	Latitude   float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude  float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude   float64 `json:"altitude"`
}

var _ decoder.UplinkFeatureGNSS = &Port180Payload{}

func (p Port180Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port180Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port180Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port180Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port180Payload) GetTTF() *time.Duration {
	return nil
}

func (p Port180Payload) GetPDOP() *float64 {
	return nil
}

func (p Port180Payload) GetSatellites() *uint8 {
	return nil
}
