package infra

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/art-frela/HW3/karpov/domain"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

const (
	templatePOSTS = "./assets/index.html"
	postID        = "id"
	//httpTimeOut   = 30 * time.Second
)

// BlogServer -
type BlogServer struct {
	log        *logrus.Logger
	mux        *chi.Mux
	controller *PostController
}

// NewBlogServer is builder for BlogServer
func NewBlogServer() *BlogServer {
	spr := NewSimplePostRepo(30)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customHTTPLogger)
	bs := &BlogServer{
		mux:        r,
		log:        logrus.New(),
		controller: NewPostController(spr),
	}
	return bs
}

// Run is running blogServer
func (bs *BlogServer) Run(hostPort string) {
	bs.registerPostRoutes()
	bs.log.Infof("http server starting on the [%s] tcp port", hostPort)
	bs.log.Fatal(http.ListenAndServe(hostPort, bs.mux))
}

func (bs *BlogServer) registerPostRoutes() {
	bs.mux.Route("/posts", func(r chi.Router) {
		r.Get("/", bs.controller.GetPosts)
		r.Get("/{"+postID+"}", bs.controller.GetOnePost)
	})
}

// [CUSTOM MIDDLEWARE]

// filterContentType - middleware to check content type as JSON
func filterContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Filtering requests by MIME type
		if r.Method == "POST" { // filter for POST request
			if r.Header.Get("Content-type") != "application/json" {
				render.Render(w, r, ErrUnsupportedFormat)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// customHTTPLogger - middleware to logrus logger
func customHTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).String()
		host, _ := os.Hostname()
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"proto":  r.Proto,
			"remote": r.RemoteAddr,
			"url":    r.RequestURI,
			//"code":     r.Response.StatusCode,
			"duration": duration,
		}).Infof("%s", host)
	})
}

// [HANDLER FUNCS]

// PostController - main controller for Posts
type PostController struct {
	PostRepo domain.PostRepository
	//CommentsRepo domain.CommentsRepository
}

// NewPostController is a builder for PostController
func NewPostController(repo domain.PostRepository) *PostController {
	pc := &PostController{
		PostRepo: repo,
	}
	return pc
}

// GetPosts - handler func for search query text at the Sites
func (pc *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	if limit == 0 {
		limit = 10
	}
	posts, err := pc.PostRepo.Find(limit, offset)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	if len(posts) == 0 {
		render.Render(w, r, ErrNotFound(err))
		return
	}
	data := templatePostsFill{
		Title: "POSTS",
		Posts: posts,
	}
	tmpl := template.Must(template.New("indexPOST").ParseFiles(templatePOSTS))
	tmpl.ExecuteTemplate(w, "indexPOST", data)
}

// GetOnePost - handler func for search query text at the Sites
func (pc *PostController) GetOnePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, postID)
	post, err := pc.PostRepo.FindByID(id)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexSinglePOST").ParseFiles(templatePOSTS))
	tmpl.ExecuteTemplate(w, "indexSinglePOST", data)
}

type templatePostsFill struct {
	Title string
	Posts []domain.PostInBlog
}

type templateOnePostFill struct {
	Title string
	Post  domain.PostInBlog
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render - implement method Render for render.Renderer
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest - wrapper for make err structure
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrServerInternal - wrapper for make err structure
func ErrServerInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound - wrapper for make err structure for empty result
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		ErrorText:      err.Error(),
	}
}

// ErrUnsupportedFormat - 415 error implementation
var ErrUnsupportedFormat = &ErrResponse{HTTPStatusCode: http.StatusUnsupportedMediaType, StatusText: "415 - Unsupported Media Type. Please send JSON"}
