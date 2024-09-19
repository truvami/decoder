package tagsl

import "time"

// +------+------+-------------------------------------------+-----------+
// | Byte | Size | Description                               | Format    |
// +------+------+-------------------------------------------+-----------+
// | 0    | 1    | Timestamp                                 | uint8     |
// | 1    | 1    | Moving                                    | uint8     |
// | 2    | 6    | MAC1                                      | uint8[6]  |
// | 8    | 1    | RSSI1                                     | int8      |
// | 9    | 6    | MAC2                                      | uint8[6]  |
// | 15   | 1    | RSSI2                                     | int8      |
// | 16   | 6    | MAC3                                      | uint8[6]  |
// | 22   | 1    | RSSI3                                     | int8      |
// | 23   | 6    | MAC4                                      | uint8[6]  |
// | 29   | 1    | RSSI4                                     | int8      |
// | 30   | 6    | MAC5                                      | uint8[6]  |
// | 36   | 1    | RSSI5                                     | int8      |
// | 37   | 6    | MAC6                                      | uint8[6]  |
// | 43   | 1    | RSSI6                                     | int8      |
// +------+------+-------------------------------------------+-----------+

type Port7Payload struct {
	Timestamp time.Time `json:"timestamp"`
	Moving    bool      `json:"moving"`
	Mac1      string    `json:"mac1"`
	Rssi1     int8      `json:"rssi1"`
	Mac2      string    `json:"mac2"`
	Rssi2     int8      `json:"rssi2"`
	Mac3      string    `json:"mac3"`
	Rssi3     int8      `json:"rssi3"`
	Mac4      string    `json:"mac4"`
	Rssi4     int8      `json:"rssi4"`
	Mac5      string    `json:"mac5"`
	Rssi5     int8      `json:"rssi5"`
	Mac6      string    `json:"mac6"`
	Rssi6     int8      `json:"rssi6"`
}
