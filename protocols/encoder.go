package main

// var msg []byte

// func main() {

// 	msg = append(msg, eventEncode("CREATE"))
// 	msg = append(msg, typeEncode(true))
// 	msg = append(msg, pathEncode("./Doc/text/file.txt")...)

// 	fmt.Println(msg)
// }

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
