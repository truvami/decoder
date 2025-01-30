package encoder

type Encoder interface {
	Encode(interface{}, int16, string) (interface{}, interface{}, error)
}
