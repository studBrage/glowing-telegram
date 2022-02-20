package cloud

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

//Virtual storage of files
//Store all files as binary sequences instead of genreating actual files
//only create file temporary when it is requested
//possibly have the created file and its binary seq stored together for easy comparison of files

var file []byte

func Cloud() {

	fmt.Print("Starter server")

	// time.Sleep(3 * time.Second)

	laddy, _ := net.ResolveTCPAddr("tcp", ":20013")
	listen, _ := net.ListenTCP("tcp", laddy)

	conn, _ := listen.AcceptTCP()
	fmt.Println("Connection established with:", conn.RemoteAddr().String())

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
			// copyToFile("txtTest")
			break
		}
		//fmt.Println(n, "bytes recieved. Local:", conn.LocalAddr().String(), " Remote:", conn.RemoteAddr().String())
		//msg := strings.Split(string(buffer[:n]), "\\x00")
		//fmt.Println()
		//fmt.Println(msg[1], ": ", msg[0])
		fmt.Println(string(buffer[:n]))
		// file = append(file, buffer[:n]...)
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
