package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CliClient struct {
	robot Robot
}

func blue(v string) string {
	return "\033[34m" + v + "\033[0m"
}
func grey(v string) string {
	return "\033[37m" + v + "\033[0m"
}

func (cc CliClient) Start() {
	reader := bufio.NewReader(os.Stdin)
	queryContext := QueryContext{
		Username: "CLI User",
		Group:    "admin", // Cli user are admin becuase it runs on the same machine as the bot
	}
	queryRobot := func(message string) {
		q := Query{
			Statement: message,
			Context:   queryContext,
		}
		if found, c := cc.robot.Query(q); found {
			for answer := range c {
				fmt.Printf("%s > %s\n", blue(cc.robot.Name()), grey(string(answer)))
			}
		}
	}
	for {
		fmt.Fprintf(cc.robot, "%s > ", blue("Admin"))
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, " \n")
		if text == "quit" {
			break
		}
		queryRobot(text)
	}

}

func NewCliClient(robot Robot) *CliClient {
	return &CliClient{
		robot: robot,
	}
}
