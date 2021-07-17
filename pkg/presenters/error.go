package presenters

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidRequest = errors.New("invalid request")
	ErrResponse		  = errors.New("error rendering response")
)

func PresentError(r *http.Request, err error) (*http.Request, interface {}) {
	switch err {
		case ErrUnauthorized:
			render.Status(r, http.StatusUnauthorized)
		case ErrForbidden:
			render.Status(r, http.StatusForbidden)
		case ErrNotFound:
			render.Status(r, http.StatusNotFound)
		case ErrInvalidRequest:
			render.Status(r, http.StatusBadRequest)
		case ErrResponse:
			render.Status(r, http.StatusUnprocessableEntity          )
		default:
			render.Status(r, http.StatusInternalServerError)
	}
	return r, map[string]string{"error": err.Error()}
}