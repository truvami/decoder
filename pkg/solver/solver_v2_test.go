package solver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

func TestMockSolverV2_Solve_ReturnsData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now()
	moving := true

	data := decoder.NewDecodedUplink(
		[]decoder.Feature{decoder.FeatureTimestamp, decoder.FeatureMoving},
		struct {
			Msg string
		}{Msg: "ok"},
	)

	m := MockSolverV2{Data: data}

	got, err := m.Solve(ctx, "AABBCC", SolverV2Options{
		DevEui:        "0102030405060708",
		UplinkCounter: 42,
		Port:          10,
		Timestamp:     &now,
		Moving:        &moving,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != data {
		t.Fatalf("expected pointer to data to be returned, got: %#v", got)
	}
	// sanity check that features on returned data are preserved
	if !got.Is(decoder.FeatureTimestamp) || !got.Is(decoder.FeatureMoving) {
		t.Fatalf("expected returned uplink to have timestamp and moving features, got: %#v", got.GetFeatures())
	}
}

func TestMockSolverV2_Solve_ReturnsError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantErr := errors.New("boom")

	m := MockSolverV2{Err: wantErr}

	got, err := m.Solve(ctx, "DEADBEEF", SolverV2Options{})
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
