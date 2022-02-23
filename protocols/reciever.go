package protocols

import "encoding/binary"

func RecieveInfo(msg []byte) ([]string, bool, []int) {
	event := eventDecode(msg[0])
	path := pathDecode(msg[7:])
	typ := typeDecode(msg[1])
	lenD := int(binary.BigEndian.Uint16(msg[2:4]))
	lenX := int(binary.BigEndian.Uint16(msg[4:6]))
	longest := int(msg[6])
	return []string{event, path}, typ, []int{lenD, lenX, longest}
}

func RecieveData(msg []byte, lenD, lenX int) (map[int]byte, []byte) {
	delta := mapDecode(msg[:lenD])
	ext := msg[lenD:]
	return delta, ext
}