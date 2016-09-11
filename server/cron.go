package server

import (
	"github.com/mrasu/malsf/structs"
	"time"
)

type Cron struct {
	interval time.Duration
	cronFn   func() (*structs.Message, error)
}

func NewCron(interval time.Duration, cronFn func() (*structs.Message, error)) *Cron {
	return &Cron{
		interval: interval,
		cronFn:   cronFn,
	}
}

func (c *Cron) Start(mch chan (*structs.Message)) chan (error) {
	ech := make(chan (error))

	go func() {
		t := time.NewTicker(c.interval)

		for {
			select {
			case <-t.C:
				message, err := c.cronFn()
				if err != nil {
					ech <- err
				} else if message == nil {
					continue
				}
				mch <- message
			}
		}
	}()

	return ech
}
