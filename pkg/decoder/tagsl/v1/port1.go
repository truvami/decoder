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
	Latitude  float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude  float64 `json:"altitude"`
	Year      uint8   `json:"year" validate:"gte=0,lte=255"`
	Month     uint8   `json:"month" validate:"gte=1,lte=12"`
	Day       uint8   `json:"day" validate:"gte=1,lte=31"`
	Hour      uint8   `json:"hour" validate:"gte=0,lte=23"`
	Minute    uint8   `json:"minute" validate:"gte=0,lte=59"`
	Second    uint8   `json:"second" validate:"gte=0,lte=59"`
}
