package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"

	"github.com/kairemor/chat-rpc/shared"
)

//QueueMessages save messages
type QueueMessages struct {
	Map map[string][]*shared.Message
	sync.Mutex
}

// ChatServer ...
type ChatServer struct {
	Port     int
	Users    []string
	Messages QueueMessages
	listener net.Listener
}

//MaxUsers at most 100 people in the chat room
const MaxUsers = 100

func (c *ChatServer) getIndexByName(name string) int {
	predicate := func(i int) bool {
		return c.Users[i] == name
	}

	return shared.MapSlice(len(c.Users), predicate)
}

//JoinRoom ...
func (c *ChatServer) JoinRoom(request *shared.UserArgs, response *shared.UserResp) error {
	if len(c.Users) == MaxUsers {
		return errors.New("No more user possible try later ")
	}

	if c.getIndexByName(request.Name) != -1 {
		fmt.Println("user with this name already exist: ", request.Name)
		response.Code = -1
	} else {
		c.Users = append(c.Users, request.Name)
		fmt.Println("join room : ", request.Name)
		response.Code = 0
	}
	return nil
}

func (c *ChatServer) delete(name string) {
	for i := 0; i < len(c.Users); i++ {
		if c.Users[i] == name {
			if i < len(c.Users)-1 {
				c.Users = append(c.Users[0:i], c.Users[i+1:]...)
			} else {
				c.Users = c.Users[:i]
			}
		}
	}
}

// DeleteUser from chat
func (c *ChatServer) DeleteUser(request *shared.DeleteUserArgs, response *shared.UserResp) error {
	if c.getIndexByName(request.Name) != -1 {
		c.delete(request.Name)
		response.Code = 0
		return nil
	}
	response.Code = -1
	return errors.New("user you try do delete don't exist")
}

//GetAllMessages find all sended messages
func (c *ChatServer) GetAllMessages(request *shared.FindAllMessagesArgs, response *shared.FindAllMessagesResp) error {
	// lock message list until we have copied messages
	c.Messages.Lock()
	defer c.Messages.Unlock()
	if messages, ok := c.Messages.Map[request.User]; !ok && c.getIndexByName(request.User) == -1 {
		return errors.New("no user")
	} else {
		response.Messages = messages
	}
	c.Messages.Map[request.User] = nil
	return nil
}

//NewMessage send
func (c *ChatServer) NewMessage(request *shared.NewMessageArgs, response *shared.SendMessageResp) error {
	fmt.Println("new message", request.Message)
	c.Messages.Lock()
	if request.Message.Receiver == "" {
		for user := range c.Messages.Map {
			c.Messages.Map[user] = append(c.Messages.Map[user], request.Message)
		}
		response.Code = 0
	} else {
		receiver := request.Message.Receiver
		if c.getIndexByName(receiver) != -1 {
			c.Messages.Map[receiver] = append(c.Messages.Map[receiver], request.Message)
			response.Code = 0
		} else {
			fmt.Println("user you try to send message don't exist ", receiver)
			response.Code = -1
		}
	}
	c.Messages.Unlock()
	return nil
}

// End chat
func (c *ChatServer) End() {
	if c.listener != nil {
		c.listener.Close()
	}
}

func (c *ChatServer) init() {
	c.Port = 4000
	c.Users = make([]string, 1, MaxUsers)
	c.Messages.Map = make(map[string][]*shared.Message, 1)
}

//Serve chat server creation
func (c *ChatServer) Serve() {
	c.init()
	rpc.Register(c)

	var err error
	c.listener, err = net.Listen("tcp", ":"+strconv.Itoa(c.Port))

	if err != nil {
		log.Fatal("server no listening error :", err)
	}

	fmt.Println("ChatServer is running port 4000")
	rpc.Accept(c.listener)
}
