package tagsl

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected interface{}
	}{
		// {
		// 	payload:  "8002cdcd1300744f5e166018040b14341a",
		// 	port:     -1,
		// 	expected: nil,
		// },
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			port:    1,
			expected: GNSSPayload{
				Moving: false,
				Lat:    47.041811,
				Lon:    7.622494,
				Alt:    5728,
				Year:   24,
				Month:  4,
				Day:    11,
				Hour:   20,
				Minute: 52,
				Second: 26,
			},
		},
		{
			payload: "808c59c3c99fc0ad",
			port:    5,
			expected: WifiPayload{
				Moving: false,
				Mac1:   "8c59c3c99fc0",
				Rssi1:  -83,
			},
		},
		{
			payload: "80e0286d8a2742a1",
			port:    5,
			expected: WifiPayload{
				Moving: false,
				Mac1:   "e0286d8a2742",
				Rssi1:  -95,
			},
		},
		{
			payload: "001f3fd57cecb4f0b0140c96bbb2e0286d8a9478b8",
			port:    5,
			expected: WifiPayload{
				Moving: false,
				Mac1:   "1f3fd57cecb4",
				Rssi1:  -16,
				Mac2:   "b0140c96bbb2",
				Rssi2:  -32,
				Mac3:   "286d8a9478b8",
				Rssi3:  0,
			},
		},
		{
			payload: "822f0101f052fab920feafd0e4158b38b9afe05994cb2f5cb2",
			port:    3,
			expected: BlePayload{
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
			payload: "01eb0101f052fab920feadd0e4158b38b9afe05994cb2f5cad",
			port:    3,
			expected: BlePayload{
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
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v-%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagSLv1Decoder()
			got, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}
		})
	}
}
