package solver

import "github.com/truvami/decoder/pkg/decoder"

type SolverV1 interface {
	Solve(payload string) (*decoder.DecodedUplink, error)
}

type MockSolverV1 struct {
	Data *decoder.DecodedUplink
	Err  error
}

func (m MockSolverV1) Solve(payload string) (*decoder.DecodedUplink, error) {
	return m.Data, m.Err
}
