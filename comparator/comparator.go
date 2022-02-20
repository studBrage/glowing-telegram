package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

// TODO
// Finish funciton for finding slice diff

func main() {

	file1, err := openFile("../destFolder/foo.txt")
	if err != nil {
		panic(err.Error())
	}
	file2, _ := openFile("../destFolder/bar.txt")

	defer file1.Close()
	defer file2.Close()

	byte1 := decodeFile(file1)
	byte2 := decodeFile(file2)

	fmt.Println(bytes.Compare(byte1, byte2))

	copyToFile("foobar.txt", byte1)
}

func openFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
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

func decodeFile(file *os.File) []byte {
	var fileRep []byte

	reader := bufio.NewReader(file)
	buf := make([]byte, 512)

	for {

		n, err := reader.Read(buf)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err.Error())
			}
			break
		}
		fileRep = append(fileRep, buf[:n]...)
		// write(conn, buf)
	}
	// fmt.Println(string(buf))

	return fileRep
}

func copyToFile(fileName string, byteRep []byte) {
	destination := fmt.Sprintf("../destFolder/%s", fileName)
	err := ioutil.WriteFile(destination, byteRep, 0644)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("File saved")
}

func compSlices(a, b []byte) (int, int, int) {
	al, bl := len(a), len(b)
	diff := al - bl
	if diff == 0 {
		return 0, len(a), diff
	} else if diff > 0 {
		return 1, len(a), int(math.Abs(float64(diff)))
	} else {
		return 2, len(b), int(math.Abs(float64(diff)))
	}
}

// func findDelta(a, b []byte) {
// 	largest, shortLen, diff := compSlices(a, b)

// 	switch largest{
// 	case 1:
// 		for i := range b{

// 		}
// 	}
// }
