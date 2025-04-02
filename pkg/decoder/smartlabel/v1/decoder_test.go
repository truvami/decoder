package smartlabel

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
		expected    interface{}
		expectedErr string
	}{
		{
			payload: "0ca90dbd",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.241,
				PhotovoltaicVoltage: 3.517,
			},
		},
		{
			payload: "0dfa0e1a",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.578,
				PhotovoltaicVoltage: 3.610,
			},
		},
		{
			payload: "0e860ef7",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.718,
				PhotovoltaicVoltage: 3.831,
			},
		},
		{
			payload: "0f501079",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.920,
				PhotovoltaicVoltage: 4.217,
			},
		},
		{
			payload: "07fa69",
			port:    2,
			expected: Port2Payload{
				Temperature: 20.42,
				Humidity:    52.5,
			},
		},
		{
			payload: "074070",
			port:    2,
			expected: Port2Payload{
				Temperature: 18.56,
				Humidity:    56.0,
			},
		},
		{
			payload: "06947d",
			port:    2,
			expected: Port2Payload{
				Temperature: 16.84,
				Humidity:    62.5,
			},
		},
		{
			payload: "04da8d",
			port:    2,
			expected: Port2Payload{
				Temperature: 12.42,
				Humidity:    70.5,
			},
		},
		{
			payload: "00012c00780f000a03e8003c0078f81c01000402",
			port:    4,
			expected: Port4Payload{
				DataRate:                   0,
				Acceleration:               false,
				Wifi:                       false,
				Gnss:                       false,
				SteadyInterval:             300,
				MovingInterval:             120,
				HeartbeatInterval:          15,
				AccelerationThreshold:      10,
				AccelerationDelay:          1000,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  120,
				TemperatureLowerThreshold:  -8,
				TemperatureUpperThreshold:  +28,
				AccessPointsThreshold:      1,
				FirmwareVersionMajor:       0,
				FirmwareVersionMinor:       4,
				FirmwareVersionPatch:       2,
			},
		},
		{
			payload: "2a0258012c1e006403e8003c012cf61e02010307",
			port:    4,
			expected: Port4Payload{
				DataRate:                   2,
				Acceleration:               true,
				Wifi:                       false,
				Gnss:                       true,
				SteadyInterval:             600,
				MovingInterval:             300,
				HeartbeatInterval:          30,
				AccelerationThreshold:      100,
				AccelerationDelay:          1000,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  300,
				TemperatureLowerThreshold:  -10,
				TemperatureUpperThreshold:  +30,
				AccessPointsThreshold:      2,
				FirmwareVersionMajor:       1,
				FirmwareVersionMinor:       3,
				FirmwareVersionPatch:       7,
			},
		},
		{
			payload: "1b04b002583c012c05dc003c0258f12304020400",
			port:    4,
			expected: Port4Payload{
				DataRate:                   3,
				Acceleration:               true,
				Wifi:                       true,
				Gnss:                       false,
				SteadyInterval:             1200,
				MovingInterval:             600,
				HeartbeatInterval:          60,
				AccelerationThreshold:      300,
				AccelerationDelay:          1500,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  600,
				TemperatureLowerThreshold:  -15,
				TemperatureUpperThreshold:  +35,
				AccessPointsThreshold:      4,
				FirmwareVersionMajor:       2,
				FirmwareVersionMinor:       4,
				FirmwareVersionPatch:       0,
			},
		},
		{
			payload: "3f0e1007087801c207d0003c04b0ec280603020c",
			port:    4,
			expected: Port4Payload{
				DataRate:                   7,
				Acceleration:               true,
				Wifi:                       true,
				Gnss:                       true,
				SteadyInterval:             3600,
				MovingInterval:             1800,
				HeartbeatInterval:          120,
				AccelerationThreshold:      450,
				AccelerationDelay:          2000,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  1200,
				TemperatureLowerThreshold:  -20,
				TemperatureUpperThreshold:  +40,
				AccessPointsThreshold:      6,
				FirmwareVersionMajor:       3,
				FirmwareVersionMinor:       2,
				FirmwareVersionPatch:       12,
			},
		},
		{
			payload:  "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:     192,
			devEui:   "10CE45FFFE00C7EC",
			expected: &exampleResponse,
		},
		{
			payload:     "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:        192,
			devEui:      "10CE45FFFE00C7ED",
			expected:    nil,
			expectedErr: "invalid character 'j' looking for beginning of value",
		},
		{
			payload:     "00",
			port:        0,
			devEui:      "",
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			payload: "0e021e0384012c01003c012c03e8",
			port:    11,
			expected: Port11ConfigurationPayload{
				Flags:              30,
				GNSSEnabled:        1,
				WiFiEnabled:        1,
				AccEnabled:         1,
				StaticSF:           "SF9",
				SteadyIntervalS:    900,
				MovingIntervalS:    300,
				HeartbeatIntervalH: 1,
				LEDBlinkIntervalS:  60,
				AccThresholdMS:     300,
				AccDelayMS:         1000,
			},
		},
		{
			payload: "11021e0384012c01003c012c03e8e43420ea",
			port:    11,
			expected: Port11ConfigurationPayload{
				Flags:              30,
				GNSSEnabled:        1,
				WiFiEnabled:        1,
				AccEnabled:         1,
				StaticSF:           "SF9",
				SteadyIntervalS:    900,
				MovingIntervalS:    300,
				HeartbeatIntervalH: 1,
				LEDBlinkIntervalS:  60,
				AccThresholdMS:     300,
				AccDelayMS:         1000,
				GitHash:            "e43420ea",
			},
		},
		{
			payload: "0a010f05095f4100000000",
			port:    11,
			expected: Port11HeartbeatPayload{
				Battery:        3.845,
				Temperature:    23.99,
				RH:             32.5,
				GNSSScansCount: 0,
				WiFiScansCount: 0,
			},
		},
		{
			payload:     "",
			port:        11,
			devEui:      "",
			expected:    nil,
			expectedErr: "data length is less than 2",
		},
		{
			payload:     "0e02f9f3eae48ae7523948d0d5xx",
			port:        11,
			devEui:      "",
			expected:    nil,
			expectedErr: "encoding/hex",
		},
		{
			payload:     "ff00",
			port:        11,
			devEui:      "",
			expected:    nil,
			expectedErr: "invalid payload for port 11",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewSmartLabelv1Decoder(middleware, WithFCount(1))
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

func TestInvalidPort(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, err := decoder.Decode("0ff0", 1, "")

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, err := decoder.Decode("0ff00ff00ff0", 1, "")

	if err == nil || err.Error() != "payload too long" {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload string
		port    int16
	}{
		{
			payload: "0f501079",
			port:    1,
		},
		{
			payload: "04da8d",
			port:    2,
		},
		{
			payload: "3f0e1007087801c207d0003c04b0ec280603020c",
			port:    4,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
			data, _ := d.Decode(test.payload, test.port, "")

			// should be able to decode base feature
			base, ok := data.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", data)
			}
			// check if it panics
			base.GetTimestamp()

			if data.Is(decoder.FeatureBattery) {
				batteryVoltage, ok := data.Data.(decoder.UplinkFeatureBattery)
				if !ok {
					t.Fatalf("expected UplinkFeatureBattery, got %T", data)
				}
				if batteryVoltage.GetBatteryVoltage() == 0 {
					t.Fatalf("expected non zero battery voltage")
				}
			}
			if data.Is(decoder.FeaturePhotovoltaic) {
				photovoltaicVoltage, ok := data.Data.(decoder.UplinkFeaturePhotovoltaic)
				if !ok {
					t.Fatalf("expected UplinkFeaturePhotovoltaic, got %T", data)
				}
				if photovoltaicVoltage.GetPhotovoltaicVoltage() == 0 {
					t.Fatalf("expected non zero photovoltaic voltage")
				}
			}
			if data.Is(decoder.FeatureTemperature) {
				temperature, ok := data.Data.(decoder.UplinkFeatureTemperature)
				if !ok {
					t.Fatalf("expected UplinkFeatureTemperature, got %T", data)
				}
				if temperature.GetTemperature() == 0 {
					t.Fatalf("expected non zero temperature")
				}
			}
			if data.Is(decoder.FeatureHumidity) {
				humidity, ok := data.Data.(decoder.UplinkFeatureHumidity)
				if !ok {
					t.Fatalf("expected UplinkFeatureHumidity, got %T", data)
				}
				if humidity.GetHumidity() == 0 {
					t.Fatalf("expected non zero humidity")
				}
			}
			if data.Is(decoder.FeatureConfig) {
				config, ok := data.Data.(decoder.UplinkFeatureConfig)
				if !ok {
					t.Fatalf("expected UplinkFeatureConfig, got %T", data)
				}
				// call functions to check if it panics
				config.GetBle()
				config.GetGnss()
				config.GetWifi()
				config.GetAcceleration()
				config.GetMovingInterval()
				config.GetSteadyInterval()
				config.GetConfigInterval()
				config.GetGnssTimeout()
				config.GetAccelerometerThreshold()
				config.GetAccelerometerDelay()
				config.GetBatteryInterval()
				config.GetRejoinInterval()
				config.GetLowLightThreshold()
				config.GetHighLightThreshold()
				config.GetLowTemperatureThreshold()
				config.GetHighTemperatureThreshold()
				config.GetAccessPointsThreshold()
				config.GetBatchSize()
				config.GetBufferSize()
			}
			if data.Is(decoder.FeatureFirmwareVersion) {
				firmwareVersion, ok := data.Data.(decoder.UplinkFeatureFirmwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureFirmwareVersion, got %T", data)
				}
				if firmwareVersion.GetFirmwareVersion() == "" {
					t.Fatalf("expected non empty firmware version")
				}
			}
		})
	}
}

func TestWithAutoPadding(t *testing.T) {
	middleware := loracloud.NewLoracloudMiddleware("access_token")

	decoder := NewSmartLabelv1Decoder(
		middleware,
		WithAutoPadding(true),
	)

	// Type assert to access the internal field
	if d, ok := decoder.(*SmartLabelv1Decoder); ok {
		if !d.autoPadding {
			t.Error("expected autoPadding to be true")
		}
	} else {
		t.Error("failed to type assert decoder")
	}
}

func TestGetPort11PayloadType(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		want        string
		expectedErr string
	}{
		{
			name:        "Empty data",
			data:        "",
			want:        "",
			expectedErr: "data length is less than 2",
		},
		{
			name:        "Single byte",
			data:        "0",
			want:        "",
			expectedErr: "data length is less than 2",
		},
		{
			name:        "Configuration payload 0E",
			data:        "0E",
			want:        "configuration",
			expectedErr: "",
		},
		{
			name:        "Configuration payload 11",
			data:        "11",
			want:        "configuration",
			expectedErr: "",
		},
		{
			name:        "Heartbeat payload",
			data:        "0A",
			want:        "heartbeat",
			expectedErr: "",
		},
		{
			name:        "Invalid payload",
			data:        "FF",
			want:        "",
			expectedErr: "invalid payload for port 11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPort11PayloadType(tt.data)
			if err != nil && tt.expectedErr == "" {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectedErr != "" {
				t.Errorf("expected error containing %q, got nil", tt.expectedErr)
			}
			if err != nil && tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectedErr, err.Error())
			}
			if got != tt.want {
				t.Errorf("getPort11PayloadType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFCount(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"), WithFCount(123))

	// cast to SmartLabelv1Decoder to access fCount
	tagXLv1Decoder := decoder.(*SmartLabelv1Decoder)
	if tagXLv1Decoder.fCount != 123 {
		t.Fatalf("expected fCount to be 123, got %v", tagXLv1Decoder.fCount)
	}
}
