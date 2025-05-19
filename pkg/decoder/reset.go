package decoder

type ResetReason string

// TODO: Add more reset reasons and descriptions, check with Viktor
const (
	// Device did reset for an unknown reason
	ResetReasonUnknown ResetReason = "unknown"
	// Device did reset because of LRR1110 failure
	ResetReasonLrr1110FailCode ResetReason = "lrr1110-failure"
	// Device did reset because of a watchdog timeout
	ResetReasonWatchdog ResetReason = "watchdog"
	// Device did reset because of a pin reset
	ResetReasonPinReset ResetReason = "pin-reset"
	// Device did reset because of a system reset
	ResetReasonSystemReset ResetReason = "system-reset"
	// Device did reset because of another reason
	ResetReasonOtherReset ResetReason = "other-reset"
	// Device did reset because of a power reset
	ResetReasonPowerReset ResetReason = "power-reset"
)
