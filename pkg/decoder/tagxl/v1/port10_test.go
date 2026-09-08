package tagxl

import (
	"testing"
	"time"

	"github.com/truvami/decoder/pkg/common"
)

func TestPort10Payload_GNSSAndBatteryMethods(t *testing.T) {
	t.Parallel()

	ttf := 25 * time.Second
	p := Port10Payload{
		Moving:     false,
		Latitude:   47.385351,
		Longitude:  8.538399,
		Altitude:   43.89,
		Timestamp:  time.Date(2024, 10, 23, 11, 11, 58, 0, time.UTC),
		Battery:    3.806,
		TTF:        &ttf,
		PDOP:       common.Float64Ptr(2.5),
		Satellites: common.Uint8Ptr(5),
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
	if got := p.GetTTF(); got == nil || *got != ttf {
		t.Fatalf("GetTTF() = %v, want %v", got, ttf)
	}
	if got := p.GetPDOP(); got == nil || *got != 2.5 {
		t.Fatalf("GetPDOP() = %v, want 2.5", got)
	}
	if got := p.GetSatellites(); got == nil || *got != 5 {
		t.Fatalf("GetSatellites() = %v, want 5", got)
	}
	if got := p.GetTimestamp(); got == nil || !got.Equal(p.Timestamp) {
		t.Fatalf("GetTimestamp() = %v, want %v", got, p.Timestamp)
	}
	if got := p.GetBatteryVoltage(); got != p.Battery {
		t.Fatalf("GetBatteryVoltage() = %v, want %v", got, p.Battery)
	}
	if got := p.GetLowBattery(); got != nil {
		t.Fatalf("GetLowBattery() = %v, want nil", got)
	}
	if got := p.IsMoving(); got {
		t.Fatalf("IsMoving() = %v, want false", got)
	}

	p.Moving = true
	if got := p.IsMoving(); !got {
		t.Fatalf("IsMoving() = %v, want true", got)
	}
}
