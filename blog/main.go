package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type WebServer struct {
	srv       *http.Server
	logger    *logrus.Logger
	BlogItems BlogItems
}

type BlogItem struct {
	ID           int64
	Title        string
	ShortContent string
	FullContent  string
}

type BlogItems struct {
	Items []BlogItem
}

type PageModel struct {
	Title string
	Data  interface{}
}

func (server *WebServer) BlogGetListHandle(w http.ResponseWriter, r *http.Request) {

	pageData := PageModel{
		Title: "Blog List",
		Data:  server.BlogItems.Items,
	}

	templ := template.Must(template.New("page").ParseFiles("./templates/blog/List.tpl", "./templates/common.tpl"))
	err := templ.ExecuteTemplate(w, "page", pageData)
	if err != nil {
		server.logger.WithError(err).Error("template")
	}
}

func (server *WebServer) BlogGetItemHandle(w http.ResponseWriter, r *http.Request) {

	postIDStr := chi.URLParam(r, "postID")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	Item, err := server.BlogItems.GetItem(postID)
	if err != nil {
		http.Error(w, "page not found", 404)
		return
	}

	pageData := PageModel{
		Title: "Blog Item",
		Data:  Item,
	}

	server.BlogItems.DeleteItem(postID)

	templ := template.Must(template.New("page").ParseFiles("./templates/blog/View.tpl", "./templates/common.tpl"))
	err = templ.ExecuteTemplate(w, "page", pageData)
	if err != nil {
		server.logger.WithError(err).Error("template")
	}
}

func (server *WebServer) BlogAddItemHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "Post" {

	}

	templ := template.Must(template.New("page").ParseFiles("./templates/blog/Add.tpl", "./templates/common.tpl"))
	err := templ.ExecuteTemplate(w, "page", server)
	if err != nil {
		server.logger.WithError(err).Error("template")
	}
}

func (server *WebServer) BlogUpdateItemHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "Post" {

	}

	templ := template.Must(template.New("page").ParseFiles("./templates/blog/Update.tpl", "./templates/common.tpl"))
	err := templ.ExecuteTemplate(w, "page", server)
	if err != nil {
		server.logger.WithError(err).Error("template")
	}
}

func (server *WebServer) BlogDeleteItemHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "Post" {

	}

	templ := template.Must(template.New("page").ParseFiles("./templates/blog/Add.tpl", "./templates/common.tpl"))
	err := templ.ExecuteTemplate(w, "page", server)
	if err != nil {
		server.logger.WithError(err).Error("template")
	}
}

func (blog *BlogItems) AddItem(item *BlogItem) {
	blog.Items = append(blog.Items, *item)
}

func (blog *BlogItems) UpdateItem(ID int64, item *BlogItem) error {
	itemKey, err := blog.FindItemKeyByID(ID)
	if err != nil {
		return err
	}
	blog.Items[itemKey] = *item
	return nil
}

func (blog *BlogItems) DeleteItem(ID int64) error {
	itemKey, err := blog.FindItemKeyByID(ID)
	if err != nil {
		return err
	}
	blog.Items = append(blog.Items[:itemKey], blog.Items[itemKey+1:]...)
	return nil
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func (blog *BlogItems) GetItem(ID int64) (BlogItem, error) {
	itemKey, err := blog.FindItemKeyByID(ID)
	if err != nil {
		return BlogItem{}, err
	}
	return blog.Items[itemKey], nil
}

func (blog *BlogItems) GetAll() ([]BlogItem, error) {
	return blog.Items, nil
}

func (blog *BlogItems) FindItemKeyByID(ID int64) (int, error) {
	for key, item := range blog.Items {
		if item.ID == ID {
			return key, nil
		}
	}
	return 0, fmt.Errorf("ID:%d not found", ID)
}

func main() {
	r := chi.NewRouter()
	logger := logrus.New()

	serv := WebServer{
		logger: logger,
	}

	serv.BlogItems.AddItem(&BlogItem{
		ID:           1,
		Title:        "News1",
		ShortContent: "Short content",
		FullContent:  "Full content",
	})

	// serv.BlogItems.Items = append(serv.BlogItems.Items, BlogItem{
	// 	ID:           1,
	// 	Title:        "News1",
	// 	ShortContent: "Short content",
	// 	FullContent:  "Full content",
	// })

	serv.BlogItems.Items = append(serv.BlogItems.Items, BlogItem{
		ID:           2,
		Title:        "News2",
		ShortContent: "Short content",
		FullContent:  "Full content",
	})

	serv.BlogItems.Items = append(serv.BlogItems.Items, BlogItem{
		ID:           3,
		Title:        "News3",
		ShortContent: "Short content",
		FullContent:  "Full content",
	})

	logger.Info(serv.BlogItems)

	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.BlogGetListHandle)

		r.Get("/view/{postID}", serv.BlogGetItemHandle)

		r.Get("/add", serv.BlogAddItemHandle)
		r.Post("/add", serv.BlogAddItemHandle)

		r.Get("/update/{postID}", serv.BlogUpdateItemHandle)
		r.Post("/update/{postID}", serv.BlogUpdateItemHandle)

		r.Get("/delete/{postID}", serv.BlogDeleteItemHandle)
		r.Post("/delete/{postID}", serv.BlogDeleteItemHandle)
	})

	logger.Info("Starting server...")
	http.ListenAndServe(":8888", r)
}
