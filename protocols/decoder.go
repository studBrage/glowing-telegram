package protocols

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

func PathDecode(path []byte) string {
	return string(path)
}

func mapDecode(mp []byte) map[int]byte {
	delta := make(map[int]byte)
	for i := 0; i < len(mp); i += 2 {
		delta[int(mp[i])] = mp[i+1]
	}
	return delta
}
