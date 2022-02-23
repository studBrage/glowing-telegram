package main

import (
	"flag"
	"fmt"
	comp "glowing-telegram/comparator"
	"glowing-telegram/monitor"
	proto "glowing-telegram/protocols"
	"io/fs"
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

func init() {
	flag.StringVar(&app, "app", "none", "Define the type of app to start")
	flag.Parse()
	// files = copyAll()
}

func main() {

	// switch app {
	// case "cloud":
	// 	fmt.Println("Dette er en cloud")
	// 	fmt.Println("------------------------------------------")
	// 	cloud.Cloud()
	// case "client":
	// 	fmt.Println("Dette er en client")
	// 	fmt.Println("------------------------------------------")
	// 	time.Sleep(8 * time.Second)
	// 	conn := client.Client(20013)
	// 	defer conn.Close()
	// 	time.Sleep(1 * time.Second)
	// 	client.WriteString(conn, "Yay det funket!")
	// case "init":
	// 	cloud := exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go", "-app=\"cloud\"").Run()
	// 	client := exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go", "-app=\"client\"").Run()
	// 	if cloud != nil || client != nil {
	// 		fmt.Println(cloud.Error())
	// 		fmt.Println(client.Error())
	// 	}
	// 	os.Exit(3)
	// }

	eventChannel = make(chan fsnotify.Event, 1)

	go monitor.Watch("destFolder", eventChannel)
	go Moni(eventChannel)

	files = copyAll()

	time.Sleep(500 * time.Millisecond)

	for p, b := range files {
		fmt.Println(p)
		fmt.Println(b)
	}
	fmt.Println("------------------------------------------")

	for {

	}

}

func Moni(c chan fsnotify.Event) {
	// var prev string
	// var current string
	var isFile bool
	for {
		select {
		case e := <-c:
			target := filepath.Base(e.Name)
			targetType := filepath.Ext(target)
			if targetType == "" {
				isFile = false
			}
			if e.Op.String() == "WRITE" && targetType == "" {
				continue
			}
			action := e.Op.String()
			path := strings.Replace(e.Name, "\\", "/", -1)

			switch action {
			case "CREATE":
				doCreate(action, isFile, path)
				msgHandler(infoMsg, dataMsg)
			case "WRITE":
				doWrite(action, isFile, path, files[path])
			case "REMOVE":
				doRemove(action, isFile, path)
			case "RENAME":
				newName := <-c
				newP := strings.Replace(newName.Name, "\\", "/", -1)
				fmt.Println("RENAMED, old name: ", path, "new name: ", newP)
				doRename(action, isFile, path, newP)
			case "CHMOD":
			default:
				continue
			}

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
}

func doWrite(e string, isFile bool, path string, backupFile []byte) {
	targetBytes := comp.DecodeFile(path)
	delta, largest, ext := comp.FindDelta(backupFile, targetBytes)
	infoMsg = proto.BuildInfo(e, isFile, path, largest, len(delta)*2, len(ext))
	dataMsg = proto.BuildData(delta, ext)

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

	monitor.RemoveWatcer(path)
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
	ep, _, _ := proto.RecieveInfo(info)
	dest := fmt.Sprintf("server/%s", ep[1])
	fmt.Println(dest)
	comp.CopyToFile(dest, data)
}
