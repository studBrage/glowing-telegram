package client

import (
	"fmt"
	"net"
	"time"
)

//store all files as binary seq for comparison
//if addition to file, send key with msg if new file or changed file
//if changed file, only send difference, server will simply add the difference to excisting binary seq
//keep log of sent files to compare with new changes

// var localDir string
// var int port := 20013
// var openFile string

// func init() {
// 	flag.StringVar(&localDir, "localDir", "path", "path to monitoring dir")
// 	flag.IntVar(&port, "port", 20013, "port to dial")
// 	flag.StringVar(&openFile, "filename", "test.txt", "filename of file to send")
// 	flag.Parse()
// }

func Client(port int) *net.TCPConn {
	server := fmt.Sprintf(":%d", port)

	fmt.Println("Connecting to port: ", fmt.Sprint(port))
	raddy, _ := net.ResolveTCPAddr("tcp", server)
	conn, _ := net.DialTCP("tcp", nil, raddy)
	fmt.Println("Connected to:", conn.RemoteAddr().String())

	// file, err := os.Open("index.jpeg")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fileType, err := getFileType(file)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Println(fileType)

	// defer file.Close()

	//reader := bufio.NewReader(file)
	//buf := make([]byte, 512)
	//
	//for {
	//
	//	_, err := reader.Read(buf)
	//
	//	if err != nil {
	//		if err != io.EOF {
	//			fmt.Println(err.Error())
	//		}
	//		break
	//	}
	//	write(conn, buf)
	//}
	//fmt.Println(string(buf))
	return conn
}

func WriteString(conn *net.TCPConn, msg string) {
	fmt.Println("Sending msg")
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
}

func WriteInfo(c *net.TCPConn, msg []byte) {
	c.Write([]byte("INFO"))
	time.Sleep(100 * time.Millisecond)
	c.Write(msg)
	time.Sleep(100 * time.Millisecond)
}

func WriteBytesChunk(c *net.TCPConn, msg []byte, chunckSize int) {
	fmt.Println("Sending msg")
	prev := 0
	till := len(msg) - chunckSize
	for prev < till {
		next := prev + chunckSize
		// fmt.Println(msg[prev:next])
		c.Write(msg[prev:next])
		prev = next
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println(msg[prev:])
	c.Write(msg[prev:])
}

func WriteFull(c *net.TCPConn, infomsg []byte, datamsg []byte, chunkSize int) {
	WriteInfo(c, infomsg)
	time.Sleep(100 * time.Millisecond)
	WriteBytesChunk(c, datamsg, chunkSize)
}
