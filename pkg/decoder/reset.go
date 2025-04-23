package decoder

type ResetReason string

// TODO: Add more reset reasons and descriptions, check with Viktor
const ResetReasonLrr1110FailCode ResetReason = "lrr1110-failure"
const ResetReasonWatchdog ResetReason = "watchdog"
const ResetReasonPinReset ResetReason = "pin-reset"
const ResetReasonSystemReset ResetReason = "system-reset"
const ResetReasonOtherReset ResetReason = "other-reset"
const ResetReasonPowerReset ResetReason = "power-reset"
