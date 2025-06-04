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
