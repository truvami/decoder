package smartlabel

// +-------+------+------------------+----------------+
// | Byte  | Size | Description      | Format         |
// +-------+------+------------------+----------------+
// | 2     | 1    | Flags              | uint8, -     |
// | 2     | 1    | GNSSEnabled        | uint8, -     |
// | 2     | 1    | WiFiEnabled        | uint8, -     |
// | 2     | 1    | AccEnabled         | uint8, -     |
// | 2     | 1    | StaticSF           | uint8, -     |
// | 3-4   | 2    | SteadyIntervalS    | uint16, s    |
// | 5-6   | 2    | MovingIntervalS    | uint16, s    |
// | 7     | 1    | HeartbeatIntervalH | uint8, h     |
// | 8-9   | 2    | LEDBlinkIntervalS  | uint16, s    |
// | 10-11 | 2    | AccThresholdMS     | uint16, ms   |
// | 12-13 | 2    | AccDelayMS         | uint16, ms   |
// | 14-17 | 4    | GitHash            | uint32, -    |
// +-------+------+------------------+----------------+

type Port11ConfigurationPayload struct {
	Flags               uint8    `json:"flags"`
	GNSSEnabled         uint8    `json:"gnssEnabled"`
	WiFiEnabled         uint8    `json:"wifiEnabled"`
	AccEnabled          uint8    `json:"accEnabled"`
	StaticSF           	string   `json:"staticSF"`
	SteadyIntervalS    	uint16   `json:"steadyIntervalS"`
	MovingIntervalS    	uint16   `json:"movingIntervalS"`
	HeartbeatIntervalH 	uint8    `json:"heartbeatIntervalH"`
	LEDBlinkIntervalS  	uint16   `json:"ledBlinkIntervalS"`
	AccThresholdMS     	uint16   `json:"accThresholdMS"`
	AccDelayMS         	uint16   `json:"accDelayMS"`
	GitHash            	string   `json:"gitHash"`
}
