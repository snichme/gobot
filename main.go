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

func main() {
	config := readConfig("./config.json")

	tasks := []Task{
		NewHiTask(),
		NewServerStatusTask(),
		NewSimpsonsTask(),
		NewJenkinsListTask(config.Settings["jenkinsURL"]),
		NewJenkinsRunBuildTask(config.Settings["jenkinsURL"]),
		NewJenkinsStatusTask(config.Settings["jenkinsURL"]),
		NewHackerNewsTopTask(),
	}
	robot := NewRobot(config, tasks)

	tcpClient := NewTCPClient(robot)
	go tcpClient.Start()
	restClient := NewRestClient(robot)
	go restClient.Start()

	// wsClient := NewWebsocketClient(robot)
	// go wsClient.Start()

	cliClient := NewCliClient(robot)
	time.Sleep(time.Second * 1)
	cliClient.Start()

	fmt.Println("Power down")
}
