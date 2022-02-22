package main

import "encoding/binary"

func buildInfo(event string, typ bool, path string, lenD int, lenX int) []byte {
	var msg []byte
	msg = append(msg, eventEncode(event))
	msg = append(msg, typeEncode(typ))
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(lenD))
	msg = append(msg, len...)
	binary.LittleEndian.PutUint16(len, uint16(lenX))
	msg = append(msg, len...)
	msg = append(msg, pathEncode(path)...)

	return msg
}

func buildData(delta []byte, ext []byte) []byte {
	return append(delta, ext...)
}
