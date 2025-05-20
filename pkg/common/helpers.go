package common

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"unsafe"

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

// Parse decodes the payload based on the provided configuration and populates the target struct
func Parse(payloadHex string, config *PayloadConfig) (any, error) {
	// Convert hex payload to bytes
	payloadBytes, err := HexStringToBytes(payloadHex)
	if err != nil {
		return nil, err
	}

	// Create an instance of the target struct
	targetValue := reflect.New(config.TargetType).Elem()

	errs := []error{}

	if len(config.Tags) != 0 {
		var index uint8 = 3
		var payloadLength uint8 = uint8(len(payloadBytes))
		for index+2 < payloadLength {
			var found bool = false
			var tag uint8 = payloadBytes[index]
			index++
			var length uint8 = payloadBytes[index]
			index++

			for _, tagConfig := range config.Tags {
				if tagConfig.Tag == tag {
					found = true

					value, err := extractFieldValue(payloadBytes, int(index), int(length), false, tagConfig.Hex)
					if err != nil {
						return nil, err
					}

					fieldValue := targetValue.FieldByName(tagConfig.Name)
					if fieldValue.IsValid() && fieldValue.CanSet() && tagConfig.Transform != nil {
						if value == nil && tagConfig.Optional {
							continue
						}

						config.Features = append(config.Features, tagConfig.Feature...)

						// transform value from pointer to value
						if fieldValue.Kind() == reflect.Pointer {
							// if fieldValue is nil set the value to nil
							if value == nil {
								fieldValue.Set(reflect.Zero(fieldValue.Type()))
								continue
							}

							transformed := tagConfig.Transform(value)

							// value is a pointer and not nil, convert to value
							// transform value from pointer to value
							// Create a new pointer of the right type
							ptr := reflect.New(fieldValue.Type().Elem())

							// Set the dereferenced value
							ptr.Elem().Set(reflect.ValueOf(transformed))

							// Set the pointer to the field
							fieldValue.Set(ptr)
							continue
						}
						fieldValue.Set(reflect.ValueOf(tagConfig.Transform(value)))
						continue
					}

					// if fieldValue is nil set the value to nil
					if value == nil && tagConfig.Optional {
						continue
					}

					if !fieldValue.IsValid() || !fieldValue.CanSet() {
						return nil, fmt.Errorf("field %s not found in target struct", tagConfig.Name)
					}

					// pointer field type e.g *uint8 -> t = uint8
					targetTypeNonPointer := fieldValue.Type().Elem()

					// transform value from pointer to value with the right type
					if fieldValue.Kind() == reflect.Pointer {
						// Create a new pointer of the right type
						ptr := reflect.New(targetTypeNonPointer)

						// Set the dereferenced value
						ptr.Elem().Set(reflect.ValueOf(convertFieldToType(value, targetTypeNonPointer)))

						// Set the pointer to the field
						fieldValue.Set(ptr)
						continue
					}

					fieldValue.Set(reflect.ValueOf(convertFieldToType(value, targetTypeNonPointer)))
				}
			}
			if found {
				index += length
			} else {
				return nil, fmt.Errorf("unknown tag %x", tag)
			}
		}

		return targetValue.Interface(), errors.Join(errs...)
	}

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

			if fieldValue.Kind() == reflect.Pointer && field.Optional {
				ptrValue := reflect.New(fieldValue.Type().Elem())
				convertedValue := convertFieldToType(value, ptrValue.Elem().Type())

				ptrValue.Elem().Set(reflect.ValueOf(convertedValue))
				fieldValue.Set(ptrValue)
			} else {
				fieldType := convertFieldToType(value, fieldValue.Type())
				fieldValue.Set(reflect.ValueOf(fieldType))
			}
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

	if len(config.Tags) != 0 {
		var minLength = 3
		if payloadLength < minLength {
			return WrapErrorWithMessage(ErrInvalidPayloadLength, ErrPayloadTooShort, fmt.Sprintf("payload length %d is less than minimum required length %d", payloadLength, minLength))
		} else {
			return nil
		}
	}

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

	var maxLength int
	for _, field := range config.Fields {
		if field.Start+field.Length > maxLength {
			maxLength = field.Start + field.Length
		}
	}
	payload := make([]byte, maxLength)

	// Encode fields into the payload
	var actualLength int
	for _, field := range config.Fields {
		fieldValue := v.FieldByName(field.Name)

		// Check if the field exists
		if !fieldValue.IsValid() {
			return "", fmt.Errorf("field %s not found in data", field.Name)
		}

		var unset bool = false
		var fieldBytes []byte
		switch fieldValue.Type() {
		case reflect.TypeOf(bool(false)):
			fieldBytes = BoolToBytes(fieldValue.Bool(), 0)
		case reflect.TypeOf([]byte{}):
			value := fieldValue.Bytes()
			unset = len(value) == 0
			fieldBytes = value
		case reflect.TypeOf(string("")):
			value, err := HexStringToBytes(fieldValue.String())
			if err != nil {
				panic(err)
			}
			unset = len(value) == 0
			fieldBytes = value
		case reflect.TypeOf(uint(0)), reflect.TypeOf(uint8(0)), reflect.TypeOf(uint16(0)), reflect.TypeOf(uint32(0)), reflect.TypeOf(uint64(0)):
			value := fieldValue.Uint()
			unset = value == 0
			fieldBytes = UintToBytes(value, field.Length)
		case reflect.TypeOf(int(0)), reflect.TypeOf(int8(0)), reflect.TypeOf(int16(0)), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0)):
			value := fieldValue.Int()
			unset = value == 0
			fieldBytes = IntToBytes(value, field.Length)
		case reflect.TypeOf(float32(0)), reflect.TypeOf(float64(0)):
			fieldBytes = FloatToBytes(fieldValue.Float(), int(fieldValue.Type().Size()))
		case reflect.TypeOf(time.Duration(0)):
			duration := fieldValue.Interface().(time.Duration).Nanoseconds()
			fieldBytes = IntToBytes(duration, int(unsafe.Sizeof(duration)))
		case reflect.TypeOf(time.Time{}):
			timestamp := fieldValue.Interface().(time.Time).Unix()
			fieldBytes = IntToBytes(timestamp, int(unsafe.Sizeof(timestamp)))
		default:
			return "", fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}

		// Apply the transform function if provided
		if field.Transform != nil {
			fieldBytes = field.Transform(fieldBytes).([]byte)
		}

		if !unset || !field.Optional {
			copy(payload[field.Start:field.Start+field.Length], fieldBytes)
			actualLength += field.Length
		}
	}

	// Convert the payload to a hexadecimal string
	return hex.EncodeToString(payload[0:actualLength]), nil
}

func BoolToBytes(value bool, bit uint8) []byte {
	if bit > 7 {
		panic("bit must be in range 0 to 7")
	}
	var b byte
	if value {
		b = 1 << bit
	}
	return []byte{b}
}

func BytesToBool(bytes []byte) bool {
	return bytes[0] == 0x01
}

func UintToBytes(value uint64, length int) []byte {
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = byte(value & 0xff)
		value >>= 8
	}
	return buf
}

func IntToBytes(value int64, length int) []byte {
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = byte(value & 0xff)
		value >>= 8
	}
	return buf
}

func BytesToInt64(bytes []byte) int64 {
	var value int64 = 0
	for i := range bytes {
		value <<= 8
		value |= int64(bytes[i])
	}
	return value
}

func FloatToBytes(value float64, length int) []byte {
	bits := math.Float64bits(value)
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = byte(bits & 0xff)
		bits >>= 8
	}
	return buf
}

func BytesToFloat32(bytes []byte) float32 {
	var bits uint32
	for i := range bytes {
		bits <<= 8
		bits |= uint32(bytes[i])
	}
	return math.Float32frombits(bits)
}

func BytesToFloat64(bytes []byte) float64 {
	var bits uint64
	for i := range bytes {
		bits <<= 8
		bits |= uint64(bytes[i])
	}
	return math.Float64frombits(bits)
}

func Uint8Ptr(value uint8) *uint8 {
	return &value
}

func Uint16Ptr(value uint16) *uint16 {
	return &value
}

func Uint32Ptr(value uint32) *uint32 {
	return &value
}

func Int8Ptr(value int8) *int8 {
	return &value
}

func BoolPtr(value bool) *bool {
	return &value
}

func StringPtr(value string) *string {
	return &value
}

func Float32Ptr(value float32) *float32 {
	return &value
}

func Float64Ptr(value float64) *float64 {
	return &value
}

func DurationPtr(duration time.Duration) *time.Duration {
	return &duration
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

func TransformPointerToValue(ptr interface{}) interface{} {
	val := reflect.ValueOf(ptr)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		return val.Elem().Interface()
	}
	panic(fmt.Sprintf("expected a pointer with not nil value, got %T", ptr))
}
