package server

import (
	"github.com/mrasu/malsf/structs"
	"golang.org/x/net/context"
)

type Receiver struct {
	name      string
	service   string
	receiveFn func(action *structs.Action) (*structs.Message, error)
}

func NewReceiver(name string, service string, receiveFn func(action *structs.Action) (*structs.Message, error)) *Receiver {
	return &Receiver{
		name:      name,
		service:   service,
		receiveFn: receiveFn,
	}
}

func (r *Receiver) Name() string {
	return r.name
}

func (r *Receiver) Service() string {
	return r.service
}

func (r *Receiver) Receive(ctx context.Context, action *structs.Action) (*structs.Message, error) {
	return r.receiveFn(action)
}
