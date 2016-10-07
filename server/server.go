package server

import (
	"fmt"
	"io"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/KanybekMomukeyev/grpc_chat/proto"
	"github.com/KanybekMomukeyev/grpc_chat/utils"
	"flag"
)

var usersLock = &sync.Mutex{}

var usersMap = make(map[string]chan pb.Message, 100) //["key":"channel of Message"]

type chatServer struct{}

func newChatServer() *chatServer {
	return &chatServer{}
}

func (s *chatServer) TransferMessage(stream pb.Chat_TransferMessageServer) error {

	clientIdentification, err := stream.Recv()

	var clientName string

	clientMailbox := make(chan pb.Message, 100)

	if err != nil {
		return err
	}

	fmt.Print("TransferMessage received message \n")
	fmt.Printf("Sender: %s\n",clientIdentification.Sender)
	fmt.Printf("Text: %s\n",clientIdentification.Text)
	fmt.Printf("Register: %t\n",clientIdentification.Register)

	if clientIdentification.Register {
		clientName = clientIdentification.Sender
		if hasListener(clientName) {
			fmt.Print("name already exists error\n")
			return fmt.Errorf("name already exists")
		}
		addListener(clientName, clientMailbox)
	} else {
		fmt.Print("need to register first error\n")
		return fmt.Errorf("need to register first")
	}

	fmt.Print("passed register\n")

	clientMessages := make(chan pb.Message, 100)
	go listenToClient(stream, clientMessages)

	for {
		select {
		case messageFromClient := <-clientMessages:
			fmt.Print("broadcast(clientName, messageFromClient)\n")
			broadcast(clientName, messageFromClient)
		case messageFromOthers := <-clientMailbox:
			fmt.Print("stream.Send(&messageFromOthers)\n")
			stream.Send(&messageFromOthers)
		}
	}
}


func addListener(name string, msgQ chan pb.Message) { // string: "channel of Message"
	usersLock.Lock()
	defer usersLock.Unlock()
	usersMap[name] = msgQ
	fmt.Print("addListener\n")
}

func removeListener(name string) {
	usersLock.Lock()
	defer usersLock.Unlock()
	delete(usersMap, name)
}

func hasListener(name string) bool {
	usersLock.Lock()
	defer usersLock.Unlock()
	_, exists := usersMap[name]
	return exists
}

func broadcast(sender string, msg pb.Message) { // string:Message
	usersLock.Lock()
	defer usersLock.Unlock()

	// for key, value := range usersMap {}
	// send message to all channels, not for own
	//

	for key, channelOfMessage := range usersMap {
		if key != sender {
			channelOfMessage <- msg
		}
	}
}

// messages <- can only be used to send Message ||--> https://golang.org/ref/spec#Channel_types
// stream Chat_TransferMessageServer used for stream receive/Sende
// this will listen from clients
func listenToClient(stream pb.Chat_TransferMessageServer, messages chan<- pb.Message) {
	for {
		msg, err := stream.Recv()

		print("stream.Recv()\n")
		s := fmt.Sprintf("msg.Sender => %s|  msg.Text => %s| msg.Register=>%t msg.Disconnect=>%t\n", msg.Sender, msg.Text, msg.Register, msg.Disconnect)
		fmt.Println(s)

		if err == io.EOF {
			// ?
		}
		if err != nil {
			// ??
		}
		messages <- *msg
	}
}

var (
	port       = flag.Int("port", 10000, "The server port")
)

func Serve(address string, secure bool) error {

	//lis, err := net.Listen("tcp", address)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	if secure {
		creds, err := credentials.NewServerTLSFromFile(utils.ConfigString("TLS_CERT"), utils.ConfigString("TLS_KEY"))
		if err != nil {
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServer(grpcServer, newChatServer())
	grpcServer.Serve(lis)
	return nil
}

//go run main.go serve 10000