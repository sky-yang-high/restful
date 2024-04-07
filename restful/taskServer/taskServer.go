package taskserver

import (
	"encoding/json"
	"net/http"
	"restful/taskstore"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// func (ts *TaskServer) CreateTaskHandler(c *gin.Context) {
// 	//gin will log the request automatically
// 	//log.Printf("handling create task request at%s\n", c.Request.URL.Path)

// 	// text := c.PostForm("text")
// 	// tags := c.PostFormArray("tags") //tags is a slice of strings
// 	// due := c.PostForm("due")        //due is a int number
// 	// t, _ := strconv.Atoi(due)
// 	// i don't know why the above code doesn't work, so i use the following code instead:
// 	text := c.Request.FormValue("text")
// 	tags := c.Request.Form["tags"]
// 	due := c.Request.FormValue("due")
// 	t, _ := strconv.Atoi(due)

//		id := ts.store.CreateTask(text, tags, time.Now().Add(time.Duration(t)*time.Minute))
//		c.JSON(http.StatusCreated, gin.H{"id": id})
//	}

// another way to implement the CreateTaskHandler function:
func (ts *TaskServer) CreateTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string   `json:"text"`
		Tags []string `json:"tags"`
		Due  int      `json:"due"`
	}
	var rt RequestTask
	//in this way, the pramer should be set in request.body with json format, instead of request.url
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ts.store.CreateTask(rt.Text, rt.Tags, time.Now().Add(time.Duration(rt.Due)*time.Minute))
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ts *TaskServer) GetTaskHandler(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	task, err := ts.store.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//gin has built-in JSON rendering
	c.JSON(http.StatusOK, task)
}

func (ts *TaskServer) GetAllTasksHandler(c *gin.Context) {
	tasks := ts.store.GetAllTasks()
	if tasks == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tasks found"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
