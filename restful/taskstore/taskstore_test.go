package taskstore

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	// create a new task
	ts := NewTaskStore()
	id := ts.CreateTask("task1", nil, time.Now())

	//get one task by id
	task, err := ts.GetTask(id)
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != id {
		t.Errorf("expected task id %d, got %d", id, task.ID)
	}
	if task.Text != "task1" {
		t.Errorf("expected task text %s,got %s", "task1", task.Text)
	}

	//get all tasks and check if the new task is there
	tasks := ts.GetAllTasks()
	if len(tasks) != 1 || tasks[0].ID != id {
		t.Errorf("expected one task with id %d, got %d tasks", id, len(tasks))
	}

	//get another task,expect an error
	_, err = ts.GetTask(id + 1)
	if err == nil {
		t.Error("expected error for non-existent task,got nil")
	}

	ts.CreateTask("task2", nil, time.Now())
	tasks = ts.GetAllTasks()
	if len(tasks) != 2 {
		t.Errorf("expected two tasks, got %d tasks", len(tasks))
	}
}

func TestDelect(t *testing.T) {
	ts := NewTaskStore()
	id1 := ts.CreateTask("task1", nil, time.Now())
	ts.CreateTask("task2", nil, time.Now())

	tasks := ts.GetAllTasks()
	if len(tasks) != 2 {
		t.Errorf("expected two tasks, got %d tasks", len(tasks))
	}

	err := ts.DeleteTask(id1)
	if err != nil {
		t.Error(err)
	}
	err = ts.DeleteTask(id1)
	if err == nil {
		t.Errorf("expected error for non-existent task,got nil")
	}

	err = ts.DeleteTask(4)
	if err == nil {
		t.Error("expected error for non-existent task,got nil")
	}

	ts.Clear()
	tasks = ts.GetAllTasks()
	if len(tasks) != 0 {
		t.Errorf("expected no tasks, got %d tasks", len(tasks))
	}
}

func TestGetTasksByTag(t *testing.T) {
	ts := NewTaskStore()
	ts.CreateTask("XY", []string{"Movies"}, time.Now())
	ts.CreateTask("YZ", []string{"Bills"}, time.Now())
	ts.CreateTask("YZR", []string{"Bills"}, time.Now())
	ts.CreateTask("YWZ", []string{"Bills", "Movies"}, time.Now())
	ts.CreateTask("WZT", []string{"Movies", "Bills"}, time.Now())

	var tests = []struct {
		tag     string
		wantNum int
	}{
		{"Movies", 3},
		{"Bills", 4},
		{"Ferrets", 0},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			numByTag := len(ts.GetTasksByTag(tt.tag))
			if numByTag != tt.wantNum {
				t.Errorf("got %v, want %v", numByTag, tt.wantNum)
			}
		})
	}

}

func TestGetTasksByDueDate(t *testing.T) {
	timeFormat := "2006-Jan-02"
	mustParseDate := func(tstr string) time.Time {
		tt, err := time.Parse(timeFormat, tstr)
		if err != nil {
			t.Fatal(err)
		}
		return tt
	}

	ts := NewTaskStore()
	ts.CreateTask("XY1", nil, mustParseDate("2020-Dec-01"))
	ts.CreateTask("XY2", nil, mustParseDate("2000-Dec-21"))
	ts.CreateTask("XY3", nil, mustParseDate("2020-Dec-01"))
	ts.CreateTask("XY4", nil, mustParseDate("2000-Dec-21"))
	ts.CreateTask("XY5", nil, mustParseDate("1991-Jan-01"))

	// Check a single task can be fetched.
	y, m, d := mustParseDate("1991-Jan-01").Date()
	tasks1 := ts.GetTasksByDueDate(y, m, d)
	if len(tasks1) != 1 {
		t.Errorf("got len=%d, want 1", len(tasks1))
	}
	if tasks1[0].Text != "XY5" {
		t.Errorf("got Text=%s, want XY5", tasks1[0].Text)
	}

	var tests = []struct {
		date    string
		wantNum int
	}{
		{"2020-Jan-01", 0},
		{"2020-Dec-01", 2},
		{"2000-Dec-21", 2},
		{"1991-Jan-01", 1},
		{"2020-Dec-21", 0},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			y, m, d := mustParseDate(tt.date).Date()
			numByDate := len(ts.GetTasksByDueDate(y, m, d))

			if numByDate != tt.wantNum {
				t.Errorf("got %v, want %v", numByDate, tt.wantNum)
			}
		})
	}
}
