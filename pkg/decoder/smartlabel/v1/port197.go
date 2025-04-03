package smartlabel

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// Uplink port: 197

// Byte 0	Byte 1	Bytes[2:7]	Byte 8	Bytes[9:14]	...
// Value	0x01	AP1 RSSI	AP1 MAC	AP2 RSSI	AP2 MAC		APN RSSI	APN MAC
// Size [Bytes]	1	1	6	1	6		1	6
// Type	UINT8	UINT8	UINT8	UINT8	UINT8		UINT8	UINT8

type Port197Payload struct {
	Mac1  string `json:"mac1"`
	Rssi1 int8   `json:"rssi1"`
	Mac2  string `json:"mac2"`
	Rssi2 int8   `json:"rssi2"`
	Mac3  string `json:"mac3"`
	Rssi3 int8   `json:"rssi3"`
	Mac4  string `json:"mac4"`
	Rssi4 int8   `json:"rssi4"`
	Mac5  string `json:"mac5"`
	Rssi5 int8   `json:"rssi5"`
	Mac6  string `json:"mac6"`
	Rssi6 int8   `json:"rssi6"`
}

var _ decoder.UplinkFeatureBase = &Port197Payload{}
var _ decoder.UplinkFeatureWiFi = &Port197Payload{}

func (p Port197Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port197Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	if p.Mac6 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac6,
			RSSI: p.Rssi6,
		})
	}

	return accessPoints
}
