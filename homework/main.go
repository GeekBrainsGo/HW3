package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

type server struct {
	lg        *logrus.Logger
	Title     string
	BlogItems blogItems
}

type blogItems []blogItem

type blogItem struct {
	Title    string
	Contents string
	Labels   []string
}

const staticDir = "./www/static"

func main() {
	r := chi.NewRouter()
	lg := logrus.New()

	serv := server{
		lg:    lg,
		Title: "Fadeev's Blog",
		BlogItems: blogItems{
			{
				Title:    "Первая запись",
				Contents: "Первая запись в блоге",
				Labels:   []string{"привет"},
			},
			{
				Title:    "Вторая запись",
				Contents: "Вторая запись в блоге",
				Labels:   []string{"два", "тест"},
			},
		},
	}

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "www/static")
	FileServer(r, "/static", http.Dir(filesDir))

	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.handleGetIndex)
		r.Get("/post/{id}", serv.handleGetPost)
	})

	lg.Info("starting the server")
	lg.Error(http.ListenAndServe(":8080", r))

}

func (serv *server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/index.html")
	data, _ := ioutil.ReadAll(file)

	indexTemplate := template.Must(template.New("index").Parse(string(data)))
	err := indexTemplate.ExecuteTemplate(w, "index", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/post.html")
	data, _ := ioutil.ReadAll(file)
	postNumberStr := chi.URLParam(r, "id")
	postNumber, _ := strconv.ParseInt(postNumberStr, 10, 64)
	indexTemplate := template.Must(template.New("index").Parse(string(data)))
	err := indexTemplate.ExecuteTemplate(w, "index", serv.BlogItems[postNumber])
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
