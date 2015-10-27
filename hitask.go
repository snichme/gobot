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
func (task HiTask) CanHandle(query Query) bool {
	return strings.Contains(query.Statement, "hello") || strings.Contains(query.Statement, "hi")
}
func (task HiTask) DoHandle(query Query) <-chan Answer {
	c1 := make(chan Answer, 1)
	c1 <- "Hi there!"
	close(c1)
	return c1
}
