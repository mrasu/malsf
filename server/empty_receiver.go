package server

import (
	"golang.org/x/net/context"
	"github.com/mrasu/malsf/structs"
	"github.com/mrasu/malsf/util"
)

type EmptyReceiver struct {
	name string
	service string
}

func NewEmptyReceiver(name string, service string) *EmptyReceiver {
	return &EmptyReceiver{
		name: name,
		service: service,
	}
}

func (e *EmptyReceiver) Name() string {
	return e.name
}

func (e *EmptyReceiver) Service() string {
	return e.service
}

func (e *EmptyReceiver) Receive(ctx context.Context, action *structs.Action) (*structs.Message, error) {
	util.LogAction(e.Name() + "<EmptyReceiver>", action.Id, action.Message)
	return nil, nil
}