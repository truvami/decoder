package tagsl

type Port128Payload struct {
	BLE                             uint8  `json:"ble" validate:"gte=0,lte=1"`
	GPS                             uint8  `json:"gps" validate:"gte=0,lte=1"`
	WIFI                            uint8  `json:"wifi" validate:"gte=0,lte=1"`
	LocalizationIntervalWhileMoving uint32 `json:"localizationIntervalWhileMoving" validate:"gte=60,lte=86400"`
	LocalizationIntervalWhileSteady uint32 `json:"localizationIntervalWhileSteady" validate:"gte=120,lte=86400"`
	HeartbeatInterval               uint32 `json:"heartbeatInterval" validate:"gte=300,lte=604800"`
	GPSTimeoutWhileWaitingForFix    uint16 `json:"gpsTimeoutWhileWaitingForFix" validate:"gte=60,lte=86400"`
	AccelerometerWakeupThreshold    uint16 `json:"accelerometerWakeupThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay              uint16 `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	BatteryKeepAliveMessageInterval uint32 `json:"batteryKeepAliveMessageInterval" validate:"gte=300,lte=604800"`
	BatchSize                       uint16 `json:"batchSize" validate:"lte=50"`
	BufferSize                      uint16 `json:"bufferSize" validate:"gte=128,lte=8128"`
}
