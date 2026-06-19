package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// Port 162 - BLE steady (with timestamp)
//
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 4    | timestamp (Unix epoch seconds)                | uint32     |
// | 4    | 1    | tag (movement / timestamp / sequence flags)   | byte       |
// | 5    | 1    | rssi beacon 1                                 | int8       |
// | 6    | 6    | mac address beacon 1                          | byte[6]    |
// | 12   | 1    | rssi beacon 2                                 | int8       |
// | 13   | 6    | mac address beacon 2                          | byte[6]    |
// | 19   | 1    | rssi beacon 3                                 | int8       |
// | 20   | 6    | mac address beacon 3                          | byte[6]    |
// | 26   | 1    | rssi beacon 4                                 | int8       |
// | 27   | 6    | mac address beacon 4                          | byte[6]    |
// | 33   | 1    | rssi beacon 5                                 | int8       |
// | 34   | 6    | mac address beacon 5                          | byte[6]    |
// | 40   | 1    | rssi beacon 6                                 | int8       |
// | 41   | 6    | mac address beacon 6                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+

type Port162Payload struct {
	Timestamp time.Time `json:"timestamp"`
	Tag       byte      `json:"tag"`
	Moving    bool      `json:"moving"` // Always false for port 162 (steady)
	Rssi1     int8      `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac1      string    `json:"mac1"`
	Rssi2     *int8     `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac2      *string   `json:"mac2"`
	Rssi3     *int8     `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac3      *string   `json:"mac3"`
	Rssi4     *int8     `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac4      *string   `json:"mac4"`
	Rssi5     *int8     `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac5      *string   `json:"mac5"`
	Rssi6     *int8     `json:"rssi6" validate:"gte=-120,lte=-20"`
	Mac6      *string   `json:"mac6"`
}

var _ decoder.UplinkFeatureBle = &Port162Payload{}
var _ decoder.UplinkFeatureMoving = &Port162Payload{}
var _ decoder.UplinkFeatureTimestamp = &Port162Payload{}
var _ decoder.UplinkFeatureBuffered = &Port162Payload{}

func (p Port162Payload) GetBufferLevel() *uint16 {
	return nil
}

func (p Port162Payload) IsBuffered() bool {
	return time.Since(p.Timestamp) > bufferedAgeThreshold
}

func (p Port162Payload) GetTimestamp() *time.Time {
	return &p.Timestamp
}

func (p Port162Payload) GetBeacons() []decoder.Beacon {
	beacons := []decoder.Beacon{}

	if p.Mac1 != "" {
		beacons = append(beacons, decoder.Beacon{
			MAC:  p.Mac1,
			RSSI: &p.Rssi1,
		})
	}

	if p.Mac2 != nil {
		beacons = append(beacons, decoder.Beacon{
			MAC:  *p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != nil {
		beacons = append(beacons, decoder.Beacon{
			MAC:  *p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != nil {
		beacons = append(beacons, decoder.Beacon{
			MAC:  *p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != nil {
		beacons = append(beacons, decoder.Beacon{
			MAC:  *p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	if p.Mac6 != nil {
		beacons = append(beacons, decoder.Beacon{
			MAC:  *p.Mac6,
			RSSI: p.Rssi6,
		})
	}

	return beacons
}

// Port 162 reports a steady (non-moving) state.
func (p Port162Payload) IsMoving() bool {
	return false
}
