package tasks

import (
	"math/rand"
	"strings"

	"github.com/snichme/gobot/types"
)

type SimpsonsTask struct {
	quotes []string
}

func NewSimpsonsTask() *SimpsonsTask {
	return &SimpsonsTask{
		quotes: []string{
			"Lisa, if I've learned anything, it's that life is just one crushing defeat after another until you just wish Flanders was dead.",
			"Sorry mom, the mob has spoken.",
			"…A little help?",
			"So I said to myself: what would God do in this situation?",
			"The goggles, they do nothing!",
			"And I'm not easily impressed — WOW, A BLUE CAR!",
			"Since the beginning of time, man has yearned to destroy the sun.",
			"Lisa, I'd like to buy your rock.",
		},
	}
}
func (tt SimpsonsTask) Name() string {
	return "Simpsons quotes"
}
func (tt SimpsonsTask) HelpText() string {
	return "Will give a simpsons quote if anyone mentions simpons"
}

func (tt SimpsonsTask) Handle(query types.Query) (bool, <-chan types.Answer) {
	if !strings.Contains(query.Statement, "simpsons") {
		return false, nil
	}
	c1 := make(chan types.Answer)
	rand.Seed(8) // Try changing this number!
	go func() {
		//time.Sleep(time.Millisecond * 200)
		c1 <- types.Answer(tt.quotes[rand.Intn(len(tt.quotes))])
		close(c1)
	}()
	return true, c1
}
