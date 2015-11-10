package main

type Client interface {
	Start()
}

type QueryContext struct {
	Username string
	Group    string
}

type Query struct {
	Statement string
	Context   QueryContext
}
