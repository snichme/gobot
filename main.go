package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func readConfig(filename string) (config RobotConfig) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Cannot parse config json: %v\n", err)
		os.Exit(1)
	}
	return
}

func getTasks() []Task {
	return []Task{
		NewHiTask(),
		NewServerStatusTask(),
		NewSimpsonsTask(),
		//NewSimpsonsTask(),
	}
}

func main() {
	config := readConfig("./config.json")
	robot := NewRobot(config, getTasks())

	tcpClient := NewTCPClient(robot)
	go tcpClient.Start()
	restClient := NewRestClient(robot)
	go restClient.Start()

	wsClient := NewWebsocketClient(robot)
	go wsClient.Start()
	cliClient := NewCliClient(robot)
	time.Sleep(time.Second * 1)
	cliClient.Start()

	fmt.Println("Power down")
}
