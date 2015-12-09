package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/snichme/gobot/types"
)

type RestClient struct {
	robot types.Robot
}

type RestResponse struct {
	RobotName string         `json:"robot"`
	Query     string         `json:"query"`
	Answers   []types.Answer `json:"answers"`
}

func restHandler(robot types.Robot) http.HandlerFunc {
	isAuthorized := func(username, password string) bool {
		users := robot.Brain().Get("users").([]map[string]string)
		for _, user := range users {
			if user["username"] == username && user["password"] == password {
				return true
			}
		}
		return false
	}

	getGroupForUser := func(username string) string {
		users := robot.Brain().Get("users").([]map[string]string)
		for _, user := range users {
			if user["username"] == username {
				return user["group"]
			}
		}
		return "guest"
	}

	return func(rw http.ResponseWriter, req *http.Request) {

		if req.Method != "POST" {
			http.Error(rw, "Not found", http.StatusNotFound)
			return
		}

		username, password, ok := req.BasicAuth()
		if ok == false || !isAuthorized(username, password) {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		req.ParseForm()
		query := req.Form.Get("query")
		if len(query) == 0 {
			http.Error(rw, "No query", http.StatusBadRequest)
			return
		}

		q := types.Query{
			Statement: query,
			Context: types.QueryContext{
				Username: username,
				Group:    getGroupForUser(username), // no auth impl yet
			},
		}
		if found, c := robot.Query(q); found {
			var answers []types.Answer
			for answer := range c {
				answers = append(answers, answer)
			}
			resp := RestResponse{
				RobotName: robot.Name(),
				Query:     query,
				Answers:   answers,
			}
			js, _ := json.Marshal(resp)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(js)
		} else {
			http.Error(rw, "Not found", http.StatusNotFound)
		}
	}
}

func (client RestClient) Start() {
	uri := "127.0.0.1:" + client.robot.Setting("http_port")
	fmt.Fprintf(client.robot, "RestClient: Listening on http://%s/rest", uri)
	http.HandleFunc("/rest", restHandler(client.robot))
	http.ListenAndServe(uri, nil)
}

func NewRestClient(robot types.Robot) *RestClient {

	return &RestClient{
		robot: robot,
	}
}
