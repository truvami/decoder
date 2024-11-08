package nomadxl

// | Byte  | Size | Description                                         | Format              |
// |-------|------|-----------------------------------------------------|---------------------|
// | 0-7   | 8    | System time (ms since reset)                        | uint64_t, ms        |
// | 8-11  | 4    | UTC Date                                            | uint32, DDMMYY      |
// | 12-15 | 4    | UTC Time                                            | uint32, HHMMSS      |
// | 16-17 | 2    | Buffer level (STA)                                  | uint16              |
// | 18-19 | 2    | Buffer level (GPS)                                  | uint16              |
// | 20-21 | 2    | Buffer level (ACC)                                  | uint16              |
// | 22-23 | 2    | Buffer level (LOG)                                  | uint16              |
// | 24-25 | 2    | Temperature                                         | int16, 0.1 °C       |
// | 26-27 | 2    | Pressure                                            | uint16, 0.1 hPa     |
// | 28-29 | 2    | Orientation X                                       | int16, mG           |
// | 30-31 | 2    | Orientation Y                                       | int16, mG           |
// | 32-33 | 2    | Orientation Z                                       | int16, mG           |
// | 34-35 | 2    | Battery voltage                                     | uint16, mV          |
// | 36    | 1    | LoRaWAN battery level (1 to 254)                    | uint8               |
// | 37    | 1    | Last TTF (time to fix)                              | uint8, s            |
// | 38-39 | 2    | NMEA sentences checksum OK                          | uint16              |
// | 40-41 | 2    | NMEA sentences checksum fail                        | uint16              |
// | 42-43 | 2    | Total GPS signal to noise (0-99 for each satellite) | uint16, C/n0 [dBHz] |
// | 44    | 1    | GPS satellite count Navstar                         | uint8               |
// | 45    | 1    | GPS satellite count Glonass                         | uint8               |
// | 46    | 1    | GPS satellite count Galileo                         | uint8               |
// | 47    | 1    | GPS satellite count Beidou                          | uint8               |
// | 48-49 | 2    | GPS dilution of precision                           | uint16, cm          |

type Port101Payload struct {
	SystemTime         int64   `json:"systemTime"`
	UTCDate            uint32  `json:"date"`
	UTCTime            uint32  `json:"time"`
	Temperature        float32 `json:"temperature"`
	Pressure           float32 `json:"pressure"`
	TimeToFix          uint8   `json:"timeToFix"`
	AccelerometerXAxis int16   `json:"accelerometerXAxis"`
	AccelerometerYAxis int16   `json:"accelerometerYAxis"`
	AccelerometerZAxis int16   `json:"accelerometerZAxis"`
	Battery            float64 `json:"battery"`
	BatteryLorawan     uint8   `json:"batteryLorawan"`
}