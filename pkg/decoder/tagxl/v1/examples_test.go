package tagxl

import (
	"context"
	"testing"
	"time"

	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
)

// dummyGNSS implements UplinkFeatureGNSS for mocking solver responses
type dummyGNSS struct{}

var _ decoder.UplinkFeatureGNSS = dummyGNSS{}

func (d dummyGNSS) GetLatitude() float64   { return 51.49278 }
func (d dummyGNSS) GetLongitude() float64  { return 0.0212 }
func (d dummyGNSS) GetAltitude() float64   { return 83.93 }
func (d dummyGNSS) GetAccuracy() *float64  { return nil }
func (d dummyGNSS) GetTTF() *time.Duration { return nil }
func (d dummyGNSS) GetPDOP() *float64      { return nil }
func (d dummyGNSS) GetSatellites() *uint8  { return nil }

// Test cases derived from user-provided device logs

// moving Wi-Fi message with timestamp (port 201, v2 RSSI+MAC with 4B timestamp)
func TestTagXL_WiFi_Port201_Timestamped(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	// 0x68BAE3AB timestamp (1757078443) + v2 + [RSSI,MAC] tuples
	payload := "68bae3ab01d3f0b0140c96bbc7e4c32a622ea4c5e0286d8a9478b4e0286d8aabfcada86e84e1a812"

	out, err := dec.Decode(context.TODO(), payload, 201)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	// Features
	if !out.Is(decoder.FeatureWiFi) || !out.Is(decoder.FeatureTimestamp) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected WiFi + Timestamp + Moving features")
	}

	ts := out.Data.(decoder.UplinkFeatureTimestamp).GetTimestamp()
	if ts == nil || ts.Unix() != 1757078443 {
		t.Fatalf("expected timestamp 1757078443, got %v", ts)
	}
	wifi := out.Data.(decoder.UplinkFeatureWiFi)
	aps := wifi.GetAccessPoints()
	if len(aps) != 5 {
		t.Fatalf("expected 5 APs, got %d", len(aps))
	}
	// Order and values based on the log
	exp := []struct {
		mac  string
		rssi int8
	}{
		{"f0b0140c96bb", -45},
		{"e4c32a622ea4", -57},
		{"e0286d8a9478", -59},
		{"e0286d8aabfc", -76},
		{"a86e84e1a812", -83},
	}
	for i, ap := range aps {
		if ap.MAC != exp[i].mac || (ap.RSSI == nil || *ap.RSSI != exp[i].rssi) {
			t.Fatalf("ap[%d] expected MAC=%s RSSI=%d, got MAC=%s RSSI=%v", i, exp[i].mac, exp[i].rssi, ap.MAC, ap.RSSI)
		}
	}
}

// Wi-Fi message without timestamp (port 197, v2 RSSI+MAC)
func TestTagXL_WiFi_Port197_NoTimestamp(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	// 0x01CFF0B0140C96BB... (version=01, five RSSI+MAC tuples)
	payload := "01cff0b0140c96bbcce4c32a622ea4c8e0286d8a9478b8e0286d8aabfcafa86e84e1a812"

	out, err := dec.Decode(context.TODO(), payload, 197)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !out.Is(decoder.FeatureWiFi) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected WiFi + Moving features")
	}
	wifi := out.Data.(decoder.UplinkFeatureWiFi)
	aps := wifi.GetAccessPoints()
	if len(aps) != 5 {
		t.Fatalf("expected 5 APs, got %d", len(aps))
	}
	exp := []struct {
		mac  string
		rssi int8
	}{
		{"f0b0140c96bb", -49},
		{"e4c32a622ea4", -52},
		{"e0286d8a9478", -56},
		{"e0286d8aabfc", -72},
		{"a86e84e1a812", -81},
	}
	for i, ap := range aps {
		if ap.MAC != exp[i].mac || (ap.RSSI == nil || *ap.RSSI != exp[i].rssi) {
			t.Fatalf("ap[%d] expected MAC=%s RSSI=%d, got MAC=%s RSSI=%v", i, exp[i].mac, exp[i].rssi, ap.MAC, ap.RSSI)
		}
	}
}

// moving Wi-Fi message without timestamp (port 198, v2 RSSI+MAC)
func TestTagXL_WiFi_Port198_NoTimestamp(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	// 0x01CDE4C32A622EA4AD726C9A74B58DAEA86E84E1A812B9E0286D8AABFCC4E0286D8A9478
	payload := "01cde4c32a622ea4ad726c9a74b58daea86e84e1a812b9e0286d8aabfcc4e0286d8a9478"

	out, err := dec.Decode(context.TODO(), payload, 198)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !out.Is(decoder.FeatureWiFi) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected WiFi + Moving features")
	}
	wifi := out.Data.(decoder.UplinkFeatureWiFi)
	aps := wifi.GetAccessPoints()
	if len(aps) != 5 {
		t.Fatalf("expected 5 APs, got %d", len(aps))
	}
	exp := []struct {
		mac  string
		rssi int8
	}{
		{"e4c32a622ea4", -51},
		{"726c9a74b58d", -83},
		{"a86e84e1a812", -82},
		{"e0286d8aabfc", -71},
		{"e0286d8a9478", -60},
	}
	for i, ap := range aps {
		if ap.MAC != exp[i].mac || (ap.RSSI == nil || *ap.RSSI != exp[i].rssi) {
			t.Fatalf("ap[%d] expected MAC=%s RSSI=%d, got MAC=%s RSSI=%v", i, exp[i].mac, exp[i].rssi, ap.MAC, ap.RSSI)
		}
	}
}

// GNSS messages (2) without timestamp and not moving (port 192, no timestamp variant)
// We use a mock solver V1 to avoid external calls and validate GNSS feature presence.
func TestTagXL_GNSS_Port192_NoTimestamp(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	ctx := context.Background()
	// Mock solver returns GNSS feature regardless of payload
	mock := solver.MockSolverV1{
		Data: decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS}, dummyGNSS{}),
	}
	dec := NewTagXLv1Decoder(ctx, mock, logger.Logger)

	payloads := []string{
		"12b319f3ba0cd9805a773bf029a2f97f9077754db825d87d456527144c8493c57b350f",
		"92b31df3ba0cd9004a751ba293dd7ab8ea26e4b490ad130aa6ecc9e0ba4a00",
	}
	for _, pl := range payloads {
		out, err := dec.Decode(ctx, pl, 192)
		if err != nil {
			t.Fatalf("decode error for payload %s: %v", pl, err)
		}
		if !out.Is(decoder.FeatureGNSS) {
			t.Fatalf("expected GNSS feature for payload %s", pl)
		}
	}
}

// moving GNSS messages (2) without timestamp (port 193, moving variant)
// Use mock solver V1 and validate GNSS feature presence.
func TestTagXL_GNSS_Port193_Moving_NoTimestamp(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	ctx := context.Background()
	mock := solver.MockSolverV1{
		Data: decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS}, dummyGNSS{}),
	}
	dec := NewTagXLv1Decoder(ctx, mock, logger.Logger)

	payloads := []string{
		"13ab251151f3ba0ed580391609b73f4e2a22e0128d15242653798a4e6056cc1d0d",
		"93ab35c650f33a0cd5004a161989f27af55aec906ed9e366120a484a2600189d00",
	}
	for _, pl := range payloads {
		out, err := dec.Decode(ctx, pl, 193)
		if err != nil {
			t.Fatalf("decode error for payload %s: %v", pl, err)
		}
		if !out.Is(decoder.FeatureGNSS) {
			t.Fatalf("expected GNSS feature for payload %s", pl)
		}
	}
}

// GNSS timestamped, not moving, rotation-triggered (port 210) — same format as 194
func TestTagXL_GNSS_Port210_Timestamped_RotationTriggered(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	log := logger.Logger
	devEui := "0011223344556677"
	fcnt := 42

	// Timestamp 0x68BAD325 => 1757074213
	secs := uint32(1757074213)
	ts := time.Unix(int64(secs), 0).UTC()

	payloads := []string{
		"68bad32509ab91418ae63a10b5004a0a3fef037ab2f06ce8e510820c1a0bdcecb49e1543fdd2f28f1c",
		"68bad32589b379e7ba0fb5006b9aaa8c8e25febf16f4e5c31d0cc8ca12a1cffdddf16c2cf82877f1edee4ecbc5ef54",
	}
	for _, payload := range payloads {
		cap := &captureSolverV2{
			resp: decoder.NewDecodedUplink(
				[]decoder.Feature{decoder.FeatureGNSS, decoder.FeatureTimestamp},
				&fakeGNSSData{lat: 47.0, lon: 8.0, alt: 10.0, ts: &ts},
			),
		}
		dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log, WithSolverV2(cap))

		ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
		ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

		out, err := dec.Decode(ctx, payload, 210)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedForwarded := payload[8:] // strip 4B timestamp (8 hex chars)
		if cap.lastPayload != expectedForwarded {
			t.Fatalf("expected forwarded payload %q, got %q", expectedForwarded, cap.lastPayload)
		}
		if cap.lastOptions.Port != 192 {
			t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
		}
		if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != false {
			t.Fatalf("expected Moving=false for port 210, got %+v", cap.lastOptions.Moving)
		}
		if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) {
			t.Fatalf("expected GNSS and Timestamp features in result")
		}
	}
}

// GNSS timestamped, moving, rotation-triggered (port 211) — same format as 195
func TestTagXL_GNSS_Port211_Timestamped_RotationTriggered(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	log := logger.Logger
	devEui := "0011223344556677"
	fcnt := 43

	// Timestamp 0x68BAD3C5 => 1757074373
	secs := uint32(1757074373)
	ts := time.Unix(int64(secs), 0).UTC()

	payloads := []string{
		"68bad3c50aabd56cb2e7ba0db5805a5ac9d4edd8de8a021b4ae2b78e8c0b8391566ab8d47d1d4c55ae794a2c2da7a637b49d32e44800",
		"68bad3c58aab4581b9e73a0eb580da120d7f85a75e770c6acad3dc2acdacbdcd576ab8147f5902557379b18d0f676a35fb9a6ae5ee03",
	}
	for _, payload := range payloads {
		cap := &captureSolverV2{
			resp: decoder.NewDecodedUplink(
				[]decoder.Feature{decoder.FeatureGNSS, decoder.FeatureTimestamp},
				&fakeGNSSData{lat: 47.0, lon: 8.0, alt: 10.0, ts: &ts},
			),
		}
		dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log, WithSolverV2(cap))

		ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
		ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

		out, err := dec.Decode(ctx, payload, 211)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedForwarded := payload[8:] // strip 4B timestamp (8 hex chars)
		if cap.lastPayload != expectedForwarded {
			t.Fatalf("expected forwarded payload %q, got %q", expectedForwarded, cap.lastPayload)
		}
		if cap.lastOptions.Port != 192 {
			t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
		}
		if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != true {
			t.Fatalf("expected Moving=true for port 211, got %+v", cap.lastOptions.Moving)
		}
		if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) {
			t.Fatalf("expected GNSS and Timestamp features in result")
		}
	}
}

// WiFi non-moving, timestamped, rotation-triggered (port 212) — same format as 200 but no Buffered
func TestTagXL_WiFi_Port212_Timestamped_RotationTriggered(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	// 0x68BAE3AB timestamp (1757078443) + v2 + [RSSI,MAC] tuples
	payload := "68bae3ab01d3f0b0140c96bbc7e4c32a622ea4c5e0286d8a9478b4e0286d8aabfcada86e84e1a812"

	out, err := dec.Decode(context.TODO(), payload, 212)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	// Features: WiFi + Timestamp + Moving, but NOT Buffered
	if !out.Is(decoder.FeatureWiFi) || !out.Is(decoder.FeatureTimestamp) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected WiFi + Timestamp + Moving features")
	}
	if out.Is(decoder.FeatureBuffered) {
		t.Fatalf("port 212 should NOT have Buffered feature")
	}

	ts := out.Data.(decoder.UplinkFeatureTimestamp).GetTimestamp()
	if ts == nil || ts.Unix() != 1757078443 {
		t.Fatalf("expected timestamp 1757078443, got %v", ts)
	}
	wifi := out.Data.(decoder.UplinkFeatureWiFi)
	aps := wifi.GetAccessPoints()
	if len(aps) != 5 {
		t.Fatalf("expected 5 APs, got %d", len(aps))
	}
	moving := out.Data.(decoder.UplinkFeatureMoving)
	if moving.IsMoving() {
		t.Fatalf("expected IsMoving=false for port 212")
	}
}

// WiFi moving, timestamped, rotation-triggered (port 213) — same format as 201 but no Buffered
func TestTagXL_WiFi_Port213_Timestamped_RotationTriggered(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	// 0x68BAE3AB timestamp (1757078443) + v2 + [RSSI,MAC] tuples
	payload := "68bae3ab01d3f0b0140c96bbc7e4c32a622ea4c5e0286d8a9478b4e0286d8aabfcada86e84e1a812"

	out, err := dec.Decode(context.TODO(), payload, 213)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	// Features: WiFi + Timestamp + Moving, but NOT Buffered
	if !out.Is(decoder.FeatureWiFi) || !out.Is(decoder.FeatureTimestamp) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected WiFi + Timestamp + Moving features")
	}
	if out.Is(decoder.FeatureBuffered) {
		t.Fatalf("port 213 should NOT have Buffered feature")
	}

	ts := out.Data.(decoder.UplinkFeatureTimestamp).GetTimestamp()
	if ts == nil || ts.Unix() != 1757078443 {
		t.Fatalf("expected timestamp 1757078443, got %v", ts)
	}
	wifi := out.Data.(decoder.UplinkFeatureWiFi)
	aps := wifi.GetAccessPoints()
	if len(aps) != 5 {
		t.Fatalf("expected 5 APs, got %d", len(aps))
	}
	moving := out.Data.(decoder.UplinkFeatureMoving)
	if !moving.IsMoving() {
		t.Fatalf("expected IsMoving=true for port 213")
	}
}

// GNSS timestamped, not moving, two-frame NAV (port 194) from logs
func TestTagXL_GNSS_Port194_Timestamped(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	log := logger.Logger
	devEui := "0011223344556677"
	fcnt := 42

	// Timestamp 0x68BAD325 => 1757074213
	secs := uint32(1757074213)
	ts := time.Unix(int64(secs), 0).UTC()

	payloads := []string{
		"68bad32509ab91418ae63a10b5004a0a3fef037ab2f06ce8e510820c1a0bdcecb49e1543fdd2f28f1c",
		"68bad32589b379e7ba0fb5006b9aaa8c8e25febf16f4e5c31d0cc8ca12a1cffdddf16c2cf82877f1edee4ecbc5ef54",
	}
	for _, payload := range payloads {
		cap := &captureSolverV2{
			resp: decoder.NewDecodedUplink(
				[]decoder.Feature{decoder.FeatureGNSS, decoder.FeatureTimestamp},
				&fakeGNSSData{lat: 47.0, lon: 8.0, alt: 10.0, ts: &ts},
			),
		}
		dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log, WithSolverV2(cap))

		ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
		ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

		out, err := dec.Decode(ctx, payload, 194)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedForwarded := payload[8:] // strip 4B timestamp (8 hex chars)
		if cap.lastPayload != expectedForwarded {
			t.Fatalf("expected forwarded payload %q, got %q", expectedForwarded, cap.lastPayload)
		}
		if cap.lastOptions.Port != 192 {
			t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
		}
		if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != false {
			t.Fatalf("expected Moving=false for port 194, got %+v", cap.lastOptions.Moving)
		}
		if cap.lastOptions.Timestamp == nil || !cap.lastOptions.Timestamp.Equal(ts) {
			t.Fatalf("expected Timestamp=%v for port 194, got %+v", ts, cap.lastOptions.Timestamp)
		}
		if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) {
			t.Fatalf("expected GNSS and Timestamp features in result")
		}
	}
}

// GNSS timestamped, moving, two-frame NAV (port 195) from logs
func TestTagXL_GNSS_Port195_Timestamped(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	log := logger.Logger
	devEui := "0011223344556677"
	fcnt := 43

	// Timestamp 0x68BAD3C5 => 1757074373
	secs := uint32(1757074373)
	ts := time.Unix(int64(secs), 0).UTC()

	payloads := []string{
		"68bad3c50aabd56cb2e7ba0db5805a5ac9d4edd8de8a021b4ae2b78e8c0b8391566ab8d47d1d4c55ae794a2c2da7a637b49d32e44800",
		"68bad3c58aab4581b9e73a0eb580da120d7f85a75e770c6acad3dc2acdacbdcd576ab8147f5902557379b18d0f676a35fb9a6ae5ee03",
	}
	for _, payload := range payloads {
		cap := &captureSolverV2{
			resp: decoder.NewDecodedUplink(
				[]decoder.Feature{decoder.FeatureGNSS, decoder.FeatureTimestamp},
				&fakeGNSSData{lat: 47.0, lon: 8.0, alt: 10.0, ts: &ts},
			),
		}
		dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log, WithSolverV2(cap))

		ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
		ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

		out, err := dec.Decode(ctx, payload, 195)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedForwarded := payload[8:] // strip 4B timestamp (8 hex chars)
		if cap.lastPayload != expectedForwarded {
			t.Fatalf("expected forwarded payload %q, got %q", expectedForwarded, cap.lastPayload)
		}
		if cap.lastOptions.Port != 192 {
			t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
		}
		if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != true {
			t.Fatalf("expected Moving=true for port 195, got %+v", cap.lastOptions.Moving)
		}
		if cap.lastOptions.Timestamp == nil || !cap.lastOptions.Timestamp.Equal(ts) {
			t.Fatalf("expected Timestamp=%v for port 195, got %+v", ts, cap.lastOptions.Timestamp)
		}
		if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) {
			t.Fatalf("expected GNSS and Timestamp features in result")
		}
	}
}
