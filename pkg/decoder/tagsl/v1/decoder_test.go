package tagsl

import (
	"context"
	"encoding/json"
	"errors"
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
		port           uint8
		skipValidation bool
		expected       any
	}{
		{
			port:    1,
			payload: "8002cdcd1300744f5e166018040b14341a",
			expected: Port1Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.041811,
				Longitude:    7.622494,
				Altitude:     572.8,
				Year:         24,
				Month:        4,
				Day:          11,
				Hour:         20,
				Minute:       52,
				Second:       26,
			},
		},
		{
			port:    1,
			payload: "8002cdcd1300744f5e166018040b14341a",
			expected: Port1Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.041811,
				Longitude:    7.622494,
				Altitude:     572.8,
				Year:         24,
				Month:        4,
				Day:          11,
				Hour:         20,
				Minute:       52,
				Second:       26,
			},
		},
		{
			port:           1,
			payload:        "8002cdcd1300744f5e166018040b14341adeadbeef",
			skipValidation: true,
			expected: Port1Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.041811,
				Longitude:    7.622494,
				Altitude:     572.8,
				Year:         24,
				Month:        4,
				Day:          11,
				Hour:         20,
				Minute:       52,
				Second:       26,
			},
		},
		{
			port:           1,
			payload:        "adfe019c3cfbc9ac0c14500a0c0e160a3b",
			skipValidation: false,
			expected: Port1Payload{
				DutyCycle:    true,
				ConfigId:     5,
				ConfigChange: true,
				Moving:       true,
				Latitude:     -33.4489,
				Longitude:    -70.6693,
				Altitude:     520,
				Year:         10,
				Month:        12,
				Day:          14,
				Hour:         22,
				Minute:       10,
				Second:       59,
			},
		},
		{
			port:           2,
			payload:        "00",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
			},
		},
		{
			port:           2,
			payload:        "01",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
			},
		},
		{
			port:           2,
			payload:        "80",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
			},
		},
		{
			port:           2,
			payload:        "81",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
			},
		},
		{
			port:           2,
			payload:        "38",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    false,
				ConfigId:     7,
				ConfigChange: false,
				Moving:       false,
			},
		},
		{
			port:           2,
			payload:        "39",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    false,
				ConfigId:     7,
				ConfigChange: false,
				Moving:       true,
			},
		},
		{
			port:           2,
			payload:        "b9",
			skipValidation: true,
			expected: Port2Payload{
				DutyCycle:    true,
				ConfigId:     7,
				ConfigChange: false,
				Moving:       true,
			},
		},
		{
			port:    3,
			payload: "01eb0101f052fab920feac",
			expected: Port3Payload{
				ScanPointer:    491,
				TotalMessages:  1,
				CurrentMessage: 1,
				Mac1:           "f052fab920fe",
				Rssi1:          -84,
			},
		},
		{
			port:    3,
			payload: "01eb0101f052fab920feae",
			expected: Port3Payload{
				ScanPointer:    491,
				TotalMessages:  1,
				CurrentMessage: 1,
				Mac1:           "f052fab920fe",
				Rssi1:          -82,
			},
		},
		{
			port:    3,
			payload: "822f0101f052fab920feafd0e4158b38b9afe05994cb2f5cb2",
			expected: Port3Payload{
				ScanPointer:    33327,
				TotalMessages:  1,
				CurrentMessage: 1,
				Mac1:           "f052fab920fe",
				Rssi1:          -81,
				Mac2:           helpers.StringPtr("d0e4158b38b9"),
				Rssi2:          helpers.Int8Ptr(-81),
				Mac3:           helpers.StringPtr("e05994cb2f5c"),
				Rssi3:          helpers.Int8Ptr(-78),
			},
		},
		{
			port:    3,
			payload: "01eb0101f052fab920feadd0e4158b38b9afe05994cb2f5cad",
			expected: Port3Payload{
				ScanPointer:    491,
				TotalMessages:  1,
				CurrentMessage: 1,
				Mac1:           "f052fab920fe",
				Rssi1:          -83,
				Mac2:           helpers.StringPtr("d0e4158b38b9"),
				Rssi2:          helpers.Int8Ptr(-81),
				Mac3:           helpers.StringPtr("e05994cb2f5c"),
				Rssi3:          helpers.Int8Ptr(-83),
			},
		},
		{
			port:    4,
			payload: "0000012c00000e1000001c200078012c05dc02020100010200002328",
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
			port:    4,
			payload: "0000003c0000012c000151800078012c05dc02020100010200005460",
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
			port:           4,
			payload:        "0000003c0000012c000151800078012c05dc02020100010200005460000a1000",
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
				BatchSize:                       helpers.Uint16Ptr(10),
				BufferSize:                      helpers.Uint16Ptr(4096),
			},
		},
		{
			port:    4,
			payload: "0000003c0000012c000151800078012c05dc02020100010200005460",
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
			port:    5,
			payload: "808c59c3c99fc0ad",
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: true,
				Mac1:      "8c59c3c99fc0",
				Rssi1:     -83,
			},
		},
		{
			port:    5,
			payload: "80e0286d8a2742a1",
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: true,
				Mac1:      "e0286d8a2742",
				Rssi1:     -95,
			},
		},
		{
			port:    5,
			payload: "001f3fd57cecb4c9b0140c96bbb2bd286d8a9478b8ad",
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "1f3fd57cecb4",
				Rssi1:     -55,
				Mac2:      helpers.StringPtr("b0140c96bbb2"),
				Rssi2:     helpers.Int8Ptr(-67),
				Mac3:      helpers.StringPtr("286d8a9478b8"),
				Rssi3:     helpers.Int8Ptr(-83),
			},
		},
		{
			port:    5,
			payload: "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8",
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -88,
				Mac2:      helpers.StringPtr("e0286d8a9478"),
				Rssi2:     helpers.Int8Ptr(-62),
				Mac3:      helpers.StringPtr("726c9a74b58d"),
				Rssi3:     helpers.Int8Ptr(-85),
				Mac4:      helpers.StringPtr("726cdac8b89d"),
				Rssi4:     helpers.Int8Ptr(-84),
				Mac5:      helpers.StringPtr("f0b0140c96bb"),
				Rssi5:     helpers.Int8Ptr(-56),
			},
		},
		{
			port:           5,
			payload:        "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8deadbeef4242d6deadbeef4242d6",
			skipValidation: false,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -88,
				Mac2:      helpers.StringPtr("e0286d8a9478"),
				Rssi2:     helpers.Int8Ptr(-62),
				Mac3:      helpers.StringPtr("726c9a74b58d"),
				Rssi3:     helpers.Int8Ptr(-85),
				Mac4:      helpers.StringPtr("726cdac8b89d"),
				Rssi4:     helpers.Int8Ptr(-84),
				Mac5:      helpers.StringPtr("f0b0140c96bb"),
				Rssi5:     helpers.Int8Ptr(-56),
				Mac6:      helpers.StringPtr("deadbeef4242"),
				Rssi6:     helpers.Int8Ptr(-42),
				Mac7:      helpers.StringPtr("deadbeef4242"),
				Rssi7:     helpers.Int8Ptr(-42),
			},
		},
		{
			port:           5,
			payload:        "00e0286d8aabfca8e0286d8a9478c2726c9a74b58dab726cdac8b89dacf0b0140c96bbc8deadbeef4242d6deadbeef4242d6deadbeefdeadbeef",
			skipValidation: true,
			expected: Port5Payload{
				Moving:    false,
				DutyCycle: false,
				Mac1:      "e0286d8aabfc",
				Rssi1:     -88,
				Mac2:      helpers.StringPtr("e0286d8a9478"),
				Rssi2:     helpers.Int8Ptr(-62),
				Mac3:      helpers.StringPtr("726c9a74b58d"),
				Rssi3:     helpers.Int8Ptr(-85),
				Mac4:      helpers.StringPtr("726cdac8b89d"),
				Rssi4:     helpers.Int8Ptr(-84),
				Mac5:      helpers.StringPtr("f0b0140c96bb"),
				Rssi5:     helpers.Int8Ptr(-56),
				Mac6:      helpers.StringPtr("deadbeef4242"),
				Rssi6:     helpers.Int8Ptr(-42),
				Mac7:      helpers.StringPtr("deadbeef4242"),
				Rssi7:     helpers.Int8Ptr(-42),
			},
		},
		{
			port:    6,
			payload: "01",
			expected: Port6Payload{
				ButtonPressed: true,
			},
		},
		{
			port:    6,
			payload: "00",
			expected: Port6Payload{
				ButtonPressed: false,
			},
		},
		{
			port:           6,
			payload:        "00deadbeef",
			skipValidation: true,
			expected: Port6Payload{
				ButtonPressed: false,
			},
		},
		{
			port:    7,
			payload: "66ec04bb80e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42ad",
			expected: Port7Payload{
				Timestamp:    time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Mac1:         "e0286d8aabfc",
				Rssi1:        -69,
				Mac2:         helpers.StringPtr("ec6c9a74b58f"),
				Rssi2:        helpers.Int8Ptr(-78),
				Mac3:         helpers.StringPtr("726c9a74b58d"),
				Rssi3:        helpers.Int8Ptr(-79),
				Mac4:         helpers.StringPtr("e0286d8a9478"),
				Rssi4:        helpers.Int8Ptr(-53),
				Mac5:         helpers.StringPtr("f0b0140c96bb"),
				Rssi5:        helpers.Int8Ptr(-46),
				Mac6:         helpers.StringPtr("260122180d42"),
				Rssi6:        helpers.Int8Ptr(-83),
			},
		},
		{
			port:           7,
			payload:        "66ec04bb01e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42addeadbeef",
			skipValidation: true,
			expected: Port7Payload{
				Timestamp:    time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Mac1:         "e0286d8aabfc",
				Rssi1:        -69,
				Mac2:         helpers.StringPtr("ec6c9a74b58f"),
				Rssi2:        helpers.Int8Ptr(-78),
				Mac3:         helpers.StringPtr("726c9a74b58d"),
				Rssi3:        helpers.Int8Ptr(-79),
				Mac4:         helpers.StringPtr("e0286d8a9478"),
				Rssi4:        helpers.Int8Ptr(-53),
				Mac5:         helpers.StringPtr("f0b0140c96bb"),
				Rssi5:        helpers.Int8Ptr(-46),
				Mac6:         helpers.StringPtr("260122180d42"),
				Rssi6:        helpers.Int8Ptr(-83),
			},
		},
		{
			port:    7,
			payload: "66ec04bb81e0286d8aabfcbb",
			expected: Port7Payload{
				Timestamp:    time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Mac1:         "e0286d8aabfc",
				Rssi1:        -69,
			},
		},
		{
			port:    7,
			payload: "66ec04bb3ce0286d8aabfcbb",
			expected: Port7Payload{
				Timestamp:    time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				DutyCycle:    false,
				ConfigId:     7,
				ConfigChange: true,
				Moving:       false,
				Mac1:         "e0286d8aabfc",
				Rssi1:        -69,
			},
		},
		{
			port:    8,
			payload: "012c141e9c455738304543434343460078012c01a8c0",
			expected: Port8Payload{
				ScanInterval:                          300,
				ScanTime:                              20,
				MaxBeacons:                            30,
				MinRssiValue:                          -100,
				AdvertisingFilter:                     "EW80ECCCCF",
				AccelerometerTriggerHoldTimer:         120,
				AccelerometerThreshold:                300,
				BLECurrentConfigurationUplinkInterval: 43200,
				ScanMode:                              1,
			},
		},
		{
			port:           8,
			payload:        "012c141e9c455738304543434343460078012c01a8c0deadbeef",
			skipValidation: true,
			expected: Port8Payload{
				ScanInterval:                          300,
				ScanTime:                              20,
				MaxBeacons:                            30,
				MinRssiValue:                          -100,
				AdvertisingFilter:                     "EW80ECCCCF",
				AccelerometerTriggerHoldTimer:         120,
				AccelerometerThreshold:                300,
				BLECurrentConfigurationUplinkInterval: 43200,
				ScanMode:                              1,
			},
		},
		{
			port:    8,
			payload: "012c141e9c455738304543434343460078012c01a8c0",
			expected: Port8Payload{
				ScanInterval:                          300,
				ScanTime:                              20,
				MaxBeacons:                            30,
				MinRssiValue:                          -100,
				AdvertisingFilter:                     "EW80ECCCCF",
				AccelerometerTriggerHoldTimer:         120,
				AccelerometerThreshold:                300,
				BLECurrentConfigurationUplinkInterval: 43200,
				ScanMode:                              1,
			},
		},
		{
			port:    10,
			payload: "0002d308b50082457f16eb66c4a5cd0ed3",
			expected: Port10Payload{
				Latitude:  47.384757,
				Longitude: 8.537471,
				Altitude:  586.7,
				Timestamp: time.Date(2024, 8, 20, 14, 18, 53, 0, time.UTC),
				Battery:   3.795,
			},
		},
		{
			port:    10,
			payload: "0002d30b070082491f11256718d9fe0ede190505",
			expected: Port10Payload{
				Latitude:   47.385351,
				Longitude:  8.538399,
				Altitude:   438.9,
				Timestamp:  time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:    3.806,
				PDOP:       helpers.Float64Ptr(2.5),
				Satellites: helpers.Uint8Ptr(5),
				TTF:        helpers.DurationPtr(time.Duration(25) * time.Second),
			},
		},
		{
			port:           10,
			payload:        "0002d30b070082491f11256718d9fe0ede190505deadbeef",
			skipValidation: true,
			expected: Port10Payload{
				Latitude:   47.385351,
				Longitude:  8.538399,
				Altitude:   438.9,
				Timestamp:  time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:    3.806,
				PDOP:       helpers.Float64Ptr(2.5),
				Satellites: helpers.Uint8Ptr(5),
				TTF:        helpers.DurationPtr(time.Duration(25) * time.Second),
			},
		},
		{
			port:    10,
			payload: "0002d30b070082491f11256718d9fe0ede",
			expected: Port10Payload{
				Latitude:  47.385351,
				Longitude: 8.538399,
				Altitude:  438.9,
				Timestamp: time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:   3.806,
			},
		},
		{
			port:    10,
			payload: "0002d30b070082491f11256718d9fe0ede",
			expected: Port10Payload{
				Latitude:  47.385351,
				Longitude: 8.538399,
				Altitude:  438.9,
				Timestamp: time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
				Battery:   3.806,
			},
		},
		{
			payload: "81ff04631cfbf094fc8e946ef132401059f00906",
			port:    10,
			expected: Port10Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     -16.4897,
				Longitude:    -68.1193,
				Altitude:     3650,
				Timestamp:    time.Date(2028, 12, 24, 20, 0, 0, 0, time.UTC),
				Battery:      4.185,
				TTF:          helpers.DurationPtr(time.Duration(240) * time.Second),
				PDOP:         helpers.Float64Ptr(4.5),
				Satellites:   helpers.Uint8Ptr(6),
			},
		},
		{
			port:    15,
			payload: "000ee5",
			expected: Port15Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				LowBattery:   false,
				Battery:      3.813,
			},
		},
		{
			port:           15,
			payload:        "800ee5deadbeef",
			skipValidation: true,
			expected: Port15Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				LowBattery:   false,
				Battery:      3.813,
			},
		},
		{
			port:    15,
			payload: "011044",
			expected: Port15Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				LowBattery:   true,
				Battery:      4.164,
			},
		},
		{
			port:    15,
			payload: "811044",
			expected: Port15Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				LowBattery:   true,
				Battery:      4.164,
			},
		},
		{
			port:    15,
			payload: "241044",
			expected: Port15Payload{
				DutyCycle:    false,
				ConfigId:     4,
				ConfigChange: true,
				LowBattery:   false,
				Battery:      4.164,
			},
		},
		{
			port:    50,
			payload: "0002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
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
				Mac2:      helpers.StringPtr("f0b0140c96bb"),
				Rssi2:     helpers.Int8Ptr(-56),
				Mac3:      helpers.StringPtr("726c9a74b58d"),
				Rssi3:     helpers.Int8Ptr(-88),
				Mac4:      helpers.StringPtr("e0286d8a9478"),
				Rssi4:     helpers.Int8Ptr(-65),
			},
		},
		{
			port:    50,
			payload: "0102d30b2a0082499c10ee66c496900ed34af0b0140c96bbb3e0286d8a9478c3fc848e9b5571c2",
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
				Mac2:      helpers.StringPtr("e0286d8a9478"),
				Rssi2:     helpers.Int8Ptr(-61),
				Mac3:      helpers.StringPtr("fc848e9b5571"),
				Rssi3:     helpers.Int8Ptr(-62),
			},
		},
		{
			port:    50,
			payload: "0102d30b2a0082499c10ee66c496900ed34af0b0140c96bbb3e0286d8a9478c3fc848e9b5571c2deadbeef4242d6",
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
				Mac2:      helpers.StringPtr("e0286d8a9478"),
				Rssi2:     helpers.Int8Ptr(-61),
				Mac3:      helpers.StringPtr("fc848e9b5571"),
				Rssi3:     helpers.Int8Ptr(-62),
				Mac4:      helpers.StringPtr("deadbeef4242"),
				Rssi4:     helpers.Int8Ptr(-42),
			},
		}, {
			port:           50,
			payload:        "0102d30b2a0082499c10ee66c496900ed34af0b0140c96bbb3e0286d8a9478c3fc848e9b5571c2deadbeef4242d6",
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
				Mac2:      helpers.StringPtr("e0286d8a9478"),
				Rssi2:     helpers.Int8Ptr(-61),
				Mac3:      helpers.StringPtr("fc848e9b5571"),
				Rssi3:     helpers.Int8Ptr(-62),
				Mac4:      helpers.StringPtr("deadbeef4242"),
				Rssi4:     helpers.Int8Ptr(-42),
			},
		},
		{
			port:    50,
			payload: "0102d30b2a0082499c10ee66c496900ed34aa1b2c3d4e5f6b8",
			expected: Port50Payload{
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
				Mac1:      "a1b2c3d4e5f6",
				Rssi1:     -72,
			},
		},
		{
			port:    50,
			payload: "0102d30b2a0082499c10ee66c496900ed34aa1b2c3d4e5f6b8",
			expected: Port50Payload{
				Moving:    true,
				DutyCycle: false,
				Latitude:  47.385386,
				Longitude: 8.538524,
				Altitude:  433.4,
				Timestamp: time.Date(2024, 8, 20, 13, 13, 52, 0, time.UTC),
				Battery:   3.795,
				TTF:       time.Duration(74) * time.Second,
				Mac1:      "a1b2c3d4e5f6",
				Rssi1:     -72,
			},
		},
		{
			payload: "81fe624b04f97b74bc07d0602b67310f2131a1b2c3d4e5f6c0",
			port:    50,
			expected: Port50Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     -27.1127,
				Longitude:    -109.3497,
				Altitude:     200,
				Timestamp:    time.Date(2021, 2, 16, 6, 33, 21, 0, time.UTC),
				Battery:      3.873,
				TTF:          time.Duration(49) * time.Second,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -64,
			},
		},
		{
			port:    51,
			payload: "0002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1fcf528f8634fb552a8db7bd6b5b9e0286d8aabfcbc",
			expected: Port51Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385504,
				Longitude:    8.53883,
				Altitude:     438.6,
				Timestamp:    time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:      3.818,
				TTF:          time.Duration(52) * time.Second,
				PDOP:         5.5,
				Satellites:   6,
				Mac1:         "726c9a74b58d",
				Rssi1:        -79,
				Mac2:         helpers.StringPtr("fcf528f8634f"),
				Rssi2:        helpers.Int8Ptr(-75),
				Mac3:         helpers.StringPtr("52a8db7bd6b5"),
				Rssi3:        helpers.Int8Ptr(-71),
				Mac4:         helpers.StringPtr("e0286d8aabfc"),
				Rssi4:        helpers.Int8Ptr(-68),
			},
		},
		{
			port:    51,
			payload: "0002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1",
			expected: Port51Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385504,
				Longitude:    8.53883,
				Altitude:     438.6,
				Timestamp:    time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:      3.818,
				TTF:          time.Duration(52) * time.Second,
				PDOP:         5.5,
				Satellites:   6,
				Mac1:         "726c9a74b58d",
				Rssi1:        -79,
			},
		},
		{
			port:    51,
			payload: "8002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1",
			expected: Port51Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385504,
				Longitude:    8.53883,
				Altitude:     438.6,
				Timestamp:    time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:      3.818,
				TTF:          time.Duration(52) * time.Second,
				PDOP:         5.5,
				Satellites:   6,
				Mac1:         "726c9a74b58d",
				Rssi1:        -79,
			},
		},
		{
			port:    51,
			payload: "0102d30ba000824ace1122671b983e0eea340b06726c9a74b58db1",
			expected: Port51Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.385504,
				Longitude:    8.53883,
				Altitude:     438.6,
				Timestamp:    time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:      3.818,
				TTF:          time.Duration(52) * time.Second,
				PDOP:         5.5,
				Satellites:   6,
				Mac1:         "726c9a74b58d",
				Rssi1:        -79,
			},
		},
		{
			port:    51,
			payload: "0002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1",
			expected: Port51Payload{
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385504,
				Longitude:    8.53883,
				Altitude:     438.6,
				Timestamp:    time.Date(2024, 10, 25, 13, 8, 14, 0, time.UTC),
				Battery:      3.818,
				TTF:          time.Duration(52) * time.Second,
				PDOP:         5.5,
				Satellites:   6,
				Mac1:         "726c9a74b58d",
				Rssi1:        -79,
			},
		},
		{
			payload: "21fcbbca14fbedc7cc00cc646ea2e40f8ca01004a1b2c3d4e5f6b8",
			port:    51,
			expected: Port51Payload{
				DutyCycle:    false,
				ConfigId:     4,
				ConfigChange: false,
				Moving:       true,
				Latitude:     -54.8019,
				Longitude:    -68.3029,
				Altitude:     20.4,
				Timestamp:    time.Date(2023, 5, 24, 23, 51, 0, 0, time.UTC),
				Battery:      3.980,
				TTF:          time.Duration(160) * time.Second,
				PDOP:         8,
				Satellites:   4,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -72,
			},
		},
		{
			port:    105,
			payload: "000166c4a5ba00e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0",
			expected: Port105Payload{
				BufferLevel: 1,
				Timestamp:   time.Date(2024, 8, 20, 14, 18, 34, 0, time.UTC),
				DutyCycle:   false,
				Moving:      false,
				Mac1:        "e0286d8aabfc",
				Rssi1:       -79,
				Mac2:        helpers.StringPtr("e0286d8a9478"),
				Rssi2:       helpers.Int8Ptr(-62),
				Mac3:        helpers.StringPtr("ec6c9a74b58f"),
				Rssi3:       helpers.Int8Ptr(-83),
				Mac4:        helpers.StringPtr("726c9a74b58d"),
				Rssi4:       helpers.Int8Ptr(-83),
				Mac5:        helpers.StringPtr("f0b0140c96bb"),
				Rssi5:       helpers.Int8Ptr(-48),
			},
		},
		{
			port:    105,
			payload: "010166c4a5ba80e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0",
			expected: Port105Payload{
				BufferLevel: 257,
				Timestamp:   time.Date(2024, 8, 20, 14, 18, 34, 0, time.UTC),
				DutyCycle:   true,
				Moving:      false,
				Mac1:        "e0286d8aabfc",
				Rssi1:       -79,
				Mac2:        helpers.StringPtr("e0286d8a9478"),
				Rssi2:       helpers.Int8Ptr(-62),
				Mac3:        helpers.StringPtr("ec6c9a74b58f"),
				Rssi3:       helpers.Int8Ptr(-83),
				Mac4:        helpers.StringPtr("726c9a74b58d"),
				Rssi4:       helpers.Int8Ptr(-83),
				Mac5:        helpers.StringPtr("f0b0140c96bb"),
				Rssi5:       helpers.Int8Ptr(-48),
			},
		},
		{
			port:    105,
			payload: "001366ee2f4d01c4eb438ddde2a504e31aea1b01a7245a4c7a0d2ec026e98d560d2ebbccd42ef92ed4ae704f5708e1d1b9",
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
				DutyCycle:   false,
				Moving:      true,
				Mac1:        "c4eb438ddde2",
				Rssi1:       -91,
				Mac2:        helpers.StringPtr("04e31aea1b01"),
				Rssi2:       helpers.Int8Ptr(-89),
				Mac3:        helpers.StringPtr("245a4c7a0d2e"),
				Rssi3:       helpers.Int8Ptr(-64),
				Mac4:        helpers.StringPtr("26e98d560d2e"),
				Rssi4:       helpers.Int8Ptr(-69),
				Mac5:        helpers.StringPtr("ccd42ef92ed4"),
				Rssi5:       helpers.Int8Ptr(-82),
				Mac6:        helpers.StringPtr("704f5708e1d1"),
				Rssi6:       helpers.Int8Ptr(-71),
			},
		},
		{
			port:           105,
			payload:        "001366ee2f4d81c4eb438ddde2a504e31aea1b01a7245a4c7a0d2ec026e98d560d2ebbccd42ef92ed4ae704f5708e1d1b9deadbeef",
			skipValidation: true,
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
				DutyCycle:   true,
				Moving:      true,
				Mac1:        "c4eb438ddde2",
				Rssi1:       -91,
				Mac2:        helpers.StringPtr("04e31aea1b01"),
				Rssi2:       helpers.Int8Ptr(-89),
				Mac3:        helpers.StringPtr("245a4c7a0d2e"),
				Rssi3:       helpers.Int8Ptr(-64),
				Mac4:        helpers.StringPtr("26e98d560d2e"),
				Rssi4:       helpers.Int8Ptr(-69),
				Mac5:        helpers.StringPtr("ccd42ef92ed4"),
				Rssi5:       helpers.Int8Ptr(-82),
				Mac6:        helpers.StringPtr("704f5708e1d1"),
				Rssi6:       helpers.Int8Ptr(-71),
			},
		},
		{
			port:    105,
			payload: "001366ee2f4d80c4eb438ddde2a5",
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
				DutyCycle:   true,
				Moving:      false,
				Mac1:        "c4eb438ddde2",
				Rssi1:       -91,
			},
		},
		{
			port:    105,
			payload: "001366ee2f4d80c4eb438ddde2a5",
			expected: Port105Payload{
				BufferLevel: 19,
				Timestamp:   time.Date(2024, 9, 21, 2, 28, 29, 0, time.UTC),
				DutyCycle:   true,
				Moving:      false,
				Mac1:        "c4eb438ddde2",
				Rssi1:       -91,
			},
		},
		{
			port:    110,
			payload: "00020002d309ae008247c5113966c45d640f7e",
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
			port:    110,
			payload: "01020002d309ae008247c5113966c45d640f7e",
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
			port:           110,
			payload:        "01020002d309ae008247c5113966c45d640f7ee72004deadbeef",
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
				TTF:         helpers.DurationPtr(time.Duration(231) * time.Second),
				PDOP:        helpers.Float64Ptr(16),
				Satellites:  helpers.Uint8Ptr(4),
			},
		},
		{
			port:    110,
			payload: "00040002d30b8c00824a35112266c45c440f832a0807",
			expected: Port110Payload{
				BufferLevel: 4,
				Moving:      false,
				DutyCycle:   false,
				Latitude:    47.385484,
				Longitude:   8.538677,
				Altitude:    438.6,
				Timestamp:   time.Date(2024, 8, 20, 9, 5, 8, 0, time.UTC),
				Battery:     3.971,
				TTF:         helpers.DurationPtr(time.Duration(42) * time.Second),
				PDOP:        helpers.Float64Ptr(4),
				Satellites:  helpers.Uint8Ptr(7),
			},
		},
		{
			port:    110,
			payload: "00040002d30b8c00824a35112266c45c440f83",
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
			port:    110,
			payload: "023881ff04631cfbf094fc8e946ef132401059f00906",
			expected: Port110Payload{
				BufferLevel:  568,
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     -16.4897,
				Longitude:    -68.1193,
				Altitude:     3650,
				Timestamp:    time.Date(2028, 12, 24, 20, 0, 0, 0, time.UTC),
				Battery:      4.185,
				TTF:          helpers.DurationPtr(time.Duration(240) * time.Second),
				PDOP:         helpers.Float64Ptr(4.5),
				Satellites:   helpers.Uint8Ptr(6),
			},
		},
		{
			port:    150,
			payload: "00020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			expected: Port150Payload{
				BufferLevel:  2,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385747,
				Longitude:    8.539271,
				Altitude:     447.3,
				Timestamp:    time.Date(2024, 8, 20, 9, 11, 41, 0, time.UTC),
				Battery:      3.969,
				TTF:          time.Duration(24) * time.Second,
				Mac1:         "e0286d8aabfc",
				Rssi1:        -87,
				Mac2:         helpers.StringPtr("f0b0140c96bb"),
				Rssi2:        helpers.Int8Ptr(-56),
				Mac3:         helpers.StringPtr("726c9a74b58d"),
				Rssi3:        helpers.Int8Ptr(-88),
				Mac4:         helpers.StringPtr("e0286d8a9478"),
				Rssi4:        helpers.Int8Ptr(-65),
			},
		},
		{
			port:    150,
			payload: "01020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			expected: Port150Payload{
				BufferLevel:  258,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385747,
				Longitude:    8.539271,
				Altitude:     447.3,
				Timestamp:    time.Date(2024, 8, 20, 9, 11, 41, 0, time.UTC),
				Battery:      3.969,
				TTF:          time.Duration(24) * time.Second,
				Mac1:         "e0286d8aabfc",
				Rssi1:        -87,
				Mac2:         helpers.StringPtr("f0b0140c96bb"),
				Rssi2:        helpers.Int8Ptr(-56),
				Mac3:         helpers.StringPtr("726c9a74b58d"),
				Rssi3:        helpers.Int8Ptr(-88),
				Mac4:         helpers.StringPtr("e0286d8a9478"),
				Rssi4:        helpers.Int8Ptr(-65),
			},
		},
		{
			port:    150,
			payload: "00000102d30a98008248b611ac66c45ed80f6b1be0286d0a6f42a3000000000044a4e0286d8a9478bff0b0140c96bbc9",
			expected: Port150Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.38524,
				Longitude:    8.538294,
				Altitude:     452.4,
				Timestamp:    time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:      3.947,
				TTF:          time.Duration(27) * time.Second,
				Mac1:         "e0286d0a6f42",
				Rssi1:        -93,
				Mac2:         helpers.StringPtr("000000000044"),
				Rssi2:        helpers.Int8Ptr(-92),
				Mac3:         helpers.StringPtr("e0286d8a9478"),
				Rssi3:        helpers.Int8Ptr(-65),
				Mac4:         helpers.StringPtr("f0b0140c96bb"),
				Rssi4:        helpers.Int8Ptr(-55),
			},
		},
		{
			port:    150,
			payload: "00000102d30a98008248b611ac66c45ed80f6b1be0286d0a6f42a3000000000044a4e0286d8a9478bff0b0140c96bbc9",
			expected: Port150Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.38524,
				Longitude:    8.538294,
				Altitude:     452.4,
				Timestamp:    time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:      3.947,
				TTF:          time.Duration(27) * time.Second,
				Mac1:         "e0286d0a6f42",
				Rssi1:        -93,
				Mac2:         helpers.StringPtr("000000000044"),
				Rssi2:        helpers.Int8Ptr(-92),
				Mac3:         helpers.StringPtr("e0286d8a9478"),
				Rssi3:        helpers.Int8Ptr(-65),
				Mac4:         helpers.StringPtr("f0b0140c96bb"),
				Rssi4:        helpers.Int8Ptr(-55),
			},
		},
		{
			port:           150,
			payload:        "00000102d30a98008248b611ac66c45ed80f6b1be0286d0a6f42a3000000000044a4e0286d8a9478bff0b0140c96bbc9deadbeef",
			skipValidation: true,
			expected: Port150Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.38524,
				Longitude:    8.538294,
				Altitude:     452.4,
				Timestamp:    time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:      3.947,
				TTF:          time.Duration(27) * time.Second,
				Mac1:         "e0286d0a6f42",
				Rssi1:        -93,
				Mac2:         helpers.StringPtr("000000000044"),
				Rssi2:        helpers.Int8Ptr(-92),
				Mac3:         helpers.StringPtr("e0286d8a9478"),
				Rssi3:        helpers.Int8Ptr(-65),
				Mac4:         helpers.StringPtr("f0b0140c96bb"),
				Rssi4:        helpers.Int8Ptr(-55),
			},
		},
		{
			port:    150,
			payload: "00000102d30a98008248b611ac66c45ed80f6b1ba1b2c3d4e5f6c2",
			expected: Port150Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.38524,
				Longitude:    8.538294,
				Altitude:     452.4,
				Timestamp:    time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:      3.947,
				TTF:          time.Duration(27) * time.Second,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -62,
			},
		},
		{
			port:    150,
			payload: "00000102d30a98008248b611ac66c45ed80f6b1ba1b2c3d4e5f6c2",
			expected: Port150Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.38524,
				Longitude:    8.538294,
				Altitude:     452.4,
				Timestamp:    time.Date(2024, 8, 20, 9, 16, 8, 0, time.UTC),
				Battery:      3.947,
				TTF:          time.Duration(27) * time.Second,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -62,
			},
		},
		{
			port:    150,
			payload: "155e81fe624b04f97b74bc07d0602b67310f2131a1b2c3d4e5f6b0",
			expected: Port150Payload{
				BufferLevel:  5470,
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     -27.1127,
				Longitude:    -109.3497,
				Altitude:     200,
				Timestamp:    time.Date(2021, 2, 16, 6, 33, 21, 0, time.UTC),
				Battery:      3.873,
				TTF:          time.Duration(49) * time.Second,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -80,
			},
		},
		{
			port:    151,
			payload: "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbce",
			expected: Port151Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385383,
				Longitude:    8.53804,
				Altitude:     488.2,
				Timestamp:    time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:      4.919,
				TTF:          time.Duration(24) * time.Second,
				PDOP:         1.5,
				Satellites:   11,
				Mac1:         "e0286d8a9478",
				Rssi1:        -53,
				Mac2:         helpers.StringPtr("f0b0140c96bb"),
				Rssi2:        helpers.Int8Ptr(-50),
			},
		},
		{
			port:    151,
			payload: "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcedeadbeef4242d6deadbeef4242d6",
			expected: Port151Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385383,
				Longitude:    8.53804,
				Altitude:     488.2,
				Timestamp:    time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:      4.919,
				TTF:          time.Duration(24) * time.Second,
				PDOP:         1.5,
				Satellites:   11,
				Mac1:         "e0286d8a9478",
				Rssi1:        -53,
				Mac2:         helpers.StringPtr("f0b0140c96bb"),
				Rssi2:        helpers.Int8Ptr(-50),
				Mac3:         helpers.StringPtr("deadbeef4242"),
				Rssi3:        helpers.Int8Ptr(-42),
				Mac4:         helpers.StringPtr("deadbeef4242"),
				Rssi4:        helpers.Int8Ptr(-42),
			},
		},
		{
			port:           151,
			payload:        "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcedeadbeef4242d6deadbeef4242d6deadbeef",
			skipValidation: true,
			expected: Port151Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385383,
				Longitude:    8.53804,
				Altitude:     488.2,
				Timestamp:    time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:      4.919,
				TTF:          time.Duration(24) * time.Second,
				PDOP:         1.5,
				Satellites:   11,
				Mac1:         "e0286d8a9478",
				Rssi1:        -53,
				Mac2:         helpers.StringPtr("f0b0140c96bb"),
				Rssi2:        helpers.Int8Ptr(-50),
				Mac3:         helpers.StringPtr("deadbeef4242"),
				Rssi3:        helpers.Int8Ptr(-42),
				Mac4:         helpers.StringPtr("deadbeef4242"),
				Rssi4:        helpers.Int8Ptr(-42),
			},
		},
		{
			port:    151,
			payload: "00000002d30b27008247b81312671bd164133718030ba1b2c3d4e5f6c0",
			expected: Port151Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385383,
				Longitude:    8.53804,
				Altitude:     488.2,
				Timestamp:    time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:      4.919,
				TTF:          time.Duration(24) * time.Second,
				PDOP:         1.5,
				Satellites:   11,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -64,
			},
		},
		{
			port:    151,
			payload: "00000002d30b27008247b81312671bd164133718030ba1b2c3d4e5f6c0",
			expected: Port151Payload{
				BufferLevel:  0,
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.385383,
				Longitude:    8.53804,
				Altitude:     488.2,
				Timestamp:    time.Date(2024, 10, 25, 17, 12, 4, 0, time.UTC),
				Battery:      4.919,
				TTF:          time.Duration(24) * time.Second,
				PDOP:         1.5,
				Satellites:   11,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -64,
			},
		},
		{
			port:    151,
			payload: "050021fcbbca14fbedc7cc00cc646ea2e40f8ca01004a1b2c3d4e5f6a8",
			expected: Port151Payload{
				BufferLevel:  1280,
				DutyCycle:    false,
				ConfigId:     4,
				ConfigChange: false,
				Moving:       true,
				Latitude:     -54.8019,
				Longitude:    -68.3029,
				Altitude:     20.4,
				Timestamp:    time.Date(2023, 5, 24, 23, 51, 0, 0, time.UTC),
				Battery:      3.980,
				TTF:          time.Duration(160) * time.Second,
				PDOP:         8,
				Satellites:   4,
				Mac1:         "a1b2c3d4e5f6",
				Rssi1:        -88,
			},
		},
		{
			port:    198,
			payload: "02",
			expected: Port198Payload{
				Reason: 2,
			},
		},
		{
			port:    198,
			payload: "04",
			expected: Port198Payload{
				Reason: 4,
			},
		},
		{
			port:    198,
			payload: "04313137",
			expected: Port198Payload{
				Reason: 4,
				Line:   helpers.StringPtr("117"),
			},
		},
		{
			port:    198,
			payload: "043131373a7372632f6770732e63",
			expected: Port198Payload{
				Reason: 4,
				Line:   helpers.StringPtr("117"),
				File:   helpers.StringPtr("src/gps.c"),
			},
		},
		{
			port:    198,
			payload: "043131373a7372632f6770732e633a6770735f73746172745f6d756c7469706c65",
			expected: Port198Payload{
				Reason:   4,
				Line:     helpers.StringPtr("117"),
				File:     helpers.StringPtr("src/gps.c"),
				Function: helpers.StringPtr("gps_start_multiple"),
			},
		},
		{
			port:    199,
			payload: "078f64e7dcffff00000f3e00127802",
			expected: Port199Payload{
				Constant: "078f64e7dcffff",
				Sequence: 3902,
				Number:   4728,
				Id:       2,
			},
		},
		{
			port:    199,
			payload: "078f64e7dcffff00000f3f00127902",
			expected: Port199Payload{
				Constant: "078f64e7dcffff",
				Sequence: 3903,
				Number:   4729,
				Id:       2,
			},
		},
		{
			port:    199,
			payload: "078f64e7dcffff00000f4000127a02",
			expected: Port199Payload{
				Constant: "078f64e7dcffff",
				Sequence: 3904,
				Number:   4730,
				Id:       2,
			},
		},
		{
			port:    199,
			payload: "078f64e7dcffff00000f4100127b02",
			expected: Port199Payload{
				Constant: "078f64e7dcffff",
				Sequence: 3905,
				Number:   4731,
				Id:       2,
			},
		},
		{
			port:    199,
			payload: "078f64e7dcffff00000f4200127c02",
			expected: Port199Payload{
				Constant: "078f64e7dcffff",
				Sequence: 3906,
				Number:   4732,
				Id:       2,
			},
		},
		{
			port:    199,
			payload: "078f64e7dcffff00000f4300127d02",
			expected: Port199Payload{
				Constant: "078f64e7dcffff",
				Sequence: 3907,
				Number:   4733,
				Id:       2,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder(WithSkipValidation(test.skipValidation))
			got, err := decoder.Decode(context.TODO(), test.payload, test.port)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got == nil || !reflect.DeepEqual(&got.Data, &test.expected) {
				t.Errorf("expected: %v\ngot: %v", test.expected, got)
			}
		})
	}

	t.Run("TestInvalidPayload", func(t *testing.T) {
		decoder := NewTagSLv1Decoder()
		_, err := decoder.Decode(context.TODO(), "", 1)
		if err == nil {
			t.Fatal("expected invalid payload")
		}
	})
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
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
			payload:  "0002d30c9300824c87117966c45dcd0f8118a1b2c3d4e5f6b8",
			port:     50,
			expected: nil,
		},
		{
			payload:  "0005f5e10000824c87117966c45dcd0f8118a1b2c3d4e5f6b8",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "0002d30c930bebc200117966c45dcd0f8118a1b2c3d4e5f6b8",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd01f418a1b2c3d4e5f6b8",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd157c18a1b2c3d4e5f6b8",
			port:     50,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f81180205a1b2c3d4e5f6a6",
			port:     51,
			expected: nil,
		},
		{
			payload:  "0005f5e10000824c87117966c45dcd0f81180205a1b2c3d4e5f6a6",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "0002d30c930bebc200117966c45dcd0f81180205a1b2c3d4e5f6a6",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd01f4180205a1b2c3d4e5f6a6",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd157c180205a1b2c3d4e5f6a6",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f81180202a1b2c3d4e5f6a6",
			port:     51,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 2),
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f8118021ca1b2c3d4e5f6a6",
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
			payload:  "00000002d30c9300824c87117966c45dcd0f812fa1b2c3d4e5f6c2",
			port:     150,
			expected: nil,
		},
		{
			payload:  "00000005f5e10000824c87117966c45dcd0f812fa1b2c3d4e5f6c2",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "00000002d30c930bebc200117966c45dcd0f812fa1b2c3d4e5f6c2",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd01f42fa1b2c3d4e5f6c2",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd157c2fa1b2c3d4e5f6c2",
			port:     150,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f0205a1b2c3d4e5f6c0",
			port:     151,
			expected: nil,
		},
		{
			payload:  "00000005f5e10000824c87117966c45dcd0f812f0205a1b2c3d4e5f6c0",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "00000002d30c930bebc200117966c45dcd0f812f0205a1b2c3d4e5f6c0",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd01f42f0205a1b2c3d4e5f6c0",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd157c2f0205a1b2c3d4e5f6c0",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f0202a1b2c3d4e5f6c0",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 2),
		},
		{
			payload:  "00000002d30c9300824c87117966c45dcd0f812f021ca1b2c3d4e5f6c0",
			port:     151,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Satellites", 28),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vValidationWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder()
			got, err := decoder.Decode(context.TODO(), test.payload, test.port)

			if err == nil && test.expected == nil {
				return
			}

			t.Logf("got %v", got)

			if err != nil && test.expected == nil || err == nil || !strings.Contains(err.Error(), test.expected.Error()) {
				t.Errorf("expected: %v\ngot: %v", test.expected, err)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode(context.TODO(), "00", 0)
	if err == nil || !errors.Is(err, helpers.ErrPortNotSupported) {
		t.Fatal("expected port not supported")
	}
}

func TestInvalidHexString(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode(context.TODO(), "xx", 6)
	if err == nil || err.Error() != "encoding/hex: invalid byte: U+0078 'x'" {
		t.Fatal("expected invalid hex byte")
	}
}

func TestFullDecode(t *testing.T) {
	tests := []struct {
		payload        string
		skipValidation bool
		expectedData   any
		port           uint8
	}{
		{
			port:    1,
			payload: "8002cdcd1300744f5e166018040b14341a",
			expectedData: Port1Payload{
				DutyCycle:    true,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       false,
				Latitude:     47.041811,
				Longitude:    7.622494,
				Altitude:     572.8,
				Year:         24,
				Month:        4,
				Day:          11,
				Hour:         20,
				Minute:       52,
				Second:       26,
			},
		},
		{
			port:    5,
			payload: "0c24651155ce55b8602232e20f52b0ac8ba91fedaaa5603197f93781a90aecdafa5fe8bc02ecdafa5fe8bd8c59c3c960f0a5",
			expectedData: Port5Payload{
				DutyCycle:    false,
				ConfigId:     1,
				ConfigChange: true,
				Moving:       false,
				Mac1:         "24651155ce55",
				Rssi1:        -72,
				Mac2:         helpers.StringPtr("602232e20f52"),
				Rssi2:        helpers.Int8Ptr(-80),
				Mac3:         helpers.StringPtr("ac8ba91fedaa"),
				Rssi3:        helpers.Int8Ptr(-91),
				Mac4:         helpers.StringPtr("603197f93781"),
				Rssi4:        helpers.Int8Ptr(-87),
				Mac5:         helpers.StringPtr("0aecdafa5fe8"),
				Rssi5:        helpers.Int8Ptr(-68),
				Mac6:         helpers.StringPtr("02ecdafa5fe8"),
				Rssi6:        helpers.Int8Ptr(-67),
				Mac7:         helpers.StringPtr("8c59c3c960f0"),
				Rssi7:        helpers.Int8Ptr(-91),
			},
		},
		{
			port:    7,
			payload: "66ec04bb00e0286d8aabfcbbec6c9a74b58fb2726c9a74b58db1e0286d8a9478cbf0b0140c96bbd2260122180d42ad",
			expectedData: Port7Payload{
				Timestamp: time.Date(2024, 9, 19, 11, 2, 19, 0, time.UTC),
				Mac1:      "e0286d8aabfc",
				Rssi1:     -69,
				Mac2:      helpers.StringPtr("ec6c9a74b58f"),
				Rssi2:     helpers.Int8Ptr(-78),
				Mac3:      helpers.StringPtr("726c9a74b58d"),
				Rssi3:     helpers.Int8Ptr(-79),
				Mac4:      helpers.StringPtr("e0286d8a9478"),
				Rssi4:     helpers.Int8Ptr(-53),
				Mac5:      helpers.StringPtr("f0b0140c96bb"),
				Rssi5:     helpers.Int8Ptr(-46),
				Mac6:      helpers.StringPtr("260122180d42"),
				Rssi6:     helpers.Int8Ptr(-83),
			},
		},
		{
			port:    105,
			payload: "0028672658500172a741b1e238b572a741b1e08bb03498b5c583e2b172a741b1e0cda772a741beed4cc472a741beef53b7",
			expectedData: Port105Payload{
				BufferLevel:  40,
				Timestamp:    time.Date(2024, 11, 2, 16, 50, 24, 0, time.UTC),
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Mac1:         "72a741b1e238",
				Rssi1:        -75,
				Mac2:         helpers.StringPtr("72a741b1e08b"),
				Rssi2:        helpers.Int8Ptr(-80),
				Mac3:         helpers.StringPtr("3498b5c583e2"),
				Rssi3:        helpers.Int8Ptr(-79),
				Mac4:         helpers.StringPtr("72a741b1e0cd"),
				Rssi4:        helpers.Int8Ptr(-89),
				Mac5:         helpers.StringPtr("72a741beed4c"),
				Rssi5:        helpers.Int8Ptr(-60),
				Mac6:         helpers.StringPtr("72a741beef53"),
				Rssi6:        helpers.Int8Ptr(-73),
			},
		},
		{
			port:    110,
			payload: "00430102d43ffa00772d870ea367250ef60eda",
			expectedData: Port110Payload{
				BufferLevel:  67,
				Timestamp:    time.Date(2024, 11, 1, 17, 25, 10, 0, time.UTC),
				DutyCycle:    false,
				ConfigId:     0,
				ConfigChange: false,
				Moving:       true,
				Latitude:     47.464442,
				Longitude:    7.810439,
				Altitude:     374.7,
				Battery:      3.802,
			},
		},
		{
			port:    198,
			payload: "01",
			expectedData: Port198Payload{
				Reason: 1,
			},
		},
		{
			port:    198,
			payload: "02",
			expectedData: Port198Payload{
				Reason: 2,
			},
		},
	}

	for _, test := range tests {
		decoder := NewTagSLv1Decoder(WithSkipValidation(test.skipValidation))
		t.Run(fmt.Sprintf("TestFullDecodeWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoded, err := decoder.Decode(context.TODO(), test.payload, test.port)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(&decoded.Data, &test.expectedData) {
				t.Errorf("expected data: %v, got: %v", test.expectedData, decoded.Data)
			}
		})
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode(context.TODO(), "deadbeef", 1)

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooShort) {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewTagSLv1Decoder()
	_, err := decoder.Decode(context.TODO(), "deadbeef4242deadbeef4242deadbeef4242", 1)

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooLong) {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload         string
		port            uint8
		allowNoFeatures bool
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			port:    1,
		},
		{
			payload: "00",
			port:    2,
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
			payload:         "012c141e9c455738304543434343460078012c01a8c0",
			port:            8,
			allowNoFeatures: true,
		},
		{
			payload: "0002d308b50082457f16eb66c4a5cd0ed3",
			port:    10,
		},
		{
			payload: "0002d308b50082457f16eb66c4a5cd0ed32a0807",
			port:    10,
		},
		{
			payload: "800ee5",
			port:    15,
		},
		{
			payload: "0002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:    50,
		},
		{
			payload: "0002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1fcf528f8634fb552a8db7bd6b5b9e0286d8aabfcbc",
			port:    51,
		},
		{
			payload: "000166c4a5ba80e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0a1b2c3d4e5f6ae",
			port:    105,
		},
		{
			payload: "00020002d309ae008247c5113966c45d640f7e",
			port:    110,
		},
		{
			payload: "00020002d309ae008247c5113966c45d640f7e2a0807",
			port:    110,
		},
		{
			payload: "00020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:    150,
		},
		{
			payload: "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcea1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:    151,
		},
		{
			payload: "043131373a7372632f6770732e633a6770735f73746172745f6d756c7469706c65",
			port:    198,
		},
		{
			payload:         "078f64e7dcffff00000f4300127d02",
			port:            199,
			allowNoFeatures: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewTagSLv1Decoder()
			decodedPayload, _ := d.Decode(context.TODO(), test.payload, test.port)

			if len(decodedPayload.GetFeatures()) == 0 && !test.allowNoFeatures {
				t.Error("expected features, got none")
			}

			if decodedPayload.Is(decoder.FeatureTimestamp) {
				timestamp, ok := decodedPayload.Data.(decoder.UplinkFeatureTimestamp)
				if !ok {
					t.Fatalf("expected UplinkFeatureTimestamp, got %T", decodedPayload)
				}
				if timestamp.GetTimestamp() == nil {
					t.Fatalf("expected non nil timestamp")
				}
			}

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
				buffered.IsBuffered()
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
				// call function to check if it panics
				batteryVoltage.GetLowBattery()
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
				config.GetAcceleration()
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
				config.GetLowTemperatureThreshold()
				config.GetHighTemperatureThreshold()
				config.GetAccessPointsThreshold()
				config.GetBatchSize()
				config.GetBufferSize()
				config.GetDataRate()
			}
			if decodedPayload.Is(decoder.FeatureConfigChange) {
				configChange, ok := decodedPayload.Data.(decoder.UplinkFeatureConfigChange)
				if !ok {
					t.Fatalf("expected UplinkFeatureConfigChange, got %T", decodedPayload)
				}
				// call functions to check if it panics
				configChange.GetConfigId()
				configChange.GetConfigChange()
			}
			if decodedPayload.Is(decoder.FeatureButton) {
				button, ok := decodedPayload.Data.(decoder.UplinkFeatureButton)
				if !ok {
					t.Fatalf("expected UplinkFeatureButton, got %T", decodedPayload)
				}
				// call function to check if it panics
				button.GetPressed()
			}
			if decodedPayload.Is(decoder.FeatureFirmwareVersion) {
				firmwareVersion, ok := decodedPayload.Data.(decoder.UplinkFeatureFirmwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureFirmwareVersion, got %T", decodedPayload)
				}
				if firmwareVersion.GetFirmwareVersion() == nil {
					t.Fatalf("expected non nil firmware version")
				}
				firmwareVersion.GetFirmwareHash()
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
			if decodedPayload.Is(decoder.FeatureResetReason) {
				resetReason, ok := decodedPayload.Data.(decoder.UplinkFeatureResetReason)
				if !ok {
					t.Fatalf("expected UplinkFeatureResetReason, got %T", decodedPayload)
				}
				// call function to check if it panics
				resetReason.GetResetReason()
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
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
			expected: []string{"\"movingInterval\": \"5m0s\"", "\"steadyInterval\": \"1h0m0s\"", "\"accelerometerThreshold\": \"300mg\"", "\"deviceState\": \"steady\""},
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
			expected: []string{"\"minRssiValue\": -100", "\"advertisingFilter\": \"EW80ECCCCF\""},
		},
		{
			payload:  "0002d308b50082457f16eb66c4a5cd0ed32a0807",
			port:     10,
			expected: []string{"\"timestamp\": \"2024-08-20T14:18:53Z\"", "\"battery\": \"3.795v\"", "\"ttf\": \"42s\"", "\"pdop\": \"4.0m\""},
		},
		{
			payload:  "800ee5",
			port:     15,
			expected: []string{"\"lowBattery\": false", "\"battery\": \"3.813v\""},
		},
		{
			payload:  "0002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:     50,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-08-20T09:11:41Z\"", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -87"},
		},
		{
			payload:  "0002d30ba000824ace1122671b983e0eea340b06726c9a74b58db1fcf528f8634fb552a8db7bd6b5b9e0286d8aabfcbc",
			port:     51,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-10-25T13:08:14Z\"", "\"pdop\": \"5.5m\"", "\"mac1\": \"726c9a74b58d\"", "\"rssi1\": -79"},
		},
		{
			payload:  "000166c4a5ba80e0286d8aabfcb1e0286d8a9478c2ec6c9a74b58fad726c9a74b58dadf0b0140c96bbd0a1b2c3d4e5f6ae",
			port:     105,
			expected: []string{"\"moving\": false", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -79"},
		},
		{
			payload:  "00020002d309ae008247c5113966c45d640f7e2e0707",
			port:     110,
			expected: []string{"\"timestamp\": \"2024-08-20T09:09:56Z\"", "\"battery\": \"3.966v\"", "\"ttf\": \"46s\"", "\"pdop\": \"3.5m\""},
		},
		{
			payload:  "00020002d30c9300824c87117966c45dcd0f8118e0286d8aabfca9f0b0140c96bbc8726c9a74b58da8e0286d8a9478bf",
			port:     150,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-08-20T09:11:41Z\"", "\"mac1\": \"e0286d8aabfc\"", "\"rssi1\": -87", "\"ttf\": \"24s\""},
		},
		{
			payload:  "00000002d30b27008247b81312671bd164133718030be0286d8a9478cbf0b0140c96bbcea1b2c3d4e5f6aea1b2c3d4e5f6ae",
			port:     151,
			expected: []string{"\"moving\": false", "\"timestamp\": \"2024-10-25T17:12:04Z\"", "\"pdop\": \"1.5m\"", "\"mac1\": \"e0286d8a9478\"", "\"rssi1\": -53", "\"ttf\": \"24s\""},
		},
		{
			payload:  "01",
			port:     198,
			expected: []string{"\"reason\": \"lrr1110-failure\""},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder()

			data, _ := decoder.Decode(context.TODO(), test.payload, test.port)

			marshaled, err := json.MarshalIndent(map[string]any{
				"data": data.Data,
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

func TestDeviceState(t *testing.T) {
	tests := []struct {
		data     Port4Payload
		expected string
	}{
		{
			data:     Port4Payload{DeviceState: 0},
			expected: "\"deviceState\":\"unknown\"",
		},
		{
			data:     Port4Payload{DeviceState: 1},
			expected: "\"deviceState\":\"moving\"",
		},
		{
			data:     Port4Payload{DeviceState: 2},
			expected: "\"deviceState\":\"steady\"",
		},
		{
			data:     Port4Payload{DeviceState: 3},
			expected: "\"deviceState\":\"unknown\"",
		},
	}

	for _, test := range tests {
		t.Run("TestDeviceState", func(t *testing.T) {
			result, err := test.data.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected err %s", err)
			}
			t.Logf("%s", result)
			if !strings.Contains(string(result), test.expected) {
				t.Errorf("expected to find %s", test.expected)
			}
		})
	}
}

func TestResetReason(t *testing.T) {
	tests := []struct {
		data     Port198Payload
		expected decoder.ResetReason
	}{
		{
			data:     Port198Payload{Reason: 0},
			expected: decoder.ResetReasonUnknown,
		},
		{
			data:     Port198Payload{Reason: 1},
			expected: decoder.ResetReasonLrr1110FailCode,
		},
		{
			data:     Port198Payload{Reason: 2},
			expected: decoder.ResetReasonPowerReset,
		},
		{
			data:     Port198Payload{Reason: 3},
			expected: decoder.ResetReasonPinReset,
		},
		{
			data:     Port198Payload{Reason: 4},
			expected: decoder.ResetReasonWatchdog,
		},
		{
			data:     Port198Payload{Reason: 5},
			expected: decoder.ResetReasonSystemReset,
		},
		{
			data:     Port198Payload{Reason: 6},
			expected: decoder.ResetReasonOtherReset,
		},
		{
			data:     Port198Payload{Reason: 7},
			expected: decoder.ResetReasonUnknown,
		},
	}

	for _, test := range tests {
		t.Run("TestResetReason", func(t *testing.T) {
			result := test.data.GetResetReason()
			if test.expected != result {
				t.Errorf("expected %v, received %v", test.expected, result)
			}
		})
	}
}
