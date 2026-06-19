package tagxl

import (
	"github.com/truvami/decoder/pkg/decoder"
)

// Port 160 - BLE steady (no timestamp)
//
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | tag (movement / timestamp / sequence flags)   | byte       |
// | 1    | 1    | rssi beacon 1                                 | int8       |
// | 2    | 6    | mac address beacon 1                          | byte[6]    |
// | 8    | 1    | rssi beacon 2                                 | int8       |
// | 9    | 6    | mac address beacon 2                          | byte[6]    |
// | 15   | 1    | rssi beacon 3                                 | int8       |
// | 16   | 6    | mac address beacon 3                          | byte[6]    |
// | 22   | 1    | rssi beacon 4                                 | int8       |
// | 23   | 6    | mac address beacon 4                          | byte[6]    |
// | 29   | 1    | rssi beacon 5                                 | int8       |
// | 30   | 6    | mac address beacon 5                          | byte[6]    |
// | 36   | 1    | rssi beacon 6                                 | int8       |
// | 37   | 6    | mac address beacon 6                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+

type Port160Payload struct {
	Tag    byte    `json:"tag"`
	Moving bool    `json:"moving"` // Always false for port 160 (steady)
	Rssi1  int8    `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac1   string  `json:"mac1"`
	Rssi2  *int8   `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac2   *string `json:"mac2"`
	Rssi3  *int8   `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac3   *string `json:"mac3"`
	Rssi4  *int8   `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac4   *string `json:"mac4"`
	Rssi5  *int8   `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac5   *string `json:"mac5"`
	Rssi6  *int8   `json:"rssi6" validate:"gte=-120,lte=-20"`
	Mac6   *string `json:"mac6"`
}

var _ decoder.UplinkFeatureBle = &Port160Payload{}
var _ decoder.UplinkFeatureMoving = &Port160Payload{}

func (p Port160Payload) GetBeacons() []decoder.Beacon {
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

// Port 160 reports a steady (non-moving) state.
func (p Port160Payload) IsMoving() bool {
	return false
}
