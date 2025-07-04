package smartlabel

import (
	"errors"
	"fmt"
	"testing"

	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder/smartlabel/v1"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		data     any
		port     uint8
		expected string
	}{
		{
			data: smartlabel.Port1Payload{
				BatteryVoltage:      3.457,
				PhotovoltaicVoltage: 1.672,
			},
			port:     1,
			expected: "0d810688",
		},
		{
			data: smartlabel.Port1Payload{
				BatteryVoltage:      3.561,
				PhotovoltaicVoltage: 1.890,
			},
			port:     1,
			expected: "0de90762",
		},
		{
			data: smartlabel.Port1Payload{
				BatteryVoltage:      3.785,
				PhotovoltaicVoltage: 2.459,
			},
			port:     1,
			expected: "0ec9099b",
		},
		{
			data: smartlabel.Port1Payload{
				BatteryVoltage:      3.817,
				PhotovoltaicVoltage: 3.322,
			},
			port:     1,
			expected: "0ee90cfa",
		},
		{
			data: smartlabel.Port2Payload{
				Temperature: 12.06,
				Humidity:    40.5,
			},
			port:     2,
			expected: "04b651",
		},
		{
			data: smartlabel.Port2Payload{
				Temperature: 8.46,
				Humidity:    57.0,
			},
			port:     2,
			expected: "034e72",
		},
		{
			data: smartlabel.Port2Payload{
				Temperature: 2.32,
				Humidity:    65.5,
			},
			port:     2,
			expected: "00e883",
		},
		{
			data: smartlabel.Port2Payload{
				Temperature: -4.98,
				Humidity:    78.0,
			},
			port:     2,
			expected: "fe0e9c",
		},
		{
			data: smartlabel.Port11Payload{
				BatteryVoltage:      3.457,
				PhotovoltaicVoltage: 1.672,
				Temperature:         12.06,
				Humidity:            40.5,
			},
			port:     11,
			expected: "0d81068804b651",
		},
		{
			data: smartlabel.Port11Payload{
				BatteryVoltage:      3.561,
				PhotovoltaicVoltage: 1.890,
				Temperature:         8.46,
				Humidity:            57.0,
			},
			port:     11,
			expected: "0de90762034e72",
		},
		{
			data: smartlabel.Port11Payload{
				BatteryVoltage:      3.785,
				PhotovoltaicVoltage: 2.459,
				Temperature:         2.32,
				Humidity:            65.5,
			},
			port:     11,
			expected: "0ec9099b00e883",
		},
		{
			data: smartlabel.Port11Payload{
				BatteryVoltage:      3.817,
				PhotovoltaicVoltage: 3.322,
				Temperature:         -4.98,
				Humidity:            78.0,
			},
			port:     11,
			expected: "0ee90cfafe0e9c",
		},
		{
			data: Port128Payload{
				DataRate:                   0,
				SteadyInterval:             21600,
				MovingInterval:             3600,
				HeartbeatInterval:          24,
				AccelerationThreshold:      300,
				AccelerationDelay:          1500,
				TemperaturePollingInterval: 900,
				TemperatureUplinkInterval:  3600,
				TemperatureUpperThreshold:  +40,
				TemperatureLowerThreshold:  -20,
				AccessPointsThreshold:      2,
			},
			port:     128,
			expected: "0054600e1018012c05dc03840e1028ec02",
		},
		{
			data: Port128Payload{
				DataRate:                   2,
				SteadyInterval:             7200,
				MovingInterval:             1800,
				HeartbeatInterval:          12,
				AccelerationThreshold:      200,
				AccelerationDelay:          1200,
				TemperaturePollingInterval: 600,
				TemperatureUplinkInterval:  1800,
				TemperatureUpperThreshold:  +30,
				TemperatureLowerThreshold:  -15,
				AccessPointsThreshold:      3,
			},
			port:     128,
			expected: "021c2007080c00c804b0025807081ef103",
		},
		{
			data: Port128Payload{
				DataRate:                   4,
				SteadyInterval:             3600,
				MovingInterval:             900,
				HeartbeatInterval:          6,
				AccelerationThreshold:      100,
				AccelerationDelay:          1000,
				TemperaturePollingInterval: 300,
				TemperatureUplinkInterval:  1200,
				TemperatureUpperThreshold:  +20,
				TemperatureLowerThreshold:  -10,
				AccessPointsThreshold:      6,
			},
			port:     128,
			expected: "040e10038406006403e8012c04b014f606",
		},
		{
			data: Port128Payload{
				DataRate:                   6,
				SteadyInterval:             1800,
				MovingInterval:             600,
				HeartbeatInterval:          4,
				AccelerationThreshold:      10,
				AccelerationDelay:          1000,
				TemperaturePollingInterval: 300,
				TemperatureUplinkInterval:  900,
				TemperatureUpperThreshold:  +20,
				TemperatureLowerThreshold:  -0,
				AccessPointsThreshold:      1,
			},
			port:     128,
			expected: "060708025804000a03e8012c0384140001",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.expected), func(t *testing.T) {
			encoder := NewSmartlabelv1Encoder()
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
	encoder := NewSmartlabelv1Encoder()
	_, err := encoder.Encode(nil, 1)
	if err == nil || err.Error() != "data must be a struct" {
		t.Fatal("expected data must be a struct")
	}
}

func TestInvalidPort(t *testing.T) {
	encoder := NewSmartlabelv1Encoder()
	_, err := encoder.Encode(nil, 0)
	if err == nil || !errors.Is(err, common.ErrPortNotSupported) {
		t.Fatal("expected port not supported")
	}
}
