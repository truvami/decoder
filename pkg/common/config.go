package common

import "reflect"

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
	Fields          []FieldConfig
	TargetType      reflect.Type
	StatusByteIndex *int // can be nil
}
