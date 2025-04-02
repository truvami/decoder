package encoder

type Encoder interface {
	Encode(any, int16, string) (any, any, error)
}
