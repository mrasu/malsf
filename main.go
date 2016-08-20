package main

import (
	"time"
	"github.com/mrasu/malsf/command"
)

func main() {
	c := client.NewCommand(client.Manager, 11110)
	c.Start()
	c2 := client.NewCommand(client.Client, 11111)
	c2.Start()

	time.Sleep(5 * time.Second)
}