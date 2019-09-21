// Package blog implement basic blog server.
package main

/*
	Basics Go.
	Rishat Ishbulatov, dated Sep 19, 2019.
	Create a route and template to display all blog posts.
	Create a route and template for viewing a specific blog post.
	Create a route and template for editing and creating material.
*/

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"

	"github.com/sirupsen/logrus"
)

var base *template.Template

// DB stands for blog database.
type DB struct {
	Title string
	Posts Posts
	sync.Mutex
}

// Posts stands for array of posts.
type Posts []Post

// Post stands for post in blog.
type Post struct {
	ID            int
	Title, Author string
	Created       time.Time
	Content       template.HTML
}

// Blog stands for blog handling multiplexer.
type Blog struct {
	*chi.Mux
	*logrus.Logger
	*DB
}

func init() {
	file, _ := ioutil.ReadFile("./www/main.html")
	base = template.Must(template.New("").Parse(string(file)))
}

func main() {
	blog := NewBlog()
	blog.Info("Server started")
	log.Fatal(http.ListenAndServe(":8000", blog))
}

// NewBlog return blog handling multiplexer.
func NewBlog() *Blog {
	mux := chi.NewRouter()
	log := logrus.New()
	blog := &Blog{
		mux,
		log,
		&DB{Title: "Awsome Blog"},
	}
	blog.Posts = Posts{
		{
			ID:      0,
			Title:   "Post One",
			Created: time.Now(),
			Author:  "Vasia Pupkine",
			Content: "This is my very first post in this awsome blog",
		},
		{
			ID:      1,
			Title:   "Post Two",
			Created: time.Now(),
			Author:  "Джон Сноу",
			Content: "Дедлайн завтра, ОГОНЬ!",
		},
	}
	mux.Route("/", func(mux chi.Router) {
		mux.Get("/", blog.Main)
		mux.Get("/edit", blog.EditPost)
		mux.Get("/edit/{id}", blog.EditPost)
		mux.Post("/edit/{id}", blog.EditPost)
		mux.Get("/posts/{id}", blog.ViewPost)
	})
	return blog
}

// Add adds new post to database and return it's id.
func (db *DB) Add(p Post) int {
	db.Lock()
	defer db.Unlock()
	db.Posts = append(db.Posts, p)
	return len(db.Posts) - 1
}

// Main handles displaying all posts in blog.
func (b *Blog) Main(w http.ResponseWriter, r *http.Request) {
	err := base.ExecuteTemplate(w, "main", b.DB)
	if err != nil {
		b.WithError(err).Error("main")
		return
	}
}

// ViewPost handles for viewing a specific blog post
func (b *Blog) ViewPost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if id >= len(b.Posts) {
		http.NotFound(w, r)
		return
	}
	err := base.ExecuteTemplate(w, "view", b.DB.Posts[id])
	if err != nil {
		b.WithError(err).Error("view")
		return
	}
}

// EditPost handles editing and creating blog's post.
func (b *Blog) EditPost(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "id")
	if len(i) == 0 {
		err := base.ExecuteTemplate(w, "edit", Post{})
		if err != nil {
			b.WithError(err).Error("addpost")
			return
		}
		return
	}
	id, _ := strconv.Atoi(i)
	if id > len(b.Posts) {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost && id < len(b.Posts) {
		err := base.ExecuteTemplate(w, "edit", b.DB.Posts[id])
		if err != nil {
			b.WithError(err).Error("editpost")
			return
		}
		return
	}
	p := Post{
		ID:      len(b.Posts),
		Title:   r.FormValue("title"),
		Author:  r.FormValue("author"),
		Created: time.Now(),
		Content: template.HTML(r.FormValue("body")),
	}
	if id == 0 {
		id = b.Add(p)
	} else {
		b.Posts[id] = p
	}
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}
