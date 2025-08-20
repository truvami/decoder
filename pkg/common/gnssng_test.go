package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeGNSSNGHeader(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		header  GNSSNGHeader
		err     error
	}{
		{
			name:    "empty",
			payload: []byte{},
			header: GNSSNGHeader{
				EndOfGroup:           false,
				ReservedForFutureUse: 0,
				GroupToken:           0,
			},
			err: ErrGNSSNGHeaderByteMissing,
		},
		{
			name:    "end of group",
			payload: []byte{0b1000_0000},
			header: GNSSNGHeader{
				EndOfGroup:           true,
				ReservedForFutureUse: 0,
				GroupToken:           0,
			},
			err: nil,
		},
		{
			name:    "group token",
			payload: []byte{0b0001_1111},
			header: GNSSNGHeader{
				EndOfGroup:           false,
				ReservedForFutureUse: 0,
				GroupToken:           0b0001_1111,
			},
			err: nil,
		},
		{
			name:    "end of group with group token (31)",
			payload: []byte{0b1000_0000 | 0b0001_1111},
			header: GNSSNGHeader{
				EndOfGroup:           true,
				ReservedForFutureUse: 0,
				GroupToken:           31,
			},
			err: nil,
		},
		{
			name:    "end of group with group token (9)",
			payload: []byte{0b1000_0000 | 0b0000_1001},
			header: GNSSNGHeader{
				EndOfGroup:           true,
				ReservedForFutureUse: 0,
				GroupToken:           9,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			header, err := DecodeGNSSNGHeader(test.payload)
			assert.Equal(t, header, test.header)
			assert.Equal(t, err, test.err)
		})
	}
}
