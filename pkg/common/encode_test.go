package common

import (
	"testing"
)

type TestStruct struct {
	Uint8Field   uint8
	Uint16Field  uint16
	Uint32Field  uint32
	Int8Field    int8
	Int16Field   int16
	Int32Field   int32
	ByteSlice    []byte
	MissingField string // This field won't be in the config
}

func TestEncodePayload(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		config         PayloadConfig
		expectedOutput string
		expectError    bool
	}{
		{
			name: "Basic encoding test",
			data: TestStruct{
				Uint8Field:  0x12,
				Uint16Field: 0x3456,
				Uint32Field: 0x789ABCDE,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Uint8Field", Start: 0, Length: 1},
					{Name: "Uint16Field", Start: 1, Length: 2},
					{Name: "Uint32Field", Start: 3, Length: 4},
				},
			},
			expectedOutput: "123456789abcde",
			expectError:    false,
		},
		{
			name: "Test with int fields",
			data: TestStruct{
				Int8Field:  -16,
				Int16Field: -1000,
				Int32Field: -70000,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Int8Field", Start: 0, Length: 1},
					{Name: "Int16Field", Start: 1, Length: 2},
					{Name: "Int32Field", Start: 3, Length: 4},
				},
			},
			expectedOutput: "f0fc18fffeee90",
			expectError:    false,
		},
		{
			name: "Test with byte slice",
			data: TestStruct{
				ByteSlice: []byte{0xAA, 0xBB, 0xCC},
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "ByteSlice", Start: 0, Length: 5}, // Longer than the slice
				},
			},
			expectedOutput: "aabbcc0000",
			expectError:    false,
		},
		{
			name: "Test with optional missing field",
			data: TestStruct{
				Uint8Field: 0x12,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Uint8Field", Start: 0, Length: 1},
					{Name: "NonExistentField", Start: 1, Length: 1, Optional: true},
				},
			},
			expectedOutput: "1200",
			expectError:    false,
		},
		{
			name: "Test with transform function",
			data: TestStruct{
				ByteSlice: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE},
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{
						Name:      "ByteSlice",
						Start:     0,
						Length:    3,
						Transform: func(v interface{}) interface{} { return v.([]byte)[:2] },
					},
				},
			},
			expectedOutput: "aabb00",
			expectError:    false,
		},
		{
			name:           "Test with non-struct data",
			data:           "not a struct",
			config:         PayloadConfig{},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with missing required field",
			data: TestStruct{},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "NonExistentField", Start: 0, Length: 1},
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with uint8 field length mismatch",
			data: TestStruct{
				Uint8Field: 0x12,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Uint8Field", Start: 0, Length: 2}, // Should be 1
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with uint16 field length mismatch",
			data: TestStruct{
				Uint16Field: 0x3456,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Uint16Field", Start: 0, Length: 3}, // Should be 2
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with uint32 field length mismatch",
			data: TestStruct{
				Uint32Field: 0x789ABCDE,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Uint32Field", Start: 0, Length: 3}, // Should be 4
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with int8 field length mismatch",
			data: TestStruct{
				Int8Field: -16,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Int8Field", Start: 0, Length: 2}, // Should be 1
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with int16 field length mismatch",
			data: TestStruct{
				Int16Field: -1000,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Int16Field", Start: 0, Length: 3}, // Should be 2
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with int32 field length mismatch",
			data: TestStruct{
				Int32Field: -70000,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "Int32Field", Start: 0, Length: 3}, // Should be 4
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with byte slice too long",
			data: TestStruct{
				ByteSlice: []byte{0xAA, 0xBB, 0xCC},
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "ByteSlice", Start: 0, Length: 2}, // Shorter than the slice
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with unsupported slice type",
			data: struct {
				IntSlice []int
			}{
				IntSlice: []int{1, 2, 3},
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "IntSlice", Start: 0, Length: 3},
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Test with unsupported field type",
			data: struct {
				FloatField float64
			}{
				FloatField: 123.456,
			},
			config: PayloadConfig{
				Fields: []FieldConfig{
					{Name: "FloatField", Start: 0, Length: 8},
				},
			},
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := EncodePayload(tt.data, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if output != tt.expectedOutput {
					t.Errorf("Expected output %s, got %s", tt.expectedOutput, output)
				}
			}
		})
	}
}

func TestEncodePayloadWithNilTransform(t *testing.T) {
	data := TestStruct{
		Uint8Field: 0x12,
	}

	config := PayloadConfig{
		Fields: []FieldConfig{
			{Name: "Uint8Field", Start: 0, Length: 1, Transform: nil},
		},
	}

	output, err := EncodePayload(data, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if output != "12" {
		t.Errorf("Expected output 12, got %s", output)
	}
}

func TestEncodePayloadWithEmptyConfig(t *testing.T) {
	data := TestStruct{}
	config := PayloadConfig{
		Fields: []FieldConfig{},
	}

	output, err := EncodePayload(data, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if output != "" {
		t.Errorf("Expected empty output, got %s", output)
	}
}
