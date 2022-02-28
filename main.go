package main

import (
	"flag"
	"fmt"
	"glowing-telegram/client"
	"glowing-telegram/cloud"
	comp "glowing-telegram/comparator"
	"glowing-telegram/monitor"
	proto "glowing-telegram/protocols"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// To run normally
//
// Server first
// run: go run main.go -app="cloud" -path=**choosen path**
//
// then client
// run: go run main.go -app="client" -path=**chosen path**
//

var app string
var eventChannel chan fsnotify.Event
var infoMsg []byte
var dataMsg []byte
var dump []byte
var files map[string][]byte
var conn *net.TCPConn
var targetPath string

func init() {
	flag.StringVar(&app, "app", "none", "Define the type of app to start")
	flag.StringVar(&targetPath, "path", "", "Directory path")
	flag.Parse()
}

func main() {

	switch app {
	case "cloud":
		incoming := make(chan []byte)
		fmt.Println("Dette er en cloud")
		fmt.Println("------------------------------------------")
		go cloud.Cloud(incoming)
		go server(incoming)

	case "client":
		fmt.Println("Dette er en client")
		fmt.Println("------------------------------------------")
		time.Sleep(2 * time.Second)
		conn = client.Client(20013)
		defer conn.Close()
		time.Sleep(1 * time.Second)

		eventChannel = make(chan fsnotify.Event, 1)
		go monitor.Watch(targetPath, eventChannel)
		go Moni(eventChannel)
		files = copyAll()
		time.Sleep(500 * time.Millisecond)

	// This case is mainly for testing the program
	case "init":
		cloud := exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go", "-app=\"cloud\"").Run()
		client := exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go", "-app=\"client\"").Run()
		if cloud != nil || client != nil {
			fmt.Println(cloud.Error())
			fmt.Println(client.Error())
		}
		os.Exit(3)
	}
	for {

	}

}

func server(incoming chan []byte) {
	tempLen := 0
	incomingLen := 0

	for {
		msg := <-incoming
		if tempLen == incomingLen {
			infoMsg = msg
			incomingLen = proto.ExtractDataLen(msg)
			tempLen = 0
			dataMsg = dump
		} else {
			dataMsg = append(dataMsg, msg...)
			tempLen = len(dataMsg)
			if tempLen == incomingLen {
				fmt.Println("File received")
				msgHandler(infoMsg, dataMsg)
			}
		}
	}
}

func Moni(c chan fsnotify.Event) {
	var isFile bool
	for {
		isFile = true
		select {
		case e := <-c:
			targetType := filepath.Ext(e.Name)
			action := e.Op.String()
			spl := strings.Split(e.Name, "\\")
			targetSpl := strings.Split(targetPath, "\\")
			spl = spl[len(targetSpl):]

			//The slash separator for path might have to be changed
			// in order to work on linux
			toServerPath := strings.Join(spl, "/")
			localPath := strings.Replace(e.Name, "\\", "/", -1)

			if targetType == "" {
				isFile = false
			}
			if action == "WRITE" && targetType == "" {
				continue
			}

			switch action {
			case "CREATE":
				doCreate(action, isFile, localPath, toServerPath)
			case "WRITE":
				doWrite(action, isFile, localPath, toServerPath, files[localPath])
			case "REMOVE":
				doRemove(action, isFile, localPath, toServerPath)
			case "RENAME":
				newName := <-c
				newSpl := strings.Split(newName.Name, "\\")

				newP := strings.Join(newSpl, "/")
				newSendPath := strings.Join(newSpl[len(targetSpl):], "/")

				doRename(action, isFile, localPath, newSendPath, newP)
			default:
				continue
			}
			client.WriteFull(conn, infoMsg, dataMsg, 1024)
			dataMsg = dump
			infoMsg = dump

		default:
			continue
		}
	}
}

func doCreate(e string, isFile bool, path, sendPath string) {
	var targetBytes []byte
	if !isFile {
		targetBytes = []byte("FOLDER")
		files[path] = targetBytes
		err := monitor.AddWatcher(path)
		if err != nil {
			fmt.Println("Error adding watcher to created dir", path)
			fmt.Println(err.Error())
		}
	} else {
		targetBytes = comp.DecodeFile(path)
		files[path] = targetBytes
	}
	infoMsg = proto.BuildInfo(e, isFile, sendPath, 0, 0, len(targetBytes))
	dataMsg = targetBytes
}

func doWrite(e string, isFile bool, path, sendPath string, backupFile []byte) {
	targetBytes := comp.DecodeFile(path)
	delta, largest, ext := comp.FindDelta(backupFile, targetBytes)
	checkerMap := map[int]byte{0: byte(0)}

	if reflect.DeepEqual(delta, checkerMap) {
		dataMsg = ext
	} else {
		dataMsg = proto.BuildData(delta, ext)
	}

	infoMsg = proto.BuildInfo(e, isFile, sendPath, largest, len(dataMsg)-len(ext), len(ext))

	files[path] = comp.UpdateChange(backupFile, largest, delta, ext)
}

func doRename(e string, isFile bool, path, sendPath string, new string) {
	infoMsg = proto.BuildInfo(e, isFile, path, 0, 0, 0)
	dataMsg = []byte(sendPath)

	monitor.RemoveWatcer(path)
	monitor.AddWatcher(new)
	files[new] = files[path]
	delete(files, path)
}

func doRemove(e string, isFile bool, path, sendPath string) {
	infoMsg = proto.BuildInfo(e, isFile, sendPath, 0, 0, 1)
	dataMsg = append(dataMsg, byte(0))
	if !isFile {
		monitor.RemoveWatcer(path)
	}
	delete(files, path)
}

func copyAll() map[string][]byte {
	allFiles := make(map[string][]byte)
	err := filepath.WalkDir("destFolder", func(path string, info fs.DirEntry, err error) error {
		pth := strings.Replace(path, "\\", "/", -1)

		if filepath.Ext(path) == "" {
			allFiles[pth] = []byte("FOLDER")
		} else {
			allFiles[pth] = comp.DecodeFile(pth)
		}
		return nil
	})

	if err != nil {
		fmt.Println("copyALL error:", err.Error())
	}
	return allFiles
}

func msgHandler(info []byte, data []byte) {
	ep, isFile, lens := proto.RecieveInfo(info)
	dest := targetPath + "/" + ep[1]

	typ := ""
	if isFile {
		typ = "FILE"
	} else {
		typ = "FOLDER"
	}
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println("Event: ", ep[0], "Type: ", typ)
	fmt.Println("Target path for event: ", dest)
	fmt.Println("Size of delta:", lens[0], "BYTES")
	fmt.Println("Size of extension:", lens[1], "BYTES")
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	switch ep[0] {
	case "CREATE":
		var err error
		if !isFile {
			err = os.Mkdir(dest, 0777)
		} else {
			err = comp.CopyToFile(dest, data)
		}
		if err != nil {
			fmt.Println("Could not create dir: ", dest)
			fmt.Println(err.Error())
		}
	case "WRITE":
		oldFile := comp.DecodeFile(dest)
		if len(oldFile) == 0 {
			comp.CopyToFile(dest, data[2:])
			break
		}
		delta, ext := proto.RecieveData(data, len(data)-lens[1], lens[1])
		updatedFile := comp.UpdateChange(oldFile, lens[2], delta, ext)
		comp.CopyToFile(dest, updatedFile)
	case "REMOVE":
		os.Remove(dest)
	case "RENAME":
		os.Rename(dest, proto.PathDecode(data))
	}
}
