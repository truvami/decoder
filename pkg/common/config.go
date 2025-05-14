package common

import (
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
)

type TagConfig struct {
	Name      string
	Tag       uint8
	Transform func(any) any
}

type FieldConfig struct {
	Name      string
	Start     int
	Length    int
	Transform func(any) any
	Optional  bool
	Hex       bool
}

// PayloadConfig defines the overall structure of the payload, including the target struct type
type PayloadConfig struct {
	Tags       []TagConfig
	Fields     []FieldConfig
	TargetType reflect.Type
	Features   []decoder.Feature
}
