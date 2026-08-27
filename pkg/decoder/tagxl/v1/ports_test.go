package tagxl

import (
	"context"
	"strings"
	"testing"

	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/solver"
)

func TestIsGnssSolverPort(t *testing.T) {
	t.Parallel()

	for _, port := range []uint8{192, 193, 194, 195, 199, 210, 211} {
		if !IsGnssSolverPort(port) {
			t.Errorf("port %d should be solver-backed", port)
		}
	}

	// Port 10 is an on-device fix; the rest are status, config, BLE and Wi-Fi.
	for _, port := range []uint8{0, 10, 150, 151, 152, 160, 161, 162, 163, 197, 198, 200, 201, 212, 213, 255} {
		if IsGnssSolverPort(port) {
			t.Errorf("port %d should not be solver-backed", port)
		}
	}
}

func TestGnssSolverPortsIsSortedAndComplete(t *testing.T) {
	t.Parallel()

	got := GnssSolverPorts()
	want := []uint8{192, 193, 194, 195, 199, 210, 211}

	if len(got) != len(want) {
		t.Fatalf("GnssSolverPorts() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("GnssSolverPorts() = %v, want %v", got, want)
		}
	}

	// The returned slice must not alias internal state.
	got[0] = 0
	if GnssSolverPorts()[0] != 192 {
		t.Fatal("GnssSolverPorts() leaked internal state to the caller")
	}
}

// The predicate must agree with what Decode actually dispatches on. Every
// solver port takes the solver branch; a non-solver port is decoded locally.
// Ports 194/195/210/211 additionally require a v2 solver, and the distinct
// error they return with only a v1 solver is itself proof they took the
// solver branch rather than falling through to getConfig.
func TestIsGnssSolverPortMatchesDecodeDispatch(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}

	d := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)

	v2Only := map[uint8]bool{194: true, 195: true, 210: true, 211: true}

	for _, port := range GnssSolverPorts() {
		_, err := d.Decode(context.TODO(), "00", port)
		if v2Only[port] {
			if err == nil || !strings.Contains(err.Error(), "without v2 solver") {
				t.Errorf("port %d: expected the v2-solver requirement from the solver branch, got %v", port, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("port %d: expected the v1 solver to handle it, got %v", port, err)
		}
	}

	// A non-solver port is rejected by getConfig, not by the solver branch.
	_, err := d.Decode(context.TODO(), "00", 0)
	if err == nil || !strings.Contains(err.Error(), "port 0 not supported") {
		t.Errorf("port 0: expected a local unsupported-port error, got %v", err)
	}
}
