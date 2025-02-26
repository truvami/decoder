package tagsl

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
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
		{
			data: Port134Payload{
				ScanInterval:                  300,
				ScanTime:                      60,
				MaxBeacons:                    8,
				MinRssi:                       -24,
				AdvertisingName:               []byte("deadbeef"),
				AccelerometerTriggerHoldTimer: 2000,
				AcceleratorThreshold:          1000,
				ScanMode:                      0,
				BleConfigUplinkInterval:       21600,
			},
			port:     134,
			expected: "012c3c08e86465616462656566000007d003e8005460",
		},
		{
			data: Port134Payload{
				ScanInterval:                  900,
				ScanTime:                      120,
				MaxBeacons:                    16,
				MinRssi:                       -20,
				AdvertisingName:               []byte("hello-world"),
				AccelerometerTriggerHoldTimer: 4000,
				AcceleratorThreshold:          2000,
				ScanMode:                      2,
				BleConfigUplinkInterval:       43200,
			},
			port:     134,
			expected: "03847810ec68656c6c6f2d776f72000fa007d002a8c0",
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

func TestInvalidData(t *testing.T) {
	encoder := NewTagSLv1Encoder()
	_, _, err := encoder.Encode(nil, 128, "")
	if err == nil || err.Error() != "data must be a struct" {
		t.Fatal("expected data must be a struct")
	}
}

func TestInvalidPort(t *testing.T) {
	encoder := NewTagSLv1Encoder()
	_, _, err := encoder.Encode(nil, 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestNewTagSLv1Encoder(t *testing.T) {
	// Test with no options
	encoder := NewTagSLv1Encoder()
	if encoder == nil {
		t.Fatal("expected encoder to be created")
	}
	
	// Test with options
	optionCalled := false
	option := func(e *TagSLv1Encoder) {
		optionCalled = true
	}
	
	encoder = NewTagSLv1Encoder(option)
	if !optionCalled {
		t.Fatal("expected option to be called")
	}
}

func TestTagSLv1EncoderWithExtraData(t *testing.T) {
	encoder := NewTagSLv1Encoder()
	
	// Test with extra data
	extraData := "extra data"
	data := Port128Payload{
		BLE: 1,
		GPS: 1,
	}
	
	_, metadata, err := encoder.Encode(data, 128, extraData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if metadata != extraData {
		t.Errorf("expected metadata to be %v, got %v", extraData, metadata)
	}
}
