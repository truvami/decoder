package tagsl

// | Byte   | Size | Description                                 | Format                                             |
// |--------|------|--------------------------------------------|-----------------------------------------------------|
// | 0-1    | 2    | Scan interval                              | uint16, s                                           |
// | 2      | 1    | Scan time                                  | uint8, s [0..180]                                   |
// | 3      | 1    | Max beacons                                | uint8                                               |
// | 4      | 1    | Min. Rssi value                            | int8                                                |
// | 5-14   | 10   | Advertising name/eddystone namespace filter | 10 x ASCII or 10 x uint8                           |
// | 15-16  | 2    | Accelerometer trigger hold timer           | uint16, s                                           |
// | 17-18  | 2    | Accelerometer threshold                    | uint16, mg                                          |
// | 19     | 1    | Scan mode                                  | 0 - no filter; 1 - advertised name filter;          |
// | 		|	   |											| 2 - eddystone namespace filter                      |
// | 20-21  | 2    | BLE current configuration uplink interval  | uint16, s                                           |

type Port8Payload struct {
	ScanInterval                          uint16 `json:"scanInterval"`
	ScanTime                              uint8  `json:"scanTime"`
	MaxBeacons                            uint8  `json:"maxBeacons"`
	MinRssiValue                          int8   `json:"minRssiValue"`
	AdvertisingFilter                     string `json:"advertisingFilter"`
	AccelerometerTriggerHoldTimer         uint16 `json:"accelerometerTriggerHoldTimer"`
	AccelerometerThreshold                uint16 `json:"accelerometerThreshold"`
	ScanMode                              uint8  `json:"scanMode"`
	BLECurrentConfigurationUplinkInterval uint16 `json:"bleCurrentConfigurationUplinkInterval"`
}
