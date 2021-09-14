package main

import "encoding/binary"

func IntToHex(number int64) []byte {
	cn := make([]byte, 8)
	binary.LittleEndian.PutUint64(cn, uint64(number))
	return cn
}
