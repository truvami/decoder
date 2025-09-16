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
	DevEui        string
	UplinkCounter uint16
	Port          uint8

	// Optional captured at timestamp of the uplink, if available.
	Timestamp *time.Time
	// Optional indicates if the device is in motion, if available.
	Moving *bool
}

type MockSolverV2 struct {
	Data *decoder.DecodedUplink
	Err  error
}

func (m MockSolverV2) Solve(ctx context.Context, payload string, options SolverV2Options) (*decoder.DecodedUplink, error) {
	return m.Data, m.Err
}
