package tagxl

import (
	"testing"
	"time"
)

func TestPort193Payload_GNSSMethodsAndMoving(t *testing.T) {
	t.Parallel()

	p := Port193Payload{
		EndOfGroup: true,
		GroupToken: 10,
		NavMessage: []byte{0x01, 0x02, 0x03},
		// Moving field is intentionally false; IsMoving() should still return true.
		Moving:    false,
		Latitude:  47.3769,
		Longitude: 8.5417,
		Altitude:  500.0,
	}

	if got := p.GetLatitude(); got != p.Latitude {
		t.Fatalf("GetLatitude() = %v, want %v", got, p.Latitude)
	}
	if got := p.GetLongitude(); got != p.Longitude {
		t.Fatalf("GetLongitude() = %v, want %v", got, p.Longitude)
	}
	if got := p.GetAltitude(); got != p.Altitude {
		t.Fatalf("GetAltitude() = %v, want %v", got, p.Altitude)
	}
	if got := p.GetAccuracy(); got != nil {
		t.Fatalf("GetAccuracy() = %v, want nil", got)
	}
	if got := p.GetTTF(); got != nil {
		t.Fatalf("GetTTF() = %v, want nil", got)
	}
	if got := p.GetPDOP(); got != nil {
		t.Fatalf("GetPDOP() = %v, want nil", got)
	}
	if got := p.GetSatellites(); got != nil {
		t.Fatalf("GetSatellites() = %v, want nil", got)
	}
	if got := p.IsMoving(); got != true {
		t.Fatalf("IsMoving() = %v, want true", got)
	}
}

// Compile-time like assertion helpers (executed to ensure method signatures remain compatible).
func TestPort193Payload_InterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var _ = func() any {
		var _ interface {
			GetLatitude() float64
			GetLongitude() float64
			GetAltitude() float64
			GetAccuracy() *float64
			GetTTF() *time.Duration
			GetPDOP() *float64
			GetSatellites() *uint8
		} = Port193Payload{}
		var _ interface {
			IsMoving() bool
		} = Port193Payload{}
		return nil
	}()
}
