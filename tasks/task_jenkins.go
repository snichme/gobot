package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/snichme/gobot/types"
)

// JenkinsRunBuildTask Task for running jobs in Jenkins
type JenkinsRunBuildTask struct {
	queryRegexp *regexp.Regexp
	jenkinsURL  string
}

// Name Returns the name of the task
func (task JenkinsRunBuildTask) Name() string {
	return "JenkinsRunBuildTask"
}

// HelpText Returns a description of what this task do
func (task JenkinsRunBuildTask) HelpText() string {
	return "Trigger a build for a job in Jenkins"
}

// Handle If query matches, grab the job name from the query and queue that build
func (task JenkinsRunBuildTask) Handle(query types.Query) (bool, <-chan types.Answer) {
	matches := task.queryRegexp.FindStringSubmatch(query.Statement)
	if matches == nil {
		return false, nil
	}
	jobName := matches[1]
	uri := fmt.Sprintf("%s/job/%s/build", task.jenkinsURL, jobName)
	body := strings.NewReader("")

	c1 := make(chan types.Answer)
	go func() {
		resp, err := http.Post(uri, " text/plain", body)
		jobURL := resp.Header.Get("Location")
		if err != nil {
			c1 <- types.Answer(err.Error())
		} else {
			var responseText string
			if resp.StatusCode == 201 {
				responseText = fmt.Sprintf("A new build for %s is queued (%s)", jobName, jobURL)
			} else if resp.StatusCode == 404 {
				responseText = fmt.Sprintf("No job with name %s exists", jobName)
			} else {
				responseText = fmt.Sprintf("Unknown response (%s) from jenkins", resp.Status)
			}
			c1 <- types.Answer(responseText)
		}
		resp.Body.Close()
		close(c1)
	}()

	return true, c1
}

// NewJenkinsRunBuildTask Creates a new JenkinsRunBuildTask
func NewJenkinsRunBuildTask(jenkinsURL string) *JenkinsRunBuildTask {
	var queryRegexp = regexp.MustCompile(`jenkins build ([a-z]+)`)
	return &JenkinsRunBuildTask{
		queryRegexp: queryRegexp,
		jenkinsURL:  jenkinsURL,
	}
}

// JenkinsListTask Task for running jobs in Jenkins
type JenkinsListTask struct {
	jenkinsURL string
}

// Name Returns the name of the task
func (task JenkinsListTask) Name() string {
	return "JenkinsListTask"
}

// HelpText Returns a description of what this task do
func (task JenkinsListTask) HelpText() string {
	return "List all available jobs in Jenkins"
}

type listTaskResponse struct {
	Jobs []map[string]string `json:"jobs"`
}

// Handle If query matches, grab the job name from the query and queue that build
func (task JenkinsListTask) Handle(query types.Query) (bool, <-chan types.Answer) {
	if query.Statement != "jenkins list" {
		return false, nil
	}
	c1 := make(chan types.Answer)
	go func() {
		var jsonResp listTaskResponse

		resp, _ := http.Get(fmt.Sprintf("%s/%s", task.jenkinsURL, "/api/json"))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		json.Unmarshal(body, &jsonResp)

		c1 <- types.Answer("Here is the jobs available in Jenkins: ")
		for _, m := range jsonResp.Jobs {
			c1 <- types.Answer(fmt.Sprintf("%s %s", m["name"], m["url"]))
		}
		close(c1)
	}()
	return true, c1
}

// NewJenkinsListTask Creates a new JenkinsRunBuildTask
func NewJenkinsListTask(jenkinsURL string) *JenkinsListTask {

	return &JenkinsListTask{
		jenkinsURL: jenkinsURL,
	}
}

// JenkinsStatusTask Task for getting the status of a job
type JenkinsStatusTask struct {
	jenkinsURL  string
	queryRegexp *regexp.Regexp
}

// Name Returns the name of the task
func (task JenkinsStatusTask) Name() string {
	return "JenkinsStatusTask"
}

// HelpText Returns a description of what this task do
func (task JenkinsStatusTask) HelpText() string {
	return "Get the status of a build in jenkins"
}

type statusTaskResponse struct {
	HealthReport []map[string]string `json:"healthReport"`
	URL          string              `json:"url"`
	Color        string              `json:"color"`
}

// Handle If query matches, grab the job name from the query and queue that build
func (task JenkinsStatusTask) Handle(query types.Query) (bool, <-chan types.Answer) {
	matches := task.queryRegexp.FindStringSubmatch(query.Statement)
	if matches == nil {
		return false, nil
	}
	jobName := matches[1]

	uri := fmt.Sprintf("%s/job/%s/api/json", task.jenkinsURL, jobName)
	c1 := make(chan types.Answer)
	go func() {
		var jsonResp statusTaskResponse

		resp, err := http.Get(uri)
		if err != nil {
			log.Fatal(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		json.Unmarshal(body, &jsonResp)

		c1 <- types.Answer(fmt.Sprintf("Status: %s", jsonResp.Color))
		c1 <- types.Answer(fmt.Sprintf("Url: %s", jsonResp.URL))
		c1 <- types.Answer(fmt.Sprintf("Health report: %s", jsonResp.HealthReport[0]["description"]))
		close(c1)
	}()
	return true, c1
}

// NewJenkinsStatusTask Creates a new JenkinsRunBuildTask
func NewJenkinsStatusTask(jenkinsURL string) *JenkinsStatusTask {
	var queryRegexp = regexp.MustCompile(`jenkins status ([a-z]+)`)

	return &JenkinsStatusTask{
		jenkinsURL:  jenkinsURL,
		queryRegexp: queryRegexp,
	}
}
