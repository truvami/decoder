package tagsl

import (
	"fmt"
	"testing"
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		data     interface{}
		port     int16
		expected string
	}{
		{
			data: Port128Payload{
				BLE:                             1,
				GPS:                             1,
				WIFI:                            1,
				LocalizationIntervalWhileMoving: 3600,
				LocalizationIntervalWhileSteady: 7200,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				BatteryKeepAliveMessageInterval: 21600,
				BatchSize:                       10,
				BufferSize:                      4096,
			},
			port:     128,
			expected: "01010100000e1000001c20000151800078012c05dc00005460000a1000",
		},
		{
			data: Port128Payload{
				BLE:                             0,
				GPS:                             1,
				WIFI:                            0,
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 300,
				HeartbeatInterval:               7200,
				GPSTimeoutWhileWaitingForFix:    60,
				AccelerometerWakeupThreshold:    200,
				AccelerometerDelay:              1000,
				BatteryKeepAliveMessageInterval: 3600,
				BatchSize:                       10,
				BufferSize:                      4096,
			},
			port:     128,
			expected: "000100000000780000012c00001c20003c00c803e800000e10000a1000",
		},
		{
			data: Port129Payload{
				TimeToBuzz: 0,
			},
			port:     129,
			expected: "00",
		},
		{
			data: Port129Payload{
				TimeToBuzz: 16,
			},
			port:     129,
			expected: "10",
		},
		{
			data: Port129Payload{
				TimeToBuzz: 32,
			},
			port:     129,
			expected: "20",
		},
		{
			data: Port131Payload{
				AccuracyEnhancement: 0,
			},
			port:     131,
			expected: "00",
		},
		{
			data: Port131Payload{
				AccuracyEnhancement: 16,
			},
			port:     131,
			expected: "10",
		},
		{
			data: Port131Payload{
				AccuracyEnhancement: 32,
			},
			port:     131,
			expected: "20",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.expected), func(t *testing.T) {
			encoder := NewTagSLv1Encoder()
			got, _, err := encoder.Encode(test.data, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != test.expected {
				t.Errorf("expected: %v\ngot: %v", test.expected, got)
			}
		})
	}
}
