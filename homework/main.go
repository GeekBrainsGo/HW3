package main

import (
	"encoding/json"
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
	Id       int64    `json:"id,omitempty"`
	Title    string   `json:"title,omitempty"`
	Contents string   `json:"contents,omitempty"`
	Labels   []string `json:"labels,omitempty"`
}

const staticDir = "./www/static"

func main() {
	r := chi.NewRouter()
	lg := logrus.New()

	serv := server{
		lg:    lg,
		Title: "Gopher's Blog",
		BlogItems: blogItems{
			{
				Id:       0,
				Title:    "Первая запись",
				Contents: "Первая запись в блоге",
				Labels:   []string{"привет"},
			},
			{
				Id:       1,
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
		r.Get("/edit/{id}", serv.handleGetEditPost)
		r.Post("/edit/{id}", serv.handlePostEditPost)
		r.Post("/create", serv.handlePostCreatePost)
		r.Post("/delete/{id}", serv.handlePostDeletePost)

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
	var searchedPost blogItem
	for _, test := range serv.BlogItems {
		if test.Id == postNumber {
			searchedPost = test
			break
		}
	}
	err := indexTemplate.ExecuteTemplate(w, "index", searchedPost)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *server) handlePostCreatePost(w http.ResponseWriter, r *http.Request) {
	var post blogItem
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		post.Id = int64(len(serv.BlogItems))
		serv.BlogItems = append(serv.BlogItems, post)
		resp, err := json.Marshal(post)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(resp)
		}
	}

}

func (serv *server) handlePostDeletePost(w http.ResponseWriter, r *http.Request) {
	postNumberStr := chi.URLParam(r, "id")
	postNumber, _ := strconv.ParseInt(postNumberStr, 10, 64)
	serv.BlogItems = append(serv.BlogItems[:postNumber], serv.BlogItems[postNumber+1:]...)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func (serv *server) handleGetEditPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/edit.html")
	data, _ := ioutil.ReadAll(file)
	postNumberStr := chi.URLParam(r, "id")
	postNumber, _ := strconv.ParseInt(postNumberStr, 10, 64)
	indexTemplate := template.Must(template.New("index").Parse(string(data)))
	var searchedPost blogItem
	for _, test := range serv.BlogItems {
		if test.Id == postNumber {
			searchedPost = test
			break
		}
	}
	err := indexTemplate.ExecuteTemplate(w, "index", searchedPost)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *server) handlePostEditPost(w http.ResponseWriter, r *http.Request) {
	var post blogItem
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		serv.BlogItems[post.Id] = post
		resp, err := json.Marshal(post)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(resp)
		}
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
