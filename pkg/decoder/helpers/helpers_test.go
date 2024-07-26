package helpers

import (
	"fmt"
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

type Port1Payload struct {
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
		TargetType: reflect.TypeOf(Port1Payload{}),
	}

	tests := []struct {
		payload  string
		config   decoder.PayloadConfig
		expected interface{}
	}{
		{
			payload: "8002cdcd1300744f5e166018040b14341a",
			config:  config,
			expected: Port1Payload{
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
			payload := decodedData.(Port1Payload)

			// Check the decoded data against the expected data using reflect.DeepEqual
			if !reflect.DeepEqual(payload, test.expected) {
				t.Fatalf("decoded data does not match expected data expected: %+v got: %+v", test.expected, payload)
			}
		})
	}
}
func TestConvertFieldToType(t *testing.T) {
	tests := []struct {
		value     interface{}
		fieldType reflect.Kind
		expected  interface{}
	}{
		{
			value:     10,
			fieldType: reflect.Int,
			expected:  10,
		},
		{
			value:     10,
			fieldType: reflect.Int8,
			expected:  int8(10),
		},
		{
			value:     10,
			fieldType: reflect.Int16,
			expected:  int16(10),
		},
		{
			value:     10,
			fieldType: reflect.Int32,
			expected:  int32(10),
		},
		{
			value:     10,
			fieldType: reflect.Int64,
			expected:  int64(10),
		},
		{
			value:     10,
			fieldType: reflect.Uint,
			expected:  uint(10),
		},
		{
			value:     10,
			fieldType: reflect.Uint8,
			expected:  uint8(10),
		},
		{
			value:     10,
			fieldType: reflect.Uint16,
			expected:  uint16(10),
		},
		{
			value:     10,
			fieldType: reflect.Uint32,
			expected:  uint32(10),
		},
		{
			value:     10,
			fieldType: reflect.Uint64,
			expected:  uint64(10),
		},
		{
			value:     10,
			fieldType: reflect.Float64,
			expected:  float64(10),
		},
		{
			value:     "hello",
			fieldType: reflect.String,
			expected:  "hello",
		},
		{
			value:     1,
			fieldType: reflect.Bool,
			expected:  true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v_%v", test.value, test.fieldType), func(t *testing.T) {
			result := convertFieldToType(test.value, test.fieldType)
			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("converted value does not match expected value expected: %v got: %v", test.expected, result)
			}
		})
	}
}

func TestInvalidPayload(t *testing.T) {
	_, err := Parse("", decoder.PayloadConfig{
		Fields: []decoder.FieldConfig{
			{Name: "Moving", Start: 0, Length: 1},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field out of bounds")
	}

	_, err = Parse("01", decoder.PayloadConfig{
		Fields: []decoder.FieldConfig{
			{Name: "Moving", Start: 0, Length: 2},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field out of bounds")
	}

	_, err = Parse("01", decoder.PayloadConfig{
		Fields: []decoder.FieldConfig{
			{Name: "Moving", Start: 10, Length: 1},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field start out of bounds")
	}
}
