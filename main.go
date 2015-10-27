package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var config map[string]string
	json.Unmarshal(file, &config)

	tasks := []Task{
		NewHiTask(),
		NewSimpsonsTask(),
		&TestTask{},
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
