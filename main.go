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

	for p, b := range files {
		fmt.Println(p)
		fmt.Println(b)
	}

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
			case "WRITE":
				doWrite(action, isFile, path, files[path])
			case "REMOVE":
				doRemove(action, isFile, path)
			case "RENAME":
				newName := <-c
				newP := strings.Replace(newName.Name, "\\", "/", -1)
				doRename(action, isFile, path, newP)
			case "CHMOD":
			default:
				continue
			}

			for p, b := range files {
				fmt.Println(p)
				fmt.Println(b)
			}
			fmt.Println("------------------------------------------")

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
	targetFile, err := comp.OpenFile(path)
	errHandler(err)
	targetBytes := comp.DecodeFile(targetFile)
	infoMsg = proto.BuildInfo(e, isFile, path, 0, 0, len(targetBytes))
	dataMsg = targetBytes
	if isFile {
		files[path] = targetBytes
		monitor.AddWatcher(path)
	} else {
		files[path] = []byte("FOLDER")
	}
}

func doWrite(e string, isFile bool, path string, backupFile []byte) {
	targetFile, err := comp.OpenFile(path)
	errHandler(err)
	targetBytes := comp.DecodeFile(targetFile)
	delta, largest, ext := comp.FindDelta(backupFile, targetBytes)
	infoMsg = proto.BuildInfo(e, isFile, path, largest, len(delta)*2, len(ext))
	dataMsg = proto.BuildData(delta, ext)

	files[path] = comp.UpdateChange(backupFile, largest, delta, ext)
}

func doRename(e string, isFile bool, path string, new string) {
	infoMsg = proto.BuildInfo(e, isFile, path, 0, 0, 0)
	dataMsg = []byte(new)

	// monitor.RemoveWatcer(path)
	// monitor.AddWatcher(new)
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
			fil, err := comp.OpenFile(pth)
			errHandler(err)
			allFiles[pth] = comp.DecodeFile(fil)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}
	return allFiles
}
