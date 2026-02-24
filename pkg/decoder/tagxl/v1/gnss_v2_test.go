package tagxl

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"go.uber.org/zap"
)

// fakeGNSSData implements minimal GNSS (and optionally Timestamp) interfaces for testing.
type fakeGNSSData struct {
	lat, lon, alt float64
	ts            *time.Time
}

var _ decoder.UplinkFeatureGNSS = &fakeGNSSData{}
var _ decoder.UplinkFeatureTimestamp = &fakeGNSSData{}

func (f *fakeGNSSData) GetLatitude() float64     { return f.lat }
func (f *fakeGNSSData) GetLongitude() float64    { return f.lon }
func (f *fakeGNSSData) GetAltitude() float64     { return f.alt }
func (f *fakeGNSSData) GetAccuracy() *float64    { return nil }
func (f *fakeGNSSData) GetTTF() *time.Duration   { return nil }
func (f *fakeGNSSData) GetPDOP() *float64        { return nil }
func (f *fakeGNSSData) GetSatellites() *uint8    { return nil }
func (f *fakeGNSSData) GetTimestamp() *time.Time { return f.ts }

// captureSolverV2 captures the last payload and options passed to Solve.
type captureSolverV2 struct {
	lastPayload string
	lastOptions solver.SolverV2Options
	resp        *decoder.DecodedUplink
	err         error
}

func (c *captureSolverV2) Solve(ctx context.Context, payload string, options solver.SolverV2Options) (*decoder.DecodedUplink, error) {
	c.lastPayload = payload
	c.lastOptions = options
	return c.resp, c.err
}

func newLogger() *zap.Logger {
	if logger.Logger == nil {
		logger.NewLogger()
	}
	return logger.Logger
}

func TestGNSS_SolverV2_192_193_NoTimestamp(t *testing.T) {
	log := newLogger()

	// header 0x80 => EndOfGroup set; payload not otherwise relevant for capture
	payload := "80abcd"
	devEui := "0011223344556677"
	fcnt := 123

	// Prepare a response with GNSS (no timestamp)
	resp := decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS}, &fakeGNSSData{
		lat: 47.0, lon: 8.0, alt: 10.0, ts: nil,
	})

	cap := &captureSolverV2{resp: resp}

	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log,
		WithSolverV2(cap),
	)

	ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
	ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

	// Port 192 -> steady (moving=false), no timestamp stripping
	out, err := dec.Decode(ctx, payload, 192)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.lastPayload != payload {
		t.Fatalf("expected payload forwarded unchanged for port 192, got %q", cap.lastPayload)
	}
	if cap.lastOptions.DevEui != devEui || cap.lastOptions.UplinkCounter != uint16(fcnt) || cap.lastOptions.Port != 192 {
		t.Fatalf("unexpected options for port 192: %+v", cap.lastOptions)
	}
	if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != false {
		t.Fatalf("expected Moving=false for port 192, got %+v", cap.lastOptions.Moving)
	}
	if cap.lastOptions.Timestamp != nil {
		t.Fatalf("expected no timestamp for port 192")
	}
	if !out.Is(decoder.FeatureGNSS) {
		t.Fatalf("expected FeatureGNSS in result")
	}

	// Port 193 -> moving (moving=true), no timestamp stripping
	out, err = dec.Decode(ctx, payload, 193)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.lastPayload != payload {
		t.Fatalf("expected payload forwarded unchanged for port 193, got %q", cap.lastPayload)
	}
	if cap.lastOptions.Port != 192 {
		t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
	}
	if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != true {
		t.Fatalf("expected Moving=true for port 193, got %+v", cap.lastOptions.Moving)
	}
	if cap.lastOptions.Timestamp != nil {
		t.Fatalf("expected no timestamp for port 193")
	}
	if !out.Is(decoder.FeatureGNSS) {
		t.Fatalf("expected FeatureGNSS in result")
	}
}

func TestGNSS_SolverV2_210_211_TimestampStrippedAndPassed(t *testing.T) {
	log := newLogger()

	// Build a payload with 4B timestamp prefix + header 0x80 after
	secs := uint32(1750000000) // some fixed time
	ts := time.Unix(int64(secs), 0).UTC()
	tsBytes := make([]byte, 4)
	// big-endian
	tsBytes[0] = byte((secs >> 24) & 0xff)
	tsBytes[1] = byte((secs >> 16) & 0xff)
	tsBytes[2] = byte((secs >> 8) & 0xff)
	tsBytes[3] = byte(secs & 0xff)

	tsHex := hex.EncodeToString(tsBytes)
	headerHex := "80"
	payloadWithTS := tsHex + headerHex + "abcd" // ts + GHDR + rest

	devEui := "0011223344556677"
	fcnt := 321

	// Response includes GNSS + Timestamp feature
	resp := decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS, decoder.FeatureTimestamp}, &fakeGNSSData{
		lat: 47.1, lon: 8.1, alt: 12.0, ts: &ts,
	})

	cap := &captureSolverV2{resp: resp}

	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log,
		WithSolverV2(cap),
	)

	ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
	ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

	// Port 210 -> steady, timestamp present and should be stripped before solve
	out, err := dec.Decode(ctx, payloadWithTS, 210)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedForwarded := headerHex + "abcd"
	if cap.lastPayload != expectedForwarded {
		t.Fatalf("expected forwarded payload %q, got %q", expectedForwarded, cap.lastPayload)
	}
	if cap.lastOptions.Port != 192 {
		t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
	}
	if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != false {
		t.Fatalf("expected Moving=false for port 210, got %+v", cap.lastOptions.Moving)
	}
	if cap.lastOptions.Timestamp == nil || !cap.lastOptions.Timestamp.Equal(ts) {
		t.Fatalf("expected Timestamp=%v for port 210, got %+v", ts, cap.lastOptions.Timestamp)
	}
	if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) {
		t.Fatalf("expected GNSS and Timestamp features in result")
	}

	// Port 211 -> moving, timestamp present and should be stripped before solve
	out, err = dec.Decode(ctx, payloadWithTS, 211)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.lastPayload != expectedForwarded {
		t.Fatalf("expected forwarded payload %q, got %q", expectedForwarded, cap.lastPayload)
	}
	if cap.lastOptions.Port != 192 {
		t.Fatalf("expected port 192, got %d", cap.lastOptions.Port)
	}
	if cap.lastOptions.Moving == nil || *cap.lastOptions.Moving != true {
		t.Fatalf("expected Moving=true for port 211, got %+v", cap.lastOptions.Moving)
	}
	if cap.lastOptions.Timestamp == nil || !cap.lastOptions.Timestamp.Equal(ts) {
		t.Fatalf("expected Timestamp=%v for port 211, got %+v", ts, cap.lastOptions.Timestamp)
	}
	if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) {
		t.Fatalf("expected GNSS and Timestamp features in result")
	}
}

func TestGNSS_SolverV2_194_195_TimestampStrippedAndPassed(t *testing.T) {
	log := newLogger()

	// Build a payload with 4B timestamp prefix + header 0x80 after
	secs := uint32(1750000000) // some fixed time
	ts := time.Unix(int64(secs), 0).UTC()
	tsBytes := make([]byte, 4)
	// big-endian
	tsBytes[0] = byte((secs >> 24) & 0xff)
	tsBytes[1] = byte((secs >> 16) & 0xff)
	tsBytes[2] = byte((secs >> 8) & 0xff)
	tsBytes[3] = byte(secs & 0xff)

	tsHex := hex.EncodeToString(tsBytes)
	headerHex := "80"
	payloadWithTS := tsHex + headerHex + "abcd" // ts + GHDR + rest

	devEui := "0011223344556677"
	fcnt := 321

	// Response includes GNSS + Timestamp feature
	resp := decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureGNSS, decoder.FeatureTimestamp}, &fakeGNSSData{
		lat: 47.1, lon: 8.1, alt: 12.0, ts: &ts,
	})

	cap := &captureSolverV2{resp: resp}

	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, log,
		WithSolverV2(cap),
	)

	ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, devEui)
	ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, fcnt)

	// Port 194 -> steady, timestamp present and should be stripped before solve
	out, err := dec.Decode(ctx, payloadWithTS, 194)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedForwarded := headerHex + "abcd"
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

	// Port 195 -> moving, timestamp present and should be stripped before solve
	out, err = dec.Decode(ctx, payloadWithTS, 195)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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
