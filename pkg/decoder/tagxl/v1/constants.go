package tagxl

import "time"

// bufferedAgeThreshold defines the age after which an uplink is considered buffered.
// Centralizing this avoids hardcoding the duration across files.
const bufferedAgeThreshold = 5 * time.Minute
