package tagxl

import (
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
	"github.com/truvami/decoder/pkg/middleware"
)

type TagXLV1 struct {
	loraClient middleware.LoracloudClient
}

func NewTagXLv1Decoder(loraClient middleware.LoracloudClient) decoder.Decoder {
	return TagXLV1{loraClient}
}

func (t TagXLV1) GetConfig(port int16) (decoder.PayloadConfig, error) {
	switch port {
	case 151:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "DeviceFlags", Start: 0, Length: 1},
				{Name: "AssetTrackingIntervals", Start: 1, Length: 4},
				{Name: "AccelerationSensor", Start: 5, Length: 4},
				{Name: "HeartbeatInterval", Start: 9, Length: 1},
				{Name: "AdvertisementFWUInterval", Start: 10, Length: 1},
				{Name: "BatteryVoltage", Start: 11, Length: 2},
				{Name: "FirmwareHash", Start: 13, Length: 4},
			},
			TargetType: reflect.TypeOf(TagXLConfig{}),
		}, nil
	case 192:
		return decoder.PayloadConfig{}, nil
	}

	return decoder.PayloadConfig{}, nil
}

func (t TagXLV1) Decode(data string, port int16, devEui string) (interface{}, error) {
	config, err := t.GetConfig(port)
	if err != nil {
		return nil, err
	}

	decodedData, err := helpers.Parse(data, config)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}

// #	Device Setting	Tag	Size	Data	Format
// 1	Device Flags	0x40	0x01	bit 0: BLE_FWU_ENABLED
// bit 1: GNSS_ENABLE
// bit 2: WIFI_ENABLE
// bit 3: Set ACCELERATION_ENABLE
// bit 4-7: RFU	bit field
// 2	Asset Tracking Intervals	0x41	0x04	data 0: MOVING_INTERVAL
// data 1: STEADY_INTERVAL	uint16_t[2]
// 3	Acceleration Sensor Settings	0x42	0x04	data 0: ACCELERATION_SENSITIVITY
// data 1: ACCELERATION_DELAY	uint16_t[2]
// 4	HEARTBEAT_INTERVAL	0x43	0x01	Heartbeat interval in hours	uint8_t
// 5	ADVERTISEMENT_FWU_INTERVAL	0x44	0x01	Value in seconds	uint8_t
// 6	Battery Voltage	0x45	0x02	Battery voltage in mV	uint16_t
// 7	Firmware Hash	0x46	0x04	First 4 bytes of SHA-1 hash of git commit	uint8_t[4]
type TagXLConfig struct {
	DeviceFlags              uint8
	AssetTrackingIntervals   [2]uint16
	AccelerationSensor       [2]uint16
	HeartbeatInterval        uint8
	AdvertisementFWUInterval uint8
	BatteryVoltage           uint16
	FirmwareHash             [4]uint8
}
