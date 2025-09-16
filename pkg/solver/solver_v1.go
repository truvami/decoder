package solver

import (
	"context"

	"github.com/truvami/decoder/pkg/decoder"
)

type SolverV1 interface {
	Solve(ctx context.Context, payload string) (*decoder.DecodedUplink, error)
}

type MockSolverV1 struct {
	Data *decoder.DecodedUplink
	Err  error
}

func (m MockSolverV1) Solve(ctx context.Context, payload string) (*decoder.DecodedUplink, error) {
	return m.Data, m.Err
}
