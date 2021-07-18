package router

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pyoko/pyogo-rest-template/pkg/models"
	"github.com/pyoko/pyogo-rest-template/pkg/presenters"
	"github.com/pyoko/pyogo-rest-template/pkg/requests"
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

// Post 
func (ro *Router) PostInit() http.Handler {
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

	if err := render.RenderList(w, r, presenters.PostListResponse(posts)); err != nil {
		r, modelErr := presenters.PresentError(r, presenters.ErrResponse)
		render.Respond(w, r, modelErr)

		return
	}
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

// User
func (ro *Router) UserInit() http.Handler {
	r := chi.NewRouter()
	r.Get("/", ro.ListUsers)
	r.Route("/{userID}", func(r chi.Router) {
		r.Use(ro.UserCtx)
		r.Get("/", ro.GetUser)
		r.Put("/", ro.UpdateUser)
		r.Delete("/", ro.DeleteUser)
	})
	r.Post("/", ro.CreateUser)

	return r
}

func (ro *Router) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ro.persistence.GetUsers()
	if err != nil {
		panic(err)
	}

	if err := render.RenderList(w, r, presenters.UserListResponse(users)); err != nil {
		r, modelErr := presenters.PresentError(r, presenters.ErrResponse)
		render.Respond(w, r, modelErr)

		return
	}
}

func (ro *Router) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error

		if userID := chi.URLParam(r, "userID"); userID != "" {
			userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
			user, err = ro.persistence.GetUserByID(userID) 
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

		ctx := context.WithValue(r.Context(), ctxKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ro *Router) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey{}).(*models.User)

	payload := render.Renderer(presenters.UserResponse(user))

	render.Render(w, r, payload)
}

func (ro *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &requests.UserRequest{}
	if err := render.Bind(r, data); err != nil {
		r, modelErr := presenters.PresentError(r, presenters.ErrInvalidRequest)
		render.Respond(w, r, modelErr)

		return
	}

	user := data.User

	tx := ro.persistence.Begin()
	ro.persistence.InsertUser(user, tx)
	tx.Commit()

	
	payload := render.Renderer(presenters.UserResponse(user))
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload)
}

func (ro *Router) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey{}).(*models.User)

	data := &requests.UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		r, modelErr := presenters.PresentError(r, presenters.ErrInvalidRequest)
		render.Respond(w, r, modelErr)

		return
	}

	user = data.User

	tx := ro.persistence.Begin()
	ro.persistence.UpdateUser(user, tx)
	tx.Commit()

	payload := render.Renderer(presenters.UserResponse(user))
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload)
}

func (ro *Router) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey{}).(*models.User)

	tx := ro.persistence.Begin()
	err := ro.persistence.DeleteUserByID(user, tx)
	tx.Commit()

	if err != nil {
		// r, modelErr := presenters.PresentError(r, presenters.ErrNotFound)
		// render.Respond(w, r, modelErr)

		log.Fatalln(err)
	}

	render.Status(r, http.StatusNoContent)
}