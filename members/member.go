package members

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"fmt"
)

type Member struct {
	name string
	addr grpc.Address
}

func NewMember(name string, addr string, port int) *Member {
	return &Member{
		name: name,
		addr: grpc.Address{
			Addr: fmt.Sprintf("%s:%d", addr, port),
		},
	}
}

func (m *Member) Addr() string {
	return m.addr.Addr
}

func (m *Member) Send(message string, messageType string) (chan *Reaction, chan error) {
	reactChan := make(chan *Reaction)
	errChan := make(chan error)
	go func() {
		conn, err := grpc.Dial(m.Addr(), grpc.WithInsecure())
		if err != nil {
			errChan <- err
		}
		defer conn.Close()

		c := NewActionServiceClient(conn)
		act := &Action{
			NodeName: m.name,
			Id: 1,
			Message: message,
			Type: messageType,
		}
		r, err := c.Notify(context.Background(), act)
		if err != nil {
			errChan <- err
			return
		}
		fmt.Println(r)
		reactChan <- r
	}()

	return reactChan, errChan
}