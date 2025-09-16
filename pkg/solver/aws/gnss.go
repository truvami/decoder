package aws

import (
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

// leapSeconds is the number of leap seconds that have accumulated since GPS epoch (Jan 6, 1980).
// As of 2025, it's 18 seconds. This value needs to be updated when new leap seconds are introduced.
const leapSeconds = 18

// GNSSCapture represents a single GNSS uplink capture, containing the hexadecimal payload
// and the time it was received.
type GNSSCapture struct {
	HexPayload string
	ReceivedAt time.Time
}

// extract30bitWords extracts 10 30-bit words from a byte slice.
// This function is crucial for parsing the raw binary data of GPS navigation messages,
// where information is often packed into 30-bit words.
func extract30bitWords(data []byte) ([]uint32, error) {
	// A typical GPS navigation page (subframe) consists of 10 30-bit words, totaling 300 bits.
	// 300 bits / 8 bits/byte = 37.5 bytes, so at least 38 bytes are needed to contain these bits.
	if len(data) < 38 {
		return nil, fmt.Errorf("need at least 38 bytes for NAV page")
	}
	words := make([]uint32, 10)
	bitPos := 0 // Tracks the current bit position across the byte slice

	for i := 0; i < 10; i++ {
		var word uint32
		// Read 30 bits to form each 30-bit word
		for b := 0; b < 30; b++ {
			byteIndex := bitPos / 8                     // Determine which byte the current bit is in
			bitIndex := 7 - (bitPos % 8)                // Determine the bit position within that byte (MSB to LSB)
			bit := (data[byteIndex] >> bitIndex) & 0x01 // Extract the bit (0 or 1)
			word = (word << 1) | uint32(bit)            // Shift the existing word left and append the new bit
			bitPos++                                    // Move to the next bit
		}
		words[i] = word
	}
	return words, nil
}

// decodeZCount extracts the Z-count (Time of Week) from the HOW (Hand Over Word)
// of a GPS navigation message subframe.
// The Z-count is a 17-bit field that represents 1.5-second intervals.
func decodeZCount(payload []byte) (float64, error) {
	words, err := extract30bitWords(payload)
	if err != nil {
		return 0, err
	}
	// The HOW word is typically the second 30-bit word (at index 1) in a GPS subframe.
	howWord := words[1]
	// The Z-count is located in bits 13 through 29 (inclusive) of the 30-bit HOW word.
	// To extract it:
	// 1. Shift the HOW word right by 13 bits to move the Z-count to the least significant bits.
	// 2. Apply a mask (0x1FFFF, which is 17 ones in binary) to isolate the 17-bit Z-count.
	zCount := (howWord >> 13) & 0x1FFFF
	// Convert the 17-bit Z-count to seconds by multiplying by the 1.5-second interval.
	towSec := float64(zCount) * 1.5
	return towSec, nil
}

// SolveCapturedAt attempts to infer the precise GNSS capture time from a slice of buffered uplinks.
// It uses the Z-count (Time of Week) from the GNSS payload and the uplink's reception time
// to estimate the actual capture time.
// It specifically looks for a Z-count fix that is up to 30 minutes older than the ReceivedAt timestamp.
func SolveCapturedAt(captures []GNSSCapture) (time.Time, error) {
	if len(captures) == 0 {
		return time.Time{}, fmt.Errorf("no captures provided")
	}

	// Internal struct to hold a candidate capture time and its error (difference from receivedAt).
	type candidate struct {
		time  time.Time
		error time.Duration
	}

	var best *candidate // Stores the best (most accurate) candidate found so far

	// The Z-count is a 17-bit counter, meaning its value repeats (rolls over)
	// every 2^17 * 1.5 seconds. This is a critical constant for disambiguating TOW.
	const rolloverIntervalSec = 196606.5 // 2^17 * 1.5 seconds

	for _, cap := range captures {
		// Decode the hexadecimal payload string into a byte slice.
		payload, err := hex.DecodeString(cap.HexPayload)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to decode hex: %w", err)
		}

		var navBytes []byte
		// Check if the payload starts with a TLV (Type-Length-Value) header (Type 0x01 for navigation data).
		if payload[0] == 0x01 {
			navLen := int(payload[1]) // Length of the navigation data
			if len(payload) < 2+navLen {
				return time.Time{}, fmt.Errorf("invalid TLV length")
			}
			navBytes = payload[2 : 2+navLen] // Extract the navigation data part
		} else {
			// If no TLV header, assume the entire payload is raw navigation data.
			navBytes = payload
		}

		// Decode the Z-count from the navigation bytes to get the Time of Week in seconds.
		towSec, err := decodeZCount(navBytes)
		if err != nil {
			return time.Time{}, fmt.Errorf("zcount decode: %w", err)
		}

		// Convert the uplink's ReceivedAt time to GPS time by adding the current leap seconds.
		// GPS time does not observe leap seconds, so this aligns it with the GPS system's internal clock.
		gpsReceived := cap.ReceivedAt.Add(time.Duration(leapSeconds) * time.Second)

		// Calculate the precise start of the GPS week (Sunday 00:00:00 UTC) for the week
		// that the gpsReceived timestamp falls into. This is crucial for correctly
		// anchoring the `towSec` within the correct GPS week.
		daysToSunday := int(gpsReceived.Weekday()) // Weekday() returns 0 for Sunday, 1 for Monday, etc.
		// Truncate to the beginning of the day (00:00:00 UTC) and then subtract days to get to Sunday.
		gpsWeekStartForReceivedAt := gpsReceived.Truncate(24 * time.Hour).Add(time.Duration(-daysToSunday) * 24 * time.Hour)

		// Calculate the Time of Week (TOW) of the gpsReceived timestamp relative to its GPS week start.
		towOfReceivedAt := gpsReceived.Sub(gpsWeekStartForReceivedAt).Seconds()

		// The `towSec` extracted from the payload is a 17-bit value and thus rolls over.
		// We need to find the integer 'n' (number of rollover cycles) such that
		// (towSec + n * rolloverIntervalSec) is closest to `towOfReceivedAt`.
		// This 'n' effectively tells us which rollover cycle within the week the capture occurred.
		bestN := 0
		minDiffTow := math.Abs(towOfReceivedAt - towSec) // Initial minimum difference (assuming n=0)

		// Iterate through a range of possible 'n' values. A GPS week (604800s) contains ~3.07 rollover intervals.
		// Checking from -3 to +3 should cover all relevant rollover possibilities around the `receivedAt`'s TOW,
		// allowing for slight discrepancies or captures that span week boundaries.
		for n := -3; n <= 3; n++ {
			currentTowCandidate := towSec + float64(n)*rolloverIntervalSec
			diff := math.Abs(towOfReceivedAt - currentTowCandidate)
			if diff < minDiffTow {
				minDiffTow = diff
				bestN = n
			}
		}

		// Now, construct the candidate GPS fix time.
		// This is done by taking the start of the GPS week (derived from receivedAt),
		// adding the payload's towSec, and then adding the best-fit rollover offset (n * rolloverIntervalSec).
		candidateFixGPS := gpsWeekStartForReceivedAt.Add(time.Duration(towSec+float64(bestN)*rolloverIntervalSec) * time.Second)

		// It's possible that the actual capture occurred in a GPS week immediately preceding or following
		// the week of `receivedAt`, even if the TOW alignment suggests the current week.
		// Therefore, we check three candidate GPS fix times: the one calculated, and one week before/after.
		candidateGPSFixTimes := []time.Time{
			candidateFixGPS,
			candidateFixGPS.Add(-7 * 24 * time.Hour), // Candidate from the previous GPS week
			candidateFixGPS.Add(7 * 24 * time.Hour),  // Candidate from the next GPS week
		}

		// Evaluate each candidate GPS fix time.
		for _, currentCandidateFix := range candidateGPSFixTimes {
			// Convert the candidate GPS time back to UTC by subtracting leap seconds.
			captureUTC := currentCandidateFix.Add(-time.Duration(leapSeconds) * time.Second)

			// Calculate the signed difference between when the uplink was received and the inferred capture time.
			// A positive `diff` means `captureUTC` is older than `receivedAt`.
			diff := cap.ReceivedAt.Sub(captureUTC)

			// Print detailed debug information for each candidate.
			fmt.Printf(
				"Uplink: %s\nDecoded TOW: %.1fs → candidate fix: %s → diff: %s (bestN: %d, currentCandidateFix: %s)\n",
				cap.HexPayload, towSec, captureUTC.Format(time.RFC3339), diff, bestN, currentCandidateFix.Format(time.RFC3339),
			)

			// Accept the candidate if it meets the criteria:
			// 1. `diff` is non-negative (captureUTC is not in the future relative to receivedAt).
			// 2. `diff` is within the allowed 30-minute margin (captureUTC is at most 30 minutes older).
			if diff >= 0 && diff <= 30*time.Minute {
				// If this is the first valid candidate, or it's better (smaller error) than the current best, update best.
				if best == nil || diff < best.error {
					best = &candidate{time: captureUTC, error: diff}
				}
			}
		}
	}

	// If an acceptable candidate capture time was found across all uplinks, return the best one.
	if best != nil {
		return best.time, nil
	}

	// Fallback: If no valid GNSS time could be inferred from any of the uplinks
	// (i.e., none met the criteria), return the ReceivedAt time of the most recently received uplink.
	// This assumes the 'captures' slice is ordered by 'ReceivedAt' time, so the last element is the newest.
	return captures[len(captures)-1].ReceivedAt, nil
}

// ExtractGNSSCaptureTime extracts the GNSS capture timestamp from the payload (hex string or []byte).
func ExtractGNSSCaptureTime(payload []byte) (utcTime time.Time, err error) {
	const leapSeconds = 18
	const gpsEpoch = "1980-01-06T00:00:00Z"
	if len(payload) < 5 {
		return time.Time{}, fmt.Errorf("payload too short")
	}
	gpsSeconds := uint32(payload[1]) | uint32(payload[2])<<8 | uint32(payload[3])<<16 | uint32(payload[4])<<24
	epoch, _ := time.Parse(time.RFC3339, gpsEpoch)
	return epoch.Add(time.Duration(int64(gpsSeconds)-leapSeconds) * time.Second), nil
}

// func main() {
// 	// Example payload (hex): 0A5E2B1C00...
// 	payloadHex := "0A5E2B1C00"
// 	payload, err := hex.DecodeString(payloadHex)
// 	if err != nil {
// 		panic(err)
// 	}
// 	gpsSeconds, utcTime, err := ExtractGNSSCaptureTime(payload)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("GNSS Capture Time (GPS seconds): %d\n", gpsSeconds)
// 	fmt.Printf("GNSS Capture Time (UTC): %s\n", utcTime.Format(time.RFC3339))
// }
