package v2

import "errors"

var (
	// Input / option validation
	ErrInvalidOptions       = errors.New("invalid solver options")
	ErrInvalidDevEui        = errors.New("invalid DevEUI (must be 16 hex chars)")
	ErrInvalidPort          = errors.New("invalid port (0-255)")
	ErrInvalidUplinkCounter = errors.New("invalid uplink counter")

	// Request/response handling
	ErrBuildRequest      = errors.New("failed to build request")
	ErrRequestFailed     = errors.New("request failed")
	ErrUnexpectedStatus  = errors.New("unexpected status code")
	ErrDecodeFailed      = errors.New("failed to decode response")
	ErrResponseInvalid   = errors.New("invalid response from LoRaCloud")
	ErrPositionInvalid   = errors.New("position resolution is invalid")
	ErrNoPosition        = errors.New("no position solution in response")
	ErrZeroCoordinates   = errors.New("position has zero coordinates (0,0)")
	ErrMissingCapturedAt = errors.New("missing captured_at (UTC) timestamp in response")
)
