package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +-----+------+------------------------------------------------+------------+
// | Tag | Size | Description                                    | Format     |
// +-----+------+------------------------------------------------+------------+
// | 40  | 1    | device flags                                   | byte       |
// |     |      | firmware upgrade flag                          | uint1      |
// |     |      | gnss flag                                      | uint1      |
// |     |      | wifi flag                                      | uint1      |
// |     |      | accelerometer flag                             | uint1      |
// |     |      | reserved                                       | uint4      |
// | 41  | 4    | moving interval                                | uint16, s  |
// |     |      | steady interval                                | uint16, s  |
// | 42  | 4    | accelerometer threshold                        | uint16, mg |
// |     |      | accelerometer delay                            | uint16, ms |
// | 43  | 1    | heartbeat interval                             | uint8, h   |
// | 44  | 1    | firmware upgrade advertisement                 | uint8, s   |
// | 45  | 2    | battery voltage                                | uint16, mv |
// | 46  | 4    | firmware hash                                  | byte[4]    |
// | 49  | 2    | reset count since flash erase                  | uint16     |
// | 4a  | 4    | reset cause register value                     | uint32     |
// | 4b  | 4    | gnss scans since reset                         | uint16     |
// |     |      | wifi scans since reset                         | uint16     |
// +-----+------+------------------------------------------------+------------+

type Port151Payload struct {
	GnssEnabled                          *bool    `json:"gnssEnabled"`
	WiFiEnabled                          *bool    `json:"wifiEnabled"`
	AccelerometerEnabled                 *bool    `json:"accelerometerEnabled"`
	LocalizationIntervalWhileMoving      *uint16  `json:"movingInterval" validate:"gte=60,lte=86400"`
	LocalizationIntervalWhileSteady      *uint16  `json:"steadyInterval" validate:"gte=120,lte=86400"`
	AccelerometerWakeupThreshold         *uint16  `json:"accelerometerWakeupThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay                   *uint16  `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	HeartbeatInterval                    *uint8   `json:"heartbeatInterval" validate:"gte=300,lte=604800"`
	AdvertisementFirmwareUpgradeInterval *uint8   `json:"advertisementFirmwareUpgradeInterval" validate:"gte=1,lte=86400"`
	Battery                              *float32 `json:"battery" validate:"gte=1,lte=5"`
	FirmwareHash                         *string  `json:"firmwareHash"`
	ResetCount                           *uint16  `json:"resetCount"`
	ResetCause                           *uint32  `json:"resetCause"`
	GnssScans                            *uint16  `json:"gnssScans"`
	WifiScans                            *uint16  `json:"wifiScans"`
}

var _ decoder.UplinkFeatureBase = &Port151Payload{}
var _ decoder.UplinkFeatureBattery = &Port151Payload{}
var _ decoder.UplinkFeatureConfig = &Port151Payload{}

func (p Port151Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port151Payload) GetBatteryVoltage() float64 {
	return float64(*p.Battery)
}

func (p Port151Payload) GetLowBattery() *bool {
	return nil
}

func (p Port151Payload) GetBle() *bool {
	return nil
}

func (p Port151Payload) GetGnss() *bool {
	return p.GnssEnabled
}

func (p Port151Payload) GetWifi() *bool {
	return p.WiFiEnabled
}

func (p Port151Payload) GetAcceleration() *bool {
	return p.AccelerometerEnabled
}

func (p Port151Payload) GetMovingInterval() *uint32 {
	if p.LocalizationIntervalWhileMoving == nil {
		return nil
	}
	movingInterval := uint32(*p.LocalizationIntervalWhileMoving)
	return &movingInterval
}

func (p Port151Payload) GetSteadyInterval() *uint32 {
	if p.LocalizationIntervalWhileSteady == nil {
		return nil
	}
	steadyInterval := uint32(*p.LocalizationIntervalWhileSteady)
	return &steadyInterval
}

func (p Port151Payload) GetConfigInterval() *uint32 {
	if p.HeartbeatInterval == nil {
		return nil
	}
	interval := uint32(*p.HeartbeatInterval) * 60 * 60
	return &interval
}

func (p Port151Payload) GetGnssTimeout() *uint16 {
	return nil
}

func (p Port151Payload) GetAccelerometerThreshold() *uint16 {
	return p.AccelerometerWakeupThreshold
}

func (p Port151Payload) GetAccelerometerDelay() *uint16 {
	return p.AccelerometerDelay
}

func (p Port151Payload) GetBatteryInterval() *uint32 {
	return nil
}

func (p Port151Payload) GetRejoinInterval() *uint32 {
	return nil
}

func (p Port151Payload) GetLowLightThreshold() *uint16 {
	return nil
}

func (p Port151Payload) GetHighLightThreshold() *uint16 {
	return nil
}

func (p Port151Payload) GetLowTemperatureThreshold() *int8 {
	return nil
}

func (p Port151Payload) GetHighTemperatureThreshold() *int8 {
	return nil
}

func (p Port151Payload) GetAccessPointsThreshold() *uint8 {
	return nil
}

func (p Port151Payload) GetBatchSize() *uint16 {
	return nil
}

func (p Port151Payload) GetBufferSize() *uint16 {
	return nil
}

func (p Port151Payload) GetDataRate() *decoder.DataRate {
	return nil
}
