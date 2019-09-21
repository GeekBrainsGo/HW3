package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

const (
	SERVER_PATH = "./web/static/"
)

// Server - Объект сервера
type Server struct {
	lg    *logrus.Logger
	Blogs BlogItems
}

// NewServer - создаёт новый экземпляр сервера
func NewServer(lg *logrus.Logger) *Server {
	return &Server{
		lg: lg,
		Blogs: BlogItems{
			{
				ID:       0,
				Title:    "Мой первый блог",
				Body:     "Первый блин комом. И зачем я это все делаю.",
				Comments: []string{"Отличная статья", "Не чего более бредового не читал", "Супер"},
			},
			{
				ID:    1,
				Title: "Эксперементы продолжаются",
				Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt.",
			},
			{
				ID:    2,
				Title: "А это что такое",
				Body:  "Здесь ничего нет",
			},
			{
				ID:    3,
				Title: "Полный бред",
				Body:  "Здесь тоже пусто",
			},
		},
	}
}

// Start - запускает сервер
func (serv *Server) Start() error {
	r := chi.NewRouter()
	r.Use(serv.RequestTracerMiddleware)
	serv.ConfigureHandlers(r)
	serv.lg.Info("Сервер запущен!")
	return http.ListenAndServe(":8080", r)
}

// ConfigureHandlers - настраивает хендлеры и их пути
func (serv *Server) ConfigureHandlers(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleBlogIndex)
		r.Route("/blog", func(r chi.Router) {
			r.Get("/{blogID}", serv.HandleBlog)
			r.Get("/del/{blogID}", serv.HandleBlogDelete)
			r.Get("/edit/{blogID}", serv.HandleEditBlog)
		})
	})
}

func (serv *Server) HandleBlogIndex(w http.ResponseWriter, r *http.Request) {

	file, _ := os.Open(SERVER_PATH + "index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleEditBlog(w http.ResponseWriter, r *http.Request) {

	blogIDStr := chi.URLParam(r, "blogID")
	blogID, _ := strconv.ParseInt(blogIDStr, 10, 64)

	file, _ := os.Open(SERVER_PATH + "edit.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("edit").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "edit", serv.Blogs[blogID])
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleBlog(w http.ResponseWriter, r *http.Request) {

	blogIDStr := chi.URLParam(r, "blogID")
	blogID, _ := strconv.ParseInt(blogIDStr, 10, 64)

	file, _ := os.Open(SERVER_PATH + "blog.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("blog").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "blog", serv.Blogs[blogID])
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func RemoveBlogSlice(slice BlogItems, start, end int64) BlogItems {
	return append(slice[:start], slice[end:]...)
}

func (serv *Server) HandleBlogDelete(w http.ResponseWriter, r *http.Request) {
	blogIDStr := chi.URLParam(r, "blogID")
	blogID, _ := strconv.ParseInt(blogIDStr, 10, 64)

	slice := RemoveBlogSlice(serv.Blogs, blogID, blogID+1)

	serv.Blogs = slice

	file, _ := os.Open(SERVER_PATH + "index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

// SendErr - отправляет ошибку пользователю и логирует её
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("Ошибка сервера")
	w.WriteHeader(code)
	errModel := ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - отправляет 500 ошибку
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}

// RequestTracerMiddleware - отслеживает и логирует входящие запросы
func (serv *Server) RequestTracerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		serv.lg.
			WithFields(map[string]interface{}{
				"url":  r.URL.String(),
				"body": string(body),
			}).
			Debug("request")
		next.ServeHTTP(w, r)
	})
}
