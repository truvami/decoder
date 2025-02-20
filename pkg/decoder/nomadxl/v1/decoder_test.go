package nomadxl

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload     string
		port        int16
		autoPadding bool
		expected    interface{}
	}{
		{
			payload:     "00000001fdd5c693000079300001b45d000000000000000000d700000000000000000b3fd724",
			port:        101,
			autoPadding: false,
			expected: Port101Payload{
				SystemTime:         8553612947,
				UTCDate:            31024,
				UTCTime:            111709,
				Temperature:        21.5,
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
			payload:     "1fdd5c693000079300001b45d000000000000000000d700000000000000000b3fd724",
			port:        101,
			autoPadding: true,
			expected: Port101Payload{
				SystemTime:         8553612947,
				UTCDate:            31024,
				UTCTime:            111709,
				Temperature:        21.5,
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
			payload:     "0000793000020152004B6076000C838C00003994",
			port:        103,
			autoPadding: false,
			expected: Port103Payload{
				UTCDate:   31024,
				UTCTime:   131410,
				Latitude:  49.39894,
				Longitude: 8.20108,
				Altitude:  147.4,
			},
		},
		{
			payload:     "793000020152004B6076000C838C00003994",
			port:        103,
			autoPadding: true,
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
			decoder := NewNomadXLv1Decoder(WithAutoPadding(test.autoPadding))
			got, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got.Data != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewNomadXLv1Decoder()
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewNomadXLv1Decoder()
	_, err := decoder.Decode("deadbeef", 101, "")

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewNomadXLv1Decoder()
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242", 101, "")

	if err == nil || err.Error() != "payload too long" {
		t.Fatal("expected error payload too long")
	}
}
