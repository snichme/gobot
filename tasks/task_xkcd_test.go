package tasks

import "testing"

var xkcd_queries = []Query{
	Query{Statement: "xkcd"},
	Query{Statement: "did someone say xkcd?"},
	Query{Statement: "xkcd is sooo good"},
}

func TestTaskXkcdShouldHandleQueries(t *testing.T) {
	task := NewXkcdTask()
	for _, query := range xkcd_queries {
		if task.CanHandle(query) == false {
			t.Errorf("Should handle \"%s\" query", query.Statement)
		}
	}
}

/*
func TestTaskXkcdShouldExecuteQueries(t *testing.T) {
	task := NewXkcdTask()
	for _, query := range xkcd_queries {
		c := task.DoHandle(query)
		for v := range c {
			fmt.Println(v)
		}
	}
}
*/
