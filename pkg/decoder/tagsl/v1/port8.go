package tagsl

import (
	"reflect"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

// | Byte   | Size | Description                                 | Format                                              |
// |--------|------|---------------------------------------------|-----------------------------------------------------|
// | 0-1    | 2    | Scan interval                               | uint16, s                                           |
// | 2      | 1    | Scan time                                   | uint8, s [0..180]                                   |
// | 3      | 1    | Max beacons                                 | uint8                                               |
// | 4      | 1    | Min. Rssi value                             | int8                                                |
// | 5-14   | 10   | Advertising name/eddystone namespace filter | 10 x ASCII or 10 x uint8                            |
// | 15-16  | 2    | Accelerometer trigger hold timer            | uint16, s                                           |
// | 17-18  | 2    | Accelerometer threshold                     | uint16, mg                                          |
// | 19     | 1    | Scan mode                                   | 0 - no filter; 1 - advertised name filter;          |
// |        |      |                                             | 2 - eddystone namespace filter                      |
// | 20-21  | 2    | BLE current configuration uplink interval   | uint16, s                                           |
// |--------|------|---------------------------------------------|-----------------------------------------------------|

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

func port8PayloadConfig() common.PayloadConfig {
	return common.PayloadConfig{
		Fields: []common.FieldConfig{
			{Name: "ScanInterval", Start: 0, Length: 2},
			{Name: "ScanTime", Start: 2, Length: 1},
			{Name: "MaxBeacons", Start: 3, Length: 1},
			{Name: "MinRssiValue", Start: 4, Length: 1},
			{Name: "AdvertisingFilter", Start: 5, Length: 10},
			{Name: "AccelerometerTriggerHoldTimer", Start: 15, Length: 2},
			{Name: "AccelerometerThreshold", Start: 17, Length: 2},
			{Name: "ScanMode", Start: 19, Length: 1},
			{Name: "BLECurrentConfigurationUplinkInterval", Start: 20, Length: 2},
		},
		TargetType: reflect.TypeOf(Port8Payload{}),
		Features:   []decoder.Feature{},
	}
}
