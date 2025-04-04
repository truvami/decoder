package common

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestInvalidHexString(t *testing.T) {
	_, err := HexStringToBytes("invalid")
	if err == nil {
		t.Fatalf("expected error while decoding hex string")
	}
}

func TestHexStringToBytes(t *testing.T) {
	_, err := HexStringToBytes("8002cdcd1300744f5e166018040b14341a")
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
}

type Port1Payload struct {
	Moving bool    `json:"moving"`
	Lat    float64 `json:"gpsLat" validate:"gte=-90,lte=90"`
	Lon    float64 `json:"gpsLon" validate:"gte=-180,lte=180"`
	Alt    float64 `json:"gpsAlt" validate:"gte=0,lte=20000"`
	Year   int     `json:"year" validate:"gte=0,lte=255"`
	Month  int     `json:"month" validate:"gte=0,lte=255"`
	Day    int     `json:"day" validate:"gte=1,lte=31"`
	Hour   int     `json:"hour" validate:"gte=0,lte=23"`
	Minute int     `json:"minute" validate:"gte=0,lte=59"`
	Second int     `json:"second" validate:"gte=0,lte=59"`
	TS     int64   `json:"ts"`
}

func TestParse(t *testing.T) {
	config := PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 0, Length: 1},
			{Name: "Lat", Start: 1, Length: 4, Transform: func(v any) any {
				return float64(v.(int)) / 1000000
			}},
			{Name: "Lon", Start: 5, Length: 4, Transform: func(v any) any {
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
		config   PayloadConfig
		expected any
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
			decodedData, err := Parse(test.payload, &test.config)
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
		value     any
		fieldType reflect.Type
		expected  any
	}{
		{
			value:     10,
			fieldType: reflect.TypeOf(int(0)),
			expected:  10,
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(int8(0)),
			expected:  int8(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(int16(0)),
			expected:  int16(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(int32(0)),
			expected:  int32(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(int64(0)),
			expected:  int64(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(uint(0)),
			expected:  uint(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(uint8(0)),
			expected:  uint8(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(uint16(0)),
			expected:  uint16(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(uint32(0)),
			expected:  uint32(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(uint64(0)),
			expected:  uint64(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(float32(0)),
			expected:  float32(10),
		},
		{
			value:     10,
			fieldType: reflect.TypeOf(float64(0)),
			expected:  float64(10),
		},
		{
			value:     "hello",
			fieldType: reflect.TypeOf(string("")),
			expected:  "hello",
		},
		{
			value:     1,
			fieldType: reflect.TypeOf(bool(false)),
			expected:  true,
		},
		{
			value:     200,
			fieldType: reflect.TypeOf(time.Duration(0)),
			expected:  time.Duration(200) * time.Nanosecond,
		},
		{
			value:     42,
			fieldType: reflect.TypeOf(time.Time{}),
			expected:  time.Date(1970, 1, 1, 0, 0, 42, 0, time.UTC),
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
	_, err := Parse("", &PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 0, Length: 1},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field out of bounds")
	}

	_, err = Parse("01", &PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 0, Length: 2},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field out of bounds")
	}

	_, err = Parse("01", &PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 10, Length: 1},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field start out of bounds")
	}
}

func TestUintToBinaryArray(t *testing.T) {
	tests := []struct {
		value    uint64
		length   int
		expected []uint8
	}{
		{
			value:    0x01,
			length:   1,
			expected: []uint8{1},
		},
		{
			value:    0x01,
			length:   2,
			expected: []uint8{0, 1},
		},
		{
			value:    0x03,
			length:   2,
			expected: []uint8{1, 1},
		},
		{
			value:    0x03,
			length:   4,
			expected: []uint8{0, 0, 1, 1},
		},
		{
			value:    0x03,
			length:   8,
			expected: []uint8{0, 0, 0, 0, 0, 0, 1, 1},
		},
		{
			value:  0x03,
			length: 16,
			expected: []uint8{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v_%v", test.value, test.length), func(t *testing.T) {
			result := UintToBinaryArray(test.value, test.length)
			for i, v := range result {
				if v != test.expected[i] {
					t.Fatalf("expected: %v got: %v", test.expected, result)
				}
			}
		})
	}
}

func TestTimePointerCompare(t *testing.T) {

	tests := []struct {
		alpha    *time.Time
		bravo    *time.Time
		expected bool
	}{
		{
			alpha:    TimePointer(42),
			bravo:    TimePointer(73),
			expected: false,
		},
		{
			alpha:    TimePointer(42.64),
			bravo:    TimePointer(73.32),
			expected: false,
		},
		{
			alpha:    TimePointer(56.64),
			bravo:    TimePointer(56.32),
			expected: false,
		},
		{
			alpha:    nil,
			bravo:    TimePointer(56.32),
			expected: false,
		},
		{
			alpha:    TimePointer(56.64),
			bravo:    nil,
			expected: false,
		},
		{
			alpha:    TimePointer(42),
			bravo:    TimePointer(42),
			expected: true,
		},
		{
			alpha:    TimePointer(42.64),
			bravo:    TimePointer(42.64),
			expected: true,
		},
		{
			alpha:    TimePointer(73.32),
			bravo:    TimePointer(73.32),
			expected: true,
		},
		{
			alpha:    nil,
			bravo:    nil,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s - %s", test.alpha, test.bravo), func(t *testing.T) {
			result := TimePointerCompare(test.alpha, test.bravo)
			if result != test.expected {
				t.Fatalf("expected %v got %v", test.expected, result)
			}
		})
	}
}
