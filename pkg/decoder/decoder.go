package decoder

type Decoder interface {
	Decode(string, int16, string) (interface{}, interface{}, error)
}
