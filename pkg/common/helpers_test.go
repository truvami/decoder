package common

import (
	"fmt"
	"math"
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
	Year   uint8   `json:"year" validate:"gte=0,lte=255"`
	Month  uint8   `json:"month" validate:"gte=0,lte=255"`
	Day    uint8   `json:"day" validate:"gte=1,lte=31"`
	Hour   uint8   `json:"hour" validate:"gte=0,lte=23"`
	Minute uint8   `json:"minute" validate:"gte=0,lte=59"`
	Second uint8   `json:"second" validate:"gte=0,lte=59"`
	TS     int64   `json:"ts"`
}

func TestDecode(t *testing.T) {
	config := PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 0, Length: 1},
			{Name: "Lat", Start: 1, Length: 4, Transform: func(v any) any {
				return float64(BytesToInt32(v.([]byte))) / 1000000
			}},
			{Name: "Lon", Start: 5, Length: 4, Transform: func(v any) any {
				return float64(BytesToInt32(v.([]byte))) / 1000000
			}},
			{Name: "Alt", Start: 9, Length: 2, Transform: func(v any) any {
				return float64(BytesToUint16(v.([]byte))) / 10
			}},
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
				Alt:    572.8,
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
			decodedData, err := Decode(StringPtr(test.payload), &test.config)
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

func TestExtractFieldValue(t *testing.T) {
	tests := []struct {
		payload     []byte
		start       int
		length      int
		optional    bool
		hexadecimal bool
		expected    any
		expectedErr string
	}{
		{
			payload:  []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:    0,
			length:   1,
			expected: []byte{0x73},
		},
		{
			payload:  []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:    0,
			length:   5,
			expected: []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
		},
		{
			payload:  []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:    1,
			length:   2,
			expected: []byte{0x6f, 0x72},
		},
		{
			payload:  []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:    2,
			length:   3,
			expected: []byte{0x72, 0x65, 0x6e},
		},
		{
			payload:  []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:    0,
			length:   -1,
			expected: []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
		},
		{
			payload:  []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:    5,
			length:   1,
			optional: true,
			expected: nil,
		},
		{
			payload:     []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:       8,
			length:      2,
			expected:    nil,
			expectedErr: "field out of bounds",
		},
		{
			payload:     []byte{0x73, 0x6f, 0x72, 0x65, 0x6e},
			start:       8,
			length:      -1,
			expected:    nil,
			expectedErr: "field start out of bounds",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v_%v_%v", test.payload, test.start, test.length), func(t *testing.T) {
			result, err := extractFieldValue(test.payload, test.start, test.length, test.optional, test.hexadecimal)
			if err != nil && err.Error() != test.expectedErr {
				t.Fatalf("expected: %s received: %s", test.expectedErr, err.Error())
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("expected: %v received: %v", test.expected, result)
			}
		})
	}
}

func TestConvertFieldValue(t *testing.T) {
	tests := []struct {
		value       any
		fieldType   reflect.Type
		expected    any
		expectedErr string
	}{
		{
			value:     []byte{0x00},
			fieldType: reflect.TypeOf(bool(false)),
			expected:  false,
		},
		{
			value:     []byte{0x01},
			fieldType: reflect.TypeOf(bool(false)),
			expected:  true,
		},
		{
			value:     []byte{0x00},
			fieldType: reflect.TypeOf(int8(0)),
			expected:  int8(0),
		},
		{
			value:     []byte{0xff},
			fieldType: reflect.TypeOf(int8(0)),
			expected:  ^int8(0),
		},
		{
			value:     []byte{0x00, 0x00},
			fieldType: reflect.TypeOf(int16(0)),
			expected:  int16(0),
		},
		{
			value:     []byte{0xff, 0xff},
			fieldType: reflect.TypeOf(int16(0)),
			expected:  ^int16(0),
		},
		{
			value:     []byte{0x00, 0x00, 0x00, 0x00},
			fieldType: reflect.TypeOf(int32(0)),
			expected:  int32(0),
		},
		{
			value:     []byte{0xff, 0xff, 0xff, 0xff},
			fieldType: reflect.TypeOf(int32(0)),
			expected:  ^int32(0),
		},
		{
			value:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			fieldType: reflect.TypeOf(int64(0)),
			expected:  int64(0),
		},
		{
			value:     []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			fieldType: reflect.TypeOf(int64(0)),
			expected:  ^int64(0),
		},
		{
			value:     []byte{0x00},
			fieldType: reflect.TypeOf(uint8(0)),
			expected:  uint8(0),
		},
		{
			value:     []byte{0xff},
			fieldType: reflect.TypeOf(uint8(0)),
			expected:  ^uint8(0),
		},
		{
			value:     []byte{0x00, 0x00},
			fieldType: reflect.TypeOf(uint16(0)),
			expected:  uint16(0),
		},
		{
			value:     []byte{0xff, 0xff},
			fieldType: reflect.TypeOf(uint16(0)),
			expected:  ^uint16(0),
		},
		{
			value:     []byte{0x00, 0x00, 0x00, 0x00},
			fieldType: reflect.TypeOf(uint32(0)),
			expected:  uint32(0),
		},
		{
			value:     []byte{0xff, 0xff, 0xff, 0xff},
			fieldType: reflect.TypeOf(uint32(0)),
			expected:  ^uint32(0),
		},
		{
			value:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			fieldType: reflect.TypeOf(uint64(0)),
			expected:  uint64(0),
		},
		{
			value:     []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			fieldType: reflect.TypeOf(uint64(0)),
			expected:  ^uint64(0),
		},
		{
			value:     "hello world",
			fieldType: reflect.TypeOf(string("")),
			expected:  "hello world",
		},
		{
			value:     "lorem ipsum dolor",
			fieldType: reflect.TypeOf(string("")),
			expected:  "lorem ipsum dolor",
		},
		{
			value:       nil,
			fieldType:   reflect.TypeOf(time.Time{}),
			expected:    nil,
			expectedErr: "unsupported field type: time.Time",
		},
		{
			value:       nil,
			fieldType:   reflect.TypeOf(time.Duration(0)),
			expected:    nil,
			expectedErr: "unsupported field type: time.Duration",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v_%v", test.value, test.fieldType), func(t *testing.T) {
			result, err := convertFieldValue(test.value, test.fieldType, nil)
			if err != nil && err.Error() != test.expectedErr {
				t.Fatalf("expected: %s received: %s", test.expectedErr, err.Error())
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("expected: %v received: %v", test.expected, result)
			}
		})
	}
}

func TestInvalidPayload(t *testing.T) {
	_, err := Decode(StringPtr(""), &PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 0, Length: 1},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field out of bounds")
	}

	_, err = Decode(StringPtr("01"), &PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 0, Length: 2},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field out of bounds")
	}

	_, err = Decode(StringPtr("01"), &PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Moving", Start: 10, Length: 1},
		},
		TargetType: reflect.TypeOf(Port1Payload{}),
	})
	if err == nil {
		t.Fatal("expected field start out of bounds")
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

func TestUnwrapError_Nil(t *testing.T) {
	var err error = nil
	result := UnwrapError(err)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

type multiError struct {
	errs []error
}

func (m multiError) Error() string {
	return "multi error"
}

func (m multiError) Unwrap() []error {
	return m.errs
}

func TestUnwrapError_WithUnwrap(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")
	merr := multiError{errs: []error{err1, err2}}
	result := UnwrapError(merr)
	if len(result) != 2 || result[0] != err1 || result[1] != err2 {
		t.Fatalf("expected [%v %v], got %v", err1, err2, result)
	}
}

func TestUnwrapError_NoUnwrap(t *testing.T) {
	err := fmt.Errorf("plain error")
	result := UnwrapError(err)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

func TestBoolToBytes(t *testing.T) {
	tests := []struct {
		value    bool
		bit      uint8
		expected []byte
	}{
		{value: false, bit: 0, expected: []byte{0x00}},
		{value: true, bit: 0, expected: []byte{0x01}},
		{value: true, bit: 1, expected: []byte{0x02}},
		{value: true, bit: 2, expected: []byte{0x04}},
		{value: true, bit: 7, expected: []byte{0x80}},
		{value: false, bit: 7, expected: []byte{0x00}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("value=%v_bit=%d", test.value, test.bit), func(t *testing.T) {
			result := BoolToBytes(test.value, test.bit)
			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, result)
			}
		})
	}

	t.Run("panic on bit > 7", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic for bit > 7")
			}
		}()
		BoolToBytes(true, 8)
	})
}

func TestBytesToInt64(t *testing.T) {
	tests := []struct {
		input    []byte
		expected int64
	}{
		{input: []byte{0x00}, expected: 0},
		{input: []byte{0x01}, expected: 1},
		{input: []byte{0x7F}, expected: 127},
		{input: []byte{0xFF}, expected: 255},
		{input: []byte{0x01, 0x00}, expected: 256},
		{input: []byte{0x12, 0x34}, expected: 0x1234},
		{input: []byte{0x00, 0x01}, expected: 1},
		{input: []byte{0xFF, 0xFF}, expected: 65535},
		{input: []byte{0x01, 0x00, 0x00}, expected: 65536},
		{input: []byte{0x00, 0x00, 0x01}, expected: 1},
		{input: []byte{0x80, 0x00, 0x00, 0x00}, expected: 0x80000000},
		{input: []byte{0x7F, 0xFF, 0xFF, 0xFF}, expected: 0x7FFFFFFF},
		{input: []byte{}, expected: 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			result := BytesToInt64(test.input)
			if result != test.expected {
				t.Fatalf("expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestFloat64ToBytes(t *testing.T) {
	tests := []struct {
		value    float64
		expected []byte
	}{
		{
			value:    0.0,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			value:    1.0,
			expected: []byte{0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			value:    -1.0,
			expected: []byte{0xbf, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			value:    123.456,
			expected: []byte{0x40, 0x5e, 0xdd, 0x2f, 0x1a, 0x9f, 0xbe, 0x77},
		},
		{
			value: math.NaN(),
			expected: func() []byte {
				// NaN can have multiple bit representations, so just check math.IsNaN after conversion
				return nil
			}(),
		},
		{
			value:    math.Inf(1),
			expected: []byte{0x7f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			value:    math.Inf(-1),
			expected: []byte{0xff, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("value=%v", test.value), func(t *testing.T) {
			result := Float64ToBytes(test.value)
			if test.expected == nil {
				// For NaN, convert back and check math.IsNaN
				bits := uint64(0)
				for _, b := range result {
					bits = (bits << 8) | uint64(b)
				}
				f := math.Float64frombits(bits)
				if !math.IsNaN(f) {
					t.Fatalf("expected NaN, got %v", f)
				}
			} else {
				if !reflect.DeepEqual(result, test.expected) {
					t.Fatalf("expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}

func TestFloat32ToBytes(t *testing.T) {
	tests := []struct {
		value    float32
		expected []byte
	}{
		{
			value:    1.0,
			expected: []byte{0x3f, 0x80, 0x00, 0x00},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("value=%v", test.value), func(t *testing.T) {
			result := Float32ToBytes(test.value)
			if test.expected == nil {
				// For NaN, convert back and check math.IsNaN
				bits := uint32(0)
				for _, b := range result {
					bits = (bits << 8) | uint32(b)
				}
				f := math.Float32frombits(bits)
				if !math.IsNaN(float64(f)) {
					t.Fatalf("expected NaN, got %v", f)
				}
			} else {
				if !reflect.DeepEqual(result, test.expected) {
					t.Fatalf("expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}
