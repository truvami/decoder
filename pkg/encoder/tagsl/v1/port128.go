package tagsl

type Port128Payload struct {
	Ble                    bool   `json:"ble"`
	Gnss                   bool   `json:"gnss"`
	Wifi                   bool   `json:"wifi"`
	MovingInterval         uint32 `json:"movingInterval" validate:"gte=60,lte=86400"`
	SteadyInterval         uint32 `json:"steadyInterval" validate:"gte=120,lte=86400"`
	ConfigInterval         uint32 `json:"configInterval" validate:"gte=300,lte=604800"`
	GnssTimeout            uint16 `json:"gnssTimeout" validate:"gte=60,lte=86400"`
	AccelerometerThreshold uint16 `json:"accelerometerThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay     uint16 `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	BatteryInterval        uint32 `json:"batteryInterval" validate:"gte=300,lte=604800"`
	BatchSize              uint16 `json:"batchSize" validate:"lte=50"`
	BufferSize             uint16 `json:"bufferSize" validate:"gte=128,lte=8128"`
}
