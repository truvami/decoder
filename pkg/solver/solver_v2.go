package solver

import (
	"context"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type SolverV2 interface {
	Solve(ctx context.Context, payload string, options SolverV2Options) (*decoder.DecodedUplink, error)
}

type SolverV2Options struct {
	DevEui string

	// UplinkCounter is the 16-bit device counter (FCntUp modulo 65536).
	// If an upstream component provides a 32-bit frame counter, truncate to the
	// lower 16 bits here and (optionally) carry the full 32-bit value via context
	// or a separate field if your solver needs it.
	UplinkCounter uint16

	Port uint8

	// Timestamp is the captured-at time of the uplink (UTC), when available.
	Timestamp *time.Time
	// Moving, when set, indicates whether the device was in motion for this uplink.
	Moving *bool
}

type MockSolverV2 struct {
	Data *decoder.DecodedUplink
	Err  error
}

func (m MockSolverV2) Solve(ctx context.Context, payload string, options SolverV2Options) (*decoder.DecodedUplink, error) {
	return m.Data, m.Err
}
