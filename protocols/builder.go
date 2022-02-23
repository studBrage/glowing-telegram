package protocols

import "encoding/binary"

func BuildInfo(event string, typ bool, path string, longest, lenD, lenX int) []byte {
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

func BuildData(delta map[int]byte, ext []byte) []byte {
	del := mapEncode(delta)
	return append(del, ext...)
}
