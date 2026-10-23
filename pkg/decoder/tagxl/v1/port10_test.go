package tagxl

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/solver"
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

func TestPort10_OptionalAccuracyAndSpeed(t *testing.T) {
	t.Parallel()

	// 20-byte uplink from older firmware, then with hAcc (12 m) and speed over ground (34 km/h) appended
	base := "0002d30b070082491f11256718d9fe0ede190505"
	tests := []struct {
		name     string
		payload  string
		accuracy *float64
		speed    *float64
	}{
		{name: "20 bytes, older firmware", payload: base},
		{name: "21 bytes, hAcc only", payload: base + "0c", accuracy: common.Float64Ptr(12)},
		{name: "22 bytes, hAcc and speed", payload: base + "0c22", accuracy: common.Float64Ptr(12), speed: common.Float64Ptr(34)},
		{name: "22 bytes, hAcc capped at 255, standing still", payload: base + "ff00", accuracy: common.Float64Ptr(255), speed: common.Float64Ptr(0)},
	}

	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := dec.Decode(context.TODO(), tt.payload, 10)
			require.NoError(t, err)
			data, ok := decoded.Data.(Port10Payload)
			require.True(t, ok, "unexpected payload type %T", decoded.Data)

			// the existing fields decode the same regardless of length
			assert.Equal(t, 47.385351, data.Latitude)
			assert.Equal(t, common.Uint8Ptr(5), data.Satellites)

			assert.Equal(t, tt.accuracy, data.Accuracy)
			assert.Equal(t, tt.accuracy, data.GetAccuracy())
			assert.Equal(t, tt.speed, data.Speed)
		})
	}
}

func TestPort10_AccuracyAndSpeedJSON(t *testing.T) {
	t.Parallel()

	p := Port10Payload{Accuracy: common.Float64Ptr(12), Speed: common.Float64Ptr(34)}
	encoded, err := json.Marshal(p)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"accuracy":"12m"`)
	assert.Contains(t, string(encoded), `"speed":"34km/h"`)

	encoded, err = json.Marshal(Port10Payload{})
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"accuracy":null`)
	assert.Contains(t, string(encoded), `"speed":null`)
}

func TestPort10_RejectsUnknownTrailingBytes(t *testing.T) {
	t.Parallel()

	// the maximum length is derived from the configured fields: a byte the decoder doesn't
	// know yet is rejected, so the decoder must be updated before firmware appends more
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	_, err := dec.Decode(context.TODO(), "0002d30b070082491f11256718d9fe0ede1905050c22ff", 10)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payload too long")
}
