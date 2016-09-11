package client

import (
	"github.com/mrasu/malsf/server"
	"github.com/mrasu/malsf/structs"
	"time"
)

type Command struct {
	port int
}

func NewCommand(port int) *Command {
	return &Command{
		port: port,
	}
}

func (c *Command) StartManager(name string, service string, receiveFn func(action *structs.Action) (*structs.Message, error)) {
	r := server.NewReceiver(name, service, receiveFn)
	c.startServer(r, make(chan (*structs.Message)))
}

func (c *Command) StartCron(name string, service string, interval time.Duration, cronFn func() (*structs.Message, error)) {
	mch := make(chan (*structs.Message))
	go func() {
		c := server.NewCron(interval, cronFn)
		ech := c.Start(mch)

		err := <-ech
		panic(err)
	}()
	c.startServer(server.NewEmptyReceiver(name, service), mch)
}

func (c *Command) startServer(s structs.ReceiverAct, mch chan (*structs.Message)) {
	go func() {
		server.StartServer(c.port, s, mch)
	}()
}
