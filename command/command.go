package client

import (
	"github.com/mrasu/malsf/server"
	"github.com/mrasu/malsf/structs"
)

type Command struct {
	port int
}

func NewCommand(port int) *Command {
	return &Command{
		port: port,
	}
}

func (c *Command) StartManager(s structs.ServerAct) {
	c.startServer(s)
}

func (c *Command) StartCron(cronAct structs.CronAct) {
	c.startServer(cronAct)

	go func() {
		c := server.NewCron(cronAct)
		panic(c.Start())
	}()
}

func (c *Command) startServer(s structs.ServerAct) {
	go func() {
		server.StartServer(c.port, s)
	}()
}