package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

const (
	STATIC_DIR = "./www/"
)

func main() {
	r := chi.NewRouter()
	lg := logrus.New()

	serv := Server{
		lg:    lg,
		Title: "BLOG",
		Posts: PostItems{
			{
				Id:		   1,
				Title:     "Пост 1",
				Text:      "Очень интересный текст",
				Labels:    []string{"путешестве", "отдых"},
			},
			{
				Id:		   2,
				Title:     "Пост 2",
				Text:      "Второй очень интересный текст",
				Labels:    []string{"домашка", "golang"},
			},
			{
				Id:		   3,
				Title:     "Пост 3",
				Text:      "Третий очень интересный текст",
				Labels:    []string{},
			},
		},
	}

	fileServer := http.FileServer(http.Dir(STATIC_DIR))
	r.Handle("/static/*", fileServer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleGetIndex)
		r.Get("/post/{id}", serv.HandleGetPost)
		r.Get("/post/create", serv.HandleGetEditPost)
		r.Post("/post/create", serv.HandleEditPost)
		r.Get("/post/{id}/edit", serv.HandleGetEditPost)
		r.Post("/post/{id}/edit", serv.HandleEditPost)
	})

	lg.Info("server is start")
	http.ListenAndServe(":8080", r)
}

type Server struct {
	lg    *logrus.Logger
	Title string
	Posts PostItems
}

type PostItems []PostItem
type PostItem struct {
	Id     	  int64
	Title     string
	Text      string
	Labels    []string
}

func (posts PostItems) PostsById(id int64) (PostItem, error) {
	for _, post := range posts {
		if post.Id == id {
			return post, nil
		}
	}
	return PostItem{}, nil
}

func (serv *Server) AddOrUpdatePost(newPost PostItem) (PostItem) {
	for key, post := range serv.Posts {
		if post.Id == newPost.Id {
			serv.Posts[key] = newPost
			return post
		}
	}

	serv.Posts = append(serv.Posts, newPost)
	return newPost
}


func (serv *Server) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/index.gohtml")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/post.gohtml")
	data, _ := ioutil.ReadAll(file)

	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	post, err := serv.Posts.PostsById(postID)
	if err != nil {
		serv.lg.WithError(err).Error("template")
		post = PostItem{}
	}

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleGetEditPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/edit_post.gohtml")
	data, _ := ioutil.ReadAll(file)

	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	post, err := serv.Posts.PostsById(postID)
	if err != nil {
		serv.lg.WithError(err).Error("template")
		post = PostItem{}
	}

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleEditPost(w http.ResponseWriter, r *http.Request) {
	/*
	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	post, err := serv.Posts.PostsById(postID)
	if err != nil {
		serv.lg.WithError(err).Error("template")
		post = PostItem{}
	}
	*/

	/*
	пока сделал передачу ID поста в json, нужно передалть для выбора поста по ID
	{"id":4, "Title":"Пост 4", "Text":"Новый очень интересный текст", "Labels":["l1","l2"]}
	*/

	decoder := json.NewDecoder(r.Body)
	var inPostItem PostItem
	err := decoder.Decode(&inPostItem)
	if err != nil {
		serv.lg.WithError(err).Error("decode edit post data")
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	newPost := serv.AddOrUpdatePost(inPostItem)
	respondWithJSON(w, http.StatusOK, newPost)
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

