// which is a RESTful API for a task store
// TODO: consider using a database(e.g. mysql,redis) instead of a map to store tasks
package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// TaskStore is a store for tasks, which is concurrent-safe.
type TaskStore struct {
	tasks  map[int]Task
	lock   sync.Mutex
	nextId int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		lock:   sync.Mutex{},
		nextId: 0,
	}
}

// CreateTask creates a new task in ts and returns its ID.
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	task := Task{
		ID:   ts.nextId,
		Text: text,
		Due:  due,
	}
	//copy tags instead of passing by reference
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.ID
}

// GetTask returns the task with the given ID. If no such task exists, it returns an error.
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	task, ok := ts.tasks[id]
	if ok {
		return task, nil
	} else {
		return Task{}, fmt.Errorf("task with ID %d not found", id)
	}
}

// GetTasksByTag returns all tasks in ts that have the given tag.
func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	//TODO: optimize this search by using a map of tags to tasks
	//for example, use a map[string][]int to store the mapping of tags to task IDs
	tasks := make([]Task, 0)
	for _, task := range ts.tasks {
		for _, t := range task.Tags {
			if t == tag {
				tasks = append(tasks, task)
				break
			}
		}
	}
	return tasks
}

// GetTasksByDueDate returns all tasks in ts that are due on the given date.
func (ts *TaskStore) GetTasksByDueDate(year int, mouth time.Month, day int) []Task {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	tasks := make([]Task, 0)
	dueDate := time.Date(year, time.Month(mouth), day, 0, 0, 0, 0, time.UTC)
	for _, task := range ts.tasks {
		if task.Due.Equal(dueDate) {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// GetAllTasks returns all tasks in ts.
func (ts *TaskStore) GetAllTasks() []Task {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	//remember to set the length as 0
	tasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// UpdateTask updates the task with the given ID with the new Task If no such task exists, it returns an error.
func (ts *TaskStore) UpdateTask(id int, task Task) error {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	_, ok := ts.tasks[id]
	if ok {
		ts.tasks[id] = task
		return nil
	} else {
		return fmt.Errorf("task with ID %d not found", id)
	}
}

// DeleteTask deletes the task with the given ID from ts. If no such task exists, it returns an error.
func (ts *TaskStore) DeleteTask(id int) error {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	_, ok := ts.tasks[id]
	if ok {
		delete(ts.tasks, id)
		return nil
	} else {
		return fmt.Errorf("task with ID %d not found", id)
	}
}

// Clear deletes all tasks from ts.
func (ts *TaskStore) Clear() {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.tasks = make(map[int]Task)
	ts.nextId = 0
}
