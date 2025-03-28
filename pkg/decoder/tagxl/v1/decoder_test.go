package tagxl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/loracloud"
)

func startMockServer() *httptest.Server {
	server := httptest.NewServer(nil)
	return server
}

func TestDecode(t *testing.T) {

	http.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// check if request body contains 10CE45FFFE00C7ED
		bodyString, _ := io.ReadAll(r.Body)
		if strings.Contains(string(bodyString), "10CE45FFFE00C7ED") {
			_, _ = w.Write([]byte("{\"invalid\": json}"))
			return
		}

		// get file from testdata
		file, err := os.Open("./response.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, _ = w.Write(data)
	})

	server := startMockServer()
	middleware := loracloud.NewLoracloudMiddleware("access_token")
	middleware.BaseUrl = server.URL
	defer server.Close()

	f, _ := os.Open("./response.json")
	var exampleResponse loracloud.UplinkMsgResponse
	d, _ := io.ReadAll(f)
	err := json.Unmarshal(d, &exampleResponse)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		payload     string
		port        int16
		devEui      string
		autoPadding bool
		expected    interface{}
		expectedErr string
	}{
		{
			payload:     "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:        192,
			devEui:      "10CE45FFFE00C7EC",
			autoPadding: false,
			expected:    &exampleResponse,
		},
		{
			payload:     "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:        192,
			devEui:      "10CE45FFFE00C7ED",
			autoPadding: false,
			expected:    nil,
			expectedErr: "invalid character 'j' looking for beginning of value",
		},
		{
			payload:     "00",
			port:        0,
			devEui:      "",
			autoPadding: false,
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			payload:     "010b0266acbcf0000000000756",
			port:        152,
			autoPadding: false,
			expected: Port152Payload{
				NewRotationState:  2,
				OldRotationState:  0,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 3, 12, 0, time.UTC).Unix()),
				NumberOfRotations: 0,
				ElapsedSeconds:    1878,
			},
		},
		{
			payload:     "10b0266acbcf0000000000756",
			port:        152,
			autoPadding: true,
			expected: Port152Payload{
				NewRotationState:  2,
				OldRotationState:  0,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 3, 12, 0, time.UTC).Unix()),
				NumberOfRotations: 0,
				ElapsedSeconds:    1878,
			},
		},
		{
			payload:     "010b1066acbe0c00a200000087",
			port:        152,
			autoPadding: false,
			expected: Port152Payload{
				NewRotationState:  0,
				OldRotationState:  1,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 7, 56, 0, time.UTC).Unix()),
				NumberOfRotations: 16.2,
				ElapsedSeconds:    135,
			},
		},
		{
			payload:     "10b1066acbe0c00a200000087",
			port:        152,
			autoPadding: true,
			expected: Port152Payload{
				NewRotationState:  0,
				OldRotationState:  1,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 7, 56, 0, time.UTC).Unix()),
				NumberOfRotations: 16.2,
				ElapsedSeconds:    135,
			},
		},
		// {
		// 	payload: "4C050145020B8C",
		// 	port:    151,
		// 	expected: Port151Payload{
		// 		Battery: 2.956,
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(middleware, WithAutoPadding(test.autoPadding), WithFCount(1))
			got, err := decoder.Decode(test.payload, test.port, test.devEui)
			if err != nil && len(test.expectedErr) == 0 {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != nil && !reflect.DeepEqual(got.Data, test.expected) && len(test.expectedErr) == 0 {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}

			if len(test.expectedErr) > 0 && err != nil && !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("expected error: %v, got: %v", test.expectedErr, err)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected error
	}{}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vValidationWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
			got, err := decoder.Decode(test.payload, test.port, "")

			if err == nil && test.expected == nil {
				return
			}

			t.Logf("got %v", got)

			if err != nil && test.expected == nil || err == nil || err.Error() != test.expected.Error() {
				t.Errorf("expected: %v\ngot: %v", test.expected, err)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("deadbeef", 152, "")

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242", 152, "")

	if err == nil || err.Error() != "payload too long" {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload        string
		port           int16
		skipValidation bool
	}{
		{
			payload: "010b0266acbcf0000000000756",
			port:    152,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewTagXLv1Decoder(
				loracloud.NewLoracloudMiddleware("apiKey"),
				WithSkipValidation(test.skipValidation),
			)
			decodedPayload, _ := d.Decode(test.payload, test.port, "")

			// should be able to decode base feature
			base, ok := decodedPayload.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", decodedPayload)
			}
			// check if it panics
			base.GetTimestamp()

			if decodedPayload.Is(decoder.FeatureGNSS) {
				gnss, ok := decodedPayload.Data.(decoder.UplinkFeatureGNSS)
				if !ok {
					t.Fatalf("expected UplinkFeatureGNSS, got %T", decodedPayload)
				}
				if gnss.GetLatitude() == 0 {
					t.Fatalf("expected non zero latitude")
				}
				if gnss.GetLongitude() == 0 {
					t.Fatalf("expected non zero longitude")
				}
				if gnss.GetAltitude() == 0 {
					t.Fatalf("expected non zero altitude")
				}
				// call function to check if it panics
				gnss.GetAltitude()
				gnss.GetPDOP()
				gnss.GetSatellites()
				gnss.GetTTF()
			}
			if decodedPayload.Is(decoder.FeatureBuffered) {
				buffered, ok := decodedPayload.Data.(decoder.UplinkFeatureBuffered)
				if !ok {
					t.Fatalf("expected UplinkFeatureBuffered, got %T", decodedPayload)
				}
				// call function to check if it panics
				buffered.GetBufferLevel()
			}
			if decodedPayload.Is(decoder.FeatureBattery) {
				batteryVoltage, ok := decodedPayload.Data.(decoder.UpLinkFeatureBattery)
				if !ok {
					t.Fatalf("expected UplinkFeatureBattery, got %T", decodedPayload)
				}
				if batteryVoltage.GetBatteryVoltage() == 0 {
					t.Fatalf("expected non zero battery voltage")
				}
			}
			if decodedPayload.Is(decoder.FeatureWiFi) {
				wifi, ok := decodedPayload.Data.(decoder.UplinkFeatureWiFi)
				if !ok {
					t.Fatalf("expected UplinkFeatureWiFi, got %T", decodedPayload)
				}
				if wifi.GetAccessPoints() == nil {
					t.Fatalf("expected non nil access points")
				}
			}
			if decodedPayload.Is(decoder.FeatureMoving) {
				moving, ok := decodedPayload.Data.(decoder.UplinkFeatureMoving)
				if !ok {
					t.Fatalf("expected UplinkFeatureMoving, got %T", decodedPayload)
				}
				// call function to check if it panics
				moving.IsMoving()
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected []string
	}{
		{
			payload:  "010b0266acbcf0000000000756",
			port:     152,
			expected: []string{"\"timestamp\": 1722596592", "\"elapsedSeconds\": 1878"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))

			data, _ := decoder.Decode(test.payload, test.port, "")

			marshaled, err := json.MarshalIndent(map[string]interface{}{
				"data":     data.Data,
				"metadata": data.Metadata,
			}, "", "   ")

			if err != nil {
				t.Fatalf("marshalling json failed because %s", err)
			}

			t.Logf("%s\n", marshaled)

			for _, value := range test.expected {
				fmt.Printf("value:%s\n", value)
				if !strings.Contains(string(marshaled), value) {
					t.Fatalf("expected to find %s\n", value)
				}
			}
		})
	}
}

func TestWithFCount(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"), WithFCount(123))

	// cast to TagXLv1Decoder to access fCount
	tagXLv1Decoder := decoder.(*TagXLv1Decoder)
	if tagXLv1Decoder.fCount != 123 {
		t.Fatalf("expected fCount to be 123, got %v", tagXLv1Decoder.fCount)
	}
}
