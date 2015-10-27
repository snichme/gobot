package main

type Query struct {
	Statement string
	Client    ClientPrint
}

type ClientPrint interface {
	Recieve(in <-chan Answer)
}

type Client interface {
	Start()
}
