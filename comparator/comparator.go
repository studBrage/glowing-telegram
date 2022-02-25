package comparator

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

func openFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFileType(path string) (string, error) {
	file, err := openFile(path)
	if err != nil {
		fmt.Println("Could not open file")
		return "", err
	}
	defer file.Close()
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

func DecodeFile(path string) []byte {
	file, err := openFile(path)
	if err != nil {
		fmt.Println("Could not open file:", err.Error())
	}
	defer file.Close()
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

func CopyToFile(destination string, byteRep []byte) error {
	// destination := fmt.Sprintf("../destFolder/%s", fileName)
	err := ioutil.WriteFile(destination, byteRep, 0644)
	if err != nil {
		return err
	}
	fmt.Println("File saved")
	return nil
}

func compSlices(a, b []byte) (int, int, int) {
	diff := len(a) - len(b)
	if diff == 0 {
		return 0, len(a), diff
	} else if diff > 0 {
		return 1, len(b), int(math.Abs(float64(diff)))
	} else {
		return 2, len(a), int(math.Abs(float64(diff)))
	}
}

func FindDelta(org, new []byte) (map[int]byte, int, []byte) {
	largest, shortLen, diff := compSlices(org, new)

	diffMap := make(map[int]byte)

	for i := 0; i < shortLen; i++ {
		if org[i] != new[i] {
			diffMap[i] = new[i]
		}
	}

	switch largest {
	case 0:
		return nil, 0, nil
	case 1:
		return diffMap, 1, org[len(org)-diff:]
	case 2:
		return diffMap, 2, new[len(new)-diff:]
	}
	return nil, -1, nil
}

func UpdateChange(target []byte, longest int, delta map[int]byte, ext []byte) []byte {
	// fmt.Println("DElta")
	// fmt.Println(delta)

	// fmt.Println("Len taget", len(target))

	for i, b := range delta {
		if i >= len(target) {
			continue
		}
		target[i] = b
	}
	if longest == 2 {
		target = append(target, ext...)
	} else if longest == 1 {
		target = target[:len(target)-len(ext)]
	}
	return target
}
