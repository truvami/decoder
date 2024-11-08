package tagsl

import "time"

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0    | 1    | Buffer level                              | uint16                  |
// | 2    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 3-6  | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 7-10 | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 11-12| 2    | Altitude                                  | uint16, 1/10 meter     |
// | 13-16| 4    | Unix timestamp                            | uint32                 |
// | 17-18| 2    | Battery voltage                           | uint16, mV             |
// | 19   | 1    | TTF                                       | uint8                  |
// | 20   | 1    | PDOP  		                             | uint8, 1/2 meter       |
// | 21   | 1    | Number of satellites                      | uint8, 		          |
// | 22-27| 6    | MAC1                                      | 6 x uint8              |
// | 28   | 1    | RSSI1                                     | int8                   |
// | …    |      |                                           |                        |
// |      | 6    | MACN                                      | 6 x uint8              |
// |      | 1    | RSSIN                                     | int8                   |
// +------+------+-------------------------------------------+------------------------+

// Timestamp for the Wi-Fi scanning is TSGNSS – TTF + 10 seconds.
type Port151Payload struct {
	BufferLevel uint16    `json:"bufferLevel"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Altitude    float64   `json:"altitude"`
	Timestamp   time.Time `json:"timestamp"`
	Battery     float64   `json:"battery"`
	TTF         uint16    `json:"ttf"`
	PDOP        float64   `json:"pdop"`
	Satellites  float64   `json:"satellites"`
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
	Mac7        string    `json:"mac7"`
	Rssi7       int8      `json:"rssi7"`
}
