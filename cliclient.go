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

func (cc CliClient) Recieve(in <-chan Answer) {
	for answer := range in {
		fmt.Printf("%s > %s\n", blue(cc.robot.Name), grey(string(answer)))
	}
}
func (cc CliClient) Start() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s > ", blue("User"))
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, " \n")
		if text == "quit" {
			break
		}
		cc.robot.Write(Query{
			Statement: text,
			Client:    cc,
		})
	}

}

func NewCliClient(robot Robot) *CliClient {
	return &CliClient{
		robot: robot,
	}
}
