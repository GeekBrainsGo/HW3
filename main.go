package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()

	serv := Server{
		Title: "TODO",
		Tasks: []string{
			"Task 1",
			"Task 2",
			"Task 3",
		},
	}

	// r.Handle("/", http.FileServer(http.Dir("./web/static/")))
	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleGetIndex)
	})

	http.ListenAndServe(":8080", r)
}

type Server struct {
	Title string
	Tasks TaskItems
}

type TaskItems []TaskItem

type TaskItem struct {
	Text      string
	Completed bool
	Labels    string
}

func (serv *Server) HandleGetIndex(w http.ResponseWriter, r *http.Request) {

	file, _ := os.Open("./web/static/index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	logrus.Info(templ.ExecuteTemplate(w, "page", serv))
}
