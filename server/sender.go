package server

import (
	"github.com/mrasu/malsf/members"
	"fmt"
	"github.com/mrasu/malsf/structs"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"sync"
	"github.com/mrasu/malsf/util"
)

type sender struct {
	Name string
	ServiceName string

	lastId int32
	mu sync.Mutex
}

func (s *sender) send(m *members.Member, message *structs.Message) (chan *structs.Reaction, chan error) {
	reactChan := make(chan *structs.Reaction)
	errChan := make(chan error)
	go func() {
		fmt.Println(m.Address())
		conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
		if err != nil {
			errChan <- err
		}
		defer conn.Close()

		c := structs.NewActionServiceClient(conn)

		act := &structs.Action{
			NodeName: s.Name,
			Service: s.ServiceName,
			Id: s.incrementId(),
			Message: message.Message,
			Type: message.MessageType,
		}
		util.LogAction(s.Name, act.Id, fmt.Sprintf("Send (%s): %s", act.Type, act.Message))
		r, err := c.Notify(context.Background(), act)
		if err != nil {
			errChan <- err
			return
		}

		reactChan <- r
	}()

	return reactChan, errChan
}

func (s *sender) sendAsync(m *members.Member, message *structs.Message) (*structs.Reaction, error) {
	rch, ech := s.send(m, message)

	select {
	case r := <-rch:
		return r, nil
	case err := <- ech:
		return nil, err
	}
}


func (s *sender) incrementId() int32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastId += 1
	fmt.Printf("%s: %d\n", s.Name, s.lastId)
	return s.lastId
}