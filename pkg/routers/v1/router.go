package router

import "github.com/pyoko/pyogo-rest-template/pkg/models"

type ctxKey struct{}

type Router struct {
	persistence *models.DB
}

func NewRouter(db *models.DB) *Router {
	return &Router{
		persistence: db,
	}
}
