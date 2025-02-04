package decoder

import "github.com/truvami/decoder/pkg/common"

type Decoder interface {
	Decode(string, int16, string) (interface{}, interface{}, error)
	DecodePosition(string, int16, string) (common.Position, interface{}, error)
	DecodeWifi(string, int16, string) (common.WifiLocation, interface{}, error)
}
