package tagsl

import (
	"fmt"
	"reflect"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/helpers"
)

type TagSLv1Decoder struct{}

func NewTagSLv1Decoder() decoder.Decoder {
	return TagSLv1Decoder{}
}

func (t TagSLv1Decoder) GetConfig(port int16) (decoder.PayloadConfig, error) {
	switch port {
	case 1:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
				{Name: "Lat", Start: 1, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Lon", Start: 5, Length: 4, Transform: func(v interface{}) interface{} {
					return float64(v.(int)) / 1000000
				}},
				{Name: "Alt", Start: 9, Length: 2},
				{Name: "Year", Start: 11, Length: 1},
				{Name: "Month", Start: 12, Length: 1},
				{Name: "Day", Start: 13, Length: 1},
				{Name: "Hour", Start: 14, Length: 1},
				{Name: "Minute", Start: 15, Length: 1},
				{Name: "Second", Start: 16, Length: 1},
			},
			TargetType: reflect.TypeOf(GNSSPayload{}),
		}, nil
	case 3:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "ScanPointer", Start: 0, Length: 2},
				{Name: "TotalMessages", Start: 2, Length: 1},
				{Name: "CurrentMessage", Start: 3, Length: 1},
				{Name: "Mac1", Start: 4, Length: 6, Optional: true},
				{Name: "Rssi1", Start: 10, Length: 1, Optional: true},
				{Name: "Mac2", Start: 11, Length: 6, Optional: true},
				{Name: "Rssi2", Start: 17, Length: 1, Optional: true},
				{Name: "Mac3", Start: 18, Length: 6, Optional: true},
				{Name: "Rssi3", Start: 24, Length: 1, Optional: true},
				{Name: "Mac4", Start: 25, Length: 6, Optional: true},
				{Name: "Rssi4", Start: 31, Length: 1, Optional: true},
				{Name: "Mac5", Start: 32, Length: 6, Optional: true},
				{Name: "Rssi5", Start: 38, Length: 1, Optional: true},
				{Name: "Mac6", Start: 39, Length: 6, Optional: true},
				{Name: "Rssi6", Start: 45, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(BlePayload{}),
		}, nil
	case 5:
		return decoder.PayloadConfig{
			Fields: []decoder.FieldConfig{
				{Name: "Moving", Start: 0, Length: 1},
				{Name: "Mac1", Start: 1, Length: 6, Optional: true},
				{Name: "Rssi1", Start: 7, Length: 1, Optional: true},
				{Name: "Mac2", Start: 8, Length: 6, Optional: true},
				{Name: "Rssi2", Start: 14, Length: 1, Optional: true},
				{Name: "Mac3", Start: 15, Length: 6, Optional: true},
				{Name: "Rssi3", Start: 21, Length: 1, Optional: true},
				{Name: "Mac4", Start: 22, Length: 6, Optional: true},
				{Name: "Rssi4", Start: 28, Length: 1, Optional: true},
				{Name: "Mac5", Start: 29, Length: 6, Optional: true},
				{Name: "Rssi5", Start: 35, Length: 1, Optional: true},
				{Name: "Mac6", Start: 36, Length: 6, Optional: true},
				{Name: "Rssi6", Start: 42, Length: 1, Optional: true},
				{Name: "Mac7", Start: 43, Length: 6, Optional: true},
				{Name: "Rssi7", Start: 49, Length: 1, Optional: true},
			},
			TargetType: reflect.TypeOf(WifiPayload{}),
		}, nil
	}

	return decoder.PayloadConfig{}, fmt.Errorf("port %v not supported", port)
}

func (t TagSLv1Decoder) Decode(data string, port int16, devEui string) (interface{}, error) {
	config, err := t.GetConfig(port)
	if err != nil {
		return nil, err
	}

	decodedData, err := helpers.Parse(data, config)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}

type GNSSPayload struct {
	Moving bool    `json:"moving"`
	Lat    float64 `json:"gps_lat"`
	Lon    float64 `json:"gps_lon"`
	Alt    float64 `json:"gps_alt"`
	Year   int     `json:"year"`
	Month  int     `json:"month"`
	Day    int     `json:"day"`
	Hour   int     `json:"hour"`
	Minute int     `json:"minute"`
	Second int     `json:"second"`
	TS     int64   `json:"ts"`
}

type AccessPoint struct {
	MAC  string  `json:"mac"`
	Rssi float32 `json:"rssi"`
}

type WifiPayload struct {
	Moving bool   `json:"moving"`
	Mac1   string `json:"mac1"`
	Rssi1  int8   `json:"rssi1"`
	Mac2   string `json:"mac2"`
	Rssi2  int8   `json:"rssi2"`
	Mac3   string `json:"mac3"`
	Rssi3  int8   `json:"rssi3"`
	Mac4   string `json:"mac4"`
	Rssi4  int8   `json:"rssi4"`
	Mac5   string `json:"mac5"`
	Rssi5  int8   `json:"rssi5"`
	Mac6   string `json:"mac6"`
	Rssi6  int8   `json:"rssi6"`
	Mac7   string `json:"mac7"`
	Rssi7  int8   `json:"rssi7"`
}

type BlePayload struct {
	ScanPointer    uint16 `json:"scanPointer"`
	TotalMessages  uint8  `json:"totalMessages"`
	CurrentMessage uint8  `json:"currentMessage"`
	Mac1           string `json:"mac1"`
	Rssi1          int8   `json:"rssi1"`
	Mac2           string `json:"mac2"`
	Rssi2          int8   `json:"rssi2"`
	Mac3           string `json:"mac3"`
	Rssi3          int8   `json:"rssi3"`
	Mac4           string `json:"mac4"`
	Rssi4          int8   `json:"rssi4"`
	Mac5           string `json:"mac5"`
	Rssi5          int8   `json:"rssi5"`
	Mac6           string `json:"mac6"`
	Rssi6          int8   `json:"rssi6"`
}
