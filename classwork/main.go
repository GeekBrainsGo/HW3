package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()
	lg := logrus.New()

	serv := Server{
		lg:    lg,
		Title: "TODO",
		Tasks: TaskItems{
			{
				Text:      "Дедлайн завтра, ОГОНЬ!",
				Completed: false,
				Labels:    []string{"срочно"},
			},
			{Text: "Проснулся", Completed: true},
			{Text: "Поел", Completed: true},
			{Text: "Уснул", Completed: false},
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleGetIndex)
		r.Post("/{taskID}/{status}", serv.HandlePostTaskStatus)
	})

	lg.Info("server is start")
	http.ListenAndServe(":8080", r)
}

type Server struct {
	lg    *logrus.Logger
	Title string
	Tasks TaskItems
}

type TaskItems []TaskItem
type TaskItem struct {
	Text      string
	Completed bool
	Labels    []string
}

func (tasks TaskItems) TasksWithStatus(completed bool) int {
	total := 0
	for _, task := range tasks {
		if task.Completed == completed {
			total++
		}
	}
	return total
}

func (tasks TaskItems) CompletePercent() float64 {
	prc := float64(tasks.TasksWithStatus(true)) / float64(len(tasks))
	return math.Round(prc * 100)
}

func (serv *Server) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandlePostTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskID")
	taskStatusStr := chi.URLParam(r, "status")

	taskID, _ := strconv.ParseInt(taskIDStr, 10, 64)
	taskStatus, _ := strconv.ParseBool(taskStatusStr)

	serv.Tasks[taskID].Completed = taskStatus

	data, _ := json.Marshal(serv.Tasks[taskID])
	w.Write(data)
	serv.lg.WithField("tasks", serv.Tasks).Info("status changed")
}
