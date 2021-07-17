package requests

import (
	"errors"
	"net/http"

	"github.com/pyoko/gorest/pkg/models"
)


type PostRequest struct {
	*models.Post
}

func (a *PostRequest) Bind(r *http.Request) error {
	if a.Post == nil {
		return errors.New("missing required post fields")
	}

	return nil
}