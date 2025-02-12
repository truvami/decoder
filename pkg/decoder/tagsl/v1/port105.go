package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/common"
)

// +-------+------+-------------------------------------------+-----------+
// | Byte  | Size | Description                               | Format    |
// +-------+------+-------------------------------------------+-----------+
// | 0     | 1    | Buffer level                              | uint16    |
// | 2     | 4    | Unix timestamp                            | uint32    |
// | 7     | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8     |
// | 7     | 6    | MAC1                                      | 6 x uint8 |
// | 13    | 1    | RSSI1                                     | int8      |
// | …     |      |                                           |           |
// |       | 6    | MACN                                      | 6 x uint8 |
// |       | 1    | RSSIN                                     | int8      |
// +-------+------+-------------------------------------------+-----------+

type Port105Payload struct {
	BufferLevel uint16    `json:"bufferLevel"`
	Timestamp   time.Time `json:"timestamp"`
	Mac1        string    `json:"mac1"`
	Rssi1       int8      `json:"rssi1"`
	Mac2        string    `json:"mac2"`
	Rssi2       int8      `json:"rssi2"`
	Mac3        string    `json:"mac3"`
	Rssi3       int8      `json:"rssi3"`
	Mac4        string    `json:"mac4"`
	Rssi4       int8      `json:"rssi4"`
	Mac5        string    `json:"mac5"`
	Rssi5       int8      `json:"rssi5"`
	Mac6        string    `json:"mac6"`
	Rssi6       int8      `json:"rssi6"`
}

var _ common.WifiLocation = &Port105Payload{}

func (p Port105Payload) GetAccessPoints() []common.WifiAccessPoint {
	var accessPoints []common.WifiAccessPoint

	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac1, p.Rssi1)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac2, p.Rssi2)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac3, p.Rssi3)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac4, p.Rssi4)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac5, p.Rssi5)
	accessPoints = common.AppendAccessPoint(accessPoints, p.Mac6, p.Rssi6)

	return accessPoints
}
