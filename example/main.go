package main

import (
	"fmt"
	"github.com/mrasu/malsf/command"
	"time"
	"os/exec"
	"strings"
	"strconv"
	"github.com/mrasu/malsf/structs"
	"golang.org/x/net/context"
)

func main() {
	fmt.Println("hello")

	c := client.NewCommand(11110)
	c.StartManager(&serverImpl{})
	c2 := client.NewCommand(11111)
	c2.StartCron(&cronImpl{})

	time.Sleep(5 * time.Second)
}

type serverImpl struct {}

func (c *serverImpl) Name() string {
	return "server01"
}

func (c *serverImpl) Service() string {
	return "server"
}

func (c *serverImpl) Tag() string {
	return "server"
}

func (s *serverImpl) Receive(ctx context.Context, action *structs.Action) (*structs.Message, error) {
	fmt.Println("work work")
	if action.Type == "Memory" {
		fmt.Println("Memory is leeking!!")
		ips := strings.Split(action.Message, ",")
		if len(ips) > 0 {

			return (&structs.Message{
				MessageType: "Kill Process",
				Message: ips[0],
			}), nil
		}
	}

	return nil, nil
}

type cronImpl struct {}

func (c *cronImpl) Name() string {
	return "consul01"
}

func (c *cronImpl) Service() string {
	return "client"
}

func(c *cronImpl) CallCron() (*structs.Message, error) {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		panic(err)
	}
	body := string(out)

	chromes := [][]string{}
	for _, line := range strings.Split(body, "\n") {
		line_contents := strings.Fields(line)
		if len(line_contents) >= 10 && strings.Contains(line_contents[10], "chrome") {
			chromes = append(chromes, line_contents)
		}
	}

	heavy_process_ids := []string{}
	for _, chrome := range chromes {
		memory, err := strconv.Atoi(chrome[5])
		//if err == nil && memory > 2000000000 {
		if err == nil && memory > 20000 {
			pid := chrome[1]
			heavy_process_ids = append(heavy_process_ids, pid)
		}
	}

	fmt.Println(heavy_process_ids)
	return &structs.Message{
		ToServices: []string{"server"},
		MessageType: "Memory",
		Message: strings.Join(heavy_process_ids, ","),
	}, nil
}

func (s *cronImpl) Receive(ctx context.Context, action *structs.Action) (*structs.Message, error) {
	fmt.Println("no work")
	return nil, nil
}