package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// Port 213: WiFi localization triggered by rotation state change (moving).
// Same binary format as port 201 (timestamped WiFi) but rotation-triggered, not buffered.
//
// Version v1 (without RSSI values):
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 4    | timestamp (Unix epoch seconds)                | uint32     |
// | 4    | 1    | version (0x00 = v1, 0x01 = v2)                | byte       |
// | 5    | 6    | mac address signal 1                          | byte[6]    |
// | 11   | 6    | mac address signal 2                          | byte[6]    |
// | 17   | 6    | mac address signal 3                          | byte[6]    |
// | 23   | 6    | mac address signal 4                          | byte[6]    |
// | 29   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+
//
// Version v2 (with RSSI values):
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 4    | timestamp (Unix epoch seconds)                | uint32     |
// | 4    | 1    | version (0x00 = v1, 0x01 = v2)                | byte       |
// | 5    | 1    | rssi signal 1                                 | int8       |
// | 6    | 6    | mac address signal 1                          | byte[6]    |
// | 12   | 1    | rssi signal 2                                 | int8       |
// | 13   | 6    | mac address signal 2                          | byte[6]    |
// | 19   | 1    | rssi signal 3                                 | int8       |
// | 20   | 6    | mac address signal 3                          | byte[6]    |
// | 26   | 1    | rssi signal 4                                 | int8       |
// | 27   | 6    | mac address signal 4                          | byte[6]    |
// | 33   | 1    | rssi signal 5                                 | int8       |
// | 34   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+

const (
	Port213HeaderLength      = 5 // minimum payload bytes (4B timestamp + 1B version)
	Port213VersionIndex      = 4 // byte offset of the version field
	Port213Version1     byte = 0x00
	Port213Version2     byte = 0x01
)

type Port213Payload struct {
	Timestamp time.Time `json:"timestamp"`
	Version   byte      `json:"version" validate:"gte=0,lte=1"`
	Moving    bool      `json:"moving"` // Always true for Port 213
	Rssi1     *int8     `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac1      string    `json:"mac1"`
	Rssi2     *int8     `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac2      *string   `json:"mac2"`
	Rssi3     *int8     `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac3      *string   `json:"mac3"`
	Rssi4     *int8     `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac4      *string   `json:"mac4"`
	Rssi5     *int8     `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac5      *string   `json:"mac5"`
}

var _ decoder.UplinkFeatureWiFi = &Port213Payload{}
var _ decoder.UplinkFeatureMoving = &Port213Payload{}
var _ decoder.UplinkFeatureTimestamp = &Port213Payload{}

func (p Port213Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port213Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	return accessPoints
}

// Port 213 is the moving rotation-triggered variant.
func (p Port213Payload) IsMoving() bool {
	return true
}
