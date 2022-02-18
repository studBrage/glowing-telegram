package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

var file []byte

func main() {

	laddy, _ := net.ResolveTCPAddr("tcp", ":20013")
	listen, _ := net.ListenTCP("tcp", laddy)

	conn, _ := listen.AcceptTCP()

	defer listen.Close()

	go read(conn)

	for {

	}
}

func read(inconn *net.TCPConn) {
	for {
		buffer := make([]byte, 1024)
		n, err := inconn.Read(buffer)
		if err != nil {
			fmt.Println("Connection lost")
			copyToFile("txtTest")
			break
		}
		//fmt.Println(n, "bytes recieved. Local:", conn.LocalAddr().String(), " Remote:", conn.RemoteAddr().String())
		//msg := strings.Split(string(buffer[:n]), "\\x00")
		//fmt.Println()
		//fmt.Println(msg[1], ": ", msg[0])
		fmt.Println(string(buffer[:n]))
		file = append(file, buffer[:n]...)
		time.Sleep(300 * time.Millisecond)
	}
}

func copyToFile(fileName string) {
	destination := fmt.Sprintf("./destFolder/%s", fileName)
	err := ioutil.WriteFile(destination, file, 0644)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("File saved")
}
