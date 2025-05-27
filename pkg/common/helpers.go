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

func extractFieldValue(payloadBytes []byte, start int, length int, optional bool, hexadecimal bool) (any, error) {
	if length == -1 {
		if start >= len(payloadBytes) && !optional {
			return nil, fmt.Errorf("field start out of bounds")
		}
		// Dynamic length: read until the end of the payload
		length = len(payloadBytes) - start
		if length == 0 {
			return nil, nil
		}
	} else if start+length > len(payloadBytes) {
		if optional {
			return nil, nil
		}
		return nil, fmt.Errorf("field out of bounds")
	}

	// Extract the field value based on its length
	var value any = payloadBytes[start : start+length]
	if hexadecimal {
		value = hex.EncodeToString(value.([]byte))
	}

	return value, nil
}

func convertFieldValue(rawValue any, fieldType reflect.Type, transform func(v any) any) (any, error) {
	var ptr bool = false
	var value any = nil
	var err error = nil

	if fieldType.Kind() == reflect.Ptr {
		ptr = true
		fieldType = fieldType.Elem()
	}

	if transform != nil {
		value = transform(rawValue)
	} else {
		switch fieldType {
		case reflect.TypeOf(bool(false)):
			value = rawValue.([]byte)[0]&0x01 == 1
		case reflect.TypeOf(int8(0)):
			value = BytesToInt8(rawValue.([]byte))
		case reflect.TypeOf(int16(0)):
			value = BytesToInt16(rawValue.([]byte))
		case reflect.TypeOf(int32(0)):
			value = BytesToInt32(rawValue.([]byte))
		case reflect.TypeOf(int64(0)):
			value = BytesToInt64(rawValue.([]byte))
		case reflect.TypeOf(uint8(0)):
			value = BytesToUint8(rawValue.([]byte))
		case reflect.TypeOf(uint16(0)):
			value = BytesToUint16(rawValue.([]byte))
		case reflect.TypeOf(uint32(0)):
			value = BytesToUint32(rawValue.([]byte))
		case reflect.TypeOf(uint64(0)):
			value = BytesToUint64(rawValue.([]byte))
		case reflect.TypeOf(string("")):
			value = fmt.Sprintf("%s", rawValue)
		default:
			err = fmt.Errorf("unsupported field type: %v", fieldType)
		}
	}

	if ptr && value != nil && err == nil {
		fieldValue := reflect.New(reflect.TypeOf(value))
		fieldValue.Elem().Set(reflect.ValueOf(value))
		value = fieldValue.Interface()
	}

	return value, err
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

// Decode decodes the payload based on the provided configuration and populates the target struct
func Decode(payloadHex *string, config *PayloadConfig) (any, error) {
	// Convert hex payload to bytes
	payloadBytes, err := HexStringToBytes(*payloadHex)
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
					if fieldValue.IsValid() && fieldValue.CanSet() {
						if value == nil && tagConfig.Optional {
							continue
						}

						convertedValue, err := convertFieldValue(value, fieldValue.Type(), tagConfig.Transform)
						if err != nil {
							return nil, err
						}

						if convertedValue != nil {
							fieldValue.Set(reflect.ValueOf(convertedValue))
						}
					}

					fieldName, ok := targetValue.Type().FieldByName(tagConfig.Name)
					if ok {
						err := validateFieldValue(fieldName, fieldValue)
						if err != nil {
							errs = append(errs, fmt.Errorf("%w for %s %v", ErrValidationFailed, fieldName.Name, fieldValue))
						}
					}
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

	for _, field := range config.Fields {
		value, err := extractFieldValue(payloadBytes, field.Start, field.Length, field.Optional, field.Hex)
		if err != nil {
			return nil, err
		}

		fieldValue := targetValue.FieldByName(field.Name)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			if value == nil && field.Optional {
				continue
			}

			convertedValue, err := convertFieldValue(value, fieldValue.Type(), field.Transform)
			if err != nil {
				return nil, err
			}

			if convertedValue != nil {
				fieldValue.Set(reflect.ValueOf(convertedValue))
			}
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
		if field.Length == -1 {
			maxLength = 50
		} else {
			maxLength = field.Start + field.Length
		}
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

		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				unset = true
				fieldBytes = make([]byte, field.Length)
			} else {
				fieldValue = fieldValue.Elem()
			}
		}

		if !unset {
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
			case reflect.TypeOf(float32(0)):
				fieldBytes = Float32ToBytes(float32(fieldValue.Float()))
			case reflect.TypeOf(float64(0)):
				fieldBytes = Float64ToBytes(fieldValue.Float())
			case reflect.TypeOf(time.Duration(0)):
				duration := fieldValue.Interface().(time.Duration).Nanoseconds()
				fieldBytes = IntToBytes(duration, int(unsafe.Sizeof(duration)))
			case reflect.TypeOf(time.Time{}):
				timestamp := fieldValue.Interface().(time.Time).Unix()
				fieldBytes = IntToBytes(timestamp, int(unsafe.Sizeof(timestamp)))
			default:
				return "", fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
			}
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

func Float32ToBytes(value float32) []byte {
	bits := math.Float32bits(value)
	length := int(reflect.TypeOf(float32(0)).Size())
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = byte(bits & 0xff)
		bits >>= 8
	}
	return buf
}

func Float64ToBytes(value float64) []byte {
	bits := math.Float64bits(value)
	length := int(reflect.TypeOf(float64(0)).Size())
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
