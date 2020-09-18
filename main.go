package main

import (
	"flag" // used to get command line arguments
	"fmt"
	"log"
	"time"

	"github.com/kairemor/chat-rpc/client"
	"github.com/kairemor/chat-rpc/server"
)

var t = time.Now().Format
var k = time.Kitchen

func main() {
	var process string
	flag.StringVar(&process, "process", "", "")
	flag.Parse()

	if process != "server" && process != "client" {
		log.Fatal("error, -process must be server or client")
	} else {
		fmt.Println("initiating", process, "process...")
		if process == "server" {
			srv := new(server.ChatServer)
			srv.Serve()
			defer srv.End()
		} else {
			cli := new(client.ChatClient)
			cli.Create()
			cli.Handle()
		}
	}
}
