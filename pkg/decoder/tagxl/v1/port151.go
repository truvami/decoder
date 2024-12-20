package tagxl

// The
// [TLV](https://lora-developers.semtech.com/documentation/tech-papers-and-guides/lora-edge-tracker-reference-design-user-guide-v1/lora-edge-tracker-reference-design-firmware/#payload-format-specification)
// payload format shall be used. It is defined as:
//
// | Tag  | Length                    | Number of Commands           | Data                   |
// |------|---------------------------|------------------------------|------------------------|
// | 0x4C | Variable number (uint8_t) | Variable number (uint8_t)    | Commands in TLV format |
//
// The settings uplink contains one or more settings in its data section. These settings are again in TLV format and are either device or runner settings:
//
// | # | Device Setting                   | Tag  | Size | Data                                                                                                                     | Format      |
// |---|----------------------------------|------|------|--------------------------------------------------------------------------------------------------------------------------|-------------|
// | 1 | Device Flags                     | 0x40 | 0x01 | bit 0: BLE_FWU_ENABLED<br/>bit 1: GNSS_ENABLE<br/>bit 2: WIFI_ENABLE<br/>bit 3: Set ACCELERATION_ENABLE<br/>bit 4-7: RFU | bit field   |
// | 2 | Asset Tracking Intervals         | 0x41 | 0x04 | data 0: MOVING_INTERVAL<br/>data 1: STEADY_INTERVAL                                                                      | uint16_t[2] |
// | 3 | Acceleration Sensor Settings     | 0x42 | 0x04 | data 0: ACCELERATION_SENSITIVITY<br/>data 1: ACCELERATION_DELAY                                                          | uint16_t[2] |
// | 4 | HEARTBEAT_INTERVAL               | 0x43 | 0x01 | Heartbeat interval in hours                                                                                              | uint8_t     |
// | 5 | ADVERTISEMENT_FWU_INTERVAL       | 0x44 | 0x01 | Value in seconds                                                                                                         | uint8_t     |
// | 6 | Battery Voltage                  | 0x45 | 0x02 | Battery voltage in mV                                                                                                    | uint16_t    |
// | 7 | Firmware Hash                    | 0x46 | 0x04 | First 4 bytes of SHA-1 hash of git commit                                                                                | uint8_t[4]  |
//
// | # | Runner Setting                   | Tag  | Size | Data                                                                                                                            | Format     |
// |---|----------------------------------|------|------|---------------------------------------------------------------------------------------------------------------------------------|------------|
// | 1 | Run Alarm                        | 0x80 | 0x02 | data 0: Duration of started alarm in minutes (min: 0, max: 255)<br/>data 1: Period of alarm beeps in seconds (min: 0, max: 255) | uint8_t[2] |
type Port151Payload struct {
	DeviceFlags              uint8    `json:"deviceFlags"`
	AssetTrackingIntervals   []uint16 `json:"assetTrackingIntervals"`
	AccelerationSensor       []uint16 `json:"accelerationSensor"`
	HeartbeatInterval        uint8    `json:"heartbeatInterval"`
	AdvertisementFwuInterval uint8    `json:"advertisementFwuInterval"`
	Battery                  float64  `json:"battery" validate:"gte=1,lte=5"`
	FirmwareHash             []uint8  `json:"firmwareHash"`
}
