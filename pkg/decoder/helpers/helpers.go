package helpers

import (
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

func hexStringToBytes(hexString string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func convertFieldToType(value interface{}, fieldType reflect.Kind) interface{} {
	switch fieldType {
	case reflect.Int:
		return int(value.(int))
	case reflect.Int8:
		return int8(value.(int))
	case reflect.Int16:
		return int16(value.(int))
	case reflect.Int32:
		return int32(value.(int))
	case reflect.Int64:
		return int64(value.(int))
	case reflect.Uint:
		return uint(value.(int))
	case reflect.Uint8:
		return uint8(value.(int))
	case reflect.Uint16:
		return uint16(value.(int))
	case reflect.Uint32:
		return uint32(value.(int))
	case reflect.Uint64:
		return uint64(value.(int))
	case reflect.Float32:
		return float32(value.(int))
	case reflect.Float64:
		return float64(value.(int))
	case reflect.String:
		return fmt.Sprintf("%v", value)
	case reflect.Bool:
		return value.(int)&0x01 == 1
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}).Kind() {
			return ParseTimestamp(value.(int))
		}
		fallthrough
	default:
		panic(fmt.Sprintf("unsupported field type: %v", fieldType))
	}
}

func extractFieldValue(payloadBytes []byte, start int, length int, optional bool) (interface{}, error) {
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
	var value interface{}
	switch length {
	case 1, 2, 4, 8:
		value = 0
		for i := 0; i < length; i++ {
			value = (value.(int) << 8) | int(payloadBytes[start+i])
		}
	default:
		// For lengths greater than 8, return the slice as a hex string
		value = hex.EncodeToString(payloadBytes[start : start+length])
	}

	return value, nil
}

// DecodeLoRaWANPayload decodes the payload based on the provided configuration and populates the target struct
func Parse(payloadHex string, config decoder.PayloadConfig) (interface{}, error) {
	// Convert hex payload to bytes
	payloadBytes, err := hexStringToBytes(payloadHex)
	if err != nil {
		return nil, err
	}

	// Create an instance of the target struct
	targetValue := reflect.New(config.TargetType).Elem()

	// Iterate over the fields in the config and extract their values
	for _, field := range config.Fields {
		start := field.Start
		length := field.Length
		optional := field.Optional

		// Extract the field value from the payload
		value, err := extractFieldValue(payloadBytes, start, length, optional)
		if err != nil {
			return nil, err
		}

		// Convert value to appropriate type and set it in the target struct
		fieldValue := targetValue.FieldByName(field.Name)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			if value == nil && optional {
				continue
			}

			log.Printf("field: %v", field.Name)
			log.Printf("value: %v", value)
			log.Printf("got: %T", value)
			log.Printf("expect: %v", fieldValue.Type().Kind())

			fieldType := convertFieldToType(value, fieldValue.Type().Kind())
			fieldValue.Set(reflect.ValueOf(fieldType))
		}

		// Apply the transform function if provided
		if field.Transform != nil {
			transformedValue := field.Transform(value)
			fieldValue.Set(reflect.ValueOf(transformedValue))
		}
	}

	return targetValue.Interface(), nil
}

func ParseDateTime(input string) (time.Time, error) {
	// The last 4 characters are always MMSS (minute and second)
	minuteSecond := input[len(input)-4:]

	// The next 6 characters are YYMM (year and month)
	yearMonth := input[len(input)-10 : len(input)-4]

	// The rest is the day (D or DD) and hour (H or HH)
	dayHour := input[:len(input)-10]

	// Identify the day and hour, which may vary in length
	var day, hour string

	// Check if the day is 2 digits or 1 digit
	if len(dayHour) == 3 {
		// Single-digit day (D), single-digit hour (H)
		day = "0" + string(dayHour[0])   // Pad the day
		hour = "0" + string(dayHour[1:]) // Pad the hour
	} else if len(dayHour) == 4 {
		// Two-digit day (DD), single-digit hour (H)
		day = string(dayHour[0:2])       // Day is two digits
		hour = "0" + string(dayHour[2:]) // Pad the hour
	} else if len(dayHour) == 5 {
		// Two-digit day (DD), two-digit hour (HH)
		day = string(dayHour[0:2])  // Day is two digits
		hour = string(dayHour[2:4]) // Hour is two digits
	}

	// Now we have day, hour, yearMonth, and minuteSecond
	// Combine everything into a full datetime string
	parsedInput := fmt.Sprintf("%s%s%s%s", yearMonth, day, hour, minuteSecond)

	// Define the layout in Go's time format for parsing
	layout := "060102150405" // YYMMDDHHMMSS (Go's reference date layout)

	// Parse the string into a time.Time object
	parsedTime, err := time.Parse(layout, parsedInput)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
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
