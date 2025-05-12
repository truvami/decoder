package common

import (
	"encoding/hex"
	"errors"
	"fmt"

	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

func HexStringToBytes(hexString string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func convertFieldToType(value any, fieldType reflect.Type) any {
	switch fieldType {
	case reflect.TypeOf(int(0)):
		return int(value.(int))
	case reflect.TypeOf(int8(0)):
		return int8(value.(int))
	case reflect.TypeOf(int16(0)):
		return int16(value.(int))
	case reflect.TypeOf(int32(0)):
		return int32(value.(int))
	case reflect.TypeOf(int64(0)):
		return int64(value.(int))
	case reflect.TypeOf(uint(0)):
		return uint(value.(int))
	case reflect.TypeOf(uint8(0)):
		return uint8(value.(int))
	case reflect.TypeOf(uint16(0)):
		return uint16(value.(int))
	case reflect.TypeOf(uint32(0)):
		return uint32(value.(int))
	case reflect.TypeOf(uint64(0)):
		return uint64(value.(int))
	case reflect.TypeOf(float32(0)):
		return float32(value.(int))
	case reflect.TypeOf(float64(0)):
		return float64(value.(int))
	case reflect.TypeOf(string("")):
		return fmt.Sprintf("%v", value)
	case reflect.TypeOf(bool(false)):
		return value.(int)&0x01 == 1
	case reflect.TypeOf(time.Duration(0)):
		return time.Duration(value.(int))
	case reflect.TypeOf(time.Time{}):
		return ParseTimestamp(value.(int))
	default:
		panic(fmt.Sprintf("unsupported field type: %v", fieldType))
	}
}

func extractFieldValue(payloadBytes []byte, start int, length int, optional bool, hexadecimal bool) (any, error) {
	if length == -1 {
		if start >= len(payloadBytes) {
			return nil, fmt.Errorf("field start out of bounds")
		}
		// Dynamic length: read until the end of the payload
		length = len(payloadBytes) - start
	} else if start+length > len(payloadBytes) {
		if optional {
			return nil, nil
		}
		return nil, fmt.Errorf("field out of bounds")
	}

	// Extract the field value based on its length
	var value any
	if hexadecimal {
		value = hex.EncodeToString(payloadBytes[start : start+length])
	} else {
		value = 0
		for i := 0; i < length; i++ {
			value = (value.(int) << 8) | int(payloadBytes[start+i])
		}
	}

	return value, nil
}

func validateFieldValue(field reflect.StructField, fieldValue reflect.Value) error {
	structType := reflect.StructOf([]reflect.StructField{field})

	structValue := reflect.New(structType).Elem()
	structValue.FieldByName(field.Name).Set(fieldValue)

	return validator.New().Struct(structValue.Interface())
}

var ErrValidationFailed = errors.New("validation failed")

func UnwrapError(err error) []error {
	var errs []error = []error{}
	if err, ok := err.(interface{ Unwrap() []error }); ok {
		errs = append(errs, err.Unwrap()...)
	}
	return errs
}

// DecodeLoRaWANPayload decodes the payload based on the provided configuration and populates the target struct
func Parse(payloadHex string, config *PayloadConfig) (any, error) {
	// Convert hex payload to bytes
	payloadBytes, err := HexStringToBytes(payloadHex)
	if err != nil {
		return nil, err
	}

	// Create an instance of the target struct
	targetValue := reflect.New(config.TargetType).Elem()

	errs := []error{}

	// Iterate over the fields in the config and extract their values
	for _, field := range config.Fields {
		start := field.Start
		length := field.Length
		optional := field.Optional
		hex := field.Hex

		// Extract the field value from the payload
		value, err := extractFieldValue(payloadBytes, start, length, optional, hex)
		if err != nil {
			return nil, err
		}

		// Convert value to appropriate type and set it in the target struct
		fieldValue := targetValue.FieldByName(field.Name)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			if value == nil && optional {
				continue
			}

			fieldType := convertFieldToType(value, fieldValue.Type())
			fieldValue.Set(reflect.ValueOf(fieldType))
		}

		// Apply the transform function if provided
		if field.Transform != nil {
			transformedValue := field.Transform(value)
			fieldValue.Set(reflect.ValueOf(transformedValue))
		}

		fieldName, ok := targetValue.Type().FieldByName(field.Name)
		if ok {
			err := validateFieldValue(fieldName, fieldValue)
			if err != nil {
				errs = append(errs, fmt.Errorf("%w for %s %v", ErrValidationFailed, fieldName.Name, fieldValue))
			}
		}
	}

	return targetValue.Interface(), errors.Join(errs...)
}

func ParseTimestamp(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0).UTC()
}

// UintToBinaryArray converts a uint64 value to a binary array of specified length.
// The value parameter represents the uint64 value to be converted.
// The length parameter specifies the length of the resulting binary array.
// The function returns a byte slice representing the binary array.
func UintToBinaryArray(value uint64, length int) []byte {
	binaryArray := make([]byte, length)
	for i := 0; i < length; i++ {
		binaryArray[length-1-i] = byte((value >> uint(i)) & 0x01)
	}
	return binaryArray
}

func HexNullPad(payload *string, config *PayloadConfig) string {
	var requiredBits = 0
	for _, field := range config.Fields {
		if !field.Optional {
			requiredBits = (field.Start + field.Length) * 8
		}
	}
	var providedBits = len(*payload) * 4

	if providedBits < requiredBits {
		var paddingBits = (requiredBits - providedBits) / 4
		*payload = strings.Repeat("0", paddingBits) + *payload
	}
	return *payload
}

func ValidateLength(payload *string, config *PayloadConfig) error {
	var payloadLength = len(*payload) / 2

	var minLength = 0
	for _, field := range config.Fields {
		if !field.Optional {
			minLength = field.Start + field.Length
		}
	}

	var maxLength = 0
	for _, field := range config.Fields {
		maxLength = field.Start + field.Length
	}

	if payloadLength < minLength {
		return WrapErrorWithMessage(ErrInvalidPayloadLength, ErrPayloadTooShort, fmt.Sprintf("payload length %d is less than minimum required length %d", payloadLength, minLength))
	}

	if payloadLength > maxLength {
		return WrapErrorWithMessage(ErrInvalidPayloadLength, ErrPayloadTooLong, fmt.Sprintf("payload length %d is greater than maximum allowed length %d", payloadLength, maxLength))
	}

	return nil
}

func Encode(data any, config PayloadConfig) (string, error) {
	v := reflect.ValueOf(data)

	// Validate input data is a struct
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("data must be a struct")
	}

	// Determine total payload length
	var length int
	for _, field := range config.Fields {
		if field.Start+field.Length > length {
			length = field.Start + field.Length
		}
	}
	payload := make([]byte, length)

	// Encode fields into the payload
	for _, field := range config.Fields {
		fieldValue := v.FieldByName(field.Name)

		// Check if the field exists
		if !fieldValue.IsValid() {
			return "", fmt.Errorf("field %s not found in data", field.Name)
		}

		// Convert the value to bytes
		var fieldBytes []byte
		switch fieldValue.Kind() {
		case reflect.Slice:
			fieldBytes = fieldValue.Bytes()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldBytes = uintToBytes(fieldValue.Uint(), field.Length)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldBytes = intToBytes(fieldValue.Int(), field.Length)
		default:
			return "", fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}

		// Apply the transform function if provided
		if field.Transform != nil {
			fieldBytes = field.Transform(fieldBytes).([]byte)
		}

		// Copy the bytes into the payload at the correct position
		copy(payload[field.Start:field.Start+field.Length], fieldBytes)
	}

	// Convert the payload to a hexadecimal string
	return hex.EncodeToString(payload), nil
}

// intToBytes converts an integer value to a byte slice
func intToBytes(value int64, length int) []byte {
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = byte(value & 0xFF)
		value >>= 8
	}
	return buf
}

// uintToBytes converts an unsigned integer value to a byte slice
func uintToBytes(value uint64, length int) []byte {
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = byte(value & 0xFF)
		value >>= 8
	}
	return buf
}

func TimePointer(timestamp float64) *time.Time {
	seconds := int64(timestamp)
	nanoseconds := int64((timestamp - float64(seconds)) * 1e9)
	time := time.Unix(seconds, nanoseconds)
	return &time
}

func TimePointerCompare(alpha *time.Time, bravo *time.Time) bool {
	if alpha == nil && bravo == nil {
		return true
	}
	if alpha == nil || bravo == nil {
		return false
	}
	return alpha.Equal(*bravo)
}
