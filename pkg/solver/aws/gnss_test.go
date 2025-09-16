package aws

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExampleSolve(t *testing.T) {
	const allowedErrorMargin = 10 * time.Minute

	// tests := []struct {
	// 	hexPayload    string
	// 	receivedAt    time.Time
	// 	expectedError error
	// 	expectedTime  time.Time
	// }{
	// 	{
	// 		hexPayload:    "90b30df6160cc1805a818ced337c546d02eaa3786da44a9f009d576a5cec7f7ce76ba541f84bdd86ef93f468fda6ab03",
	// 		receivedAt:    time.Date(2025, time.July, 4, 14, 49, 16, 0, time.UTC), // 2025-07-04 14:49:16
	// 		expectedError: nil,
	// 		expectedTime:  time.Date(2025, time.July, 4, 14, 43, 34, 0, time.UTC), // 2025-07-04 14:43:34
	// 	},
	// 	{
	// 		hexPayload:    "85abf56669e2160cc5805ad84fa786edd04bc025c78f9f119bb6ec0a546a98b952b95eab9b76b9cb96120106e89b92e01c00",
	// 		receivedAt:    time.Date(2025, time.July, 4, 9, 33, 39, 0, time.UTC), // 2025-07-04 11:33:39
	// 		expectedError: nil,
	// 		expectedTime:  time.Date(2025, time.July, 4, 9, 16, 1, 0, time.UTC), // 2025-07-04 11:16:01
	// 	},
	// }

	tests := []struct {
		uplinks       []GNSSCapture
		expectedTime  time.Time
		expectedError error
	}{
		// {
		// 	uplinks: []GNSSCapture{
		// 		{
		// 			HexPayload: "85abf56669e2160cc5805ad84fa786edd04bc025c78f9f119bb6ec0a546a98b952b95eab9b76b9cb96120106e89b92e01c00",
		// 			ReceivedAt: time.Date(2025, time.July, 4, 9, 33, 39, 0, time.UTC), // 2025-07-04 11:33:39
		// 		},
		// 		{
		// 			HexPayload: "86ab856155e3160bc5805a018047e6db16741024a7c4a16943644e2b555a80e8a4dcf8a9a792bd7756c366e4e94a52a03400",
		// 			ReceivedAt: time.Date(2025, time.July, 4, 9, 26, 10, 0, time.UTC), // 2025-07-04 11:26:10
		// 		},
		// 		{
		// 			HexPayload: "06ab01ec4ce3960bc580da0320428301d12b314db503c9504313245c565a2865a909a26b0aa63436621a1b7cc75b5bf4e901",
		// 			ReceivedAt: time.Date(2025, time.July, 4, 9, 26, 5, 0, time.UTC), // 2025-07-04 11:26:05
		// 		},
		// 	},
		// 	expectedTime:  time.Date(2025, time.July, 4, 9, 16, 1, 0, time.UTC), // 2025-07-04 11:16:01
		// 	expectedError: nil,
		// },
		// {
		// 	uplinks: []GNSSCapture{
		// 		{
		// 			HexPayload: "07abf1a0c9e3160bc5805a91049830bf4af3b25bb06d4447402239dc566a6cc980992bc59316a6231e956481e6a61f1a7e03",
		// 			ReceivedAt: time.Date(2025, time.July, 4, 9, 34, 3, 0, time.UTC), // 2025-07-04 11:34:03
		// 		},
		// 		{
		// 			HexPayload: "87ab052ed2e3960ac5006b40e063ad4ad1efe249b0341fdb7e1b491d493157428ed0d8e0537be45d14ec252943c589e7ae00",
		// 			ReceivedAt: time.Date(2025, time.July, 4, 9, 34, 8, 0, time.UTC), // 2025-07-04 11:34:08
		// 		},
		// 	},
		// 	expectedTime:  time.Date(2025, time.July, 4, 9, 16, 1, 0, time.UTC), // 2025-07-04 11:16:01
		// 	expectedError: nil,
		// },
		{
			uplinks: []GNSSCapture{
				{
					HexPayload: "98abc1e9c5f51ad71802eb020af4ae7c155b8480b2c90b6eb22ad49e3e6521e089d0d078cc92dc7ba74b5d440ed41c584200",
					ReceivedAt: time.Date(2025, time.July, 11, 12, 43, 4, 0, time.UTC), // 2025-07-11 14:43:04
				},
				{
					HexPayload: "18bbc55f5caf798d21b02ea03cedca57558ef8d76d246d366b75aca9546ff99b080d8dc62cdbb775509f8a01d91ccd9605",
					ReceivedAt: time.Date(2025, time.July, 11, 12, 42, 52, 0, time.UTC), // 2025-07-11 14:42:52
				},
			},
			expectedTime:  time.Date(2025, time.July, 11, 12, 42, 39, 0, time.UTC), // 2025-07-11 14:42:39
			expectedError: nil,
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("TestSolveCapturedAt#%d", index), func(t *testing.T) {
			extractedTime, err := SolveCapturedAt(test.uplinks)
			if test.expectedError != nil {
				assert.Error(t, err, "expected error to match")
				assert.Equal(t, test.expectedError.Error(), err.Error(), "expected error message to match")
				return
			}
			assert.NoError(t, err, "expected no error during Solve")

			delta := test.expectedTime.Sub(extractedTime)
			if delta < 0 {
				delta = -delta
			}

			assert.LessOrEqual(t, delta, allowedErrorMargin, "expected time to be within allowed error margin")
			for _, uplink := range test.uplinks {
				t.Logf("Uplink: %s\n", uplink.HexPayload)
			}
			t.Logf("Expected time: %v\n\n", test.expectedTime)
			t.Logf("Extracted time: %v\n", extractedTime)
		})
	}
}

func TestExtractGNSSCaptureTime(t *testing.T) {
	const gpsUtcOffset = 18 // Adjust as appropriate
	var gpsEpoch = time.Date(1980, 1, 6, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		hexInput string

		wantUTCTime time.Time
		wantErr     bool
	}{
		{
			name:        "valid payload, gpsSeconds = 100000",
			hexInput:    "00a0860100",
			wantUTCTime: gpsEpoch.Add(time.Duration(100000-gpsUtcOffset) * time.Second),
			wantErr:     false,
		},
		{
			name:        "valid payload, gpsSeconds = 123456789",
			hexInput:    "0015CD5B07",
			wantUTCTime: gpsEpoch.Add(time.Duration(123456789-gpsUtcOffset) * time.Second),
			wantErr:     false,
		},
		{
			name:        "payload too short",
			hexInput:    "01020304",
			wantUTCTime: time.Time{},
			wantErr:     true,
		},
		{
			name:        "payload with extra bytes",
			hexInput:    "0001020304FFEE",
			wantUTCTime: gpsEpoch.Add(time.Duration(0x04030201-gpsUtcOffset) * time.Second),
			wantErr:     false,
		},
		{
			name:        "90b30df6160cc180",
			hexInput:    "90b30df6160cc180",
			wantUTCTime: time.Date(2025, time.July, 4, 14, 43, 34, 0, time.UTC), // 2025-07-04 14:43:34
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Payload hex: %s\n", tt.hexInput)
			t.Logf("Expected UTC time: %v\n", tt.wantUTCTime)

			payload, err := hex.DecodeString(tt.hexInput)
			if err != nil {
				t.Fatalf("failed to decode hex input: %v", err)
			}

			gotUTCTime, err := ExtractGNSSCaptureTime(payload)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.wantUTCTime.Equal(gotUTCTime), "expected %v, got %v", tt.wantUTCTime, gotUTCTime)
			}
		})
	}
}
