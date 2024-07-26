package tagsl

// +------+------+-------------------------------------------+------------------------+
// | Byte | Size | Description                               | Format                 |
// +------+------+-------------------------------------------+------------------------+
// | 0    | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 1-4  | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 5-8  | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 9-10 | 2    | Altitude                                  | uint16, 1/100 meter    |
// | 11   | 1    | Year                                      | uint8, year after 2000 |
// | 12   | 1    | Month                                     | uint8, [1..12]         |
// | 13   | 1    | Day                                       | uint8, [1..31]         |
// | 14   | 1    | Hour                                      | [0..23]                |
// | 15   | 1    | Minute                                    | [0..59]                |
// | 16   | 1    | Second                                    | [0..59]                |
// +------+------+-------------------------------------------+------------------------+

type Port1Payload struct {
	Moving    bool    `json:"moving"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Year      uint8   `json:"year"`
	Month     uint8   `json:"month"`
	Day       uint8   `json:"day"`
	Hour      uint8   `json:"hour"`
	Minute    uint8   `json:"minute"`
	Second    uint8   `json:"second"`
}
