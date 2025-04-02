package smartlabel

import (
	"fmt"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+---------------------------------------------+--------------+
// | Byte | Size | Description                                 | Format       |
// +------+------+---------------------------------------------+--------------+
// | 0    | 1    | flags data rate[0:2] acc[3] wifi[4] gnss[5] | byte         |
// | 1    | 2    | steady interval in seconds                  | uint16       |
// | 3    | 2    | moving interval in seconds                  | uint16       |
// | 5    | 1    | config interval in seconds                  | uint8        |
// | 6    | 2    | acceleration threshold                      | uint16, mg   |
// | 8    | 2    | acceleration delay                          | uint16, ms   |
// | 10   | 2    | temperature sensor polling interval         | uint16, s    |
// | 12   | 2    | temperature uplink hold interval            | uint16, s    |
// | 14   | 1    | temperature lower threshold                 | int8, C      |
// | 15   | 1    | temperature upper threshold                 | int8, C      |
// | 16   | 1    | minimal number of access points             | uint8        |
// | 17   | 3    | firmware version major minor patch          | uint8[3]     |
// +------+------+---------------------------------------------+--------------+

type Port4Payload struct {
	DataRate                   uint8  `json:"dataRate" validate:"gte=0,lte=7"`
	Acceleration               bool   `json:"acceleration"`
	Wifi                       bool   `json:"wifi"`
	Gnss                       bool   `json:"gnss"`
	SteadyInterval             uint16 `json:"steadyInterval"`
	MovingInterval             uint16 `json:"movingInterval"`
	HeartbeatInterval          uint8  `json:"heartbeatInterval"`
	AccelerationThreshold      uint16 `json:"accelerationThreshold"`
	AccelerationDelay          uint16 `json:"accelerationDelay"`
	TemperaturePollingInterval uint16 `json:"temperaturePollingInterval"`
	TemperatureUplinkInterval  uint16 `json:"temperatureUplinkInterval"`
	TemperatureLowerThreshold  int8   `json:"temperatureLowerThreshold"`
	TemperatureUpperThreshold  int8   `json:"temperatureUpperThreshold"`
	AccessPointsThreshold      uint8  `json:"accessPointsThreshold"`
	FirmwareVersionMajor       uint8  `json:"firmwareVersionMajor"`
	FirmwareVersionMinor       uint8  `json:"firmwareVersionMinor"`
	FirmwareVersionPatch       uint8  `json:"firmwareVersionPatch"`
}

var _ decoder.UplinkFeatureBase = &Port4Payload{}
var _ decoder.UplinkFeatureConfig = &Port4Payload{}
var _ decoder.UplinkFeatureFirmwareVersion = &Port4Payload{}

func (p Port4Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port4Payload) GetBle() *bool {
	return nil
}

func (p Port4Payload) GetGnss() *bool {
	return &p.Gnss
}

func (p Port4Payload) GetWifi() *bool {
	return &p.Wifi
}

func (p Port4Payload) GetAcceleration() *bool {
	return &p.Acceleration
}

func (p Port4Payload) GetMovingInterval() *uint32 {
	movingInterval := uint32(p.MovingInterval)
	return &movingInterval
}

func (p Port4Payload) GetSteadyInterval() *uint32 {
	steadyInterval := uint32(p.SteadyInterval)
	return &steadyInterval
}

func (p Port4Payload) GetConfigInterval() *uint32 {
	return nil
}

func (p Port4Payload) GetGnssTimeout() *uint16 {
	return nil
}

func (p Port4Payload) GetAccelerometerThreshold() *uint16 {
	return &p.AccelerationThreshold
}

func (p Port4Payload) GetAccelerometerDelay() *uint16 {
	return &p.AccelerationDelay
}

func (p Port4Payload) GetBatteryInterval() *uint32 {
	return nil
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
	return &p.TemperatureLowerThreshold
}

func (p Port4Payload) GetHighTemperatureThreshold() *int8 {
	return &p.TemperatureUpperThreshold
}

func (p Port4Payload) GetAccessPointsThreshold() *uint8 {
	return &p.AccessPointsThreshold
}

func (p Port4Payload) GetBatchSize() *uint16 {
	return nil
}

func (p Port4Payload) GetBufferSize() *uint16 {
	return nil
}

func (p Port4Payload) GetFirmwareVersion() string {
	return fmt.Sprintf("%d.%d.%d", p.FirmwareVersionMajor, p.FirmwareVersionMinor, p.FirmwareVersionPatch)
}
