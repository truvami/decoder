package common

// Documentation:
// https://web.archive.org/web/20250820082113/https://www.semtech.com/loracloud-documentation/mdmsvc.html#lora-edge-gnss-ng-nav-group-positioning-protocol

type GNSSNGHeader struct {
	EndOfGroup           bool
	ReservedForFutureUse uint8
	GroupToken           uint8
}

const (
	GNSSNGHeaderEndOfGroupMask           = 0b1000_0000
	GNSSNGHeaderReservedForFutureUseMask = 0b0100_0000
	GNSSNGHeaderGroupTokenMask           = 0b0011_1111
)

func DecodeGNSSNGHeader(payload []byte) (GNSSNGHeader, error) {
	if len(payload) < 1 {
		return GNSSNGHeader{}, ErrGNSSNGHeaderByteMissing
	}

	headerByte := payload[0]
	header := GNSSNGHeader{
		EndOfGroup:           (headerByte & GNSSNGHeaderEndOfGroupMask) != 0,
		ReservedForFutureUse: (headerByte & GNSSNGHeaderReservedForFutureUseMask) >> 6,
		GroupToken:           (headerByte & GNSSNGHeaderGroupTokenMask),
	}

	return header, nil
}
