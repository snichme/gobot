package main

import "strings"

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
func (task HiTask) Handle(query Query) (bool, <-chan Answer) {

	if !(strings.Contains(query.Statement, "hello") || strings.Contains(query.Statement, "hi")) {
		return false, nil
	}
	c1 := make(chan Answer, 1)
	go func() {
		c1 <- "Hi there!"
		close(c1)
	}()
	return true, c1
}

func (task HiTask) CanHandle(query Query) bool {
	return strings.Contains(query.Statement, "hello") || strings.Contains(query.Statement, "hi")
}

func (task HiTask) DoHandle(query Query) <-chan Answer {
	c1 := make(chan Answer, 1)
	c1 <- "Hi there!"
	close(c1)
	return c1
}
