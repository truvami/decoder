package smartlabel

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

// Moving / Steady
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | tag (0x10 = moving, 0x18 = steady)            | byte       |
// | 1    | 1    | rssi signal 1                                 | int8       |
// | 2    | 6    | mac address signal 1                          | byte[6]    |
// | 8    | 1    | rssi signal 2                                 | int8       |
// | 9    | 6    | mac address signal 2                          | byte[6]    |
// | 15   | 1    | rssi signal 3                                 | int8       |
// | 16   | 6    | mac address signal 3                          | byte[6]    |
// | 22   | 1    | rssi signal 4                                 | int8       |
// | 23   | 6    | mac address signal 4                          | byte[6]    |
// | 29   | 1    | rssi signal 5                                 | int8       |
// | 30   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+
//
// Moving / Steady with timestamp
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | tag (0x14 = moving, 0x1C = steady)            | byte       |
// | 1    | 4    | timestamp                                     | uint32     |
// | 5    | 1    | rssi signal 1                                 | int8       |
// | 6    | 6    | mac address signal 1                          | byte[6]    |
// | 12   | 1    | rssi signal 2                                 | int8       |
// | 13   | 6    | mac address signal 2                          | byte[6]    |
// | 19   | 1    | rssi signal 3                                 | int8       |
// | 20   | 6    | mac address signal 3                          | byte[6]    |
// | 26   | 1    | rssi signal 4                                 | int8       |
// | 27   | 6    | mac address signal 4                          | byte[6]    |
// | 33   | 1    | rssi signal 5                                 | int8       |
// | 34   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+
//
// Moving / Steady with sequence number
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | tag (0x80 = moving, 0x88 = steady)            | byte       |
// | 1    | 2    | sequence number                               | uint16     |
// | 3    | 1    | rssi signal 1                                 | int8       |
// | 4    | 6    | mac address signal 1                          | byte[6]    |
// | 10   | 1    | rssi signal 2                                 | int8       |
// | 11   | 6    | mac address signal 2                          | byte[6]    |
// | 17   | 1    | rssi signal 3                                 | int8       |
// | 18   | 6    | mac address signal 3                          | byte[6]    |
// | 24   | 1    | rssi signal 4                                 | int8       |
// | 25   | 6    | mac address signal 4                          | byte[6]    |
// | 31   | 1    | rssi signal 5                                 | int8       |
// | 32   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+
//
// Moving / Steady with sequence number and timestamp
// +------+------+-----------------------------------------------+------------+
// | Byte | Size | Description                                   | Format     |
// +------+------+-----------------------------------------------+------------+
// | 0    | 1    | tag (0x84 = moving, 0x8C = steady)            | byte       |
// | 1    | 2    | sequence number                               | uint16     |
// | 3    | 4    | timestamp                                     | uint32     |
// | 7    | 1    | rssi signal 1                                 | int8       |
// | 8    | 6    | mac address signal 1                          | byte[6]    |
// | 14   | 1    | rssi signal 2                                 | int8       |
// | 15   | 6    | mac address signal 2                          | byte[6]    |
// | 21   | 1    | rssi signal 3                                 | int8       |
// | 22   | 6    | mac address signal 3                          | byte[6]    |
// | 28   | 1    | rssi signal 4                                 | int8       |
// | 29   | 6    | mac address signal 4                          | byte[6]    |
// | 35   | 1    | rssi signal 5                                 | int8       |
// | 36   | 6    | mac address signal 5                          | byte[6]    |
// +------+------+-----------------------------------------------+------------+

const (
	Port190Moving             = 0x10
	Port190MovingTimestamp    = 0x14
	Port190Steady             = 0x18
	Port190SteadyTimestamp    = 0x1c
	Port190SeqMoving          = 0x80
	Port190SeqMovingTimestamp = 0x84
	Port190SeqSteady          = 0x88
	Port190SeqSteadyTimestamp = 0x8c
)

type Port190Payload struct {
	Tag            byte       `json:"tag" validate:"gte=0,lte=1"`
	SequenceNumber *uint16    `json:"sequenceNumber"`
	Timestamp      *time.Time `json:"timestamp"`
	Rssi1          *int8      `json:"rssi1" validate:"gte=-120,lte=-20"`
	Mac1           string     `json:"mac1"`
	Rssi2          *int8      `json:"rssi2" validate:"gte=-120,lte=-20"`
	Mac2           *string    `json:"mac2"`
	Rssi3          *int8      `json:"rssi3" validate:"gte=-120,lte=-20"`
	Mac3           *string    `json:"mac3"`
	Rssi4          *int8      `json:"rssi4" validate:"gte=-120,lte=-20"`
	Mac4           *string    `json:"mac4"`
	Rssi5          *int8      `json:"rssi5" validate:"gte=-120,lte=-20"`
	Mac5           *string    `json:"mac5"`
}

var _ decoder.UplinkFeatureWiFi = &Port190Payload{}
var _ decoder.UplinkFeatureMoving = &Port190Payload{}
var _ decoder.UplinkFeatureSequenceNumber = &Port190Payload{}
var _ decoder.UplinkFeatureTimestamp = &Port190Payload{}

func (p Port190Payload) GetAccessPoints() []decoder.AccessPoint {
	accessPoints := []decoder.AccessPoint{}

	if p.Mac1 != "" {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  p.Mac1,
			RSSI: p.Rssi1,
		})
	}

	if p.Mac2 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac2,
			RSSI: p.Rssi2,
		})
	}

	if p.Mac3 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac3,
			RSSI: p.Rssi3,
		})
	}

	if p.Mac4 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac4,
			RSSI: p.Rssi4,
		})
	}

	if p.Mac5 != nil {
		accessPoints = append(accessPoints, decoder.AccessPoint{
			MAC:  *p.Mac5,
			RSSI: p.Rssi5,
		})
	}

	return accessPoints
}

func (p Port190Payload) IsMoving() bool {
	switch p.Tag {
	case Port190Moving:
		return true
	case Port190MovingTimestamp:
		return true
	case Port190SeqMoving:
		return true
	case Port190SeqMovingTimestamp:
		return true
	}

	return false
}

func (p Port190Payload) GetTimestamp() *time.Time {
	return p.Timestamp
}

func (p Port190Payload) GetSequenceNumber() *uint16 {
	return p.SequenceNumber
}
