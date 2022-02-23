package protocols

// var msg []byte

// func main() {
// 	msg := []byte{1, 1, 46, 47, 68, 111, 99, 47, 116, 101, 120, 116, 47, 102, 105, 108, 101, 46, 116, 120, 116}

// 	event := eventDecode(msg[0])
// 	typ := typeDecode(msg[1])
// 	path := pathDecode(msg[2:])

// 	fmt.Println(event)
// 	fmt.Println(typ)
// 	fmt.Println(path)

// 	ttt := []byte{6, 0, 10, 44, 12, 97, 15, 142}
// 	mp := mapDecode(ttt[2:])
// 	fmt.Println(mp)
// }

func eventDecode(event byte) string {
	switch event {
	case 0:
		return "Unrecognized event"
	case 1:
		return "CREATE"
	case 2:
		return "WRITE"
	case 3:
		return "REMOVE"
	case 4:
		return "RENAME"
	case 5:
		return "CHMOD"
	default:
		return "NONE"
	}
}

func typeDecode(typ byte) bool {
	switch int(typ) {
	case 0:
		return false
	case 1:
		return true
	default:
		return false
	}
}

func pathDecode(path []byte) string {
	return string(path)
}

func mapDecode(mp []byte) map[int]byte {
	delta := make(map[int]byte)
	for i := 0; i < len(mp); i += 2 {
		delta[int(mp[i])] = mp[i+1]
	}
	return delta
}
