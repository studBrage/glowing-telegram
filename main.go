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
		infoChan := make(chan bool)
		fmt.Println("Dette er en cloud")
		fmt.Println("------------------------------------------")
		go cloud.Cloud(infoChan, incoming)
		go server(infoChan, incoming)
		fmt.Println("Its a go!")

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

	// eventChannel = make(chan fsnotify.Event, 1)

	// go monitor.Watch("destFolder", eventChannel)
	// go Moni(eventChannel)

	// time.Sleep(500 * time.Millisecond)

	// for p, b := range files {
	// 	fmt.Println(p)
	// 	fmt.Println(b)
	// }
	// fmt.Println("------------------------------------------")

	for {

	}

}

func server(infoChan chan bool, incoming chan []byte) {
	current := "IDLE"
	next := ""
	lenChecker := 0
	var sliceDumper []byte
	for {
		isInfo := <-infoChan
		fmt.Println("Fra main:", isInfo)
		fmt.Println("State:", current)
		if isInfo {
			next = "INFO"
		}
		switch current {
		case "INFO":
			// fmt.Println("Getting info")
			infoMsg = <-incoming
			fmt.Println("Info recieved")
			fmt.Println(infoMsg)
			next = "DATA"
		case "DATA":
			lenChecker = proto.ExtractDataLen(infoMsg)
			// isFile := proto.ExtractType(infoMsg)
			chunk := <-incoming
			// fmt.Println("CHUNK:", chunk)
			dataMsg = append(dataMsg, chunk...)
			// fmt.Println("LEN DATA:", len(dataMsg))
			// fmt.Println("DATAMSG:", dataMsg)
			// fmt.Println("LENCHECK:", lenChecker)
			// if !isFile {
			// 	msgHandler(infoMsg)
			// }
			if len(dataMsg) == lenChecker {
				next = "HANDLE"
			}
		case "HANDLE":
			msgHandler(infoMsg, dataMsg)
			infoMsg = sliceDumper
			dataMsg = sliceDumper
			fmt.Println("empty infomsg:", infoMsg)
			next = "IDLE"
		case "IDLE":
			current = next
			continue
		}
		fmt.Println("NEXT:", next)
		current = next

		//fmt.Println("Skjer det ingenting her eller?")

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
			// for p, b := range files {
			// 	fmt.Println(p)
			// 	fmt.Println(b)
			// }
			// fmt.Println("------------------------------------------")

			// switch action {
			// case "RENAME":
			// 	newName := <-c
			// 	<-c
			// 	<-c
			// 	_, fil := filepath.Split(newName.Name)
			// 	fmt.Println("Rename", fil, "to", fil)
			// case "CREATE":
			// 	if targetType == "" {
			// 		fmt.Println("Created new folder", target)
			// 	} else {
			// 		fmt.Println("Created file ", target)
			// 	}
			// }

			// fmt.Println(e.Op)
			// if filepath.Ext(fil) == "" {
			// 	fmt.Println(" MAPPE", e.Name)
			// } else {
			// 	fmt.Println(fil, " path:", e.Name)
			// }
			// fmt.Println()
		default:
			client.WriteString(conn, "NONE")
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
	infoMsg = proto.BuildInfo(e, isFile, path, largest, len(delta), len(ext))
	dataMsg = proto.BuildData(delta, ext)

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

func errHandler(err error) {
	if err != nil {
		panic(err.Error())
	}
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
