package decoder

type DataRate string

const DataRateBlazing DataRate = "blazing"
const DataRateFast DataRate = "fast"
const DataRateQuick DataRate = "quick"
const DataRateModerate DataRate = "moderate"
const DataRateSlow DataRate = "slow"
const DataRateGlacial DataRate = "glacial"

const DataRateAutomaticNarrow DataRate = "automatic-narrow"
const DataRateAutomaticWide DataRate = "automatic-wide"

const DataRateUnknown DataRate = "unknown"

// TagXL specific data rates
const DataRateTagXLDR5 DataRate = "dr5-sf7"        // 0: DR5 (EU868 SF7)
const DataRateTagXLDR4 DataRate = "dr4-sf8"        // 1: DR4 (EU868 SF8)
const DataRateTagXLDR3 DataRate = "dr3-sf9"        // 2: DR3 (EU868 SF9, US915 SF7)
const DataRateTagXLDR2 DataRate = "dr2-sf10"       // 3: DR2 (EU868 SF10, US915 SF8)
const DataRateTagXLDR1 DataRate = "dr1-sf11"       // 4: DR1 (EU868 SF11, US915 SF9)
const DataRateTagXLDR0 DataRate = "dr0-sf12"       // 5: DR0 (EU868 SF12)
const DataRateTagXLDR1To3 DataRate = "dr1-3-array" // 6: DR1-3 array (EU868 SF9-11, US915 SF7-9)
const DataRateTagXLADR DataRate = "adr"            // 7: ADR (SF7-12) for EU868
