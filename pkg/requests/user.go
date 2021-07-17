package requests

import (
	"errors"
	"net/http"

	"github.com/pyoko/gorest/pkg/models"
)


type UserRequest struct {
	*models.User
}

func (a *UserRequest) Bind(r *http.Request) error {
	if a.User == nil {
		return errors.New("missing required user fields")
	}

	return nil
}