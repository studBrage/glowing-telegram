package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var localDir string
var port int
var openFile string

func init() {
	flag.StringVar(&localDir, "localDir", "path", "path to monitoring dir")
	flag.IntVar(&port, "port", 20013, "port to dial")
	flag.StringVar(&openFile, "filename", "test.txt", "filename of file to send")
	flag.Parse()
}

func main() {
	//server := fmt.Sprintf(":%d", port)

	//raddy, _ := net.ResolveTCPAddr("tcp", server)
	//conn, _ := net.DialTCP("tcp", nil, raddy)

	file, err := os.Open("index.jpeg")
	if err != nil {
		panic(err.Error())
	}

	fileType, err := getFileType(file)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(fileType)

	defer file.Close()
	//defer conn.Close()

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
}

func write(conn *net.TCPConn, msg []byte) {
	fmt.Println("Sending msg")
	conn.Write(msg)
	time.Sleep(100 * time.Millisecond)
}

func getFileType(file *os.File) (string, error) {
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(buffer[:n])
	fmt.Println(fileType)
	types := strings.Split(fileType, "/")

	return types[1], nil
}
