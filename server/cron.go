package server

import (
	"time"
	"github.com/mrasu/malsf/discover"
	"fmt"
	"github.com/mrasu/malsf/structs"
)

type Cron struct {
	interval int
	cronAct structs.CronAct
	*sender
}

func NewCron(cronAct structs.CronAct) *Cron {
	return &Cron{
		sender: &sender{
			Name: cronAct.Name(),
			ServiceName: cronAct.Service(),
		},
		interval: 1,
		cronAct: cronAct,
	}
}

func (c *Cron) Start() error {
	t := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-t.C:
			message, err := c.cronAct.CallCron()
			if err != nil {
				return err
			} else if message == nil {
				continue
			}

			n := discover.NewNodeDiscoverer()
			for _, service := range message.ToServices {
				members, err := n.GetMembersByTag(service)
				if err != nil {
					return err
				}

				for _, member := range members {
					fmt.Println(member.Address())
					rch, ech := c.send(member, message)

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
}
