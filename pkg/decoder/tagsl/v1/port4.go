package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +-------+------+-------------------------------------------+------------+
// | Byte  | Size | Description                               | Format     |
// +-------+------+-------------------------------------------+------------+
// | 0-3   | 4    | Localization interval while moving, IM    | uint32, s  |
// | 4-7   | 4    | Localization interval while steady, IS    | uint32, s  |
// | 8-11  | 4    | Heartbeat interval, IH                    | uint32, s  |
// | 12-13 | 2    | GPS timeout while waiting for fix         | uint16, s  |
// | 14-15 | 2    | Accelerometer wakeup threshold            | uint16, mg |
// | 16-17 | 2    | Accelerometer delay                       | uint16, ms |
// | 18    | 1    | Device state (moving = 1, steady = 2)     | uint8      |
// | 19-21 | 3    | Firmware version (major,;minor; patch)    | 3 x uint8  |
// | 22-23 | 2    | Hardware version (type; revision)         | 2 x uint8  |
// | 24-27 | 4    | Battery “keep-alive” message interval, IB | uint32, s  |
// +-------+------+-------------------------------------------+------------+

type Port4Payload struct {
	LocalizationIntervalWhileMoving uint32 `json:"localizationIntervalWhileMoving" validate:"gte=60,lte=86400"`
	LocalizationIntervalWhileSteady uint32 `json:"localizationIntervalWhileSteady" validate:"gte=120,lte=86400"`
	HeartbeatInterval               uint32 `json:"heartbeatInterval" validate:"gte=300,lte=604800"`
	GPSTimeoutWhileWaitingForFix    uint16 `json:"gpsTimeoutWhileWaitingForFix" validate:"gte=60,lte=86400"`
	AccelerometerWakeupThreshold    uint16 `json:"accelerometerWakeupThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay              uint16 `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	DeviceState                     uint8  `json:"deviceState"`
	FirmwareVersionMajor            uint8  `json:"firmwareVersionMajor"`
	FirmwareVersionMinor            uint8  `json:"firmwareVersionMinor"`
	FirmwareVersionPatch            uint8  `json:"firmwareVersionPatch"`
	HardwareVersionType             uint8  `json:"hardwareVersionType"`
	HardwareVersionRevision         uint8  `json:"hardwareVersionRevision"`
	BatteryKeepAliveMessageInterval uint32 `json:"batteryKeepAliveMessageInterval" validate:"gte=300,lte=604800"`
	BatchSize                       uint16 `json:"batchSize" validate:"lte=50"`
	BufferSize                      uint16 `json:"bufferSize" validate:"gte=128,lte=8128"`
}

var _ decoder.UplinkFeatureBase = &Port4Payload{}

func (p Port4Payload) GetTimestamp() *time.Time {
	return nil
}
