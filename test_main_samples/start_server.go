package main

import (
	"log"
	"github.com/KanybekMomukeyev/grpc_chat/server"
)

func main() {
	err := server.Serve("10000", false)
	if err != nil {
		log.Fatalln("Error: %s", err.Error())
	}
}
