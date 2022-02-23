package monitor

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func init() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Watch(dir string, c chan fsnotify.Event) {

	defer watcher.Close()

	fmt.Println("Watching following directories")
	fmt.Println("======================================")
	err := filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			pth := strings.Replace(path, "\\", "/", -1)
			fmt.Println(pth)
			return watcher.Add(pth)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("======================================")

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// fmt.Printf("EVENT! %#v\n", event)
				// fmt.Println(event.String())
				// fmt.Println(event.Op.String() == "RENAME")
				c <- event
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err.Error())
			}
		}
	}()

	<-done
}

func AddWatcher(path string) error {
	return watcher.Add(path)
}

func RemoveWatcer(path string) error {
	return watcher.Remove(path)
}
