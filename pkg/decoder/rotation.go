package decoder

type RotationState string

const (
	RotationStateUndefined RotationState = "undefined"
	RotationStatePouring   RotationState = "pouring"
	RotationStateMixing    RotationState = "mixing"
	RotationStateError     RotationState = "error"
)
