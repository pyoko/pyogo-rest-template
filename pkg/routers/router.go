package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/pyoko/pyogo-rest-template/pkg/models"
	v1 "github.com/pyoko/pyogo-rest-template/pkg/routers/v1"
)

type Router struct {
	persistence *db.DB
}

func NewRouter(db *db.DB) *Router {
	return &Router{
		persistence: db,
	}
}

func (ro *Router) Init() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		v1Router := v1.NewRouter(ro.persistence)
		r.Mount("/posts", v1Router.PostInit())
		r.Mount("/users", v1Router.UserInit())
	})

	return r
}