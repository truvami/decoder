package common

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"unsafe"

	"reflect"
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

func Decode(payloadHex *string, config *PayloadConfig) (any, error) {
	payloadBytes, err := HexStringToBytes(*payloadHex)
	if err != nil {
		return nil, err
	}

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
					config.Features = append(config.Features, tagConfig.Feature)

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

func insertFieldBytes(fieldValue reflect.Value, length int, transform func(v any) any) (bool, []byte, error) {
	var null bool = false
	var set bool = false
	var bytes []byte
	var err error = nil

	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			null = true
			bytes = make([]byte, length)
		} else {
			fieldValue = fieldValue.Elem()
		}
	}

	if !null {
		switch fieldValue.Type() {
		case reflect.TypeOf(bool(false)):
			set = true
			bytes = BoolToBytes(fieldValue.Bool(), 0)
		case reflect.TypeOf(int8(0)), reflect.TypeOf(int16(0)), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0)):
			value := fieldValue.Int()
			set = value != 0
			bytes = IntToBytes(value, length)
		case reflect.TypeOf(uint8(0)), reflect.TypeOf(uint16(0)), reflect.TypeOf(uint32(0)), reflect.TypeOf(uint64(0)):
			value := fieldValue.Uint()
			set = value != 0
			bytes = UintToBytes(value, length)
		case reflect.TypeOf(float32(0)):
			value := float32(fieldValue.Float())
			set = value != 0
			bytes = Float32ToBytes(value)
		case reflect.TypeOf(float64(0)):
			value := fieldValue.Float()
			set = value != 0
			bytes = Float64ToBytes(value)
		case reflect.TypeOf([]byte{}):
			value := fieldValue.Bytes()
			set = len(value) != 0
			bytes = value
		case reflect.TypeOf(string("")):
			value := fieldValue.String()
			set = len(value) != 0
			bytes, err = HexStringToBytes(value)
		case reflect.TypeOf(time.Time{}):
			value := fieldValue.Interface().(time.Time).Unix()
			set = value != 0
			bytes = IntToBytes(value, int(unsafe.Sizeof(value)))
		case reflect.TypeOf(time.Duration(0)):
			value := fieldValue.Interface().(time.Duration).Nanoseconds()
			set = value != 0
			bytes = IntToBytes(value, int(unsafe.Sizeof(value)))
		default:
			err = fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}
	}

	if !null && transform != nil && err == nil {
		bytes = transform(bytes).([]byte)
	}

	return set, bytes, err
}

func Encode(data any, config PayloadConfig) (string, error) {
	v := reflect.ValueOf(data)

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

	var actualLength int
	for _, field := range config.Fields {
		fieldValue := v.FieldByName(field.Name)

		if !fieldValue.IsValid() {
			return "", fmt.Errorf("field %s not found in data", field.Name)
		}

		set, bytes, err := insertFieldBytes(fieldValue, field.Length, field.Transform)
		if err != nil {
			return "", err
		}

		if set || !field.Optional {
			copy(payload[field.Start:field.Start+field.Length], bytes)
			actualLength += field.Length
		}
	}

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
