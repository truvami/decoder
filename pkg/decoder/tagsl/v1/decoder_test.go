package tagsl

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected interface{}
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			port:    1,
			expected: Port1Payload{
				Moving:    false,
				Latitude:  47.041811,
				Longitude: 7.622494,
				Altitude:  5728,
				Year:      24,
				Month:     4,
				Day:       11,
				Hour:      20,
				Minute:    52,
				Second:    26,
			},
		},
		{
			payload: "00",
			port:    2,
			expected: Port2Payload{
				Moving: false,
			},
		},
		{
			payload: "01",
			port:    2,
			expected: Port2Payload{
				Moving: true,
			},
		},
		{
			payload: "808c59c3c99fc0ad",
			port:    5,
			expected: Port5Payload{
				Moving: false,
				Mac1:   "8c59c3c99fc0",
				Rssi1:  -83,
			},
		},
		{
			payload: "80e0286d8a2742a1",
			port:    5,
			expected: Port5Payload{
				Moving: false,
				Mac1:   "e0286d8a2742",
				Rssi1:  -95,
			},
		},
		{
			payload: "001f3fd57cecb4f0b0140c96bbb2e0286d8a9478b8",
			port:    5,
			expected: Port5Payload{
				Moving: false,
				Mac1:   "1f3fd57cecb4",
				Rssi1:  -16,
				Mac2:   "b0140c96bbb2",
				Rssi2:  -32,
				Mac3:   "286d8a9478b8",
				Rssi3:  0,
			},
		},
		{
			payload: "822f0101f052fab920feafd0e4158b38b9afe05994cb2f5cb2",
			port:    3,
			expected: Port3Payload{
				ScanPointer:    33327,
				TotalMessages:  1,
				CurrentMessage: 1,
				Mac1:           "f052fab920fe",
				Rssi1:          -81,
				Mac2:           "d0e4158b38b9",
				Rssi2:          -81,
				Mac3:           "e05994cb2f5c",
				Rssi3:          -78,
			},
		},
		{
			payload: "01eb0101f052fab920feadd0e4158b38b9afe05994cb2f5cad",
			port:    3,
			expected: Port3Payload{
				ScanPointer:    491,
				TotalMessages:  1,
				CurrentMessage: 1,
				Mac1:           "f052fab920fe",
				Rssi1:          -83,
				Mac2:           "d0e4158b38b9",
				Rssi2:          -81,
				Mac3:           "e05994cb2f5c",
				Rssi3:          -83,
			},
		},
		{
			payload: "0000012c00000e1000001c200078012c05dc02020100010200002328",
			port:    4,
			expected: Port4Payload{
				LocalizationIntervalWhileMoving: 300,
				LocalizationIntervalWhileSteady: 3600,
				HeartbeatInterval:               7200,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				DeviceState:                     2,
				FirmwareVersionMajor:            2,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            0,
				HardwareVersionType:             1,
				HardwareVersionRevision:         2,
				BatteryKeepAliveMessageInterval: 9000,
			},
		},
		{
			payload: "01",
			port:    6,
			expected: Port6Payload{
				ButtonPressed: true,
			},
		},
		{
			payload: "00",
			port:    6,
			expected: Port6Payload{
				ButtonPressed: false,
			},
		},
		{
			payload: "800ee5",
			port:    15,
			expected: Port15Payload{
				LowBattery:     false,
				BatteryVoltage: 3.813,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder()
			got, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode("00", 0, "")
	if err == nil {
		t.Fatal("expected port not supported")
	}
}
