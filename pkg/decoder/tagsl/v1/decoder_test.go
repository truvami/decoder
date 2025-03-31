package tagsl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload        string
		port           int16
		autoPadding    bool
		skipValidation bool
		expected       interface{}
	}{
		{
			payload:     "8002cdcd1300744f5e166018040b14341a",
			port:        1,
			autoPadding: false,
			expected: Port1Payload{
				Moving:    false,
				DutyCycle: true,
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
				Moving:    false,
				DutyCycle: true,
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
			payload:        "8002cdcd1300744f5e166018040b14341adeadbeef",
			port:           1,
			autoPadding:    false,
			skipValidation: true,
			expected: Port1Payload{
				Moving:    false,
				DutyCycle: true,
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
			payload:        "00",
			port:           2,
			autoPadding:    false,
			skipValidation: true,
			expected:       Port2Payload{},
		},
		{
			payload:        "01",
			port:           2,
			autoPadding:    false,
			skipValidation: true,
			expected:       Port2Payload{},
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
			payload:        "0000003c0000012c000151800078012c05dc02020100010200005460000a1000",
			port:           4,
			autoPadding:    false,
			skipValidation: true,
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
				BatchSize:                       10,
				BufferSize:                      4096,
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
				Moving:    false,
				DutyCycle: true,
				Mac1:      "8c59c3c99fc0",
				Rssi1:     -83,
			},
		},
		{
			payload:     "80e0286d8a2742a1",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: true,
				Mac1:      "e0286d8a2742",
				Rssi1:     -95,
			},
		},
		{
			payload:     "001f3fd57cecb4f0b0140c96bbb2e0286d8a9478b8",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "1f3fd57cecb4",
				Rssi1:     -16,
				Mac2:      "b0140c96bbb2",
				Rssi2:     -32,
				Mac3:      "286d8a9478b8",
				Rssi3:     0,
			},
		},
		{
			payload:     "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8",
			port:        5,
			autoPadding: false,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -88,
				Mac2:      "e0286d8a9478",
				Rssi2:     -62,
				Mac3:      "726c9a74b58d",
				Rssi3:     -85,
				Mac4:      "726cdac8b89d",
				Rssi4:     -84,
				Mac5:      "f0b0140c96bb",
				Rssi5:     -56,
			},
		},
		{
			payload:        "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8deadbeef4242d6deadbeef4242d6",
			port:           5,
			autoPadding:    false,
			skipValidation: false,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -88,
				Mac2:      "e0286d8a9478",
				Rssi2:     -62,
				Mac3:      "726c9a74b58d",
				Rssi3:     -85,
				Mac4:      "726cdac8b89d",
				Rssi4:     -84,
				Mac5:      "f0b0140c96bb",
				Rssi5:     -56,
				Mac6:      "deadbeef4242",
				Rssi6:     -42,
				Mac7:      "deadbeef4242",
				Rssi7:     -42,
			},
		},
		{
			payload:        "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8deadbeef4242d6deadbeef4242d6deadbeefdeadbeef",
			port:           5,
			autoPadding:    false,
			skipValidation: true,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -88,
				Mac2:      "e0286d8a9478",
				Rssi2:     -62,
				Mac3:      "726c9a74b58d",
				Rssi3:     -85,
				Mac4:      "726cdac8b89d",
				Rssi4:     -84,
				Mac5:      "f0b0140c96bb",
				Rssi5:     -56,
				Mac6:      "deadbeef4242",
				Rssi6:     -42,
				Mac7:      "deadbeef4242",
				Rssi7:     -42,
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
			payload:        "00deadbeef",
			port:           6,
			autoPadding:    false,
			skipValidation: true,
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
				Moving:    false,
				DutyCycle: false,
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
			payload:        "66ec04bb00e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42addeadbeef",
			port:           7,
			autoPadding:    false,
			skipValidation: true,
			expected: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				Moving:    false,
				DutyCycle: false,
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
				Moving:    false,
				DutyCycle: false,
			},
		},
		{
			payload:     "66ec04bb00",
			port:        7,
			autoPadding: true,
			expected: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				Moving:    false,
				DutyCycle: false,
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
			payload:        "012c141e9c455738304543434343460078012c01a8c0deadbeef",
			port:           8,
			autoPadding:    false,
			skipValidation: true,
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
				TTF:        time.Duration(25) * time.Second,
			},
		},
		{
			payload:        "0002d30b070082491f11256718d9fe0ede190505deadbeef",
			port:           10,
			autoPadding:    false,
			skipValidation: true,
			expected: Port10Payload{
				Latitude:   47.385351,
				Longitude:  8.538399,
				Altitude:   438.9,
				Timestamp:  time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:    3.806,
				PDOP:       2.5,
				Satellites: 5,
				TTF:        time.Duration(25) * time.Second,
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
			payload:        "800ee5deadbeef",
			port:           15,
			autoPadding:    false,
			skipValidation: true,
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
				Moving:    false,
				DutyCycle: false,
				Latitude:  47.385747,
				Longitude: 8.539271,
				Altitude:  447.3,
				Timestamp: time.Date(2024, 8, 20, 9, 11, 41, 0, time.UTC),
				Battery:   3.969,
				TTF:       time.Duration(24) * time.Second,
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
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
				Mac1:      "f0b0140c96bb",
				Rssi1:     -77,
				Mac2:      "e0286d8a9478",
				Rssi2:     -61,
				Mac3:      "fc848e9b5571",
				Rssi3:     -62,
			},
		},
		{
			payload:     "0102d30b2a0082499c10ee66c496900ed34af0b0140c96bbb3e0286d8a9478c3fc848e9b5571c2deadbeef4242d6deadbeef4242d6deadbeef4242d6",
			port:        50,
			autoPadding: false,
			expected: Port50Payload{
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
				Mac1:      "f0b0140c96bb",
				Rssi1:     -77,
				Mac2:      "e0286d8a9478",
				Rssi2:     -61,
				Mac3:      "fc848e9b5571",
				Rssi3:     -62,
				Mac4:      "deadbeef4242",
				Rssi4:     -42,
				Mac5:      "deadbeef4242",
				Rssi5:     -42,
				Mac6:      "deadbeef4242",
				Rssi6:     -42,
			},
		}, {
			payload:        "0102d30b2a0082499c10ee66c496900ed34af0b0140c96bbb3e0286d8a9478c3fc848e9b5571c2deadbeef4242d6deadbeef4242d6deadbeef4242d6deadbeef",
			port:           50,
			autoPadding:    false,
			skipValidation: true,
			expected: Port50Payload{
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
				Mac1:      "f0b0140c96bb",
				Rssi1:     -77,
				Mac2:      "e0286d8a9478",
				Rssi2:     -61,
				Mac3:      "fc848e9b5571",
				Rssi3:     -62,
				Mac4:      "deadbeef4242",
				Rssi4:     -42,
				Mac5:      "deadbeef4242",
				Rssi5:     -42,
				Mac6:      "deadbeef4242",
				Rssi6:     -42,
			},
		},
		{
			payload:     "0102d30b2a0082499c10ee66c496900ed34a",
			port:        50,
			autoPadding: false,
			expected: Port50Payload{
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
			},
		},
		{
			payload:     "102d30b2a0082499c10ee66c496900ed34a",
			port:        50,
			autoPadding: true,
			expected: Port50Payload{
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
			},
		},
		{
			payload:     "0002D30BA000824ACE1122671B983E0EEA340B06726C9A74B58DB1FCF528F8634FB552A8DB7BD6B5B9E0286D8AABFCBC",
			port:        51,
			autoPadding: false,
			expected: Port51Payload{
				Moving:     false,
				DutyCycle:  false,
				Latitude:   47.385504,
				Longitude:  8.53883,
				Altitude:   438.6,
				Timestamp:  time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:    3.818,
				TTF:        time.Duration(52) * time.Second,
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
				Moving:     false,
				DutyCycle:  false,
				Latitude:   47.385504,
				Longitude:  8.53883,
				Altitude:   438.6,
				Timestamp:  time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:    3.818,
				TTF:        time.Duration(52) * time.Second,
				PDOP:       5.5,
				Satellites: 6,
			},
		},
		{
			payload:     "2d30ba000824ace1122671b983e0eea340b06",
			port:        51,
			autoPadding: true,
			expected: Port51Payload{
				Moving:     false,
				DutyCycle:  false,
				Latitude:   47.385504,
				Longitude:  8.53883,
				Altitude:   438.6,
				Timestamp:  time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:    3.818,
				TTF:        time.Duration(52) * time.Second,
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
				Moving:      false,
				DutyCycle:   true,
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
				Moving:      false,
				DutyCycle:   true,
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
				Moving:      false,
				DutyCycle:   true,
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
			payload:        "001366ee2f4d00c4eb438ddde2a504e31aea1b01a7245a4c7a0d2ec026e98d560d2ebbccd42ef92ed4ae704f5708e1d1b9deadbeef",
			port:           105,
			autoPadding:    false,
			skipValidation: true,
			expected: Port105Payload{
				BufferLevel: 19,
				Moving:      false,
				DutyCycle:   true,
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
				Moving:      false,
				DutyCycle:   true,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
			},
		},
		{
			payload:     "1366ee2f4d00c4",
			port:        105,
			autoPadding: true,
			expected: Port105Payload{
				BufferLevel: 19,
				Moving:      false,
				DutyCycle:   true,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
			},
		},
		{
			payload:     "00020002d309ae008247c5113966c45d640f7e",
			port:        110,
			autoPadding: false,
			expected: Port110Payload{
				BufferLevel: 2,
				Moving:      false,
				DutyCycle:   false,
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
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385006,
				Longitude:   8.538053,
				Altitude:    440.9,
				Timestamp:   time.Date(2024, 8, 20, 9, 9, 56, 0, time.UTC),
				Battery:     3.966,
			},
		},
		{
			payload:        "01020002d309ae008247c5113966c45d640f7edeadbeef",
			port:           110,
			autoPadding:    false,
			skipValidation: true,
			expected: Port110Payload{
				BufferLevel: 258,
				Moving:      false,
				DutyCycle:   false,
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
				Moving:      false,
				DutyCycle:   false,
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
				Moving:      false,
				DutyCycle:   false,
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
				TTF:         time.Duration(24) * time.Second,
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
				TTF:         time.Duration(24) * time.Second,
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
				Moving:      true,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         time.Duration(27) * time.Second,
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
			payload:     "00000102d30a98008248b611ac66c45ed80f6b1be0286d0a6f42a3000000000044a4e0286d8a9478bff0b0140c96bbc9deadbeef4242d6deadbeef4242d6",
			port:        150,
			autoPadding: false,
			expected: Port150Payload{
				BufferLevel: 0,
				Moving:      true,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         time.Duration(27) * time.Second,
				Mac1:        "e0286d0a6f42",
				Rssi1:       -93,
				Mac2:        "000000000044",
				Rssi2:       -92,
				Mac3:        "e0286d8a9478",
				Rssi3:       -65,
				Mac4:        "f0b0140c96bb",
				Rssi4:       -55,
				Mac5:        "deadbeef4242",
				Rssi5:       -42,
				Mac6:        "deadbeef4242",
				Rssi6:       -42,
			},
		},
		{
			payload:        "00000102d30a98008248b611ac66c45ed80f6b1be0286d0a6f42a3000000000044a4e0286d8a9478bff0b0140c96bbc9deadbeef4242d6deadbeef4242d6deadbeef",
			port:           150,
			autoPadding:    false,
			skipValidation: true,
			expected: Port150Payload{
				BufferLevel: 0,
				Moving:      true,
				DutyCycle:   false,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         time.Duration(27) * time.Second,
				Mac1:        "e0286d0a6f42",
				Rssi1:       -93,
				Mac2:        "000000000044",
				Rssi2:       -92,
				Mac3:        "e0286d8a9478",
				Rssi3:       -65,
				Mac4:        "f0b0140c96bb",
				Rssi4:       -55,
				Mac5:        "deadbeef4242",
				Rssi5:       -42,
				Mac6:        "deadbeef4242",
				Rssi6:       -42,
			},
		},
		{
			payload:     "00000102d30a98008248b611ac66c45ed80f6b1b",
			port:        150,
			autoPadding: false,
			expected: Port150Payload{
				BufferLevel: 0,
				Moving:      true,
				DutyCycle:   false,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         time.Duration(27) * time.Second,
			},
		},
		{
			payload:     "102d30a98008248b611ac66c45ed80f6b1b",
			port:        150,
			autoPadding: true,
			expected: Port150Payload{
				BufferLevel: 0,
				Moving:      true,
				DutyCycle:   false,
				Latitude:    47.38524,
				Longitude:   8.538294,
				Altitude:    452.4,
				Timestamp:   time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:     3.947,
				TTF:         time.Duration(27) * time.Second,
			},
		},
		{
			payload:     "00000002D30B27008247B81312671BD164133718030BE0286D8A9478CBF0B0140C96BBCE",
			port:        151,
			autoPadding: false,
			expected: Port151Payload{
				BufferLevel: 0,
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     4.919,
				TTF:         time.Duration(24) * time.Second,
				Mac1:        "e0286d8a9478",
				Rssi1:       -53,
				Mac2:        "f0b0140c96bb",
				Rssi2:       -50,
				PDOP:        1.5,
				Satellites:  11,
			},
		},
		{
			payload:     "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcedeadbeef4242d6deadbeef4242d6deadbeef4242d6deadbeef4242d6",
			port:        151,
			autoPadding: false,
			expected: Port151Payload{
				BufferLevel: 0,
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     4.919,
				TTF:         time.Duration(24) * time.Second,
				PDOP:        1.5,
				Satellites:  11,
				Mac1:        "e0286d8a9478",
				Rssi1:       -53,
				Mac2:        "f0b0140c96bb",
				Rssi2:       -50,
				Mac3:        "deadbeef4242",
				Rssi3:       -42,
				Mac4:        "deadbeef4242",
				Rssi4:       -42,
				Mac5:        "deadbeef4242",
				Rssi5:       -42,
				Mac6:        "deadbeef4242",
				Rssi6:       -42,
			},
		},
		{
			payload:        "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcedeadbeef4242d6deadbeef4242d6deadbeef4242d6deadbeef4242d6deadbeef",
			port:           151,
			autoPadding:    false,
			skipValidation: true,
			expected: Port151Payload{
				BufferLevel: 0,
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     4.919,
				TTF:         time.Duration(24) * time.Second,
				PDOP:        1.5,
				Satellites:  11,
				Mac1:        "e0286d8a9478",
				Rssi1:       -53,
				Mac2:        "f0b0140c96bb",
				Rssi2:       -50,
				Mac3:        "deadbeef4242",
				Rssi3:       -42,
				Mac4:        "deadbeef4242",
				Rssi4:       -42,
				Mac5:        "deadbeef4242",
				Rssi5:       -42,
				Mac6:        "deadbeef4242",
				Rssi6:       -42,
			},
		},
		{
			payload:     "00000002d30b27008247b81312671bd164133718030b",
			port:        151,
			autoPadding: false,
			expected: Port151Payload{
				BufferLevel: 0,
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     4.919,
				TTF:         time.Duration(24) * time.Second,
				PDOP:        1.5,
				Satellites:  11,
			},
		},
		{
			payload:     "2d30b27008247b81312671bd164133718030b",
			port:        151,
			autoPadding: true,
			expected: Port151Payload{
				BufferLevel: 0,
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385383,
				Longitude:   8.53804,
				Altitude:    488.2,
				Timestamp:   time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:     4.919,
				TTF:         time.Duration(24) * time.Second,
				PDOP:        1.5,
				Satellites:  11,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder(WithAutoPadding(test.autoPadding), WithSkipValidation(test.skipValidation))
			got, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got == nil || got.Data != test.expected {
				t.Errorf("expected: %v\ngot: %v", test.expected, got)
			}
		})
	}

	t.Run("TestInvalidPayload", func(t *testing.T) {
		decoder := NewTagSLv1Decoder()
		_, err := decoder.Decode("", 1, "")
		if err == nil {
			t.Fatal("expected invalid payload")
		}
	})
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected error
	}{
		{
			payload:  "8002cdcd1300744f5e166018040b14341a",
			port:     1,
			expected: nil,
		},
		{
			payload:  "8005f5e10000744f5e166018040b14341a",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "8002cdcd130bebc200166018040b14341a",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "8002cdcd1300744f5e166018490b14341a",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Month", 73),
		},
		{
			payload:  "8002cdcd1300744f5e166018044914341a",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Day", 73),
		},
		{
			payload:  "8002cdcd1300744f5e166018040b49341a",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Hour", 73),
		},
		{
			payload:  "8002cdcd1300744f5e166018040b14491a",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Minute", 73),
		},
		{
			payload:  "8002cdcd1300744f5e166018040b143449",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Second", 73),
		},
		{
			payload:  "0002d30b070082491f11256718d9fe0ede190505",
			port:     10,
			expected: nil,
		},
		{
			payload:  "0005f5e1000082491f11256718d9fe0ede190505",
			port:     10,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "0002d30b070bebc20011256718d9fe0ede190505",
			port:     10,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "0002d30b070082491f11256718d9fe01f4190505",
			port:     10,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "0002d30b070082491f11256718d9fe157c190505",
			port:     10,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "0002d30b070082491f11256718d9fe0ede190502",
			port:     10,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 2),
		},
		{
			payload:  "0002d30b070082491f11256718d9fe0ede19051c",
			port:     10,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 28),
		},
		{
			payload:  "001044",
			port:     15,
			expected: nil,
		},
		{
			payload:  "0001f4",
			port:     15,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "00157c",
			port:     15,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f8118",
			port:     50,
			expected: nil,
		},
		{
			payload:  "0005f5e10000824c87117966c45dcd0f8118",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "0002d30c930bebc200117966c45dcd0f8118",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd01f418",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd157c18",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f81180205",
			port:     51,
			expected: nil,
		},
		{
			payload:  "0005f5e10000824c87117966c45dcd0f81180205",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "0002d30c930bebc200117966c45dcd0f81180205",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd01f4180205",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd157c180205",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f81180202",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 2),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f8118021c",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 28),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f81",
			port:     110,
			expected: nil,
		},
		{
			payload:  "00000005f5e10000824c87117966c45dcd0f81",
			port:     110,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "00000002d30c930bebc200117966c45dcd0f81",
			port:     110,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd01f4",
			port:     110,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd157c",
			port:     110,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f",
			port:     150,
			expected: nil,
		},
		{
			payload:  "00000005f5e10000824c87117966c45dcd0f812f",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "00000002d30c930bebc200117966c45dcd0f812f",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd01f42f",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd157c2f",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f0205",
			port:     151,
			expected: nil,
		},
		{
			payload:  "00000005f5e10000824c87117966c45dcd0f812f0205",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "00000002d30c930bebc200117966c45dcd0f812f0205",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd01f42f0205",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd157c2f0205",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f0202",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 2),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f021c",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 28),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vValidationWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder()
			got, err := decoder.Decode(test.payload, test.port, "")

			if err == nil && test.expected == nil {
				return
			}

			t.Logf("got %v", got)

			if err != nil && test.expected == nil || err == nil || err.Error() != test.expected.Error() {
				t.Errorf("expected: %v\ngot: %v", test.expected, err)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestInvalidHexString(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode("xx", 6, "")
	if err == nil || err.Error() != "encoding/hex: invalid byte: U+0078 'x'" {
		t.Fatal("expected invalid hex byte")
	}
}

func TestParseStatusByte(t *testing.T) {
	tests := []struct {
		input    byte
		expected Status
	}{
		{
			input: 0x80,
			expected: Status{
				DutyCycle:           true,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
		},
		{
			input: 0xFF,
			expected: Status{
				DutyCycle:           true,
				ConfigChangeId:      15,
				ConfigChangeSuccess: true,
				Moving:              true,
			},
		},
		{
			input: 0x00,
			expected: Status{
				DutyCycle:           false,
				ConfigChangeId:      0,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
		},
		{
			input: 0x4A,
			expected: Status{
				DutyCycle:           false,
				ConfigChangeId:      9,
				ConfigChangeSuccess: false,
				Moving:              false,
			},
		},
		{
			input: 0x8D,
			expected: Status{
				DutyCycle:           true,
				ConfigChangeId:      1,
				ConfigChangeSuccess: true,
				Moving:              true,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestStatusByteInput%v", test.input), func(t *testing.T) {
			got := parseStatusByte(test.input)
			if got != test.expected {
				t.Errorf("expected: %v\ngot: %v", test.expected, got)
			}
		})
	}
}

func TestFullDecode(t *testing.T) {
	tests := []struct {
		payload        string
		autoPadding    bool
		skipValidation bool
		expectedData   interface{}
		expectedStatus *Status
		port           int16
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			expectedData: Port1Payload{
				Moving:    false,
				DutyCycle: true,
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
				Moving:      true,
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
			payload: "0028672658500172a741b1e238b572a741b1e08bb03498b5c583e2b172a741b1e0cda772a741beed4cc472a741beef53b7",
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

	for _, test := range tests {
		decoder := NewTagSLv1Decoder(WithAutoPadding(test.autoPadding), WithSkipValidation(test.skipValidation))
		t.Run(fmt.Sprintf("TestFullDecodeWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoded, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if decoded.Data != test.expectedData {
				t.Errorf("expected data: %v, got: %v", test.expectedData, decoded.Data)
			}

			if decoded.Metadata != nil && test.expectedStatus == nil {
				t.Errorf("expected status to be nil, got: %v", decoded.Metadata)
			}
			if decoded.Metadata == nil && test.expectedStatus != nil {
				t.Errorf("expected status: %v, got: nil", test.expectedStatus)
			}

			// check if status is equal to expected status using reflect.DeepEqual
			if reflect.DeepEqual(decoded.Metadata, *test.expectedStatus) == false {
				t.Errorf("expected status: %v, got: %v", *test.expectedStatus, decoded.Metadata)
			}
		})
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode("deadbeef", 1, "")

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242", 1, "")

	if err == nil || err.Error() != "payload too long" {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload        string
		port           int16
		skipValidation bool
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			port:    1,
		},
		{
			payload:        "00",
			port:           2,
			skipValidation: true,
		},
		{
			payload: "822f0101f052fab920feafd0e4158b38b9afe05994cb2f5cb2a1b2c3d4e5f6aea1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:    3,
		},
		{
			payload: "0000012c00000e1000001c200078012c05dc02020100010200002328",
			port:    4,
		},
		{
			payload: "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8deadbeef4242d6deadbeef4242d6",
			port:    5,
		},
		{
			payload: "01",
			port:    6,
		},
		{
			payload: "66ec04bb00e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42ad",
			port:    7,
		},
		{
			payload: "012c141e9c455738304543434343460078012c01a8c0",
			port:    8,
		},
		{
			payload: "0002d308b50082457f16eb66c4a5cd0ed3",
			port:    10,
		},
		{
			payload: "800ee5",
			port:    15,
		},
		{
			payload: "0002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bfa1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:    50,
		},
		{
			payload: "0002D30BA000824ACE1122671B983E0EEA340B06726C9A74B58DB1FCF528F8634FB552A8DB7BD6B5B9E0286D8AABFCBC",
			port:    51,
		},
		{
			payload: "000166c4a5ba00e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0a1b2c3d4e5f6ae",
			port:    105,
		},
		{
			payload: "00020002d309ae008247c5113966c45d640f7e",
			port:    110,
		},
		{
			payload: "00020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bfa1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:    150,
		},
		{
			payload: "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcea1b2c3d4e5f6aea1b2c3d4e5f6aea1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:    151,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewTagSLv1Decoder(
				WithSkipValidation(test.skipValidation),
			)
			decodedPayload, _ := d.Decode(test.payload, test.port, "")

			// should be able to decode base feature
			base, ok := decodedPayload.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", decodedPayload)
			}
			// check if it panics
			base.GetTimestamp()

			if decodedPayload.Is(decoder.FeatureGNSS) {
				gnss, ok := decodedPayload.Data.(decoder.UplinkFeatureGNSS)
				if !ok {
					t.Fatalf("expected UplinkFeatureGNSS, got %T", decodedPayload)
				}
				if gnss.GetLatitude() == 0 {
					t.Fatalf("expected non zero latitude")
				}
				if gnss.GetLongitude() == 0 {
					t.Fatalf("expected non zero longitude")
				}
				if gnss.GetAltitude() == 0 {
					t.Fatalf("expected non zero altitude")
				}
				// call function to check if it panics
				gnss.GetAltitude()
				gnss.GetPDOP()
				gnss.GetSatellites()
				gnss.GetTTF()
				gnss.GetAccuracy()
			}
			if decodedPayload.Is(decoder.FeatureBuffered) {
				buffered, ok := decodedPayload.Data.(decoder.UplinkFeatureBuffered)
				if !ok {
					t.Fatalf("expected UplinkFeatureBuffered, got %T", decodedPayload)
				}
				// call function to check if it panics
				buffered.GetBufferLevel()
			}
			if decodedPayload.Is(decoder.FeatureBattery) {
				batteryVoltage, ok := decodedPayload.Data.(decoder.UplinkFeatureBattery)
				if !ok {
					t.Fatalf("expected UplinkFeatureBattery, got %T", decodedPayload)
				}
				if batteryVoltage.GetBatteryVoltage() == 0 {
					t.Fatalf("expected non zero battery voltage")
				}
			}
			if decodedPayload.Is(decoder.FeatureWiFi) {
				wifi, ok := decodedPayload.Data.(decoder.UplinkFeatureWiFi)
				if !ok {
					t.Fatalf("expected UplinkFeatureWiFi, got %T", decodedPayload)
				}
				if wifi.GetAccessPoints() == nil {
					t.Fatalf("expected non nil access points")
				}
			}
			if decodedPayload.Is(decoder.FeatureMoving) {
				moving, ok := decodedPayload.Data.(decoder.UplinkFeatureMoving)
				if !ok {
					t.Fatalf("expected UplinkFeatureMoving, got %T", decodedPayload)
				}
				// call function to check if it panics
				moving.IsMoving()
			}
			if decodedPayload.Is(decoder.FeatureDutyCycle) {
				dutyCycle, ok := decodedPayload.Data.(decoder.UplinkFeatureDutyCycle)
				if !ok {
					t.Fatalf("expected UplinkFeatureDutyCycle, got %T", decodedPayload)
				}
				// call function to check if it panics
				dutyCycle.IsDutyCycle()
			}
			if decodedPayload.Is(decoder.FeatureConfig) {
				config, ok := decodedPayload.Data.(decoder.UplinkFeatureConfig)
				if !ok {
					t.Fatalf("expected UplinkFeatureConfig, got %T", decodedPayload)
				}
				// call functions to check if it panics
				config.GetBle()
				config.GetGnss()
				config.GetWifi()
				config.GetMovingInterval()
				config.GetSteadyInterval()
				config.GetConfigInterval()
				config.GetGnssTimeout()
				config.GetAccelerometerThreshold()
				config.GetAccelerometerDelay()
				config.GetBatteryInterval()
				config.GetRejoinInterval()
				config.GetLowLightThreshold()
				config.GetHighLightThreshold()
				config.GetBatchSize()
				config.GetBufferSize()
			}
			if decodedPayload.Is(decoder.FeatureFirmwareVersion) {
				firmwareVersion, ok := decodedPayload.Data.(decoder.UplinkFeatureFirmwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureFirmwareVersion, got %T", decodedPayload)
				}
				if firmwareVersion.GetFirmwareVersion() == "" {
					t.Fatalf("expected non empty firmware version")
				}
			}
			if decodedPayload.Is(decoder.FeatureHardwareVersion) {
				hardwareVersion, ok := decodedPayload.Data.(decoder.UplinkFeatureHardwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureHardwareVersion, got %T", decodedPayload)
				}
				if hardwareVersion.GetHardwareVersion() == "" {
					t.Fatalf("expected non empty hardware version")
				}
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected []string
	}{
		{
			payload:  "8002cdcd1300744f5e166018040b14341a",
			port:     1,
			expected: []string{"\"latitude\": 47.041811", "\"longitude\": 7.622494", "\"altitude\": 572.8"},
		},
		{
			payload:  "822f0101f052fab920feafd0e4158b38b9afe05994cb2f5cb2a1b2c3d4e5f6aea1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:     3,
			expected: []string{"\"scanPointer\": 33327", "\"mac1\": \"f052fab920fe\"", "\"rssi1\": -81"},
		},
		{
			payload:  "0000012c00000e1000001c200078012c05dc02020100010200002328",
			port:     4,
			expected: []string{"\"batteryKeepAliveMessageInterval\": 9000"},
		},
		{
			payload:  "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8deadbeef4242d6deadbeef4242d6",
			port:     5,
			expected: []string{"\"moving\": false", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -88"},
		},
		{
			payload:  "01",
			port:     6,
			expected: []string{"\"buttonPressed\": true"},
		},
		{
			payload:  "66ec04bb00e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42ad",
			port:     7,
			expected: []string{"\"timestamp\": \"2024-09-19T11:02:19Z\"", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -69"},
		},
		{
			payload:  "012c141e9c455738304543434343460078012c01a8c0",
			port:     8,
			expected: []string{"\"minRssiValue\": -100", "\"advertisingFilter\": \"4048812220199682886\""},
		},
		{
			payload:  "0002d308b50082457f16eb66c4a5cd0ed3",
			port:     10,
			expected: []string{"\"timestamp\": \"2024-08-20T14:18:53Z\"", "\"battery\": 3.795", "\"ttf\": \"0s\""},
		},
		{
			payload:  "800ee5",
			port:     15,
			expected: []string{"\"lowBattery\": false", "\"battery\": 3.813"},
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bfa1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:     50,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-08-20T09:11:41Z\"", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -87"},
		},
		{
			payload:  "0002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1fcf528f8634fb552a8db7bd6b5b9e0286d8aabfcbc",
			port:     51,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-10-25T13:08:14Z\"", "\"pdop\": 5.5", "\"mac1\": \"726c9a74b58d\"", "\"rssi1\": -79"},
		},
		{
			payload:  "000166c4a5ba00e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0a1b2c3d4e5f6ae",
			port:     105,
			expected: []string{"\"moving\": false", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -79"},
		},
		{
			payload:  "00020002d309ae008247c5113966c45d640f7e",
			port:     110,
			expected: []string{"\"timestamp\": \"2024-08-20T09:09:56Z\"", "\"battery\": 3.966"},
		},
		{
			payload:  "00020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bfa1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:     150,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-08-20T09:11:41Z\"", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -87", "\"ttf\": \"24s\""},
		},
		{
			payload:  "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcea1b2c3d4e5f6aea1b2c3d4e5f6aea1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:     151,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-10-25T17:12:04Z\"", "\"pdop\": 1.5", "\"mac1\": \"e0286d8a9478\"", "\"rssi1\": -53", "\"ttf\": \"24s\""},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder()

			data, _ := decoder.Decode(test.payload, test.port, "")

			marshaled, err := json.MarshalIndent(map[string]interface{}{
				"data":     data.Data,
				"metadata": data.Metadata,
			}, "", "   ")

			if err != nil {
				t.Fatalf("marshalling json failed because %s", err)
			}

			t.Logf("%s\n", marshaled)

			for _, value := range test.expected {
				fmt.Printf("value:%s\n", value)
				if !strings.Contains(string(marshaled), value) {
					t.Fatalf("expected to find %s\n", value)
				}
			}
		})
	}
}
