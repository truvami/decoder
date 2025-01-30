package tagsl

type Port134Payload struct {
	ScanInterval                  uint16 `json:"scanInterval"`
	ScanTime                      uint8  `json:"scanTime" validate:"lte=180"`
	MaxBeacons                    uint8  `json:"maxBeacons"`
	MinRssi                       int8   `json:"minRssi"`
	AdvertisingName               []byte `json:"advertisingName" validate:"max=9"`
	AccelerometerTriggerHoldTimer uint16 `json:"accelerometerTriggerHoldTimer"`
	AcceleratorThreshold          uint16 `json:"acceleratorThreshold"`
	ScanMode                      uint8  `json:"scanMode" validate:"lte=2"`
	BleConfigUplinkInterval       uint16 `json:"bleConfigUplinkInterval"`
}
