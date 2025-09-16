package solver

import (
	"context"
	"errors"
	"testing"

	"github.com/truvami/decoder/pkg/decoder"
)

func TestMockSolverV1_Solve_ReturnsData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	data := decoder.NewDecodedUplink(
		[]decoder.Feature{decoder.FeatureBattery},
		struct {
			Status string
		}{Status: "ok"},
	)

	m := MockSolverV1{Data: data}

	got, err := m.Solve(ctx, "01020304")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != data {
		t.Fatalf("expected pointer to data to be returned, got: %#v", got)
	}
	if !got.Is(decoder.FeatureBattery) {
		t.Fatalf("expected returned uplink to have battery feature")
	}
}

func TestMockSolverV1_Solve_ReturnsError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantErr := errors.New("oops")

	m := MockSolverV1{Err: wantErr}

	got, err := m.Solve(ctx, "BADF00D")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != wantErr.Error() {
		t.Fatalf("unexpected error, want: %v, got: %v", wantErr, err)
	}
	if got != nil {
		t.Fatalf("expected nil data when error is returned, got: %#v", got)
	}
}
