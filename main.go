package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/snichme/gobot/clients"
	"github.com/snichme/gobot/tasks"
	"github.com/snichme/gobot/types"
)

func readConfig(filename string) (config types.RobotConfig) {
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

	availTasks := []types.Task{
		tasks.NewHiTask(),
		tasks.NewServerStatusTask(),
		tasks.NewSimpsonsTask(),
		tasks.NewXkcdTask(),
		tasks.NewJenkinsListTask(config.Settings["jenkinsURL"]),
		tasks.NewJenkinsRunBuildTask(config.Settings["jenkinsURL"]),
		tasks.NewJenkinsStatusTask(config.Settings["jenkinsURL"]),
		tasks.NewHackerNewsTopTask(),
	}
	robot := NewRobot(config, availTasks)

	tcpClient := clients.NewTCPClient(robot)
	go tcpClient.Start()
	restClient := clients.NewRestClient(robot)
	go restClient.Start()

	// ircClient := clients.NewIRCClient(robot)
	// go ircClient.Start()
	// wsClient := NewWebsocketClient(robot)
	// go wsClient.Start()

	cliClient := clients.NewCliClient(robot)
	time.Sleep(time.Second * 1)
	cliClient.Start()
	fmt.Println("Power down")
}
