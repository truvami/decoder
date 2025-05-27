package common

func BytesToInt8(bytes []byte) int8 {
	var value int8 = 0
	for i := range bytes {
		value = int8(bytes[i])
	}
	return value
}

func BytesToInt16(bytes []byte) int16 {
	var value int16 = 0
	for i := range bytes {
		value <<= 8
		value |= int16(bytes[i])
	}
	return value
}

func BytesToInt32(bytes []byte) int32 {
	var value int32 = 0
	for i := range bytes {
		value <<= 8
		value |= int32(bytes[i])
	}
	return value
}

func BytesToInt64(bytes []byte) int64 {
	var value int64 = 0
	for i := range bytes {
		value <<= 8
		value |= int64(bytes[i])
	}
	return value
}

func BytesToUint8(bytes []byte) uint8 {
	var value uint8 = 0
	for i := range bytes {
		value = uint8(bytes[i])
	}
	return value
}

func BytesToUint16(bytes []byte) uint16 {
	var value uint16 = 0
	for i := range bytes {
		value <<= 8
		value |= uint16(bytes[i])
	}
	return value
}

func BytesToUint32(bytes []byte) uint32 {
	var value uint32 = 0
	for i := range bytes {
		value <<= 8
		value |= uint32(bytes[i])
	}
	return value
}

func BytesToUint64(bytes []byte) uint64 {
	var value uint64 = 0
	for i := range bytes {
		value <<= 8
		value |= uint64(bytes[i])
	}
	return value
}
