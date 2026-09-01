package tagsl

import (
	"reflect"

	"github.com/truvami/decoder/pkg/common"
)

type Port128Payload struct {
	Ble                    bool   `json:"ble"`
	Gnss                   bool   `json:"gnss"`
	Wifi                   bool   `json:"wifi"`
	MovingInterval         uint32 `json:"movingInterval" validate:"gte=5,lte=86400"`
	SteadyInterval         uint32 `json:"steadyInterval" validate:"gte=120,lte=86400"`
	ConfigInterval         uint32 `json:"configInterval" validate:"gte=300,lte=604800"`
	GnssTimeout            uint16 `json:"gnssTimeout" validate:"gte=60,lte=86400"`
	AccelerometerThreshold uint16 `json:"accelerometerThreshold" validate:"gte=10,lte=8000"`
	AccelerometerDelay     uint16 `json:"accelerometerDelay" validate:"gte=1000,lte=10000"`
	BatteryInterval        uint32 `json:"batteryInterval" validate:"gte=300,lte=604800"`
	BatchSize              uint16 `json:"batchSize" validate:"lte=50"`
	BufferSize             uint16 `json:"bufferSize" validate:"gte=128,lte=8128"`
}

func Port128PayloadConfig() common.PayloadConfig {
	return common.PayloadConfig{
		Fields: []common.FieldConfig{
			{Name: "Ble", Start: 0, Length: 1},
			{Name: "Gnss", Start: 1, Length: 1},
			{Name: "Wifi", Start: 2, Length: 1},
			{Name: "MovingInterval", Start: 3, Length: 4},
			{Name: "SteadyInterval", Start: 7, Length: 4},
			{Name: "ConfigInterval", Start: 11, Length: 4},
			{Name: "GnssTimeout", Start: 15, Length: 2},
			{Name: "AccelerometerThreshold", Start: 17, Length: 2},
			{Name: "AccelerometerDelay", Start: 19, Length: 2},
			{Name: "BatteryInterval", Start: 21, Length: 4},
			{Name: "BatchSize", Start: 25, Length: 2, Optional: true},
			{Name: "BufferSize", Start: 27, Length: 2, Optional: true},
		},
		TargetType: reflect.TypeOf(Port128Payload{}),
	}
}
