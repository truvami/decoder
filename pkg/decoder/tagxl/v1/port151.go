package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port151Payload struct {
	Battery                              *float32 `json:"battery" validate:"gte=1,lte=5"`
	GnssScans                            *uint16  `json:"gnssScans"`
	WifiScans                            *uint16  `json:"wifiScans"`
	LocalizationIntervalWhileMoving      *uint16  `json:"movingInterval" validate:"gte=60,lte=86400"`
	LocalizationIntervalWhileSteady      *uint16  `json:"steadyInterval" validate:"gte=120,lte=86400"`
	AccelerometerWakeupThreshold         *uint16  `json:"accelerometerWakeupThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay                   *uint16  `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	HeartbeatInterval                    *uint8   `json:"heartbeatInterval" validate:"gte=300,lte=604800"`
	GnssEnabled                          *bool    `json:"gnssEnabled"`
	WiFiEnabled                          *bool    `json:"wifiEnabled"`
	AccelerometerEnabled                 *bool    `json:"accelerometerEnabled"`
	AdvertisementFirmwareUpgradeInterval *uint8   `json:"advertisementFirmwareUpgradeInterval" validate:"gte=1,lte=86400"`
	FirmwareHash                         *string  `json:"firmwareHash"`
	ResetCount                           *uint16  `json:"resetCount"`
	ResetCause                           *uint32  `json:"resetCause"`
}

var _ decoder.UplinkFeatureBase = &Port151Payload{}
var _ decoder.UplinkFeatureBattery = &Port151Payload{}
var _ decoder.UplinkFeatureConfig = &Port151Payload{}

func (p Port151Payload) GetTimestamp() *time.Time { // coverage-ignore
	return nil
}

func (p Port151Payload) GetBatteryVoltage() float64 {
	return float64(*p.Battery)
}

func (p Port151Payload) GetLowBattery() *bool { // coverage-ignore
	return nil
}

func (p Port151Payload) GetBle() *bool { // coverage-ignore
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
	if p.LocalizationIntervalWhileMoving == nil { // coverage-ignore
		return nil
	}
	movingInterval := uint32(*p.LocalizationIntervalWhileMoving)
	return &movingInterval
}

func (p Port151Payload) GetSteadyInterval() *uint32 {
	if p.LocalizationIntervalWhileSteady == nil { // coverage-ignore
		return nil
	}
	steadyInterval := uint32(*p.LocalizationIntervalWhileSteady)
	return &steadyInterval
}

func (p Port151Payload) GetConfigInterval() *uint32 {
	if p.HeartbeatInterval == nil { // coverage-ignore
		return nil
	}
	interval := uint32(*p.HeartbeatInterval)
	return &interval
}

func (p Port151Payload) GetGnssTimeout() *uint16 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetAccelerometerThreshold() *uint16 {
	return p.AccelerometerWakeupThreshold
}

func (p Port151Payload) GetAccelerometerDelay() *uint16 {
	return p.AccelerometerDelay
}

func (p Port151Payload) GetBatteryInterval() *uint32 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetRejoinInterval() *uint32 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetLowLightThreshold() *uint16 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetHighLightThreshold() *uint16 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetLowTemperatureThreshold() *int8 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetHighTemperatureThreshold() *int8 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetAccessPointsThreshold() *uint8 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetBatchSize() *uint16 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetBufferSize() *uint16 { // coverage-ignore
	return nil
}

func (p Port151Payload) GetDataRate() *decoder.DataRate { // coverage-ignore
	return nil
}
