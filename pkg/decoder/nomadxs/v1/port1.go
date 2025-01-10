package nomadxs

// +-------+------+-------------------------------------------+------------------------+
// | Byte  | Size | Description                               | Format                 |
// +-------+------+-------------------------------------------+------------------------+
// | 0     | 1    | Status[6:2] + Moving flag[0] (moving = 1) | uint8                  |
// | 1-4   | 4    | Latitude                                  | int32, 1/1’000’000 deg |
// | 5-8   | 4    | Longitude                                 | int32, 1/1’000’000 deg |
// | 9-10  | 2    | Altitude                                  | uint16, 1/100 meter    |
// | 11    | 1    | Year                                      | uint8, year after 2000 |
// | 12    | 1    | Month                                     | uint8, [1..12]         |
// | 13    | 1    | Day                                       | uint8, [1..31]         |
// | 14    | 1    | Hour                                      | [1..23]                |
// | 15    | 1    | Minute                                    | [1..59]                |
// | 16    | 1    | Second                                    | [1..59]                |
// | 17    | 1    | Time to fix                               | uint8, second(s)       |
// | 18-19 | 2    | Ambient light                             | uint16, Lux            |
// | 20-21 | 2    | Accelerometer X-axis                      | int16, mg              |
// | 22-23 | 2    | Accelerometer Y-axis                      | int16, mg              |
// | 24-25 | 2    | Accelerometer Z-axis                      | int16, mg              |
// | 26-27 | 2    | Temperature*                              | int16, 0.1 °C          |
// | 28-29 | 2    | Pressure*                                 | uint16, 0.1 hPa        |
// | 30-31 | 2    | Gyroscope* X-axis                         | int16, 0.1 dps         |
// | 32-33 | 2    | Gyroscope* Y-axis                         | int16, 0.1 dps         |
// | 34-35 | 2    | Gyroscope* Z-axis                         | int16, 0.1 dps         |
// | 36-37 | 2    | Magnetometer* X-axis                      | int16, mgauss          |
// | 38-39 | 2    | Magnetometer* Y-axis                      | int16, mgauss          |
// | 40-41 | 2    | Magnetometer* Z-axis                      | int16, mgauss          |
// +-------+------+-------------------------------------------+------------------------+

type Port1Payload struct {
	Moving             bool    `json:"moving"`
	Latitude           float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude          float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Altitude           float64 `json:"altitude"`
	Year               uint8   `json:"year" validate:"gte=0,lte=255"`
	Month              uint8   `json:"month" validate:"gte=1,lte=12"`
	Day                uint8   `json:"day" validate:"gte=1,lte=31"`
	Hour               uint8   `json:"hour" validate:"gte=0,lte=23"`
	Minute             uint8   `json:"minute" validate:"gte=0,lte=59"`
	Second             uint8   `json:"second" validate:"gte=0,lte=59"`
	TimeToFix          uint8   `json:"timeToFix"`
	AmbientLight       uint16  `json:"ambientLight"`
	AccelerometerXAxis int16   `json:"accelerometerXAxis"`
	AccelerometerYAxis int16   `json:"accelerometerYAxis"`
	AccelerometerZAxis int16   `json:"accelerometerZAxis"`
	Temperature        float32 `json:"temperature" validate:"gte=-20,lte=60"`
	Pressure           float32 `json:"pressure" validate:"gte=0,lte=1100"`
	GyroscopeXAxis     float32 `json:"gyroscopeXAxis"`
	GyroscopeYAxis     float32 `json:"gyroscopeYAxis"`
	GyroscopeZAxis     float32 `json:"gyroscopeZAxis"`
	MagnetometerXAxis  float32 `json:"magnetometerXAxis"`
	MagnetometerYAxis  float32 `json:"magnetometerYAxis"`
	MagnetometerZAxis  float32 `json:"magnetometerZAxis"`
}
