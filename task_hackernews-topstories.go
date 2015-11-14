package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//HackerNewsTopTask Get the top 5 stories on HackerNews
type HackerNewsTopTask struct {
}

//Name Get the name of the task
func (task HackerNewsTopTask) Name() string {
	return "HackerNewsTask"
}

//HelpText Get a description of the task
func (task HackerNewsTopTask) HelpText() string {
	return "Get the top 5 stories on HackerNews"
}

func (task HackerNewsTopTask) getTopIds() ([]int, error) {
	var jsonResp []int
	uri := "https://hacker-news.firebaseio.com/v0/topstories.json"
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func (task HackerNewsTopTask) getStories(storyIds []int) ([]string, error) {
	var uriTemplate = "https://hacker-news.firebaseio.com/v0/item/%d.json"
	var jsonResp struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	}
	var stories []string
	for _, id := range storyIds {
		uri := fmt.Sprintf(uriTemplate, id)
		resp, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			return nil, err
		}
		stories = append(stories, fmt.Sprintf("%s %s", jsonResp.Title, jsonResp.Url))
	}
	return stories, nil
}

//Handle Run the task if query matches
func (task HackerNewsTopTask) Handle(query Query) (bool, <-chan Answer) {
	if !strings.Contains(query.Statement, "hackernews top") && !strings.Contains(query.Statement, "hn top") {
		return false, nil
	}
	c1 := make(chan Answer)
	go func() {
		ids, err := task.getTopIds()
		if err != nil {
			log.Fatal(err)
			return
		}
		stories, err := task.getStories(ids[:5])
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, story := range stories {
			c1 <- Answer(story)
		}
		close(c1)
	}()
	return true, c1
}

//NewHackerNewsTopTask Get a new HackerNewsTopTask
func NewHackerNewsTopTask() *HackerNewsTopTask {
	return &HackerNewsTopTask{}
}
