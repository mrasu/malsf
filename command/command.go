package client

import "github.com/mrasu/malsf/server"

type Command struct {
	t ClientType
	port int
}

func NewCommand(t ClientType, port int) *Command {
	return &Command{
		t: t,
		port: port,
	}
}

type ClientType int

const (
	Manager ClientType = iota
	Client
)

func (c *Command) Start() {
	switch c.t {
	case Manager:
		c.startAsManager()
	case Client:
		c.startAsClient()
	}
}

func (c *Command) startAsManager() {
	c.startServer()
}

func (c *Command) startAsClient() {
	c.startServer()

	go func() {
		c := server.NewTick()
		panic(c.Start())
	}()
}

func (c *Command) startServer() {
	go func() {
		server.StartServer(c.port)
	}()
}