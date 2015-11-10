package main

import "testing"

func TestRobotAccessWithoutPermissions(t *testing.T) {
	taskName := "TaskName"
	qc := QueryContext{
		Username: "TestUser",
		Group:    "TestGroup",
	}
	r := Robot{}
	if r.HasAccess(taskName, qc) == false {
		t.Errorf("If no permissions given to robot all tasks should be open")
	}
}

func TestRobotAccessHaveCorrectgroup(t *testing.T) {
	taskName := "TaskName"
	qc := QueryContext{
		Username: "TestUser",
		Group:    "TTGroup",
	}
	r := Robot{
		permissions: map[string][]string{
			"TaskName": []string{"TTGroup"},
		},
	}
	if r.HasAccess(taskName, qc) == false {
		t.Errorf("User should have access if is in correct group")
	}
}
func TestRobotAccessDoesntHaveCorrectgroup(t *testing.T) {
	taskName := "TaskName"
	qc := QueryContext{
		Username: "TestUser",
		Group:    "TestGroup",
	}
	r := Robot{
		permissions: map[string][]string{
			"TaskName": []string{"TTGroup"},
		},
	}
	if r.HasAccess(taskName, qc) == true {
		t.Errorf("User should not have access if is in wrong group")
	}
}
func TestRobotAccessAsAdmin(t *testing.T) {
	taskName := "TaskName"
	qc := QueryContext{
		Username: "TestUser",
		Group:    "admin",
	}
	r := Robot{
		permissions: map[string][]string{
			"TaskName": []string{"Somegroup"},
		},
	}
	if r.HasAccess(taskName, qc) == false {
		t.Errorf("Admin should always have access")
	}
}
