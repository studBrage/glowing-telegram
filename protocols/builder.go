package protocols

import "encoding/binary"

func BuildInfo(event string, typ bool, path string, longest, lenD, lenX int) []byte {
	var msg []byte
	// msg has form:
	// [0]event[1 byte] - [1]type[1 byte] - [2:4]lenght delta[2 byte]
	//  - [4:6]lenght extension[2 byte] - [6]lenght of longest[1 byte]
	//  - [7:]path[rest of msg]
	msg = append(msg, eventEncode(event))
	msg = append(msg, typeEncode(typ))
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(lenD))
	msg = append(msg, len...)
	binary.LittleEndian.PutUint16(len, uint16(lenX))
	msg = append(msg, len...)
	msg = append(msg, byte(longest))
	msg = append(msg, pathEncode(path)...)

	return msg
}

func BuildData(delta map[int]byte, ext []byte) []byte {
	del := mapEncode(delta)
	return append(del, ext...)
}
