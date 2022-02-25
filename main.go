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

// ------------------------------------
// TODO
//
// Directory monitor DONE
// File comparing DONE
// Byte storage of files
// Proper file transfer
// File storage in cloud
//
// ------------------------------------

var app string
var eventChannel chan fsnotify.Event
var infoMsg []byte
var dataMsg []byte
var dump []byte
var files map[string][]byte
var conn *net.TCPConn

func init() {
	flag.StringVar(&app, "app", "none", "Define the type of app to start")
	flag.Parse()
	// files = copyAll()
}

func main() {

	switch app {
	case "cloud":
		incoming := make(chan []byte)
		fmt.Println("Dette er en cloud")
		fmt.Println("------------------------------------------")
		go cloud.Cloud(incoming)
		go server(incoming)
		// fmt.Println("Its a go!")

	case "client":
		fmt.Println("Dette er en client")
		fmt.Println("------------------------------------------")
		time.Sleep(2 * time.Second)
		conn = client.Client(20013)
		defer conn.Close()
		time.Sleep(1 * time.Second)

		eventChannel = make(chan fsnotify.Event, 1)
		go monitor.Watch("destFolder", eventChannel)
		go Moni(eventChannel)
		files = copyAll()
		time.Sleep(500 * time.Millisecond)

		// client.WriteString(conn, "Yay det funket!")
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
		// fmt.Println("Message recieved:", msg)
		if tempLen == incomingLen {
			infoMsg = msg
			// fmt.Println("Infomsg:", infoMsg)
			incomingLen = proto.ExtractDataLen(msg)
			// fmt.Println("Len of incoming msg:", incomingLen)
			tempLen = 0
			dataMsg = dump
		} else {
			dataMsg = append(dataMsg, msg...)
			// fmt.Println("DataMsg:", dataMsg)
			// fmt.Println("Len of currently recieved msg:", len(dataMsg))
			tempLen = len(dataMsg)
			if tempLen == incomingLen {
				fmt.Println("File received")
				msgHandler(infoMsg, dataMsg)
			}
		}
	}
}

func Moni(c chan fsnotify.Event) {
	// var prev string
	// var current string
	var isFile bool
	for {
		isFile = true
		select {
		case e := <-c:
			target := filepath.Base(e.Name)
			targetType := filepath.Ext(target)
			// fmt.Println("Target type:", targetType)
			action := e.Op.String()
			path := strings.Replace(e.Name, "\\", "/", -1)
			if targetType == "" {
				isFile = false
			}
			if action == "WRITE" && targetType == "" {
				continue
			}

			switch action {
			case "CREATE":
				doCreate(action, isFile, path)
				// msgHandler(infoMsg, dataMsg)
			case "WRITE":
				doWrite(action, isFile, path, files[path])
				// msgHandler(infoMsg, dataMsg)
			case "REMOVE":
				doRemove(action, isFile, path)
				// msgHandler(infoMsg, dataMsg)
			case "RENAME":
				newName := <-c
				newP := strings.Replace(newName.Name, "\\", "/", -1)
				// fmt.Println("RENAMED, old name: ", path, "new name: ", newP)
				// msgHandler(infoMsg, dataMsg)
				doRename(action, isFile, path, newP)
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

func doCreate(e string, isFile bool, path string) {
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
	infoMsg = proto.BuildInfo(e, isFile, path, 0, 0, len(targetBytes))
	dataMsg = targetBytes

	// fmt.Println("doCREARE msg")
	// fmt.Println(infoMsg)
	// fmt.Println(dataMsg)
}

func doWrite(e string, isFile bool, path string, backupFile []byte) {
	targetBytes := comp.DecodeFile(path)
	delta, largest, ext := comp.FindDelta(backupFile, targetBytes)
	checkerMap := map[int]byte{0: byte(0)}
	// fmt.Println("Delta in doWrite func:", delta)
	if reflect.DeepEqual(delta, checkerMap) {
		dataMsg = ext
	} else {
		dataMsg = proto.BuildData(delta, ext)
	}

	infoMsg = proto.BuildInfo(e, isFile, path, largest, len(dataMsg)-len(ext), len(ext))
	// fmt.Println("doWRITE msg")
	// fmt.Println(infoMsg)
	// fmt.Println(dataMsg)

	files[path] = comp.UpdateChange(backupFile, largest, delta, ext)
}

func doRename(e string, isFile bool, path string, new string) {
	infoMsg = proto.BuildInfo(e, isFile, path, 0, 0, 0)
	dataMsg = []byte(new)

	monitor.RemoveWatcer(path)
	monitor.AddWatcher(new)
	files[new] = files[path]
	delete(files, path)
}

func doRemove(e string, isFile bool, path string) {
	infoMsg = proto.BuildInfo(e, isFile, path, 0, 0, 1)
	dataMsg = append(dataMsg, byte(0))
	if !isFile {
		monitor.RemoveWatcer(path)
	}
	delete(files, path)
}

func copyAll() map[string][]byte {
	allFiles := make(map[string][]byte)
	err := filepath.WalkDir("destFolder", func(path string, info fs.DirEntry, err error) error {
		// fmt.Printf("visited file or dir: %q", path)
		// fmt.Println("   ", filepath.Ext(path))
		pth := strings.Replace(path, "\\", "/", -1)
		// fmt.Println(pth)

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
	dest := fmt.Sprintf("server/%s", ep[1])
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

	fmt.Println("Msg sent as chunks of 512 bytes")
	fmt.Println("==========================================")
	// client.WriteBytes512(data)
	fmt.Println("==========================================")

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
