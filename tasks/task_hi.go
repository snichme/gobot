package tasks

import (
	"strings"

	"github.com/snichme/gobot/types"
)

type HiTask struct {
}

func NewHiTask() *HiTask {
	return &HiTask{}
}

func (task HiTask) Name() string {
	return "HiTask"
}
func (task HiTask) HelpText() string {
	return "Welcomes you"
}
func (task HiTask) Handle(query types.Query) (bool, <-chan types.Answer) {

	if !(strings.Contains(query.Statement, "hello") || strings.Contains(query.Statement, "hi")) {
		return false, nil
	}
	c1 := make(chan types.Answer, 1)
	go func() {
		c1 <- "Hi there!"
		close(c1)
	}()
	return true, c1
}
