/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

const (
	servAddr     = ":8080"
	templateExt  = "*.gohtml"
	templatePath = "templates"
	postsURL     = "/posts"
	editURL      = "/edit"
	deleteURL    = "/delete"
	createURL    = "/create"
	staticPath   = "/static"
)

// Post is the base post type
type Post struct {
	ID      int
	Title   string
	Date    string // todo change to time.Time?
	Summary string
	Body    template.HTML
}

type dbPosts map[string]Post // todo - need xSQL storage

// Handler is the global handlers struct
type Handler struct {
	posts  dbPosts
	tGlob  *template.Template
	globID int
	mux    sync.Mutex
}

func main() {

	// new handler struct
	handlers := &Handler{
		posts: dbPosts{},
		tGlob: template.Must(template.ParseGlob(path.Join(templatePath, templateExt))),
	}
	handlers.initPosts()

	// prepare routes & middleware
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/", handlers.mainPage)
	mux.Route(postsURL, func(r chi.Router) {
		r.Get("/{id}", handlers.postPage)
		r.Get("/", handlers.mainPage)
		r.Route(deleteURL, func(r chi.Router) {
			r.Post("/{id}", handlers.deletePostPage)
		})
		r.Route(createURL, func(r chi.Router) {
			r.Get("/", handlers.createPostPageForm)
			r.Post("/", handlers.createPostPage)
		})
		r.Route(editURL, func(r chi.Router) {
			r.Get("/{id}", handlers.editPostPageForm)
			r.Post("/{id}", handlers.editPostPage)
		})
	})

	// custom server is for custom parameters & graceful shutdown
	srv := &http.Server{Addr: servAddr, Handler: mux}

	// static files
	mux.Handle(staticPath+"/*", http.StripPrefix(staticPath, http.FileServer(http.Dir("."+staticPath))))

	// graceful shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill cannot be trapped anyway!
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println("Signal received:", <-shutdown)
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error while shutdown server:", err)
		}
	}()

	fmt.Println("Starting server at:", servAddr)
	log.Printf("Shutdown server at: %s\n%v", servAddr, srv.ListenAndServe())
}
