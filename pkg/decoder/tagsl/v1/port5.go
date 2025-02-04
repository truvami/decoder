package tagsl

import (
	"github.com/truvami/decoder/pkg/common"
)

// +------+------+-------------------------------------------+-----------+
// | Byte | Size | Description                               | Format    |
// +------+------+-------------------------------------------+-----------+
// | 0    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8     |
// | 1-6  | 6    | MAC1                                      | 6 x uint8 |
// | 7    | 1    | RSSI1                                     | int8      |
// | …    |      |                                           |           |
// |      | 6    | MACN                                      | 6 x uint8 |
// |      | 1    | RSSIN                                     | int8      |
// +------+------+-------------------------------------------+-----------+

type Port5Payload struct {
	Mac1  string `json:"mac1"`
	Rssi1 int8   `json:"rssi1"`
	Mac2  string `json:"mac2"`
	Rssi2 int8   `json:"rssi2"`
	Mac3  string `json:"mac3"`
	Rssi3 int8   `json:"rssi3"`
	Mac4  string `json:"mac4"`
	Rssi4 int8   `json:"rssi4"`
	Mac5  string `json:"mac5"`
	Rssi5 int8   `json:"rssi5"`
	Mac6  string `json:"mac6"`
	Rssi6 int8   `json:"rssi6"`
	Mac7  string `json:"mac7"`
	Rssi7 int8   `json:"rssi7"`
}

var _ common.WifiLocation = &Port5Payload{}

func (p Port5Payload) GetAccessPoints() []common.WifiAccessPoint {
	var accessPoints []common.WifiAccessPoint

	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac1, p.Rssi1)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac2, p.Rssi2)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac3, p.Rssi3)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac4, p.Rssi4)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac5, p.Rssi5)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac6, p.Rssi6)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac7, p.Rssi7)

	return accessPoints
}
