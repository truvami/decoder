package helpers

import (
	"reflect"
	"testing"

	"github.com/truvami/decoder/pkg/decoder"
)

func TestInvalidHexString(t *testing.T) {
	_, err := hexStringToBytes("invalid")
	if err == nil {
		t.Fatalf("expected error while decoding hex string")
	}
}

func TestHexStringToBytes(t *testing.T) {
	_, err := hexStringToBytes("8002cdcd1300744f5e166018040b14341a")
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
}

type GNSSPayload struct {
	Moving bool    `json:"moving"`
	Lat    float64 `json:"gps_lat"`
	Lon    float64 `json:"gps_lon"`
	Alt    float64 `json:"gps_alt"`
	Year   int     `json:"year"`
	Month  int     `json:"month"`
	Day    int     `json:"day"`
	Hour   int     `json:"hour"`
	Minute int     `json:"minute"`
	Second int     `json:"second"`
	TS     int64   `json:"ts"`
}

func TestParse(t *testing.T) {
	config := decoder.PayloadConfig{
		Fields: []decoder.FieldConfig{
			{Name: "Moving", Start: 0, Length: 1},
			{Name: "Lat", Start: 1, Length: 4, Transform: func(v interface{}) interface{} {
				return float64(v.(int)) / 1000000
			}},
			{Name: "Lon", Start: 5, Length: 4, Transform: func(v interface{}) interface{} {
				return float64(v.(int)) / 1000000
			}},
			{Name: "Alt", Start: 9, Length: 2},
			{Name: "Year", Start: 11, Length: 1},
			{Name: "Month", Start: 12, Length: 1},
			{Name: "Day", Start: 13, Length: 1},
			{Name: "Hour", Start: 14, Length: 1},
			{Name: "Minute", Start: 15, Length: 1},
			{Name: "Second", Start: 16, Length: 1},
		},
		TargetType: reflect.TypeOf(GNSSPayload{}),
	}

	tests := []struct {
		payload  string
		config   decoder.PayloadConfig
		expected interface{}
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			config:  config,
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
	}

	for _, test := range tests {
		t.Run(test.payload, func(t *testing.T) {
			decodedData, err := Parse(test.payload, test.config)
			if err != nil {
				t.Fatalf("error decoding payload: %v", err)
			}

			// Type assert to Payload
			payload := decodedData.(GNSSPayload)

			// Check the decoded data against the expected data using reflect.DeepEqual
			if !reflect.DeepEqual(payload, test.expected) {
				t.Fatalf("decoded data does not match expected data expected: %+v got: %+v", test.expected, payload)
			}
		})
	}
}
