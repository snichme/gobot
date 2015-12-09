package types

type Answer string

type Task interface {
	Name() string
	HelpText() string
	Handle(query Query) (bool, <-chan Answer)
}

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

type RobotConfig struct {
	Settings    map[string]string   `json:"settings"`
	Permissions map[string][]string `json:"permissions"`
}

type RobotBrain interface {
	Get(key string) interface{}
	Set(key string, value interface{}) bool
}

type Robot interface {
	Name() string
	HasAccess(taskName string, context QueryContext) bool
	Query(s Query) (bool, <-chan Answer)
	//Query(s types.Query) (bool, <-chan types.Answer)
	Brain() RobotBrain
	Write(p []byte) (n int, err error)
	Setting(key string) string
}
