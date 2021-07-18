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