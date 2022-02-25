package cloud

import (
	"fmt"
	"net"
	"time"
)

//Virtual storage of files
//Store all files as binary sequences instead of genreating actual files
//only create file temporary when it is requested
//possibly have the created file and its binary seq stored together for easy comparison of files

func Cloud(dataChan chan []byte) {

	fmt.Print("Starter server")

	// time.Sleep(3 * time.Second)

	laddy, _ := net.ResolveTCPAddr("tcp", ":20013")
	listen, _ := net.ListenTCP("tcp", laddy)

	conn, _ := listen.AcceptTCP()
	fmt.Println("Connection established with:", conn.RemoteAddr().String())

	defer listen.Close()

	go read(conn, dataChan)
	for {

	}
}

func read(inconn *net.TCPConn, dataChan chan []byte) {
	for {
		// fmt.Println("Toppen av read")
		buffer := make([]byte, 1024)
		n, err := inconn.Read(buffer)
		if err != nil {
			fmt.Println("Connection lost")
			// copyToFile("txtTest")
			break
		}
		// fmt.Println(n, "bytes recieved")
		// fmt.Println(buffer)
		// fmt.Println(buffer[:n])
		dataChan <- buffer[:n]

		time.Sleep(40 * time.Millisecond)
	}
}
