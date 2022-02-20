package main

import (
	"flag"
	"fmt"
	"glowing-telegram/client"
	cloud "glowing-telegram/cloud"
	"os"
	"os/exec"
	"time"
)

// ------------------------------------
// TODO
//
// Directory monitor
// File comparing
// Proper file transfer
// File storage in cloud
//
// ------------------------------------

var app string

func init() {
	flag.StringVar(&app, "app", "none", "Define the type of app to start")
	flag.Parse()
}

func main() {

	switch app {
	case "cloud":
		fmt.Println("Dette er en cloud")
		fmt.Println("------------------------------------------")
		cloud.Cloud()
	case "client":
		fmt.Println("Dette er en client")
		fmt.Println("------------------------------------------")
		time.Sleep(8 * time.Second)
		conn := client.Client(20013)
		defer conn.Close()
		time.Sleep(1 * time.Second)
		client.WriteString(conn, "Yay det funket!")
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
