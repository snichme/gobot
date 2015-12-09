package tasks

import "testing"

var serverstatus_queries = []Query{
	Query{Statement: "whats up?"},
	Query{Statement: "Whats up?"},
	Query{Statement: "how are you?"},
	Query{Statement: "Good evening, how are you today?"},
}

func TestTaskServerStatusShouldHandleQueries(t *testing.T) {
	task := NewServerStatusTask()

	for _, query := range serverstatus_queries {
		if task.CanHandle(query) == false {
			t.Errorf("Should handle \"%s\" query", query.Statement)
		}
	}
}

/*
func TestTaskServerStatusShouldExecuteQueries(t *testing.T) {
	task := NewServerStatusTask()
	for _, query := range serverstatus_queries {
		c := task.DoHandle(query)
		for v := range c {
			fmt.Println(v)
		}
	}
}
*/
