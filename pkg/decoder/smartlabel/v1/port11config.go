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
	GNSSEnabled         uint8    `json:"gnss_enabled"`
	WiFiEnabled         uint8    `json:"wifi_enabled"`
	AccEnabled          uint8    `json:"acc_enabled"`
	StaticSF           	string    `json:"static_sf"`
	SteadyIntervalS    	uint16   `json:"steady_interval_s"`
	MovingIntervalS    	uint16   `json:"moving_interval_s"`
	HeartbeatIntervalH 	uint8    `json:"heartbeat_interval_h"`
	LEDBlinkIntervalS  	uint16   `json:"led_blink_interval_s"`
	AccThresholdMS     	uint16   `json:"acc_threshold_ms"`
	AccDelayMS         	uint16   `json:"acc_delay_ms"`
	GitHash            	string   `json:"git_hash"`
}
