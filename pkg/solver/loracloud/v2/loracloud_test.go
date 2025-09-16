package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"go.uber.org/zap"
)

// mockNonNestedServer returns an httptest server that emulates the non-Traxmate LoRaCloud response:
// { "result": { ... UplinkMsgResponse.Result ... } }
func mockNonNestedServer(t *testing.T, positionValid bool, lat, lon float64, captureUTC float64, statusCode int) *httptest.Server {
	t.Helper()

	type resultStruct struct {
		Deveui          string `json:"deveui"`
		PendingRequests struct {
			Requests []any `json:"requests"`
			ID       int   `json:"id"`
			Updelay  int   `json:"updelay"`
			Upcount  int   `json:"upcount"`
		} `json:"pending_requests"`
		InfoFields  struct{} `json:"info_fields"`
		LogMessages []any    `json:"log_messages"`
		Fports      struct {
			Dmport     int `json:"dmport"`
			Gnssport   int `json:"gnssport"`
			Wifiport   int `json:"wifiport"`
			Fragport   int `json:"fragport"`
			Streamport int `json:"streamport"`
			Gnssngport int `json:"gnssngport"`
		} `json:"fports"`
		Dnlink            any   `json:"dnlink"`
		FulfilledRequests []any `json:"fulfilled_requests"`
		CancelledRequests []any `json:"cancelled_requests"`
		File              any   `json:"file"`
		StreamRecords     any   `json:"stream_records"`
		PositionSolution  struct {
			Llh             []float64 `json:"llh"`
			Accuracy        float64   `json:"accuracy"`
			Ecef            []float64 `json:"ecef"`
			Gdop            float64   `json:"gdop"`
			CaptureTimeGps  float64   `json:"capture_time_gps"`
			CaptureTimeUtc  float64   `json:"capture_time_utc"`
			CaptureTimesGps []float64 `json:"capture_times_gps"`
			CaptureTimesUtc []float64 `json:"capture_times_utc"`
			Timestamp       float64   `json:"timestamp"`
			AlgorithmType   string    `json:"algorithm_type"`
		} `json:"position_solution"`
		Operation string `json:"operation"`
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/device/send" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		resp := struct {
			Result resultStruct `json:"result"`
		}{}

		// minimal fields
		resp.Result.Deveui = "00-11-22-33-44-55-66-77"
		resp.Result.Operation = "gnss"
		if positionValid {
			resp.Result.PositionSolution.Llh = []float64{lat, lon, 10}
			resp.Result.PositionSolution.CaptureTimeUtc = captureUTC
			resp.Result.PositionSolution.Accuracy = 5
			resp.Result.PositionSolution.Gdop = 1.2
			resp.Result.PositionSolution.Timestamp = captureUTC
			resp.Result.PositionSolution.AlgorithmType = "gnss"
		} else {
			// invalid: zero coords, no capture time
			resp.Result.PositionSolution.Llh = []float64{0, 0, 0}
			resp.Result.PositionSolution.CaptureTimeUtc = 0
			resp.Result.PositionSolution.Accuracy = 0
			resp.Result.PositionSolution.Gdop = 0
			resp.Result.PositionSolution.Timestamp = 0
			resp.Result.PositionSolution.AlgorithmType = "gnss"
		}

		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func newLogger(t *testing.T) *zap.Logger {
	t.Helper()
	logger, _ := zap.NewDevelopment()
	return logger
}

// Helper to run Solve with a given server and options
func runSolve(t *testing.T, srv *httptest.Server, opts solver.SolverV2Options, payload string, clientOpts ...LoracloudClientOptions) (*decoder.DecodedUplink, error) {
	t.Helper()
	ctx := context.Background()
	logger := newLogger(t)

	// Force BaseUrl to test server
	clientOpts = append(clientOpts, WithBaseUrl(srv.URL))

	c, err := NewLoracloudClient(ctx, "Bearer test-token", logger, clientOpts...)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return c.Solve(ctx, payload, opts)
}

func TestSolve_Valid_NoOptionalFeatures(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	// First byte 0x80 -> EndOfGroup=true
	payload := "80"

	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 10,
		Port:          150,
		// No timestamp, no moving
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	// Expect GNSS feature only
	if !out.Is(decoder.FeatureGNSS) {
		t.Fatalf("expected FeatureGNSS to be set")
	}
	if out.Is(decoder.FeatureTimestamp) || out.Is(decoder.FeatureMoving) || out.Is(decoder.FeatureBuffered) {
		t.Fatalf("unexpected optional features set")
	}

	// Data should implement GNSS interface
	gnss, ok := out.Data.(decoder.UplinkFeatureGNSS)
	if !ok {
		t.Fatalf("data does not implement UplinkFeatureGNSS")
	}
	if gnss.GetLatitude() == 0 || gnss.GetLongitude() == 0 {
		t.Fatalf("unexpected zero GNSS coordinates")
	}
}

func TestSolve_WithTimestampBuffered(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	// EndOfGroup
	payload := "80"

	// 2 minutes ago, default threshold is 1 minute => buffered
	ts := time.Now().Add(-2 * time.Minute)

	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 11,
		Port:          150,
		Timestamp:     &ts,
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureTimestamp) || !out.Is(decoder.FeatureBuffered) {
		t.Fatalf("expected GNSS, Timestamp, and Buffered features to be set")
	}

	// Interfaces
	tsIF, ok := out.Data.(decoder.UplinkFeatureTimestamp)
	if !ok || tsIF.GetTimestamp() == nil {
		t.Fatalf("expected UplinkFeatureTimestamp to be implemented")
	}

	bufIF, ok := out.Data.(decoder.UplinkFeatureBuffered)
	if !ok || !bufIF.IsBuffered() {
		t.Fatalf("expected UplinkFeatureBuffered to be implemented and buffered")
	}
}

func TestSolve_WithMoving(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	payload := "80"
	mv := true

	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 12,
		Port:          150,
		Moving:        &mv,
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	if !out.Is(decoder.FeatureGNSS) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected GNSS and Moving features to be set")
	}
	mvIF, ok := out.Data.(decoder.UplinkFeatureMoving)
	if !ok || !mvIF.IsMoving() {
		t.Fatalf("expected UplinkFeatureMoving to be implemented and true")
	}
}

func TestSolve_InvalidDevEui_Wrapped(t *testing.T) {
	// server won't be reached, but provide a dummy to satisfy client creation
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	payload := "80"
	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "001122334455667", // 15 chars, invalid
		UplinkCounter: 1,
		Port:          1,
	}, payload)
	if err == nil {
		t.Fatalf("expected error for invalid DevEUI")
	}
	if !errors.Is(err, ErrInvalidOptions) {
		t.Fatalf("expected ErrInvalidOptions, got %v", err)
	}
	if !errors.Is(err, ErrInvalidDevEui) {
		t.Fatalf("expected ErrInvalidDevEui in wrap chain, got %v", err)
	}
}

func TestSolve_ResponseDecodeError_Wrapped(t *testing.T) {
	// Server returns invalid JSON body for decode
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/device/send" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{not-json"))
	}))
	defer srv.Close()

	payload := "80"
	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 13,
		Port:          150,
	}, payload)
	if err == nil {
		t.Fatalf("expected decode error")
	}
	if !errors.Is(err, ErrDecodeFailed) {
		t.Fatalf("expected ErrDecodeFailed, got %v", err)
	}
}

func TestSolve_PositionInvalid_Wrapped(t *testing.T) {
	// invalid position (0,0) and missing timestamp leads v1 to ErrPositionResolutionIsEmpty due to EndOfGroup
	srv := mockNonNestedServer(t, false, 0, 0, 0, http.StatusOK)
	defer srv.Close()

	payload := "80"
	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 14,
		Port:          150,
	}, payload)
	if err == nil {
		t.Fatalf("expected position invalid error")
	}
	if !errors.Is(err, ErrPositionInvalid) {
		t.Fatalf("expected ErrPositionInvalid, got %v", err)
	}
}

func TestSolve_BufferedThreshold_Configurable(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	payload := "80"
	// threshold 5 minutes
	threshold := 5 * time.Minute

	// 3 minutes ago -> NOT buffered
	tsRecent := time.Now().Add(-3 * time.Minute)
	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 15,
		Port:          150,
		Timestamp:     &tsRecent,
	}, payload, WithBufferedThreshold(threshold))
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}
	if !out.Is(decoder.FeatureTimestamp) {
		t.Fatalf("expected FeatureTimestamp")
	}
	if out.Is(decoder.FeatureBuffered) {
		t.Fatalf("did not expect FeatureBuffered at 3m with 5m threshold")
	}

	// 6 minutes ago -> buffered
	tsOld := time.Now().Add(-6 * time.Minute)
	out2, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 16,
		Port:          150,
		Timestamp:     &tsOld,
	}, payload, WithBufferedThreshold(threshold))
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}
	if !out2.Is(decoder.FeatureBuffered) {
		t.Fatalf("expected FeatureBuffered at 6m with 5m threshold")
	}
}

// Ensure server receives a valid JSON body with deveui formatted and uplink fields and auto timestamp set
func TestServerReceivesValidRequest(t *testing.T) {
	var receivedBody bytes.Buffer
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/device/send" {
			http.NotFound(w, r)
			return
		}
		_, _ = receivedBody.ReadFrom(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"deveui": "00-11-22-33-44-55-66-77",
				"position_solution": map[string]any{
					"llh":              []float64{47, 8, 10},
					"accuracy":         5.0,
					"gdop":             1.2,
					"capture_time_utc": float64(time.Now().UTC().Unix()),
					"timestamp":        float64(time.Now().UTC().Unix()),
					"algorithm_type":   "gnss",
				},
				"operation": "gnss",
			},
		})
	}))
	defer srv.Close()

	payload := "80"
	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 17,
		Port:          150,
		// No Timestamp provided -> client should auto-set current time in uplink.timestamp
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	if receivedBody.Len() == 0 {
		t.Fatalf("expected server to receive a JSON body")
	}

	// Validate request content structure
	var req map[string]any
	if err := json.Unmarshal(receivedBody.Bytes(), &req); err != nil {
		t.Fatalf("failed to decode received request body: %v", err)
	}

	deveuiVal, ok := req["deveui"].(string)
	if !ok || len(deveuiVal) != len("00-11-22-33-44-55-66-77") {
		t.Fatalf("expected hyphenated DevEUI, got: %#v", req["deveui"])
	}

	uplink, ok := req["uplink"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'uplink' object in request")
	}
	ts, ok := uplink["timestamp"].(float64)
	if !ok || ts == 0 {
		t.Fatalf("expected non-zero uplink.timestamp to be set automatically")
	}
}

// Additional coverage tests

// Not EndOfGroup (payload header 0x00) -> v1 does not enforce position validity.
// v2 should not return error; it should simply avoid GNSS feature if position invalid,
// while still applying optional features as provided.
func TestSolve_NotEndOfGroup_NoGNSSFeature_ButOptionalApplied(t *testing.T) {
	srv := mockNonNestedServer(t, false, 0, 0, 0, http.StatusOK)
	defer srv.Close()

	payload := "00" // EndOfGroup=false
	mv := true
	now := time.Now()

	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 18,
		Port:          150,
		Timestamp:     &now,
		Moving:        &mv,
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	// GNSS should not be set due to invalid coordinates and lack of EndOfGroup enforcement
	if out.Is(decoder.FeatureGNSS) {
		t.Fatalf("did not expect FeatureGNSS for invalid position when EndOfGroup=false")
	}
	// Optional features should still be applied
	if !out.Is(decoder.FeatureTimestamp) || !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected FeatureTimestamp and FeatureMoving to be set")
	}

	// Interfaces
	if _, ok := out.Data.(decoder.UplinkFeatureTimestamp); !ok {
		t.Fatalf("expected UplinkFeatureTimestamp to be implemented")
	}
	if mvIF, ok := out.Data.(decoder.UplinkFeatureMoving); !ok || !mvIF.IsMoving() {
		t.Fatalf("expected UplinkFeatureMoving implemented and true")
	}
}

// Moving + Timestamp buffered => all optional interfaces on, plus GNSS when valid.
func TestSolve_MovingAndTimestampBuffered_BothInterfaces(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	payload := "80"
	oldTs := time.Now().Add(-10 * time.Minute)
	mv := true

	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 19,
		Port:          150,
		Timestamp:     &oldTs,
		Moving:        &mv,
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	if !(out.Is(decoder.FeatureGNSS) && out.Is(decoder.FeatureTimestamp) && out.Is(decoder.FeatureMoving) && out.Is(decoder.FeatureBuffered)) {
		t.Fatalf("expected GNSS, Timestamp, Moving, Buffered features")
	}

	if _, ok := out.Data.(decoder.UplinkFeatureBuffered); !ok {
		t.Fatalf("expected UplinkFeatureBuffered to be implemented")
	}
}

// Invalid payload should be wrapped as ErrInvalidOptions
func TestSolve_InvalidPayload_Wrapped(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 20,
		Port:          150,
	}, "")
	if err == nil {
		t.Fatalf("expected error for empty payload")
	}
	if !errors.Is(err, ErrInvalidOptions) {
		t.Fatalf("expected ErrInvalidOptions, got %v", err)
	}
}

// Unexpected status should be wrapped as ErrUnexpectedStatus
func TestSolve_UnexpectedStatus_Wrapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/device/send" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("oops"))
	}))
	defer srv.Close()

	payload := "80"
	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 21,
		Port:          150,
	}, payload)
	if err == nil {
		t.Fatalf("expected error for unexpected status")
	}
	if !errors.Is(err, ErrUnexpectedStatus) {
		t.Fatalf("expected ErrUnexpectedStatus, got %v", err)
	}
}

// New client with Semtech base should trigger shutdown warning path; ensure construction succeeds.
func TestNewClient_SemtechShutdownWarn(t *testing.T) {
	ctx := context.Background()
	logger := newLogger(t)
	_, err := NewLoracloudClient(ctx, "Bearer test-token", logger, WithBaseUrl(SemtechLoRaCloudBaseUrl))
	if err != nil {
		t.Fatalf("expected no error constructing client with Semtech base: %v", err)
	}
}

// Request failure path: invalid base URL should cause request error and be wrapped as ErrRequestFailed.
func TestSolve_RequestFailed_Wrapped(t *testing.T) {
	ctx := context.Background()
	logger := newLogger(t)
	c, err := NewLoracloudClient(ctx, "Bearer test-token", logger, WithBaseUrl("http://127.0.0.1:0"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	payload := "80"
	_, err = c.Solve(ctx, payload, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 22,
		Port:          150,
	})
	if err == nil {
		t.Fatalf("expected request failed error")
	}
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("expected ErrRequestFailed, got %v", err)
	}
}

// Moving provided as false should still set FeatureMoving and implement UplinkFeatureMoving=false.
func TestSolve_WithMovingFalse(t *testing.T) {
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	payload := "80"
	mv := false

	out, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "0011223344556677",
		UplinkCounter: 23,
		Port:          150,
		Moving:        &mv,
	}, payload)
	if err != nil {
		t.Fatalf("Solve returned error: %v", err)
	}

	if !out.Is(decoder.FeatureMoving) {
		t.Fatalf("expected FeatureMoving set")
	}
	if mvIF, ok := out.Data.(decoder.UplinkFeatureMoving); !ok || mvIF.IsMoving() {
		t.Fatalf("expected UplinkFeatureMoving implemented and false")
	}
}

// DevEui with invalid hex characters (length 16) should error via hex decode.
func TestSolve_InvalidDevEui_NonHex(t *testing.T) {
	// server dummy
	srv := mockNonNestedServer(t, true, 47.0, 8.0, float64(time.Now().UTC().Unix()), http.StatusOK)
	defer srv.Close()

	payload := "80"
	_, err := runSolve(t, srv, solver.SolverV2Options{
		DevEui:        "00112233445566ZZ", // invalid hex
		UplinkCounter: 24,
		Port:          150,
	}, payload)
	if err == nil {
		t.Fatalf("expected error for invalid hex DevEUI")
	}
	if !errors.Is(err, ErrInvalidOptions) || !errors.Is(err, ErrInvalidDevEui) {
		t.Fatalf("expected ErrInvalidOptions and ErrInvalidDevEui, got %v", err)
	}
}
