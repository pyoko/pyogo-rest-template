package router

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pyoko/gorest/pkg/models"
	"github.com/pyoko/gorest/pkg/presenters"
	"github.com/pyoko/gorest/pkg/requests"
)

type Router struct {
	persistence *models.DB
}

func NewRouter(db *models.DB) *Router {
	return &Router{
		persistence: db,
	}
}

type ctxKey struct{}

func (ro *Router) Init() http.Handler {
	r := chi.NewRouter()
	r.Get("/", ro.ListPosts)
	r.Route("/{postID}", func(r chi.Router) {
		r.Use(ro.PostCtx)
		r.Get("/", ro.GetPost)
		r.Put("/", ro.UpdatePost)
		r.Delete("/", ro.DeletePost)
	})
	r.Post("/", ro.CreatePost)

	return r
}

func (ro *Router) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := ro.persistence.GetPosts()
	if err != nil {
		panic(err)
	}

	// Start streaming data from the channel.
	render.Respond(w, r, posts)
}

func (ro *Router) PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var post *models.Post
		var err error

		if postID := chi.URLParam(r, "postID"); postID != "" {
			postID, _ := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
			post, err = ro.persistence.GetPostByID(postID) 
		} else {
			r, modelErr := presenters.PresentError(r, presenters.ErrNotFound)
			render.Respond(w, r, modelErr)

			panic(err)
		}

		if err != nil {
			r, modelErr := presenters.PresentError(r, presenters.ErrNotFound)
			render.Respond(w, r, modelErr)

			panic(err)
		}

		ctx := context.WithValue(r.Context(), ctxKey{}, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ro *Router) GetPost(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value(ctxKey{}).(*models.Post)

	payload := render.Renderer(presenters.PostResponse(post))

	render.Render(w, r, payload)
}

func (ro *Router) CreatePost(w http.ResponseWriter, r *http.Request) {
	data := &requests.PostRequest{}
	if err := render.Bind(r, data); err != nil {
		r, modelErr := presenters.PresentError(r, presenters.ErrInvalidRequest)
		render.Respond(w, r, modelErr)

		return
	}

	post := data.Post

	tx := ro.persistence.Begin()
	ro.persistence.InsertPost(post, tx)
	tx.Commit()

	
	payload := render.Renderer(presenters.PostResponse(post))
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload)
}

func (ro *Router) UpdatePost(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value(ctxKey{}).(*models.Post)

	data := &requests.PostRequest{Post: post}
	if err := render.Bind(r, data); err != nil {
		r, modelErr := presenters.PresentError(r, presenters.ErrInvalidRequest)
		render.Respond(w, r, modelErr)

		return
	}

	post = data.Post

	tx := ro.persistence.Begin()
	ro.persistence.UpdatePost(post, tx)
	tx.Commit()

	payload := render.Renderer(presenters.PostResponse(post))
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload)
}

func (ro *Router) DeletePost(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value(ctxKey{}).(*models.Post)

	tx := ro.persistence.Begin()
	err := ro.persistence.DeletePostByID(post, tx)
	tx.Commit()

	if err != nil {
		// r, modelErr := presenters.PresentError(r, presenters.ErrNotFound)
		// render.Respond(w, r, modelErr)

		log.Fatalln(err)
	}

	render.Status(r, http.StatusNoContent)
}