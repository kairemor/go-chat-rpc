package shared

import (
	"fmt"
)

// MapSlice ...
func MapSlice(length int, check func(i int) bool) int {
	for i := 0; i < length; i++ {
		if check(i) {
			return i
		}
	}
	return -1
}

// Message ...
type Message struct {
	Sender, Receiver, Message, Time string
}

//UserArgs ...
type UserArgs struct {
	Name string
}

//UserResp ...
type UserResp struct {
	Code int
}

// DeleteUserArgs ...
type DeleteUserArgs struct {
	Name string
}

// DeleteUserResp ...
type DeleteUserResp struct {
	Code int
}

// NewMessageArgs ...
type NewMessageArgs struct {
	Message *Message
}

//NewMessageResp ...
type NewMessageResp struct {
	Code int
}

//FindAllMessagesArgs ...
type FindAllMessagesArgs struct {
	User string
}

//FindAllMessagesResp ...
type FindAllMessagesResp struct {
	Messages []*Message // messages unseen by user
}

//String...
func (m *Message) String() string {
	return fmt.Sprintf("[%s @ %s]: %s", m.Sender, m.Time, m.Message)
}
