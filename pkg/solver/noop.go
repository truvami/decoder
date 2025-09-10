package solver

import (
	"context"

	"github.com/truvami/decoder/pkg/decoder"
)

// Use this solver when you want to skip the position resolution and return a noop result.
type NoopSolver struct {
}

func (n NoopSolver) Solve(ctx context.Context, payload string) (*decoder.DecodedUplink, error) {
	result := decoder.NewDecodedUplink([]decoder.Feature{}, []any{})
	return result, nil
}
