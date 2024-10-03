package nomadxl

// | Byte  | Size | Description | Format               |
// |-------|------|-------------|----------------------|
// | 0-3   | 4    | UTC Date    | uint32, DDMMYY       |
// | 4-7   | 4    | UTC Time    | uint32, HHMMSS       |
// | 8-11  | 4    | Latitude    | int32, 1/100'000 deg |
// | 12-15 | 4    | Longitude   | int32, 1/100'000 deg |
// | 16-19 | 4    | Altitude    | int32, 1/100 m       |

type Port1Payload struct {
	UTCDate				uint32 	`json:"date"`
	UTCTime				uint32 	`json:"time"`
	Latitude			float64 `json:"latitude"`
	Longitude			float64 `json:"longitude"`
	Altitude			float64 `json:"altitude"`
}
