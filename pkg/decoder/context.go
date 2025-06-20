package decoder

type DecoderContextKey string

const (
	// DEVEUI_CONTEXT_KEY is the context key used to store and retrieve the device EUI (Extended Unique Identifier)
	// in the application context. This key is typically used when passing the devEUI value through context objects.
	DEVEUI_CONTEXT_KEY DecoderContextKey = "devEui"

	// FCNT_CONTEXT_KEY is the context key used to store and retrieve the frame count in the application context.
	FCNT_CONTEXT_KEY DecoderContextKey = "fCount"

	// PORT_CONTEXT_KEY is the context key used to store and retrieve the port number in the application context.
	PORT_CONTEXT_KEY DecoderContextKey = "port"
)
