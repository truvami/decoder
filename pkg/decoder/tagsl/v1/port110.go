package tagsl

import "time"

// +-------+------+-------------------------------------------+------------------------+
// | Byte  | Size | Description                               | Format                 |
// +-------+------+-------------------------------------------+------------------------+
// | 0     | 2    | Buffer level                              | uint8                  |
// | 2     | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 3-6   | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 7-10  | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 11-12 | 2    | Altitude                                  | uint16, 1/10 meter     |
// | 11-14 | 4    | Unix timestamp                            | uint32                 |
// | 17-18 | 2    | Battery voltage                           | uint16, mV             |
// +-------+------+-------------------------------------------+------------------------+

type Port110Payload struct {
	BufferLevel uint8     `json:"bufferLevel"`
	Moving      bool      `json:"moving"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Altitude    float64   `json:"altitude"`
	Timestamp   time.Time `json:"timestamp"`
	Battery     float64   `json:"battery"`
}