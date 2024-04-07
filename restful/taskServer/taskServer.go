package taskserver

import (
	"encoding/json"
	"log"
	"net/http"
	"restful/taskstore"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

	r.ParseForm()
	text := r.FormValue("text")
	tags := r.Form["tags"]    //tags is a slice of strings
	due := r.FormValue("due") //due is a int number
	t, _ := strconv.Atoi(due)

	id := ts.store.CreateTask(text, tags, time.Now().Add(time.Duration(t)*time.Minute))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(strconv.Itoa(id)))
}

func (ts *TaskServer) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task request at%s\n", r.URL.Path)

	//this is the official way to get the id parameter, but i use the gorilla mux library
	// sid := r.PathValue("id")
	// if sid == "" {
	// 	http.Error(w, "id parameter is missing", http.StatusBadRequest)
	// 	return
	// }

	//id, err := strconv.Atoi(sid)
	//another way to get the id parameter
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
