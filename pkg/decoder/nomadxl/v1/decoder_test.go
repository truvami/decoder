package nomadxl

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
			payload: "00000001fdd5c693000079300001b45d000000000000000000f600000000000000000b3fd7249f4a00420000000000001501",
			port:    101,
			expected: Port101Payload{
				SystemTime:         8553612947,
				UTCDate:            31024,
				UTCTime:            111709,
				Temperature:        24.6,
				Pressure:           0,
				TimeToFix:          36,
				AccelerometerXAxis: 0,
				AccelerometerYAxis: 0,
				AccelerometerZAxis: 0,
				Battery:            2.879,
				BatteryLorawan:     215,
			},
		},
		{
			payload: "0000793000020152004B6076000C838C00003994",
			port:    103,
			expected: Port103Payload{
				UTCDate:   31024,
				UTCTime:   131410,
				Latitude:  49.39894,
				Longitude: 8.20108,
				Altitude:  147.4,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXLv1Decoder()
			got, _, err := decoder.Decode(test.payload, test.port, "", false)
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
	decoder := NewNomadXLv1Decoder()
	_, _, err := decoder.Decode("00", 0, "", false)
	if err == nil {
		t.Fatal("expected port not supported")
	}
}
