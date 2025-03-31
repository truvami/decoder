package tagxl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | flags ble[0] gnss[1] wifi[2] acc[3] rfu [4:7] | byte       |
// | 1    | 2    | moving interval in seconds                    | uint16     |
// | 3    | 2    | steady interval in seconds                    | uint16     |
// | 5    | 2    | acceleration threshold in milli-g             | uint16     |
// | 7    | 2    | acceleration delay in milliseconds            | uint16     |
// | 9    | 1    | heartbeat interval in seconds                 | uint8      |
// | 10   | 1    | fwu advertisement interval in seconds         | uint8      |
// | 11   | 2    | battery voltage in millie-volts               | uint16     |
// | 13   | 4    | first 4 bytes of sha-1 git commit             | byte[4]    |
// | 17   | 2    | resets since flash erase                      | uint16     |
// | 19   | 4    | reset cause register state                    | uint32     |
// | 23   | 2    | gnss scans since reset                        | uint16     |
// | 25   | 2    | wifi scans since reset                        | uint16     |
// +------+------+-----------------------------------------------+------------+

type Port151Payload struct {
	Ble                      bool    `json:"ble"`
	Gnss                     bool    `json:"gnss"`
	Wifi                     bool    `json:"wifi"`
	Acceleration             bool    `json:"acceleration"`
	Rfu                      uint8   `json:"rfu" validate:"gte=0,lte=15"`
	MovingInterval           uint16  `json:"movingInterval"`
	SteadyInterval           uint16  `json:"steadyInterval"`
	AccelerationThreshold    uint16  `json:"accelerationThreshold"`
	AccelerationDelay        uint16  `json:"accelerationDelay"`
	HeartbeatInterval        uint8   `json:"heartbeatInterval"`
	FwuAdvertisementInterval uint8   `json:"fwuAdvertisementInterval"`
	BatteryVoltage           float32 `json:"batteryVoltage" validate:"gte=1,lte=5"`
	FirmwareHash             string  `json:"firmwareHash"`
	ResetCount               uint16  `json:"resetCount"`
	ResetCause               uint32  `json:"resetCause"`
	GnssScans                uint16  `json:"gnssScans"`
	WifiScans                uint16  `json:"wifiScans"`
}

var _ decoder.UplinkFeatureBase = &Port151Payload{}
var _ decoder.UplinkFeatureBattery = &Port151Payload{}
var _ decoder.UplinkFeatureConfig = &Port151Payload{}

func (p Port151Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port151Payload) GetBatteryVoltage() float64 {
	return float64(p.BatteryVoltage)
}

func (p Port151Payload) GetBle() *bool {
	return &p.Ble
}

func (p Port151Payload) GetGnss() *bool {
	return &p.Gnss
}

func (p Port151Payload) GetWifi() *bool {
	return &p.Wifi
}

func (p Port151Payload) GetMovingInterval() *uint32 {
	movingInterval := uint32(p.MovingInterval)
	return &movingInterval
}

func (p Port151Payload) GetSteadyInterval() *uint32 {
	steadyInterval := uint32(p.SteadyInterval)
	return &steadyInterval
}

func (p Port151Payload) GetConfigInterval() *uint32 {
	return nil
}

func (p Port151Payload) GetGnssTimeout() *uint16 {
	return nil
}

func (p Port151Payload) GetAccelerometerThreshold() *uint16 {
	return &p.AccelerationThreshold
}

func (p Port151Payload) GetAccelerometerDelay() *uint16 {
	return &p.AccelerationDelay
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

func (p Port151Payload) GetBatchSize() *uint16 {
	return nil
}

func (p Port151Payload) GetBufferSize() *uint16 {
	return nil
}
