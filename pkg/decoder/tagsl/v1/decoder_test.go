package tagsl

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload     string
		port        int16
		autoPadding bool
		expected    interface{}
	}{
		{
			payload:     "8002cdcd1300744f5e166018040b14341a",
			port:        1,
			autoPadding: false,
			expected: Port1Payload{
				Latitude:  47.041811,
				Longitude: 7.622494,
				Altitude:  572.8,
				Year:      24,
				Month:     4,
				Day:       11,
				Hour:      20,
				Minute:    52,
				Second:    26,
			},
		},
		{
			payload:     "8002cdcd1300744f5e166018040b14341a",
			port:        1,
			autoPadding: true,
			expected: Port1Payload{
				Latitude:  47.041811,
				Longitude: 7.622494,
				Altitude:  572.8,
				Year:      24,
				Month:     4,
				Day:       11,
				Hour:      20,
				Minute:    52,
				Second:    26,
			},
		},
		{
			payload:     "00",
			port:        2,
			autoPadding: false,
			expected:    Port2Payload{},
		},
		{
			payload:     "01",
			port:        2,
			autoPadding: false,
			expected:    Port2Payload{},
		},
		{
			payload:     "822f0101f052fab920feafd0e4158b38b9afe05994cb2f5cb2",
			port:        3,
			autoPadding: false,
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
			payload:     "01eb0101f052fab920feadd0e4158b38b9afe05994cb2f5cad",
			port:        3,
			autoPadding: false,
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
			payload:     "01eb0101",
			port:        3,
			autoPadding: false,
			expected: Port3Payload{
				ScanPointer:    491,
				TotalMessages:  1,
				CurrentMessage: 1,
			},
		},
		{
			payload:     "1eb0101",
			port:        3,
			autoPadding: true,
			expected: Port3Payload{
				ScanPointer:    491,
				TotalMessages:  1,
				CurrentMessage: 1,
			},
		},
		{
			payload:     "0000012c00000e1000001c200078012c05dc02020100010200002328",
			port:        4,
			autoPadding: false,
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
			payload:     "0000003c0000012c000151800078012c05dc02020100010200005460",
			port:        4,
			autoPadding: false,
			expected: Port4Payload{
				LocalizationIntervalWhileMoving: 60,
				LocalizationIntervalWhileSteady: 300,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				DeviceState:                     2,
				FirmwareVersionMajor:            2,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            0,
				HardwareVersionType:             1,
				HardwareVersionRevision:         2,
				BatteryKeepAliveMessageInterval: 21600,
			},
		},
		{
			payload:     "3c0000012c000151800078012c05dc02020100010200005460",
			port:        4,
			autoPadding: true,
			expected: Port4Payload{
				LocalizationIntervalWhileMoving: 60,
				LocalizationIntervalWhileSteady: 300,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				DeviceState:                     2,
				FirmwareVersionMajor:            2,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            0,
				HardwareVersionType:             1,
				HardwareVersionRevision:         2,
				BatteryKeepAliveMessageInterval: 21600,
			},
		},
		{
			payload:     "808c59c3c99fc0ad",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Mac1:  "8c59c3c99fc0",
				Rssi1: -83,
			},
		},
		{
			payload:     "80e0286d8a2742a1",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Mac1:  "e0286d8a2742",
				Rssi1: -95,
			},
		},
		{
			payload:     "001f3fd57cecb4f0b0140c96bbb2e0286d8a9478b8",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Mac1:  "1f3fd57cecb4",
				Rssi1: -16,
				Mac2:  "b0140c96bbb2",
				Rssi2: -32,
				Mac3:  "286d8a9478b8",
				Rssi3: 0,
			},
		},
		{
			payload:     "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Mac1:  "e0286d8aabfc",
				Rssi1: -88,
				Mac2:  "e0286d8a9478",
				Rssi2: -62,
				Mac3:  "726c9a74b58d",
				Rssi3: -85,
				Mac4:  "726cdac8b89d",
				Rssi4: -84,
				Mac5:  "f0b0140c96bb",
				Rssi5: -56,
			},
		},
		{
			payload:     "00",
			port:        5,
			autoPadding: false,
			expected:    Port5Payload{},
		},
		{
			payload:     "",
			port:        5,
			autoPadding: true,
			expected:    Port5Payload{},
		},
		{
			payload:     "01",
			port:        6,
			autoPadding: false,
			expected: Port6Payload{
				ButtonPressed: true,
			},
		},
		{
			payload:     "1",
			port:        6,
			autoPadding: true,
			expected: Port6Payload{
				ButtonPressed: true,
			},
		},
		{
			payload:     "00",
			port:        6,
			autoPadding: false,
			expected: Port6Payload{
				ButtonPressed: false,
			},
		},
		{
			payload:     "66ec04bb00e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42ad",
			port:        7,
			autoPadding: false,
			expected: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				Mac1:      "e0286d8aabfc",
				Rssi1:     -69,
				Mac2:      "ec6c9a74b58f",
				Rssi2:     -78,
				Mac3:      "726c9a74b58d",
				Rssi3:     -79,
				Mac4:      "e0286d8a9478",
				Rssi4:     -53,
				Mac5:      "f0b0140c96bb",
				Rssi5:     -46,
				Mac6:      "260122180d42",
				Rssi6:     -83,
			},
		},
		{
			payload:     "66ec04bb00",
			port:        7,
			autoPadding: false,
			expected: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
			},
		},
		{
			payload:     "66ec04bb00",
			port:        7,
			autoPadding: true,
			expected: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
			},
		},
		{
			payload:     "012c141e9c455738304543434343460078012c01a8c0",
			port:        8,
			autoPadding: false,
			expected: Port8Payload{
				ScanInterval:                          300,
				ScanTime:                              20,
				MaxBeacons:                            30,
				MinRssiValue:                          -100,
				AdvertisingFilter:                     "4048812220199682886",
				AccelerometerTriggerHoldTimer:         120,
				AccelerometerThreshold:                300,
				BLECurrentConfigurationUplinkInterval: 43200,
				ScanMode:                              1,
			},
		},
		{
			payload:     "12c141e9c455738304543434343460078012c01a8c0",
			port:        8,
			autoPadding: true,
			expected: Port8Payload{
				ScanInterval:                          300,
				ScanTime:                              20,
				MaxBeacons:                            30,
				MinRssiValue:                          -100,
				AdvertisingFilter:                     "4048812220199682886",
				AccelerometerTriggerHoldTimer:         120,
				AccelerometerThreshold:                300,
				BLECurrentConfigurationUplinkInterval: 43200,
				ScanMode:                              1,
			},
		},
		{
			payload:     "0002d308b50082457f16eb66c4a5cd0ed3",
			port:        10,
			autoPadding: false,
			expected: Port10Payload{
				Latitude:  47.384757,
				Longitude: 8.537471,
				Altitude:  586.7,
				Timestamp: time.Date(2024, 8, 20, 14, 18, 53, 0, time.UTC),
				Battery:   3.795,
			},
		},
		{
			payload:     "0002D30B070082491F11256718D9FE0EDE190505",
			port:        10,
			autoPadding: false,
			expected: Port10Payload{
				Latitude:   47.385351,
				Longitude:  8.538399,
				Altitude:   438.9,
				Timestamp:  time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:    3.806,
				PDOP:       2.5,
				Satellites: 5,
				TTF:        25,
			},
		},
		{
			payload:     "0002d30b070082491f11256718d9fe0ede",
			port:        10,
			autoPadding: false,
			expected: Port10Payload{
				Latitude:  47.385351,
				Longitude: 8.538399,
				Altitude:  438.9,
				Timestamp: time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:   3.806,
			},
		},
		{
			payload:     "2d30b070082491f11256718d9fe0ede",
			port:        10,
			autoPadding: true,
			expected: Port10Payload{
				Latitude:  47.385351,
				Longitude: 8.538399,
				Altitude:  438.9,
				Timestamp: time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:   3.806,
			},
		},
		{
			payload:     "800ee5",
			port:        15,
			autoPadding: false,
			expected: Port15Payload{
				LowBattery: false,
				Battery:    3.813,
			},
		},
		{
			payload:     "001044",
			port:        15,
			autoPadding: false,
			expected: Port15Payload{
				LowBattery: false,
				Battery:    4.164,
			},
		},
		{
			payload:     "1044",
			port:        15,
			autoPadding: true,
			expected: Port15Payload{
				LowBattery: false,
				Battery:    4.164,
			},
		},
		{
			payload:     "0002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:        50,
			autoPadding: false,
			expected: Port50Payload{
				Latitude:  47.385747,
				Longitude: 8.539271,
				Altitude:  447.3,
				Timestamp: time.Date(2024, 8, 20, 9, 11, 41, 0, time.UTC),
				Battery:   3.969,
				TTF:       24,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -87,
				Mac2:      "f0b0140c96bb",
				Rssi2:     -56,
				Mac3:      "726c9a74b58d",
				Rssi3:     -88,
				Mac4:      "e0286d8a9478",
				Rssi4:     -65,
			},
		},
		{
			payload:     "0102d30b2a0082499c10ee66c496900ed34af0b0140c96bbb3e0286d8a9478c3fc848e9b5571c2",
			port:        50,
			autoPadding: false,
			expected: Port50Payload{
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       74,
				Mac1:      "f0b0140c96bb",
				Rssi1:     -77,
				Mac2:      "e0286d8a9478",
				Rssi2:     -61,
				Mac3:      "fc848e9b5571",
				Rssi3:     -62,
			},
		},
		{
			payload:     "0102d30b2a0082499c10ee66c496900ed34a",
			port:        50,
			autoPadding: false,
			expected: Port50Payload{
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       74,
			},
		},
		{
			payload:     "102d30b2a0082499c10ee66c496900ed34a",
			port:        50,
			autoPadding: true,
			expected: Port50Payload{
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       74,
			},
		},
		{
			payload:     "0002D30BA000824ACE1122671B983E0EEA340B06726C9A74B58DB1FCF528F8634FB552A8DB7BD6B5B9E0286D8AABFCBC",
			port:        51,
			autoPadding: false,
			expected: Port51Payload{
				Latitude:   47.385504,
				Longitude:  8.53883,
				Altitude:   438.6,
				Timestamp:  time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:    3.818,
				TTF:        52,
				Mac1:       "726c9a74b58d",
				Rssi1:      -79,
				Mac2:       "fcf528f8634f",
				Rssi2:      -75,
				Mac3:       "52a8db7bd6b5",
				Rssi3:      -71,
				Mac4:       "e0286d8aabfc",
				Rssi4:      -68,
				PDOP:       5.5,
				Satellites: 6,
			},
		},
		{
			payload:     "0002d30ba000824ace1122671b983e0eea340b06",
			port:        51,
			autoPadding: false,
			expected: Port51Payload{
				Latitude:   47.385504,
				Longitude:  8.53883,
				Altitude:   438.6,
				Timestamp:  time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:    3.818,
				TTF:        52,
				PDOP:       5.5,
				Satellites: 6,
			},
		},
		{
			payload:     "2d30ba000824ace1122671b983e0eea340b06",
			port:        51,
			autoPadding: true,
			expected: Port51Payload{
				Latitude:   47.385504,
				Longitude:  8.53883,
				Altitude:   438.6,
				Timestamp:  time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:    3.818,
				TTF:        52,
				PDOP:       5.5,
				Satellites: 6,
			},
		},
		{
			payload:     "000166c4a5ba00e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0",
			port:        105,
			autoPadding: false,
			expected: Port105Payload{
				BufferLevel: 1,
				Timestamp:   time.Date(2024, 8, 20, 14, 18, 34, 0, time.UTC),
				Mac1:        "e0286d8aabfc",
				Rssi1:       -79,
				Mac2:        "e0286d8a9478",
				Rssi2:       -62,
				Mac3:        "ec6c9a74b58f",
				Rssi3:       -83,
				Mac4:        "726c9a74b58d",
				Rssi4:       -83,
				Mac5:        "f0b0140c96bb",
				Rssi5:       -48,
			},
		},
		{
			payload:     "010166c4a5ba00e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0",
			port:        105,
			autoPadding: false,
			expected: Port105Payload{
				BufferLevel: 257,
				Timestamp:   time.Date(2024, 8, 20, 14, 18, 34, 0, time.UTC),
				Mac1:        "e0286d8aabfc",
				Rssi1:       -79,
				Mac2:        "e0286d8a9478",
				Rssi2:       -62,
				Mac3:        "ec6c9a74b58f",
				Rssi3:       -83,
				Mac4:        "726c9a74b58d",
				Rssi4:       -83,
				Mac5:        "f0b0140c96bb",
				Rssi5:       -48,
			},
		},
		{
			payload:     "001366ee2f4d00c4eb438ddde2a504e31aea1b01a7245a4c7a0d2ec026e98d560d2ebbccd42ef92ed4ae704f5708e1d1b9",
			port:        105,
			autoPadding: false,
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
				Mac1:        "c4eb438ddde2",
				Rssi1:       -91,
				Mac2:        "04e31aea1b01",
				Rssi2:       -89,
				Mac3:        "245a4c7a0d2e",
				Rssi3:       -64,
				Mac4:        "26e98d560d2e",
				Rssi4:       -69,
				Mac5:        "ccd42ef92ed4",
				Rssi5:       -82,
				Mac6:        "704f5708e1d1",
				Rssi6:       -71,
			},
		},
		{
			payload:     "001366ee2f4d00c4",
			port:        105,
			autoPadding: false,
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
			},
		},
		{
			payload:     "1366ee2f4d00c4",
			port:        105,
			autoPadding: true,
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
			},
		},
		{
			payload:     "00020002d309ae008247c5113966c45d640f7e",
			port:        110,
			autoPadding: false,
			expected: Port110Payload{
				BufferLevel: 2,
				Latitude:    47.385006,
				Longitude:   8.538053,
				Altitude:    440.9,
				Timestamp:   time.Date(2024, 8, 20, 9, 9, 56, 0, time.UTC),
				Battery:     3.966,
			},
		},
		{
			payload:     "01020002d309ae008247c5113966c45d640f7e",
			port:        110,
			autoPadding: false,
			expected: Port110Payload{
				BufferLevel: 258,
				Latitude:    47.385006,
				Longitude:   8.538053,
				Altitude:    440.9,
				Timestamp:   time.Date(2024, 8, 20, 9, 9, 56, 0, time.UTC),
				Battery:     3.966,
			},
		},
		{
			payload:     "00040002d30b8c00824a35112266c45c440f83",
			port:        110,
			autoPadding: false,
			expected: Port110Payload{
				BufferLevel: 4,
				Latitude:    47.385484,
				Longitude:   8.538677,
				Altitude:    438.6,
				Timestamp:   time.Date(2024, 8, 20, 9, 5, 8, 0, time.UTC),
				Battery:     3.971,
			},
		},
		{
			payload:     "40002d30b8c00824a35112266c45c440f83",
			port:        110,
			autoPadding: true,
			expected: Port110Payload{
				BufferLevel: 4,
				Latitude:    47.385484,
				Longitude:   8.538677,
				Altitude:    438.6,
				Timestamp:   time.Date(2024, 8, 20, 9, 5, 8, 0, time.UTC),
				Battery:     3.971,
			},
		},
		{
			payload:     "00020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:        150,
			autoPadding: false,
			expected: Port150Payload{
				BufferLevel: 2,
				Latitude:    47.385747,
				Longitude:   8.539271,
				Altitude:    447.3,
				Timestamp:   time.Date(2024, 8, 20, 9, 11, 41, 0, time.UTC),
				Battery:     3.969,
				TTF:         24,
				Mac1:        "e0286d8aabfc",
				Rssi1:       -87,
				Mac2:        "f0b0140c96bb",
				Rssi2:       -56,
				Mac3:        "726c9a74b58d",
				Rssi3:       -88,
				Mac4:        "e0286d8a9478",
				Rssi4:       -65,
			},
		},
		{
			payload:     "01020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:        150,
			autoPadding: false,
			expected: Port150Payload{
				BufferLevel: 258,
				Latitude:    47.385747,
				Longitude:   8.539271,
				Altitude:    447.3,
				Timestamp:   time.Date(2024, 8, 20, 9, 11, 41, 0, time.UTC),
				Battery:     3.969,
				TTF:         24,
				Mac1:        "e0286d8aabfc",
				Rssi1:       -87,
				Mac2:        "f0b0140c96bb",
				Rssi2:       -56,
				Mac3:        "726c9a74b58d",
				Rssi3:       -88,
				Mac4:        "e0286d8a9478",
				Rssi4:       -65,
			},
		},
		{
			payload:     "00000102d30a98008248b611ac66c45ed80f6b1be0286d0a6f42a3000000000044a4e0286d8a9478bff0b0140c96bbc9",
			port:        150,
			autoPadding: false,
			expected: Port150Payload{
				BufferLevel: 0,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         27,
				Mac1:        "e0286d0a6f42",
				Rssi1:       -93,
				Mac2:        "000000000044",
				Rssi2:       -92,
				Mac3:        "e0286d8a9478",
				Rssi3:       -65,
				Mac4:        "f0b0140c96bb",
				Rssi4:       -55,
			},
		},
		{
			payload:     "00000102d30a98008248b611ac66c45ed80f6b1b",
			port:        150,
			autoPadding: false,
			expected: Port150Payload{
				BufferLevel: 0,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         27,
			},
		},
		{
			payload:     "102d30a98008248b611ac66c45ed80f6b1b",
			port:        150,
			autoPadding: true,
			expected: Port150Payload{
				BufferLevel: 0,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         27,
			},
		},
		{
			payload:     "00000002D30B27008247B81312671BD1641EEF18030BE0286D8A9478CBF0B0140C96BBCE",
			port:        151,
			autoPadding: false,
			expected: Port151Payload{
				BufferLevel: 0,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     7.919,
				TTF:         24,
				Mac1:        "e0286d8a9478",
				Rssi1:       -53,
				Mac2:        "f0b0140c96bb",
				Rssi2:       -50,
				PDOP:        1.5,
				Satellites:  11,
			},
		},
		{
			payload:     "00000002d30b27008247b81312671bd1641eef18030b",
			port:        151,
			autoPadding: false,
			expected: Port151Payload{
				BufferLevel: 0,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     7.919,
				TTF:         24,
				PDOP:        1.5,
				Satellites:  11,
			},
		},
		{
			payload:     "2d30b27008247b81312671bd1641eef18030b",
			port:        151,
			autoPadding: true,
			expected: Port151Payload{
				BufferLevel: 0,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     7.919,
				TTF:         24,
				PDOP:        1.5,
				Satellites:  11,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder(WithAutoPadding(test.autoPadding))
			got, _, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != test.expected {
				t.Errorf("expected: %v\ngot: %v", test.expected, got)
			}
		})
	}

	t.Run("TestInvalidPort", func(t *testing.T) {
		decoder := NewTagSLv1Decoder()
		_, _, err := decoder.Decode("00", 0, "")
		if err == nil {
			t.Fatal("expected port not supported")
		}
	})

	t.Run("TestInvalidPayload", func(t *testing.T) {
		decoder := NewTagSLv1Decoder()
		_, _, err := decoder.Decode("", 1, "")
		if err == nil {
			t.Fatal("expected invalid payload")
		}
	})
}

func TestInvalidPort(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, _, err := decoder.Decode("00", 0, "")
	if err == nil {
		t.Fatal("expected port not supported")
	}
}

func TestParseStatusByte(t *testing.T) {
	tests := []struct {
		input    byte
		expected Status
		err      error
	}{
		{
			input: 0x80,
			expected: Status{
				DutyCycle:           true,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
			err: nil,
		},
		{
			input: 0xFF,
			expected: Status{
				DutyCycle:           true,
				ConfigChangeId:      15,
				ConfigChangeSuccess: true,
				Moving:              true,
			},
			err: nil,
		},
		{
			input: 0x00,
			expected: Status{
				DutyCycle:           false,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
			err: nil,
		},
		{
			input: 0x4A,
			expected: Status{
				DutyCycle:           false,
				ConfigChangeId:      9,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
			err: nil,
		},
		{
			input: 0x8D,
			expected: Status{
				DutyCycle:           true,
				ConfigChangeId:      1,
				ConfigChangeSuccess: true,
				Moving:              true,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestStatusByteInput%v", test.input), func(t *testing.T) {
			got, err := parseStatusByte(test.input)
			if err != nil && test.err == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && test.err != nil {
				t.Fatalf("expected error: %v, got: %v", test.err, err)
			}
			if err != nil && test.err != nil && err.Error() != test.err.Error() {
				t.Fatalf("expected error: %v, got: %v", test.err, err)
			}
			if got != test.expected {
				t.Errorf("expected: %v\ngot: %v", test.expected, got)
			}
		})
	}
}

func TestFullDecode(t *testing.T) {
	tests := []struct {
		payload        string
		expectedData   interface{}
		expectedStatus *Status
		port           int16
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			expectedData: Port1Payload{
				Latitude:  47.041811,
				Longitude: 7.622494,
				Altitude:  572.8,
				Year:      24,
				Month:     4,
				Day:       11,
				Hour:      20,
				Minute:    52,
				Second:    26,
			},
			expectedStatus: &Status{
				DutyCycle:           true,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
			port: 1,
		},
		{
			payload: "66ec04bb00e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42ad",
			expectedData: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				Mac1:      "e0286d8aabfc",
				Rssi1:     -69,
				Mac2:      "ec6c9a74b58f",
				Rssi2:     -78,
				Mac3:      "726c9a74b58d",
				Rssi3:     -79,
				Mac4:      "e0286d8a9478",
				Rssi4:     -53,
				Mac5:      "f0b0140c96bb",
				Rssi5:     -46,
				Mac6:      "260122180d42",
				Rssi6:     -83,
			},
			expectedStatus: &Status{
				DutyCycle:           false,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
			port: 7,
		},
		{
			payload: "0c24651155ce55b8602232e20f52b0ac8ba91fedaaa5603197f93781a90aecdafa5fe8bc02ecdafa5fe8bd8c59c3c960f0a5",
			expectedData: Port5Payload{
				Mac1:  "24651155ce55",
				Rssi1: -72,
				Mac2:  "602232e20f52",
				Rssi2: -80,
				Mac3:  "ac8ba91fedaa",
				Rssi3: -91,
				Mac4:  "603197f93781",
				Rssi4: -87,
				Mac5:  "0aecdafa5fe8",
				Rssi5: -68,
				Mac6:  "02ecdafa5fe8",
				Rssi6: -67,
				Mac7:  "8c59c3c960f0",
				Rssi7: -91,
			},
			expectedStatus: &Status{
				DutyCycle:           false,
				ConfigChangeId:      1,
				ConfigChangeSuccess: true,
				Moving:              false,
			},
			port: 5,
		},
		{
			payload: "00430102d43ffa00772d870ea367250ef60eda",
			expectedData: Port110Payload{
				BufferLevel: 67,
				Latitude:    47.464442,
				Longitude:   7.810439,
				Altitude:    374.7,
				Timestamp:   time.Date(2024, 11, 1, 17, 25, 10, 0, time.UTC),
				Battery:     3.802,
			},
			expectedStatus: &Status{
				DutyCycle:           false,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              true,
			},
			port: 110,
		},
		{
			payload: "0028672658500172a741b1e238b572a741b1e08bb03498b5c583e2b172a741b1e0cda772a741beed4cc472a741beef53b772a741b1dd0000",
			expectedData: Port105Payload{
				BufferLevel: 40,
				Mac1:        "72a741b1e238",
				Rssi1:       -75,
				Mac2:        "72a741b1e08b",
				Rssi2:       -80,
				Mac3:        "3498b5c583e2",
				Rssi3:       -79,
				Mac4:        "72a741b1e0cd",
				Rssi4:       -89,
				Mac5:        "72a741beed4c",
				Rssi5:       -60,
				Mac6:        "72a741beef53",
				Rssi6:       -73,
				Timestamp:   time.Date(2024, 11, 2, 16, 50, 24, 0, time.UTC),
			},
			expectedStatus: &Status{
				DutyCycle:           false,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              true,
			},
			port: 105,
		},
	}

	decoder := NewTagSLv1Decoder()
	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFullDecodeWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			data, status, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if data != test.expectedData {
				t.Errorf("expected data: %v, got: %v", test.expectedData, data)
			}

			if status != nil && test.expectedStatus == nil {
				t.Errorf("expected status to be nil, got: %v", status)
			}
			if status == nil && test.expectedStatus != nil {
				t.Errorf("expected status: %v, got: nil", test.expectedStatus)
			}

			// check if status is equal to expected status using reflect.DeepEqual
			if reflect.DeepEqual(status, *test.expectedStatus) == false {
				t.Errorf("expected status: %v, got: %v", *test.expectedStatus, status)
			}
		})
	}
}
