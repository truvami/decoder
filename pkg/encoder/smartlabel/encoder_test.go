package smartlabel

import (
	"fmt"
	"testing"

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
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.expected), func(t *testing.T) {
			encoder := NewSmartlabelv1Encoder()
			got, err := encoder.Encode(test.data, test.port)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != test.expected {
				t.Errorf("expected: %v", test.expected)
				t.Errorf("received: %v", got)
			}
		})
	}
}
