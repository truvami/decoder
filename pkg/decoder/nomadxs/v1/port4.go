package nomadxs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

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
	AccelerometerWakeupThreshold    uint16 `json:"accelerometerWakeUpThreshold"`
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

func (p Port4Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MovingInterval         string `json:"movingInterval"`
		SteadyInterval         string `json:"steadyInterval"`
		ConfigInterval         string `json:"configInterval"`
		BatteryInterval        string `json:"batteryInterval"`
		RejoinInterval         string `json:"rejoinInterval"`
		GnssTimeout            string `json:"gnssTimeout"`
		AccelerometerThreshold string `json:"accelerometerThreshold"`
		AccelerometerDelay     string `json:"accelerometerDelay"`
		AccuracyEnhancement    string `json:"accuracyEnhancement"`
		LightLowerThreshold    string `json:"lightLowerThreshold"`
		LightUpperThreshold    string `json:"lightUpperThreshold"`
		FirmwareVersion        string `json:"firmwareVersion"`
		HardwareVersion        string `json:"hardwareVersion"`
	}{
		MovingInterval:         (time.Duration(p.LocalizationIntervalWhileMoving) * time.Second).String(),
		SteadyInterval:         (time.Duration(p.LocalizationIntervalWhileSteady) * time.Second).String(),
		ConfigInterval:         (time.Duration(p.HeartbeatInterval) * time.Second).String(),
		BatteryInterval:        (time.Duration(p.BatteryKeepAliveMessageInterval) * time.Second).String(),
		RejoinInterval:         (time.Duration(p.ReJoinInterval) * time.Second).String(),
		GnssTimeout:            (time.Duration(p.GPSTimeoutWhileWaitingForFix) * time.Second).String(),
		AccelerometerThreshold: fmt.Sprintf("%dmg", p.AccelerometerWakeupThreshold),
		AccelerometerDelay:     (time.Duration(p.AccelerometerDelay) * time.Millisecond).String(),
		AccuracyEnhancement:    (time.Duration(p.AccuracyEnhancement) * time.Second).String(),
		LightLowerThreshold:    fmt.Sprintf("%dlux", p.LightLowerThreshold),
		LightUpperThreshold:    fmt.Sprintf("%dlux", p.LightUpperThreshold),
		FirmwareVersion:        fmt.Sprintf("%d.%d.%d", p.FirmwareVersionMajor, p.FirmwareVersionMinor, p.FirmwareVersionPatch),
		HardwareVersion:        fmt.Sprintf("%d.%d", p.HardwareVersionType, p.HardwareVersionRevision),
	})
}

var _ decoder.UplinkFeatureBase = &Port4Payload{}
var _ decoder.UplinkFeatureConfig = &Port4Payload{}
var _ decoder.UplinkFeatureFirmwareVersion = &Port4Payload{}
var _ decoder.UplinkFeatureHardwareVersion = &Port4Payload{}

func (p Port4Payload) GetTimestamp() *time.Time {
	return nil
}

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
	return &p.ReJoinInterval
}

func (p Port4Payload) GetLowLightThreshold() *uint16 {
	return &p.LightLowerThreshold
}

func (p Port4Payload) GetHighLightThreshold() *uint16 {
	return &p.LightUpperThreshold
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
	return nil
}

func (p Port4Payload) GetBufferSize() *uint16 {
	return nil
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
