package structs

import (
	"golang.org/x/net/context"
)

type CronAct interface{
	Name() string
	Service() string
	CallCron() (*Message, error)
	Receive(ctx context.Context, action *Action) (*Message, error)
}
