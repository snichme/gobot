package main

import (
	"fmt"
	"time"
)

type Robot struct {
	Name   string
	tasks  []Task
	config map[string]string
}

func (l Robot) Write(s Query) {
	for _, task := range l.tasks {
		if task.CanHandle(s) {
			c := task.DoHandle(s)
			s.Client.Recieve(c)
			return
		}
	}
}

func main() {

	tasks := []Task{
		NewHiTask(),
		NewSimpsonsTask(),
		&TestTask{},
	}
	config := map[string]string{
		"TCP_PORT": "3333",
	}

	robot := Robot{
		Name:   "GoBot",
		tasks:  tasks,
		config: config,
	}

	tcpClient := NewTCPClient(robot)
	go tcpClient.Start()

	cliClient := NewCliClient(robot)
	time.Sleep(time.Second * 1)
	cliClient.Start()

	fmt.Println("Power down")
}
