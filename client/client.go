package client

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kairemor/chat-rpc/shared"
)

var t = time.Now().Format
var k = time.Kitchen

// ChatClient type
type ChatClient struct {
	Name string
	Conn *rpc.Client
}

// Globale variables
const (
	JoinRoom   = "Server.JoinRoom"
	Delete     = "Server.DeleteUser"
	FindAll    = "Server.GetAllMessages"
	NewMessage = "Server.NewMessage"
	Exit       = "/exit"
	Help       = "/help"
	WhiteSpace = ""
	host       = "localhost:"
	port       = "4000"
)

// End chat
func (c *ChatClient) End() {
	if c.Conn != nil {
		args := &shared.DeleteUserArgs{Name: c.Name}
		resp := &shared.UserResp{}
		err := c.Conn.Call(FindAll, args, resp)
		if err != nil {
			log.Fatal("error removing user")
		}
		c.Conn.Close()
	}
}

// Create new
func (c *ChatClient) Create() {
	if c.Conn == nil {
		var err error
		c.Conn, err = rpc.Dial("tcp", host+port) // register RPC c
		if err != nil {
			log.Fatal("c error on Dial")
		}
	}
}

//FindAllMessages ...
func (c *ChatClient) FindAllMessages(DoneChan chan int) {
	go func() {
		<-DoneChan
		return
	}()
	for {
		args := &shared.FindAllMessagesArgs{User: c.Name}
		resp := &shared.FindAllMessagesResp{}
		err := c.Conn.Call(FindAll, args, resp)
		if err != nil {
			log.Fatal(err)
		}
		for _, message := range resp.Messages {
			fmt.Println(message)
		}
	}
}

// NewMessage to chat
func (c *ChatClient) NewMessage(message *shared.Message) {
	args := &shared.NewMessageArgs{Message: message}
	resp := &shared.NewMessageResp{}
	err := c.Conn.Call(NewMessage, args, resp)
	if err != nil {
		log.Fatal(err)
	}
	if resp.Code == -1 {
		fmt.Printf("no user %s \n", message.Receiver)
	}
}

//Handle ...
func (c *ChatClient) Handle() {
	DoneChan := make(chan int)
	defer c.End()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.register(&wg)
	wg.Wait()

	go c.FindAllMessages(DoneChan)
	go c.listen(DoneChan)
	<-DoneChan
}

func (c *ChatClient) register(wg *sync.WaitGroup) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Put your name: ")
		name, _ := reader.ReadString('\n')
		nameNoLine := name[:len(name)-1]
		if strings.IndexByte(nameNoLine, ' ') != -1 {
			fmt.Println("put a valid name")
			continue
		}
		fmt.Printf("registering user: %s", name)
		args := &shared.UserArgs{Name: nameNoLine}
		resp := &shared.UserResp{}
		err := c.Conn.Call(JoinRoom, args, resp)
		if err != nil {
			log.Fatal(err)
		}
		if resp.Code == -1 {
			fmt.Printf("no registered %s is taken\n", name[:len(name)-1])
			continue
		} else {
			fmt.Printf("welcome %s", name)
			c.Name = nameNoLine
			break
		}
	}
	wg.Done()
}

func (c *ChatClient) listen(DoneChan chan int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		messageNoLine := message[:len(message)-1]
		switch messageNoLine {
		case WhiteSpace:
			continue
		case Help:
			fmt.Println("Write a message to send room")
			fmt.Println("send PM message by -> @receiver <message>")
			fmt.Println("to exit this chat type -> /exit")
		case Exit:
			DoneChan <- 1
			break
		default:
			if messageNoLine[0] == '@' {
				directMessage := messageNoLine[1:]
				endOfReceiver := strings.IndexByte(directMessage, ' ')
				if endOfReceiver == -1 || endOfReceiver == len(directMessage)-1 {
					fmt.Println("message no valid no content")
					continue
				}

				receiver := directMessage[:endOfReceiver]
				message := directMessage[endOfReceiver+1:]
				checkIsWhitespace := func(s string) bool {
					for i := 0; i < len(s); i++ {
						if s[i] != ' ' {
							return false
						}
					}
					return true
				}
				if !checkIsWhitespace(message) {
					MessageStruct := &shared.Message{Sender: c.Name, Receiver: receiver, Message: message, Time: t(k)}
					c.NewMessage(MessageStruct)
				} else {
					fmt.Println("message must contain non-whitespace characters")
				}
			} else {
				MessageStruct := &shared.Message{Sender: c.Name, Receiver: "", Message: messageNoLine, Time: t(k)}
				c.NewMessage(MessageStruct)
			}
		}
	}
}
