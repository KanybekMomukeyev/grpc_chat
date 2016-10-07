package main

import (
"bufio"
"fmt"
"os"

"google.golang.org/grpc"
"google.golang.org/grpc/credentials"

"golang.org/x/net/context"
pb "github.com/KanybekMomukeyev/grpc_chat/proto"
)

func listenToClient(sendQ chan pb.Message, reader *bufio.Reader, name string) {
	for {
		msg, _ := reader.ReadString('\n')
		sendQ <- pb.Message{Sender: name, Text: msg}
	}
}

func receiveMessages(stream pb.Chat_TransferMessageClient, mailbox chan pb.Message) {
	for {
		msg, _ := stream.Recv()
		mailbox <- *msg
	}
}

var (
	serverAddr = "localhost:10000"
)

func Connect(address, nickname string, secure bool) error {
	var opts []grpc.DialOption
	if secure {
		var sn string
		//TODO: server hotname override
		var creds credentials.TransportCredentials
		//var creds credentials.TransportAuthenticator
		//TODO: if self-signed
		creds = credentials.NewClientTLSFromCert(nil, sn)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	//conn, err := grpc.Dial(address, opts...)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		return err
	}

	defer conn.Close()

	client := pb.NewChatClient(conn)

	stream, err := client.TransferMessage(context.Background())
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	clientName, err := reader.ReadString('\n')

	if err != nil {
		return err
	}
	stream.Send(&pb.Message{Sender: clientName, Register: true})

	mailBox := make(chan pb.Message, 100)
	go receiveMessages(stream, mailBox)

	sendQ := make(chan pb.Message, 100)
	go listenToClient(sendQ, reader, clientName)

	for {
		select {
		case toSend := <-sendQ:
			stream.Send(&toSend)

		case received := <-mailBox:
			fmt.Printf("%s> %s", received.Sender, received.Text)
		}
	}
	return nil
}

//go run main.go connect koke no

func main() {
	Connect("","Koke",false)
}

