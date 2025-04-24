package decoder

type ResetReason string

// TODO: Add more reset reasons and descriptions, check with Viktor
const (
	// Device did reset for an unknown reason
	ResetReasonUnknown ResetReason = "UNKNOWN"
	// Device did reset because of LRR1110 failure
	ResetReasonLrr1110FailCode ResetReason = "LRR1110_FAILURE"
	// Device did reset because of a watchdog timeout
	ResetReasonWatchdog ResetReason = "WATCHDOG"
	// Device did reset because of a pin reset
	ResetReasonPinReset ResetReason = "PIN_RESET"
	// Device did reset because of a system reset
	ResetReasonSystemReset ResetReason = "SYSTEM_RESET"
	// Device did reset because of another reason
	ResetReasonOtherReset ResetReason = "OTHER_RESET"
	// Device did reset because of a power reset
	ResetReasonPowerReset ResetReason = "POWER_RESET"
)
