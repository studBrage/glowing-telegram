package protocols

import (
	"encoding/binary"
)

func eventEncode(event string) byte {
	switch event {
	case "CREATE":
		return byte(1)
	case "WRITE":
		return byte(2)
	case "REMOVE":
		return byte(3)
	case "RENAME":
		return byte(4)
	case "CHMOD":
		return byte(5)
	default:
		return byte(0)
	}
}

func typeEncode(typ bool) byte {
	if typ {
		return byte(1)
	} else {
		return byte(0)
	}
}

func pathEncode(path string) []byte {
	return []byte(path)
}

func mapEncode(delta map[int]byte) []byte {
	lenComp := make([]byte, 2)
	var deltaSeq []byte
	for i, b := range delta {
		diff := []byte{byte(i), b}
		deltaSeq = append(deltaSeq, diff...)
	}
	binary.LittleEndian.PutUint16(lenComp, uint16(len(deltaSeq)))
	lenComp = append(lenComp, deltaSeq...)
	return lenComp
}

func extEncode(ext []byte) []byte {
	lenComp := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenComp, uint16(len(ext)))
	return append(lenComp, ext...)
}
