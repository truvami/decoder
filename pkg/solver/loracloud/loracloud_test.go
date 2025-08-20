package loracloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"go.uber.org/zap"
)

func startMockServer(handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

func TestPost(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := startMockServer(mux)
	middleware, err := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample(), WithBaseUrl(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer server.Close()

	// Test case 1: Successful request
	url := fmt.Sprintf("%v/success", server.URL)
	body := []byte(`{"key": "value"}`)

	response, err := middleware.post(url, body)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected response: %v, got: %v", http.StatusOK, response)
	}

	// Test case 2: Request with error
	url = fmt.Sprintf("%v/error", server.URL)
	body = []byte(`{"key": "value}`)

	response, err = middleware.post(url, body)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected response: %v, got: %v", http.StatusInternalServerError, response)
	}
}

func TestDeliverUplinkMessage(t *testing.T) {
	t.Run("Successful request", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"result": {
					"deveui": "01-23-45-67-89-AB-CD-EF",
					"pending_requests": {
						"requests": [],
						"id": 1,
						"updelay": 0,
						"upcount": 0
					},
					"info_fields": {},
					"log_messages": [],
					"fports": {
						"dmport": 1,
						"gnssport": 2,
						"wifiport": 3,
						"fragport": 4,
						"streamport": 5,
						"gnssngport": 6
					},
					"position_solution": {
							"algorithm_type": "gnssng",
							"llh": [51.49278, 53.0212, 0],
							"accuracy": 20.7,
							"gdop": 2.48,
							"capture_time_utc": 1722433373.18046
					},
					"operation": "gnss"
				}
			}`))
		})

		server := startMockServer(mux)
		middleware, err := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample(), WithBaseUrl(server.URL))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
		}

		response, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if response.Result.Deveui != "0123456789ABCDEF" {
			t.Errorf("expected deveui: %v, got: %v", "0123456789ABCDEF", response.Result.Deveui)
		}
	})

	t.Run("Validation error", func(t *testing.T) {
		server := startMockServer(nil)
		middleware, err := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample(), WithBaseUrl(server.URL))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
		}

		_, err = middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !strings.Contains(err.Error(), "error validating uplink message") {
			t.Errorf("expected validation error, got: %v", err)
		}
	})

	t.Run("Unexpected status code", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors": ["Invalid request"]}`))
		})

		server := startMockServer(mux)
		middleware, err := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample(), WithBaseUrl(server.URL))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
		}

		_, err = middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !errors.Is(err, ErrUnexpectedStatusCode) {
			t.Errorf("expected status code error, got: %v", err)
		}
	})

	t.Run("Error decoding response", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid-json`))
		})

		server := startMockServer(mux)
		middleware, err := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample(), WithBaseUrl(server.URL))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
		}

		_, err = middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !errors.Is(err, ErrDecodingResponse) {
			t.Errorf("expected decoding error, got: %v", err)
		}
	})
}

func TestResponseVariants(t *testing.T) {
	type Expected = struct {
		timestamp *time.Time
		latitude  float64
		longitude float64
		altitude  float64
	}
	var tests = []struct {
		name     string
		result   []byte
		expected Expected
		err      error
	}{
		{
			name: "normal response",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"algorithm_type": "gnssng",
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": 1722433373.18046
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{
				timestamp: common.TimePointer(1722433373.18046),
				latitude:  51.49278,
				longitude: 0.0212,
				altitude:  83.93,
			},
			err: nil,
		},
		{
			name: "llh empty array",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"algorithm_type": "gnssng",
						"llh": [],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": 1722433373.18046
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{},
			err:      ErrPositionResolutionIsEmpty,
		},
		{
			name: "captured at null and no algorithm type",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": null
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{},
			err:      ErrPositionResolutionIsEmpty,
		},
		{
			name: "captured at null and gnss ng algorithm type",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"algorithm_type": "gnssng",
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": null,
						"capture_times_utc": [1722433364.06164, 1722433373.18046, null]
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{
				timestamp: common.TimePointer(1722433373.18046),
				latitude:  51.49278,
				longitude: 0.0212,
				altitude:  83.93,
			},
		},
		{
			name: "captured at null and gnss ng algorithm type",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"algorithm_type": "gnssng",
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": 1722433364.06164,
						"capture_times_utc": [1722433364.06164, 1722433373.18046, null]
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{
				timestamp: common.TimePointer(1722433373.18046),
				latitude:  51.49278,
				longitude: 0.0212,
				altitude:  83.93,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(test.result)
			})

			server := startMockServer(mux)
			middleware, err := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample(), WithBaseUrl(server.URL))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer server.Close()

			devEui := "b2e6876e64be9692"
			uplinkMsg := UplinkMsg{
				MsgType: "uplink",
				FCount:  42,
				Port:    192,
				Payload: "8c9e50de366a460e8a70fe72e04445db95d1eca8dcdac252",
			}

			response, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.timestamp, response.GetTimestamp())
				assert.Equal(t, test.expected.latitude, response.GetLatitude())
				assert.Equal(t, test.expected.longitude, response.GetLongitude())
				assert.Equal(t, test.expected.altitude, response.GetAltitude())
			}
		})
	}
}

func TestValidateContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
	}{
		{
			name:    "missing port",
			ctx:     context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef"),
			wantErr: ErrContextPortNotFound,
		},
		{
			name: "missing devEui",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				return ctx
			}(),
			wantErr: ErrContextDevEuiNotFound,
		},
		{
			name: "missing fCount",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				return ctx
			}(),
			wantErr: ErrContextFCountNotFound,
		},
		{
			name: "invalid port type",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, 1) // int instead of uint8
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 0)
				return ctx
			}(),
			wantErr: ErrContextPortNotFound,
		},
		{
			name: "invalid devEui length",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcde") // 15 chars
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 0)
				return ctx
			}(),
			wantErr: ErrContextInvalidDevEui,
		},
		{
			name: "invalid devEui non-hex",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdeg") // 'g' is not hex
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 0)
				return ctx
			}(),
			wantErr: ErrContextInvalidDevEui,
		},
		{
			name: "invalid fCount negative",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, -1)
				return ctx
			}(),
			wantErr: ErrContextInvalidFCount,
		},
		{
			name: "valid context",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(10))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 42)
				return ctx
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContext(tt.ctx)
			if tt.wantErr == nil {
				assert.NoError(t, err, "expected no error but got: %v", err)
				return
			}

			assert.Error(t, err, "expected error but got none")
			assert.ErrorIs(t, err, tt.wantErr, "expected error to match")

		})
	}
}

func TestWithBaseUrl(t *testing.T) {
	// Create a LoracloudClient with a default BaseUrl
	client := LoracloudClient{
		BaseUrl: "https://default.url",
	}

	// Use WithBaseUrl to set a new BaseUrl
	newUrl := "https://custom.url"
	option := WithBaseUrl(newUrl)
	option(&client)

	if client.BaseUrl != newUrl {
		t.Errorf("expected BaseUrl to be %q, got %q", newUrl, client.BaseUrl)
	}
}

func TestIsSemtechLoRaCloudShutdown(t *testing.T) {
	shutdownDate := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		baseUrl     string
		currentTime time.Time
		err         error
	}{
		{
			name:        "after Semtech shutdown date - non-Semtech base URL",
			baseUrl:     TraxmateLoRaCloudBaseUrl,
			currentTime: shutdownDate.Add(time.Hour),
			err:         nil,
		},
		{
			name:        "before Semtech shutdown date",
			baseUrl:     SemtechLoRaCloudBaseUrl,
			currentTime: shutdownDate.Add(-time.Hour),
			err:         nil,
		},
		{
			name:        "exactly at Semtech shutdown date",
			baseUrl:     SemtechLoRaCloudBaseUrl,
			currentTime: shutdownDate,
			err:         nil,
		},
		{
			name:        "after Semtech shutdown date",
			baseUrl:     SemtechLoRaCloudBaseUrl,
			currentTime: shutdownDate.Add(time.Hour),
			err:         ErrSemtechLoRaCloudShutdown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTimeNow := func() time.Time {
				return tt.currentTime
			}

			client, err := NewLoracloudClient(
				context.Background(),
				"access_token",
				zap.NewNop(),
				WithBaseUrl(tt.baseUrl),
				WithTimeNow(mockTimeNow),
			)

			if tt.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				err = client.isSemtechLoRaCloudShutdown()
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithTimeNow(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	mockTimeNow := func() time.Time {
		return fixedTime
	}

	client := LoracloudClient{}
	option := WithTimeNow(mockTimeNow)
	option(&client)

	// Test that the time function was set correctly
	assert.Equal(t, fixedTime, client.timeNow())
}
