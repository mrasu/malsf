package structs

import (
	"golang.org/x/net/context"
)

type ServerAct interface{
	Name() string
	Service() string
	Receive(ctx context.Context, action *Action) (*Message, error)
}
