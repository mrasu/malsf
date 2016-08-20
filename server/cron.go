package server

import (
	"time"
	"github.com/mrasu/malsf/discover"
	"fmt"
)

type Cron struct {
	interval int
}

func NewTick() *Cron {
	return &Cron{
		interval: 1,
	}
}

func (c *Cron) Start() error {
	t := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-t.C:
			n := discover.NewNodeDiscoverer()
			members, err := n.GetMembersByTag("client")
			if err != nil {
				return err
			}
			for _, member := range members {
				fmt.Println(member.Addr())
				rch, ech := member.Send("hello world", "Memory")

				select {
				case <-rch:
					//do nothing
				case err := <- ech:
					return err
				}
			}
		}
	}
}