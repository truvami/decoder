package nomadxs

// +-------+------+-------------------------------------------+------------------+
// | Byte  | Size | Description                               | Format           |
// +-------+------+-------------------------------------------+------------------+
// | 0-3   | 4    | Localization interval while moving, IM    | uint32, s        |
// | 4-7   | 4    | Localization interval while steady, IS    | uint32, s        |
// | 8-11  | 4    | Config/Status interval, IC                | uint32, s        |
// | 12-13 | 2    | GPS timeout while waiting for fix         | uint16, s        |
// | 14-15 | 2    | Accelerometer wakeup threshold            | uint16, mg       |
// | 16-17 | 2    | Accelerometer delay                       | uint16, ms       |
// | 18-20 | 3    | Firmware version (major,;minor; patch)    | 3 x uint8        |
// | 21-22 | 2    | Hardware version (type; revision)         | 2 x uint8        |
// | 23-26 | 4    | Battery “keep-alive” message interval, IB | uint32, s        |
// | 27-30 | 4    | Re-Join interval in case of Join Failed   | uint32, s        |
// | 31    | 1    | Accuracy enhancement                      | uint8, s [0..59] |
// | 32-33 | 2    | Light lower threshold                     | uint16, Lux      |
// | 34-35 | 2    | Light upper threshold                     | uint16, Lux      |
// +-------+------+-------------------------------------------+------------------+

type Port4Payload struct {
	LocalizationIntervalWhileMoving uint32 `json:"localizationIntervalWhileMoving"`
	LocalizationIntervalWhileSteady uint32 `json:"localizationIntervalWhileSteady"`
	HeartbeatInterval               uint32 `json:"heartbeatInterval"`
	GPSTimeoutWhileWaitingForFix    uint16 `json:"gpsTimeoutWhileWaitingForFix"`
	AccelerometerWakeupThreshold    uint16 `json:"accelerometerWakeupThreshold"`
	AccelerometerDelay              uint16 `json:"accelerometerDelay"`
	FirmwareVersionMajor            uint8  `json:"firmwareVersionMajor"`
	FirmwareVersionMinor            uint8  `json:"firmwareVersionMinor"`
	FirmwareVersionPatch            uint8  `json:"firmwareVersionPatch"`
	HardwareVersionType             uint8  `json:"hardwareVersionType"`
	HardwareVersionRevision         uint8  `json:"hardwareVersionRevision"`
	BatteryKeepAliveMessageInterval uint32 `json:"batteryKeepAliveMessageInterval"`
	ReJoinInterval                  uint32 `json:"reJoinInterval"`
	AccuracyEnhancement             uint8  `json:"accuracyEnhancement"`
	LightLowerThreshold             uint16 `json:"lightLowerThreshold"`
	LightUpperThreshold             uint16 `json:"lightUpperThreshold"`
}
