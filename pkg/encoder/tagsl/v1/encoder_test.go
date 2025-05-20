package tagsl

import (
	"errors"
	"fmt"
	"testing"

	helpers "github.com/truvami/decoder/pkg/common"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		data     any
		port     uint8
		expected string
	}{
		{
			data: Port128Payload{
				Ble:                    true,
				Gnss:                   true,
				Wifi:                   true,
				MovingInterval:         3600,
				SteadyInterval:         7200,
				ConfigInterval:         86400,
				GnssTimeout:            120,
				AccelerometerThreshold: 300,
				AccelerometerDelay:     1500,
				BatteryInterval:        21600,
				BatchSize:              10,
				BufferSize:             4096,
			},
			port:     128,
			expected: "01010100000e1000001c20000151800078012c05dc00005460000a1000",
		},
		{
			data: Port128Payload{
				Ble:                    false,
				Gnss:                   true,
				Wifi:                   false,
				MovingInterval:         120,
				SteadyInterval:         300,
				ConfigInterval:         7200,
				GnssTimeout:            60,
				AccelerometerThreshold: 200,
				AccelerometerDelay:     1000,
				BatteryInterval:        3600,
				BatchSize:              10,
				BufferSize:             4096,
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
			data: Port130Payload{
				EraseFlash: false,
			},
			port:     130,
			expected: "00",
		},
		{
			data: Port130Payload{
				EraseFlash: true,
			},
			port:     130,
			expected: "de",
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
			data: Port132Payload{
				EraseFlash: false,
			},
			port:     132,
			expected: "00",
		},
		{
			data: Port132Payload{
				EraseFlash: true,
			},
			port:     132,
			expected: "00",
		},
		{
			data: Port134Payload{
				ScanInterval:            300,
				ScanTime:                60,
				MaxBeacons:              8,
				MinRssi:                 -24,
				AdvertisingName:         []byte("deadbeef"),
				AccelerometerDelay:      2000,
				AccelerometerThreshold:  1000,
				ScanMode:                0,
				BleConfigUplinkInterval: 21600,
			},
			port:     134,
			expected: "012c3c08e86465616462656566000007d003e8005460",
		},
		{
			data: Port134Payload{
				ScanInterval:            900,
				ScanTime:                120,
				MaxBeacons:              16,
				MinRssi:                 -20,
				AdvertisingName:         []byte("hello-world"),
				AccelerometerDelay:      4000,
				AccelerometerThreshold:  2000,
				ScanMode:                2,
				BleConfigUplinkInterval: 43200,
			},
			port:     134,
			expected: "03847810ec68656c6c6f2d776f72000fa007d002a8c0",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.expected), func(t *testing.T) {
			encoder := NewTagSLv1Encoder()
			got, err := encoder.Encode(test.data, test.port)
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
	_, err := encoder.Encode(nil, 128)
	if err == nil || err.Error() != "data must be a struct" {
		t.Fatal("expected data must be a struct")
	}
}

func TestInvalidPort(t *testing.T) {
	encoder := NewTagSLv1Encoder()
	_, err := encoder.Encode(nil, 0)
	if err == nil || !errors.Is(err, helpers.ErrPortNotSupported) {
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

	encoder.Encode(nil, 0)
}
