package nomadxs

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
			payload: "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f09a71d2e",
			port:    1,
			expected: Port1Payload{
				Moving:             false,
				Year:               24,
				Month:              7,
				Day:                25,
				Hour:               20,
				Minute:             38,
				Second:             7,
				Latitude:           46.407935,
				Longitude:          6.21577,
				Altitude:           478.8,
				TimeToFix:          36,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
				Temperature:        24.71,
				Pressure:           7470,
			},
		},
		{
			payload: "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f09a71d2e000000000000",
			port:    1,
			expected: Port1Payload{
				Moving:             false,
				Year:               24,
				Month:              7,
				Day:                25,
				Hour:               20,
				Minute:             38,
				Second:             7,
				Latitude:           46.407935,
				Longitude:          6.21577,
				Altitude:           478.8,
				TimeToFix:          36,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
				Temperature:        24.71,
				Pressure:           7470,
				GyroscopeXAxis:     0.0,
				GyroscopeYAxis:     0.0,
				GyroscopeZAxis:     0.0,
				MagnetometerXAxis:  0.0,
				MagnetometerYAxis:  0.0,
				MagnetometerZAxis:  0.0,
			},
		},
		{
			payload: "0000007800000708000151800078012c05dc000100010100000258000002580500000000",
			port:    4,
			expected: Port4Payload{
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 1800,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				FirmwareVersionMajor:            0,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            0,
				BatteryKeepAliveMessageInterval: 600,
				HardwareVersionType:             1,
				HardwareVersionRevision:         1,
				ReJoinInterval:                  600,
				AccuracyEnhancement:             5,
				LightLowerThreshold:             0,
				LightUpperThreshold:             0,
			},
		},
		{
			payload: "010df6",
			port:    15,
			expected: Port15Payload{
				LowBattery:     true,
				BatteryVoltage: 3.574,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXSv1Decoder()
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
	decoder := NewNomadXSv1Decoder()
	_, err := decoder.Decode("00", 0, "")
	if err == nil {
		t.Fatal("expected port not supported")
	}
}
