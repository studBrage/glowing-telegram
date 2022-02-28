package protocols

import "encoding/binary"

func RecieveInfo(msg []byte) ([]string, bool, []int) {
	event := eventDecode(msg[0])
	path := PathDecode(msg[7:])
	typ := typeDecode(msg[1])
	lenD := int(binary.BigEndian.Uint16(msg[2:4]))
	lenX := int(binary.BigEndian.Uint16(msg[4:6]))
	longest := int(msg[6])
	return []string{event, path}, typ, []int{lenD, lenX, longest}
}

func RecieveData(msg []byte, lenD, lenX int) (map[int]byte, []byte) {
	ext := msg[lenD:]
	if lenD == 0 {
		return nil, ext
	}
	delta := mapDecode(msg[:lenD])
	return delta, ext
}

func ExtractDataLen(info []byte) int {
	lenD := int(binary.BigEndian.Uint16(info[2:4]))
	lenX := int(binary.BigEndian.Uint16(info[4:6]))
	return lenD + lenX
}

func ExtractType(info []byte) bool {
	return typeDecode(info[1])
}
