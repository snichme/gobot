package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/snichme/gobot/types"
)

type XkcdTask struct {
	queryRegexp *regexp.Regexp
}

type xkcdResponse struct {
	Title    string `json:"safe_title"`
	ImageUrl string `json:"img"`
}

func NewXkcdTask() *XkcdTask {
	var queryRegexp = regexp.MustCompile(`xkcd`)
	return &XkcdTask{
		queryRegexp: queryRegexp,
	}
}

func (task XkcdTask) Name() string {
	return "Xkcd"
}
func (task XkcdTask) HelpText() string {
	return "Gives you a xkcd ref if someone mentions xkcd"
}

func (task XkcdTask) Handle(query types.Query) (bool, <-chan types.Answer) {
	if !task.queryRegexp.MatchString(query.Statement) {
		return false, nil
	}
	c1 := make(chan types.Answer)

	go func(uri string) {
		defer close(c1)
		resp, err := http.Get(uri)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		//fmt.Printf("resp %s", body)

		var jsonResp xkcdResponse
		err = json.Unmarshal(body, &jsonResp)
		if err != nil {
			return
		}

		c1 <- types.Answer(fmt.Sprintf("You mentioned xkcd, here is a link to the latest: %s (%s)", jsonResp.ImageUrl, jsonResp.Title))

	}("http://xkcd.com/info.0.json")

	return true, c1
}
