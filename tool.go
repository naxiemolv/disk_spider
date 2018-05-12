package disk_spider

import (
	"bytes"
	"encoding/binary"

)

func bytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func bigEndianPack(len uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, len)
	return b
}

func appendBigEndian(buff []byte, len uint32) []byte {
	return bytesCombine(buff, bigEndianPack(len))
}

func uint32BigEndianBytesToInt(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}

