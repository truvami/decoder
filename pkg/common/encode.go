package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
)

// EncodePayload encodes a struct into a byte array according to the provided payload configuration
func EncodePayload(data interface{}, config PayloadConfig) (string, error) {
	// Validate that data is a struct
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("data must be a struct")
	}

	// Calculate the total length of the payload
	totalLength := 0
	for _, field := range config.Fields {
		totalLength += field.Length
	}

	// Create a byte array to hold the encoded payload
	payload := make([]byte, totalLength)

	// Encode each field
	for _, field := range config.Fields {
		// Get the field value from the struct
		fieldValue := v.FieldByName(field.Name)
		if !fieldValue.IsValid() {
			if field.Optional {
				// Skip optional fields that don't exist
				continue
			}
			return "", fmt.Errorf("field %s not found in struct", field.Name)
		}

		// Apply any transformation function
		if field.Transform != nil {
			fieldValue = reflect.ValueOf(field.Transform(fieldValue.Interface()))
		}

		// Encode the field value into the payload
		switch fieldValue.Kind() {
		case reflect.Uint8:
			if field.Length != 1 {
				return "", fmt.Errorf("field %s is uint8 but length is %d", field.Name, field.Length)
			}
			payload[field.Start] = uint8(fieldValue.Uint())
		case reflect.Uint16:
			if field.Length != 2 {
				return "", fmt.Errorf("field %s is uint16 but length is %d", field.Name, field.Length)
			}
			binary.BigEndian.PutUint16(payload[field.Start:field.Start+field.Length], uint16(fieldValue.Uint()))
		case reflect.Uint32:
			if field.Length != 4 {
				return "", fmt.Errorf("field %s is uint32 but length is %d", field.Name, field.Length)
			}
			binary.BigEndian.PutUint32(payload[field.Start:field.Start+field.Length], uint32(fieldValue.Uint()))
		case reflect.Int8:
			if field.Length != 1 {
				return "", fmt.Errorf("field %s is int8 but length is %d", field.Name, field.Length)
			}
			payload[field.Start] = uint8(fieldValue.Int())
		case reflect.Int16:
			if field.Length != 2 {
				return "", fmt.Errorf("field %s is int16 but length is %d", field.Name, field.Length)
			}
			binary.BigEndian.PutUint16(payload[field.Start:field.Start+field.Length], uint16(fieldValue.Int()))
		case reflect.Int32:
			if field.Length != 4 {
				return "", fmt.Errorf("field %s is int32 but length is %d", field.Name, field.Length)
			}
			binary.BigEndian.PutUint32(payload[field.Start:field.Start+field.Length], uint32(fieldValue.Int()))
		case reflect.Slice:
			// Handle byte slices ([]byte)
			if fieldValue.Type().Elem().Kind() == reflect.Uint8 {
				bytes := fieldValue.Bytes()
				if len(bytes) > field.Length {
					return "", fmt.Errorf("field %s is too long: %d > %d", field.Name, len(bytes), field.Length)
				}
				copy(payload[field.Start:field.Start+len(bytes)], bytes)
			} else {
				return "", fmt.Errorf("unsupported slice type for field %s", field.Name)
			}
		default:
			return "", fmt.Errorf("unsupported type for field %s: %s", field.Name, fieldValue.Kind())
		}
	}

	// Convert the payload to a hexadecimal string
	return hex.EncodeToString(payload), nil
}
