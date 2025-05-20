package encoder

type Encoder interface {
	Encode(any, uint8) (any, error)
}
