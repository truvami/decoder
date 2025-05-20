package tagsl

import (
	"errors"
	"fmt"
	"testing"
	"time"

	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		data     any
		port     uint8
		expected string
	}{
		{
			data: tagsl.Port1Payload{
				Moving:    false,
				Latitude:  46.63858,
				Longitude: 10.39973,
				Altitude:  2909,
				Year:      16,
				Month:     8,
				Day:       32,
				Hour:      12,
				Minute:    56,
				Second:    34,
			},
			port:     1,
			expected: "0002c7a5f4009eaff271a21008200c3822",
		},
		{
			data: tagsl.Port1Payload{
				Moving:    true,
				Latitude:  46.60000,
				Longitude: 10.41667,
				Altitude:  3033,
				Year:      20,
				Month:     4,
				Day:       45,
				Hour:      46,
				Minute:    32,
				Second:    9,
			},
			port:     1,
			expected: "0102c70f40009ef21e767a14042d2e2009",
		},
		{
			data: tagsl.Port4Payload{
				LocalizationIntervalWhileMoving: 60,
				LocalizationIntervalWhileSteady: 120,
				HeartbeatInterval:               3600,
				GPSTimeoutWhileWaitingForFix:    180,
				AccelerometerWakeupThreshold:    20,
				AccelerometerDelay:              2000,
				DeviceState:                     1,
				FirmwareVersionMajor:            4,
				FirmwareVersionMinor:            5,
				FirmwareVersionPatch:            2,
				HardwareVersionType:             2,
				HardwareVersionRevision:         7,
				BatteryKeepAliveMessageInterval: 600,
				BatchSize:                       16,
				BufferSize:                      4096,
			},
			port:     4,
			expected: "0000003c0000007800000e1000b4001407d00104050202070000025800101000",
		},
		{
			data: tagsl.Port4Payload{
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 600,
				HeartbeatInterval:               3600,
				GPSTimeoutWhileWaitingForFix:    240,
				AccelerometerWakeupThreshold:    10,
				AccelerometerDelay:              1000,
				DeviceState:                     0,
				FirmwareVersionMajor:            4,
				FirmwareVersionMinor:            2,
				FirmwareVersionPatch:            0,
				HardwareVersionType:             2,
				HardwareVersionRevision:         4,
				BatteryKeepAliveMessageInterval: 900,
				BatchSize:                       32,
				BufferSize:                      8128,
			},
			port:     4,
			expected: "000000780000025800000e1000f0000a03e80004020002040000038400201fc0",
		},
		{
			data: tagsl.Port5Payload{
				Moving: false,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
			},
			port:     5,
			expected: "00e1c384f2ab5dc0",
		},
		{
			data: tagsl.Port5Payload{
				Moving: true,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
				Mac2:   "f6a8d09c3e72",
				Rssi2:  -72,
			},
			port:     5,
			expected: "01e1c384f2ab5dc0f6a8d09c3e72b8",
		},
		{
			data: tagsl.Port5Payload{
				Moving: true,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
				Mac2:   "f6a8d09c3e72",
				Rssi2:  -72,
				Mac3:   "d9475b62c801",
				Rssi3:  -80,
			},
			port:     5,
			expected: "01e1c384f2ab5dc0f6a8d09c3e72b8d9475b62c801b0",
		},
		{
			data: tagsl.Port5Payload{
				Moving: true,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
				Mac2:   "f6a8d09c3e72",
				Rssi2:  -72,
				Mac3:   "d9475b62c801",
				Rssi3:  -80,
				Mac4:   "0a3ed14bf69c",
				Rssi4:  -88,
			},
			port:     5,
			expected: "01e1c384f2ab5dc0f6a8d09c3e72b8d9475b62c801b00a3ed14bf69ca8",
		},
		{
			data: tagsl.Port5Payload{
				Moving: true,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
				Mac2:   "f6a8d09c3e72",
				Rssi2:  -72,
				Mac3:   "d9475b62c801",
				Rssi3:  -80,
				Mac4:   "0a3ed14bf69c",
				Rssi4:  -88,
				Mac5:   "bc2e90da1473",
				Rssi5:  -96,
			},
			port:     5,
			expected: "01e1c384f2ab5dc0f6a8d09c3e72b8d9475b62c801b00a3ed14bf69ca8bc2e90da1473a0",
		},
		{
			data: tagsl.Port5Payload{
				Moving: true,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
				Mac2:   "f6a8d09c3e72",
				Rssi2:  -72,
				Mac3:   "d9475b62c801",
				Rssi3:  -80,
				Mac4:   "0a3ed14bf69c",
				Rssi4:  -88,
				Mac5:   "bc2e90da1473",
				Rssi5:  -96,
				Mac6:   "3f5a7cc0b6e8",
				Rssi6:  -104,
			},
			port:     5,
			expected: "01e1c384f2ab5dc0f6a8d09c3e72b8d9475b62c801b00a3ed14bf69ca8bc2e90da1473a03f5a7cc0b6e898",
		},
		{
			data: tagsl.Port5Payload{
				Moving: true,
				Mac1:   "e1c384f2ab5d",
				Rssi1:  -64,
				Mac2:   "f6a8d09c3e72",
				Rssi2:  -72,
				Mac3:   "d9475b62c801",
				Rssi3:  -80,
				Mac4:   "0a3ed14bf69c",
				Rssi4:  -88,
				Mac5:   "bc2e90da1473",
				Rssi5:  -96,
				Mac6:   "3f5a7cc0b6e8",
				Rssi6:  -104,
				Mac7:   "a4d38e27fc60",
				Rssi7:  -112,
			},
			port:     5,
			expected: "01e1c384f2ab5dc0f6a8d09c3e72b8d9475b62c801b00a3ed14bf69ca8bc2e90da1473a03f5a7cc0b6e898a4d38e27fc6090",
		},
		{
			data: tagsl.Port6Payload{
				ButtonPressed: false,
			},
			port:     6,
			expected: "00",
		},
		{
			data: tagsl.Port6Payload{
				ButtonPressed: true,
			},
			port:     6,
			expected: "01",
		},
		{
			data: tagsl.Port7Payload{
				Timestamp: time.Date(1984, 4, 19, 0, 0, 0, 0, time.UTC),
				Moving:    false,
				Mac1:      "fa6d293c851b",
				Rssi1:     -48,
			},
			port:     7,
			expected: "1ae4790000fa6d293c851bd0",
		},
		{
			data: tagsl.Port7Payload{
				Timestamp: time.Date(1996, 7, 3, 0, 0, 0, 0, time.UTC),
				Moving:    true,
				Mac1:      "fa6d293c851b",
				Rssi1:     -48,
				Mac2:      "0e42c97a1f64",
				Rssi2:     -56,
			},
			port:     7,
			expected: "31d9b80001fa6d293c851bd00e42c97a1f64c8",
		},
		{
			data: tagsl.Port7Payload{
				Timestamp: time.Date(2004, 12, 24, 0, 0, 0, 0, time.UTC),
				Moving:    false,
				Mac1:      "fa6d293c851b",
				Rssi1:     -48,
				Mac2:      "0e42c97a1f64",
				Rssi2:     -56,
				Mac3:      "b3885e902da7",
				Rssi3:     -64,
			},
			port:     7,
			expected: "41cb5c0000fa6d293c851bd00e42c97a1f64c8b3885e902da7c0",
		},
		{
			data: tagsl.Port7Payload{
				Timestamp: time.Date(2011, 5, 31, 0, 0, 0, 0, time.UTC),
				Moving:    true,
				Mac1:      "fa6d293c851b",
				Rssi1:     -48,
				Mac2:      "0e42c97a1f64",
				Rssi2:     -56,
				Mac3:      "b3885e902da7",
				Rssi3:     -64,
				Mac4:      "4cd29176ab0f",
				Rssi4:     -72,
			},
			port:     7,
			expected: "4de42f8001fa6d293c851bd00e42c97a1f64c8b3885e902da7c04cd29176ab0fb8",
		},
		{
			data: tagsl.Port7Payload{
				Timestamp: time.Date(2018, 8, 28, 0, 0, 0, 0, time.UTC),
				Moving:    false,
				Mac1:      "fa6d293c851b",
				Rssi1:     -48,
				Mac2:      "0e42c97a1f64",
				Rssi2:     -56,
				Mac3:      "b3885e902da7",
				Rssi3:     -64,
				Mac4:      "4cd29176ab0f",
				Rssi4:     -72,
				Mac5:      "a81b3def09cd",
				Rssi5:     -80,
			},
			port:     7,
			expected: "5b84908000fa6d293c851bd00e42c97a1f64c8b3885e902da7c04cd29176ab0fb8a81b3def09cdb0",
		},
		{
			data: tagsl.Port7Payload{
				Timestamp: time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
				Moving:    true,
				Mac1:      "fa6d293c851b",
				Rssi1:     -48,
				Mac2:      "0e42c97a1f64",
				Rssi2:     -56,
				Mac3:      "b3885e902da7",
				Rssi3:     -64,
				Mac4:      "4cd29176ab0f",
				Rssi4:     -72,
				Mac5:      "a81b3def09cd",
				Rssi5:     -80,
				Mac6:      "3fe478115062",
				Rssi6:     -88,
			},
			port:     7,
			expected: "696ec58001fa6d293c851bd00e42c97a1f64c8b3885e902da7c04cd29176ab0fb8a81b3def09cdb03fe478115062a8",
		},
		{
			data: tagsl.Port10Payload{
				Moving:     false,
				Latitude:   46.5372,
				Longitude:  8.1286,
				Altitude:   4274,
				Timestamp:  time.Date(2002, 5, 10, 0, 0, 0, 0, time.UTC),
				Battery:    3.780,
				TTF:        time.Duration(24) * time.Second,
				PDOP:       1.0,
				Satellites: 8,
			},
			port:     10,
			expected: "0002c619f0007c0858a6f43cdb0d800ec4180208",
		},
		{
			data: tagsl.Port10Payload{
				Moving:     true,
				Latitude:   45.9763,
				Longitude:  7.6586,
				Altitude:   4478,
				Timestamp:  time.Date(2015, 7, 14, 0, 0, 0, 0, time.UTC),
				Battery:    3.835,
				TTF:        time.Duration(37) * time.Second,
				PDOP:       2.5,
				Satellites: 9,
			},
			port:     10,
			expected: "0102bd8aec0074dc68aeec55a451000efb250509",
		},
		{
			data: tagsl.Port10Payload{
				Moving:     false,
				Latitude:   -3.0674,
				Longitude:  37.3556,
				Altitude:   5895,
				Timestamp:  time.Date(1995, 10, 1, 0, 0, 0, 0, time.UTC),
				Battery:    3.895,
				TTF:        time.Duration(73) * time.Second,
				PDOP:       6.5,
				Satellites: 10,
			},
			port:     10,
			expected: "00ffd131f8023a0050e646306dda000f37490d0a",
		},
		{
			data: tagsl.Port10Payload{
				Moving:     true,
				Latitude:   -15.5656,
				Longitude:  -72.6467,
				Altitude:   6425,
				Timestamp:  time.Date(1990, 12, 1, 0, 0, 0, 0, time.UTC),
				Battery:    3.960,
				TTF:        time.Duration(192) * time.Second,
				PDOP:       21.0,
				Satellites: 12,
			},
			port:     10,
			expected: "01ff127ce0fbab7fd4fafa2756f2800f78c02a0c",
		},
		{
			data: tagsl.Port15Payload{
				LowBattery: false,
				Battery:    3.895,
			},
			port:     15,
			expected: "000f37",
		},
		{
			data: tagsl.Port15Payload{
				LowBattery: true,
				Battery:    3.920,
			},
			port:     15,
			expected: "010f50",
		},
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
