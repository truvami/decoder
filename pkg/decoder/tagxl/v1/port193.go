package tagxl

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

type Port193Payload struct {
	EndOfGroup bool    `json:"endOfGroup"`
	GroupToken uint8   `json:"groupToken" validate:"gte=2, lte=31"`
	NavMessage []byte  `json:"navMessage"`
	Moving     bool    `json:"moving"` // always true for port 193
	Latitude   float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude  float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude   float64 `json:"altitude"`
}

var _ decoder.UplinkFeatureGNSS = &Port193Payload{}
var _ decoder.UplinkFeatureMoving = &Port193Payload{}

func (p Port193Payload) GetLatitude() float64 {
	return p.Latitude
}

func (p Port193Payload) GetLongitude() float64 {
	return p.Longitude
}

func (p Port193Payload) GetAltitude() float64 {
	return p.Altitude
}

func (p Port193Payload) GetAccuracy() *float64 {
	return nil
}

func (p Port193Payload) GetTTF() *time.Duration {
	return nil
}

func (p Port193Payload) GetPDOP() *float64 {
	return nil
}

func (p Port193Payload) GetSatellites() *uint8 {
	return nil
}

func (p Port193Payload) IsMoving() bool {
	return true
}
