package taskserver

import (
	"encoding/json"
	"log"
	"net/http"
	"restful/taskstore"
	"strconv"
	"time"
)

type TaskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *TaskServer {
	return &TaskServer{
		store: taskstore.NewTaskStore(),
	}
}

// renderJSON is a helper function to render JSON data to the response writer
func renderJSON(w http.ResponseWriter, data interface{}) {
	jsdata, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsdata)
}

func (ts *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling create task request at%s\n", r.URL.Path)

	que := r.URL.Query()

	Text := que.Get("text")
	//TODO: maybe the time should be in the query string as well
	t, _ := strconv.Atoi(que.Get("due"))
	Due := time.Now().Add(time.Duration(t) * time.Minute)
	Tags := que["tags"]

	id := ts.store.CreateTask(Text, Tags, Due)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(strconv.Itoa(id)))
}

func (ts *TaskServer) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task request at%s\n", r.URL.Path)

	sid := r.URL.Query().Get("id")
	if sid == "" {
		http.Error(w, "id parameter is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderJSON(w, task)
}

func (ts *TaskServer) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all tasks request at%s\n", r.URL.Path)

	tasks := ts.store.GetAllTasks()
	if tasks == nil {
		http.Error(w, "no tasks found", http.StatusNotFound)
		return
	}

	renderJSON(w, tasks)
}
