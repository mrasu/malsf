package main

import (
	"fmt"
	"github.com/mrasu/malsf/command"
	"time"
	"os/exec"
	"strings"
	"github.com/mrasu/malsf/structs"
	"os"
	"encoding/json"
)

func main() {
	arg := os.Args[1]
	fmt.Printf("arg: %s\n", arg)

	if arg == "m" {
		c := client.NewCommand(11110)
		c.StartManager("server01", "server", ReceiveToServer)
		ccc := make(chan(int))
		<- ccc
	}
	if arg == "c" {
		c2 := client.NewCommand(11111)
		c2.StartCron("consul01", "client", 1 * time.Second, ReceiveTime)
		time.Sleep(4 * time.Second)
	}
}

func ReceiveToServer(action *structs.Action) (*structs.Message, error) {
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

func ReceiveTime() (*structs.Message, error) {
	out, err := exec.Command("python3", "cron_task.py").Output()
	if err != nil {
		panic(err)
	}

	m := &structs.Message{}
	err = json.Unmarshal(out, m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message: %s\n", m)
	return m, nil
}
