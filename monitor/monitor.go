package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Hello")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer watcher.Close()

	done := make(chan bool)
	eventChan := make(chan string)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// fmt.Printf("EVENT! %#v\n", event)
				// fmt.Println(event.String())
				eventChan <- event.String()
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err.Error())
			}
		}
	}()
	if err := watcher.Add("../destFolder"); err != nil {
		fmt.Println(err.Error())
	}

	<-done
}
