package encoder

type Encoder interface {
	Encode(any, uint8, string) (any, any, error)
}
