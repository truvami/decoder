package tagxl

import (
	"github.com/truvami/decoder/pkg/decoder"
)

// +-----+------+------------------------------------------------+------------+
// | Tag | Size | Description                                    | Format     |
// +-----+------+------------------------------------------------+------------+
// | 40  | 1    | device flags                                   | byte       |
// |     |      | reserved                                       | uint4      |
// |     |      | accelerometer flag                             | uint1      |
// |     |      | wifi flag                                      | uint1      |
// |     |      | gnss flag                                      | uint1      |
// |     |      | firmware upgrade flag                          | uint1      |
// | 41  | 4    | moving interval                                | uint16, s  |
// |     |      | steady interval                                | uint16, s  |
// | 42  | 4    | accelerometer threshold                        | uint16, mg |
// |     |      | accelerometer delay                            | uint16, ms |
// | 43  | 1    | heartbeat interval                             | uint8, h   |
// | 44  | 1    | firmware upgrade advertisement                 | uint8, s   |
// | 45  | 2    | battery voltage                                | uint16, mv |
// | 46  | 4    | firmware hash                                  | byte[4]    |
// | 47  | 1    | rotation flags                                 | byte       |
// |     |      | reserved                                       | uint6      |
// |     |      | rotation invert                                | uint1      |
// |     |      | rotation confirmed                             | uint1      |
// | 49  | 2    | reset count since flash erase                  | uint16     |
// | 4a  | 4    | reset cause register value                     | uint32     |
// | 4b  | 4    | gnss scans since reset                         | uint16     |
// |     |      | wifi scans since reset                         | uint16     |
// | 4e  | 1    | Data rate setting (0-7)                        | uint8      |
// |     |      |   0: DR5 (EU868 SF7)                           |            |
// |     |      |   1: DR4 (EU868 SF8)                           |            |
// |     |      |   2: DR3 (EU868 SF9, US915 SF7)                |            |
// |     |      |   3: DR2 (EU868 SF10, US915 SF8)               |            |
// |     |      |   4: DR1 (EU868 SF11, US915 SF9)               |            |
// |     |      |   5: DR0 (EU868 SF12)                          |            |
// |     |      |   6: DR1-3 array (EU868 SF9-11, US915 SF7-9)   |            |
// |     |      |   7: ADR (SF7-12) for EU868                    |            |
// |     |      | See: https://docs.truvami.com/docs/Devices/tag%20XL%20/Payload%20Format%20%20tag%20XL/#settings-downlink
// +-----+------+------------------------------------------------+------------+

type Port151Payload struct {
	AccelerometerEnabled                 *bool             `json:"accelerometerEnabled"`
	WifiEnabled                          *bool             `json:"wifiEnabled"`
	GnssEnabled                          *bool             `json:"gnssEnabled"`
	FirmwareUpgrade                      *bool             `json:"firmwareUpgrade"`
	LocalizationIntervalWhileMoving      *uint16           `json:"movingInterval" validate:"gte=60,lte=86400"`
	LocalizationIntervalWhileSteady      *uint16           `json:"steadyInterval" validate:"gte=120,lte=86400"`
	AccelerometerWakeupThreshold         *uint16           `json:"accelerometerWakeupThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay                   *uint16           `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	HeartbeatInterval                    *uint8            `json:"heartbeatInterval" validate:"gte=0,lte=168"`
	AdvertisementFirmwareUpgradeInterval *uint8            `json:"advertisementFirmwareUpgradeInterval" validate:"gte=1,lte=86400"`
	Battery                              *float32          `json:"battery" validate:"gte=1,lte=5"`
	FirmwareHash                         *string           `json:"firmwareHash"`
	RotationInvert                       *bool             `json:"rotationInvert"`
	RotationConfirmed                    *bool             `json:"rotationConfirmed"`
	ResetCount                           *uint16           `json:"resetCount"`
	ResetCause                           *uint32           `json:"resetCause"`
	GnssScans                            *uint16           `json:"gnssScans"`
	WifiScans                            *uint16           `json:"wifiScans"`
	DataRate                             *decoder.DataRate `json:"dataRate"`
}

var _ decoder.UplinkFeatureBattery = &Port151Payload{}
var _ decoder.UplinkFeatureConfig = &Port151Payload{}
var _ decoder.UplinkFeatureFirmwareVersion = &Port151Payload{}

func (p Port151Payload) GetBatteryVoltage() float64 {
	if p.Battery == nil {
		return 0
	}
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
	return p.WifiEnabled
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
	return p.DataRate
}

func (p Port151Payload) GetFirmwareHash() *string {
	return p.FirmwareHash
}

func (p Port151Payload) GetFirmwareVersion() *string {
	return nil
}

// DataRateFromUint8 converts a uint8 data rate value to the corresponding TagXL DataRate enum.
func DataRateFromUint8(value uint8) decoder.DataRate {
	switch value {
	case 0:
		return decoder.DataRateTagXLDR5
	case 1:
		return decoder.DataRateTagXLDR4
	case 2:
		return decoder.DataRateTagXLDR3
	case 3:
		return decoder.DataRateTagXLDR2
	case 4:
		return decoder.DataRateTagXLDR1
	case 5:
		return decoder.DataRateTagXLDR0
	case 6:
		return decoder.DataRateTagXLDR1To3
	case 7:
		return decoder.DataRateTagXLADR
	default:
		return decoder.DataRate("unknown")
	}
}
