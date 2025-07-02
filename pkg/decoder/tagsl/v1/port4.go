package tagsl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/common"
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
	LocalizationIntervalWhileMoving uint32  `json:"localizationIntervalWhileMoving" validate:"gte=60,lte=86400"`
	LocalizationIntervalWhileSteady uint32  `json:"localizationIntervalWhileSteady" validate:"gte=120,lte=86400"`
	HeartbeatInterval               uint32  `json:"heartbeatInterval" validate:"gte=300,lte=604800"`
	GPSTimeoutWhileWaitingForFix    uint16  `json:"gpsTimeoutWhileWaitingForFix" validate:"gte=60,lte=86400"`
	AccelerometerWakeupThreshold    uint16  `json:"accelerometerWakeupThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay              uint16  `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	DeviceState                     uint8   `json:"deviceState"`
	FirmwareVersionMajor            uint8   `json:"firmwareVersionMajor"`
	FirmwareVersionMinor            uint8   `json:"firmwareVersionMinor"`
	FirmwareVersionPatch            uint8   `json:"firmwareVersionPatch"`
	HardwareVersionType             uint8   `json:"hardwareVersionType"`
	HardwareVersionRevision         uint8   `json:"hardwareVersionRevision"`
	BatteryKeepAliveMessageInterval uint32  `json:"batteryKeepAliveMessageInterval" validate:"gte=300,lte=604800"`
	BatchSize                       *uint16 `json:"batchSize" validate:"lte=50"`
	BufferSize                      *uint16 `json:"bufferSize" validate:"gte=128,lte=8128"`
}

func (p Port4Payload) MarshalJSON() ([]byte, error) {
	deviceState := func() string {
		switch p.DeviceState {
		case 1:
			return "moving"
		case 2:
			return "steady"
		default:
			return "unknown"
		}
	}()
	return json.Marshal(&struct {
		MovingInterval         string  `json:"movingInterval"`
		SteadyInterval         string  `json:"steadyInterval"`
		ConfigInterval         string  `json:"configInterval"`
		BatteryInterval        string  `json:"batteryInterval"`
		GnssTimeout            string  `json:"gnssTimeout"`
		AccelerometerThreshold string  `json:"accelerometerThreshold"`
		AccelerometerDelay     string  `json:"accelerometerDelay"`
		DeviceState            string  `json:"deviceState"`
		FirmwareVersion        string  `json:"firmwareVersion"`
		HardwareVersion        string  `json:"hardwareVersion"`
		BatchSize              *uint16 `json:"batchSize"`
		BufferSize             *uint16 `json:"bufferSize"`
	}{
		MovingInterval:         (time.Duration(p.LocalizationIntervalWhileMoving) * time.Second).String(),
		SteadyInterval:         (time.Duration(p.LocalizationIntervalWhileSteady) * time.Second).String(),
		ConfigInterval:         (time.Duration(p.HeartbeatInterval) * time.Second).String(),
		BatteryInterval:        (time.Duration(p.BatteryKeepAliveMessageInterval) * time.Second).String(),
		GnssTimeout:            (time.Duration(p.GPSTimeoutWhileWaitingForFix) * time.Second).String(),
		AccelerometerThreshold: fmt.Sprintf("%dmg", p.AccelerometerWakeupThreshold),
		AccelerometerDelay:     (time.Duration(p.AccelerometerDelay) * time.Millisecond).String(),
		DeviceState:            deviceState,
		FirmwareVersion:        fmt.Sprintf("%d.%d.%d", p.FirmwareVersionMajor, p.FirmwareVersionMinor, p.FirmwareVersionPatch),
		HardwareVersion:        fmt.Sprintf("%d.%d", p.HardwareVersionType, p.HardwareVersionRevision),
		BatchSize:              p.BatchSize,
		BufferSize:             p.BufferSize,
	})
}

var _ decoder.UplinkFeatureMoving = &Port4Payload{}
var _ decoder.UplinkFeatureConfig = &Port4Payload{}
var _ decoder.UplinkFeatureFirmwareVersion = &Port4Payload{}
var _ decoder.UplinkFeatureHardwareVersion = &Port4Payload{}

func (p Port4Payload) GetBle() *bool {
	return nil
}

func (p Port4Payload) GetGnss() *bool {
	return nil
}

func (p Port4Payload) GetWifi() *bool {
	return nil
}

func (p Port4Payload) GetAcceleration() *bool {
	return nil
}

func (p Port4Payload) GetMovingInterval() *uint32 {
	return &p.LocalizationIntervalWhileMoving
}

func (p Port4Payload) GetSteadyInterval() *uint32 {
	return &p.LocalizationIntervalWhileSteady
}

func (p Port4Payload) GetConfigInterval() *uint32 {
	return &p.HeartbeatInterval
}

func (p Port4Payload) GetGnssTimeout() *uint16 {
	return &p.GPSTimeoutWhileWaitingForFix
}

func (p Port4Payload) GetAccelerometerThreshold() *uint16 {
	return &p.AccelerometerWakeupThreshold
}

func (p Port4Payload) GetAccelerometerDelay() *uint16 {
	return &p.AccelerometerDelay
}

func (p Port4Payload) GetBatteryInterval() *uint32 {
	return &p.BatteryKeepAliveMessageInterval
}

func (p Port4Payload) GetRejoinInterval() *uint32 {
	return nil
}

func (p Port4Payload) GetLowLightThreshold() *uint16 {
	return nil
}

func (p Port4Payload) GetHighLightThreshold() *uint16 {
	return nil
}

func (p Port4Payload) GetLowTemperatureThreshold() *int8 {
	return nil
}

func (p Port4Payload) GetHighTemperatureThreshold() *int8 {
	return nil
}

func (p Port4Payload) GetAccessPointsThreshold() *uint8 {
	return nil
}

func (p Port4Payload) GetBatchSize() *uint16 {
	return p.BatchSize
}

func (p Port4Payload) GetBufferSize() *uint16 {
	return p.BufferSize
}

func (p Port4Payload) GetDataRate() *decoder.DataRate {
	return nil
}

func (p Port4Payload) GetFirmwareHash() *string {
	return nil
}

func (p Port4Payload) GetFirmwareVersion() *string {
	return common.StringPtr(fmt.Sprintf("%d.%d.%d", p.FirmwareVersionMajor, p.FirmwareVersionMinor, p.FirmwareVersionPatch))
}

func (p Port4Payload) GetHardwareVersion() string {
	return fmt.Sprintf("%d.%d", p.HardwareVersionType, p.HardwareVersionRevision)
}

func (p Port4Payload) IsMoving() bool {
	return p.DeviceState == 0x01
}
