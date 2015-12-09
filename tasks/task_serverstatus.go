package tasks

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/snichme/gobot/types"
)

type ServerStatusTask struct {
	queryRegexp *regexp.Regexp
}

func NewServerStatusTask() *ServerStatusTask {
	var queryRegexp = regexp.MustCompile(`([wW]hats up\?)|(how are you\??)`)

	return &ServerStatusTask{
		queryRegexp: queryRegexp,
	}
}

func (task ServerStatusTask) Name() string {
	return "ServerStatusTask"
}
func (task ServerStatusTask) HelpText() string {
	return "Tells you how the server where the bot lives are doing"
}

func (task ServerStatusTask) Handle(query types.Query) (bool, <-chan types.Answer) {
	if !task.queryRegexp.MatchString(query.Statement) {
		return false, nil
	}
	c1 := make(chan types.Answer)
	go func(cmd string) {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		c1 <- types.Answer(fmt.Sprintf("My status is: %s", strings.Trim(string(out), " \n")))
		close(c1)
	}("uptime")
	return true, c1

}
