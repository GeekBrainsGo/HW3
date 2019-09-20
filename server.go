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
	Pg    Page
	lg    *logrus.Logger
	Blogs BlogItems
}

// NewServer - создаёт новый экземпляр сервера
func NewServer(lg *logrus.Logger) *Server {
	return &Server{
		Pg: Page{},
		lg: lg,
		Blogs: BlogItems{
			{
				Title:    "Мой первый блог",
				Body:     "Первый блин комом. И зачем я это все делаю.",
				Comments: []string{"Отличная статья", "Не чего более бредового не читал", "Супер"},
			},
			{
				Title: "Эксперементы продолжаются",
				Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt.",
			},
			{Title: "А это что такое",
				Body: "Здесь ничего нет",
			},
			{Title: "Полный бред",
				Body: "Здесь тоже пусто",
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
		})
	})
}

func (serv *Server) HandleBlogIndex(w http.ResponseWriter, r *http.Request) {

	// serv.Pg = Page{
	// 	Title:   "Мой личный блог",
	// 	Content: "Главная страница",
	// }

	file, _ := os.Open(SERVER_PATH + "index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	// templ := template.Must(template.ParseFiles(SERVER_PATH+"index.html", SERVER_PATH+"head.html"))
	err := templ.ExecuteTemplate(w, "page", serv)
	// err := templ.ExecuteTemplate(w, SERVER_PATH+"index.html", serv)
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
				"url": r.URL.String(),
				// "cookie": r.Header.Get("Cookie"),
				"body": string(body),
			}).
			Debug("request")
		next.ServeHTTP(w, r)
	})
}
