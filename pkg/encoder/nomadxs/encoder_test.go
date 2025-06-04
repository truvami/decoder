package nomadxs

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		data     any
		port     uint8
		expected string
	}{
		{
			data: nomadxs.Port1Payload{
				Moving:             false,
				Latitude:           27.9881,
				Longitude:          86.9250,
				Altitude:           4848,
				Year:               78,
				Month:              5,
				Day:                8,
				Hour:               10,
				Minute:             20,
				Second:             30,
				TimeToFix:          time.Duration(16) * time.Second,
				AmbientLight:       40887,
				AccelerometerXAxis: 42,
				AccelerometerYAxis: 73,
				AccelerometerZAxis: 420,
				Temperature:        -22.98,
				Pressure:           314.3,
			},
			port:     1,
			expected: "0001ab1084052e5ec8bd604e05080a141e109fb7002a004901a4f7060c47",
		},
		{
			data: nomadxs.Port1Payload{
				Moving:             true,
				Latitude:           35.3606,
				Longitude:          138.7274,
				Altitude:           3776,
				Year:               21,
				Month:              2,
				Day:                23,
				Hour:               8,
				Minute:             30,
				Second:             0,
				TimeToFix:          time.Duration(34) * time.Second,
				AmbientLight:       65023,
				AccelerometerXAxis: 21,
				AccelerometerYAxis: 37,
				AccelerometerZAxis: 42,
				Temperature:        8.42,
				Pressure:           641.6,
			},
			port:     1,
			expected: "01021b8f580844cfe89380150217081e0022fdff00150025002a034a1910",
		},
		{
			data: nomadxs.Port1Payload{
				Moving:             false,
				Latitude:           -32.6532,
				Longitude:          -70.0109,
				Altitude:           6461,
				Year:               85,
				Month:              1,
				Day:                16,
				Hour:               11,
				Minute:             17,
				Second:             51,
				TimeToFix:          time.Duration(67) * time.Second,
				AmbientLight:       34509,
				AccelerometerXAxis: 56,
				AccelerometerYAxis: 112,
				AccelerometerZAxis: 3200,
				Temperature:        -2.05,
				Pressure:           518.3,
			},
			port:     1,
			expected: "00fe0dc070fbd3b7ecfc625501100b11334386cd003800700c80ff33143f",
		},
		{
			data: nomadxs.Port1Payload{
				Moving:             true,
				Latitude:           5.2270,
				Longitude:          -60.7500,
				Altitude:           2810,
				Year:               92,
				Month:              7,
				Day:                15,
				Hour:               20,
				Minute:             40,
				Second:             5,
				TimeToFix:          time.Duration(239) * time.Second,
				AmbientLight:       6581,
				AccelerometerXAxis: 82,
				AccelerometerYAxis: 96,
				AccelerometerZAxis: 120,
				Temperature:        17.45,
				Pressure:           738.0,
			},
			port:     1,
			expected: "01004fc1f8fc6107506dc45c070f142805ef19b500520060007806d11cd4",
		},
		{
			data: nomadxs.Port4Payload{
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 120,
				HeartbeatInterval:               3600,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    10,
				AccelerometerDelay:              1000,
				FirmwareVersionMajor:            0,
				FirmwareVersionMinor:            5,
				FirmwareVersionPatch:            21,
				HardwareVersionType:             0,
				HardwareVersionRevision:         5,
				BatteryKeepAliveMessageInterval: 900,
				ReJoinInterval:                  300,
				AccuracyEnhancement:             4,
				LightLowerThreshold:             32,
				LightUpperThreshold:             128,
			},
			port:     4,
			expected: "000000780000007800000e100078000a03e80005150005000003840000012c0400200080",
		},
		{
			data: nomadxs.Port4Payload{
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 300,
				HeartbeatInterval:               7200,
				GPSTimeoutWhileWaitingForFix:    240,
				AccelerometerWakeupThreshold:    100,
				AccelerometerDelay:              1500,
				FirmwareVersionMajor:            0,
				FirmwareVersionMinor:            7,
				FirmwareVersionPatch:            3,
				HardwareVersionType:             0,
				HardwareVersionRevision:         8,
				BatteryKeepAliveMessageInterval: 1800,
				ReJoinInterval:                  300,
				AccuracyEnhancement:             8,
				LightLowerThreshold:             64,
				LightUpperThreshold:             256,
			},
			port:     4,
			expected: "000000780000012c00001c2000f0006405dc0007030008000007080000012c0800400100",
		},
		{
			data: nomadxs.Port4Payload{
				LocalizationIntervalWhileMoving: 300,
				LocalizationIntervalWhileSteady: 1800,
				HeartbeatInterval:               21600,
				GPSTimeoutWhileWaitingForFix:    480,
				AccelerometerWakeupThreshold:    200,
				AccelerometerDelay:              2000,
				FirmwareVersionMajor:            1,
				FirmwareVersionMinor:            0,
				FirmwareVersionPatch:            8,
				HardwareVersionType:             1,
				HardwareVersionRevision:         2,
				BatteryKeepAliveMessageInterval: 3600,
				ReJoinInterval:                  600,
				AccuracyEnhancement:             16,
				LightLowerThreshold:             128,
				LightUpperThreshold:             512,
			},
			port:     4,
			expected: "0000012c000007080000546001e000c807d0010008010200000e10000002581000800200",
		},
		{
			data: nomadxs.Port4Payload{
				LocalizationIntervalWhileMoving: 600,
				LocalizationIntervalWhileSteady: 3600,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    480,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              4000,
				FirmwareVersionMajor:            1,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            2,
				HardwareVersionType:             1,
				HardwareVersionRevision:         6,
				BatteryKeepAliveMessageInterval: 3600,
				ReJoinInterval:                  600,
				AccuracyEnhancement:             24,
				LightLowerThreshold:             128,
				LightUpperThreshold:             2048,
			},
			port:     4,
			expected: "0000025800000e100001518001e0012c0fa0010102010600000e10000002581800800800",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.expected), func(t *testing.T) {
			encoder := NewNomadXSv1Encoder()
			received, err := encoder.Encode(test.data, test.port)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if received != test.expected {
				t.Errorf("expected: %v\n", test.expected)
				t.Errorf("received: %v\n", received)
			}
		})
	}
}

func TestInvalidData(t *testing.T) {
	encoder := NewNomadXSv1Encoder()
	_, err := encoder.Encode(nil, 1)
	if err == nil || err.Error() != "data must be a struct" {
		t.Fatal("expected data must be a struct")
	}
}

func TestInvalidPort(t *testing.T) {
	encoder := NewNomadXSv1Encoder()
	_, err := encoder.Encode(nil, 0)
	if err == nil || !errors.Is(err, common.ErrPortNotSupported) {
		t.Fatal("expected port not supported")
	}
}
