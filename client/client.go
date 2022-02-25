package client

import (
	"fmt"
	"net"
	"time"
)

func Client(port int) *net.TCPConn {
	server := fmt.Sprintf(":%d", port)

	fmt.Println("Connecting to port: ", fmt.Sprint(port))
	raddy, _ := net.ResolveTCPAddr("tcp", server)
	conn, _ := net.DialTCP("tcp", nil, raddy)
	fmt.Println("Connected to:", conn.RemoteAddr().String())

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

// func WriteInfo(c *net.TCPConn, msg []byte) {
// 	c.Write([]byte("INFO"))
// 	time.Sleep(100 * time.Millisecond)
// 	c.Write(msg)
// 	time.Sleep(100 * time.Millisecond)
// }

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
	// fmt.Println(msg[prev:])
	c.Write(msg[prev:])
}

func WriteFull(c *net.TCPConn, infomsg []byte, datamsg []byte, chunkSize int) {
	c.Write(infomsg)
	time.Sleep(100 * time.Millisecond)
	WriteBytesChunk(c, datamsg, chunkSize)
}
