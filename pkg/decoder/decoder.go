package decoder

import "reflect"

// FieldConfig defines the structure of a single field in the payload
type FieldConfig struct {
	Name      string
	Start     int
	Length    int
	Transform func(interface{}) interface{}
	Optional  bool
	Hex       bool
}

// PayloadConfig defines the overall structure of the payload, including the target struct type
type PayloadConfig struct {
	Fields     []FieldConfig
	TargetType reflect.Type
}

type Decoder interface {
	Decode(string, int16, string) (interface{}, interface{}, error)
}
