package presenters

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pyoko/pyogo-rest-template/pkg/models"
)

type User struct {
	*models.User

	// Additional fields.
	SelfURL string `json:"SelfURL" xml:"SelfURL"`
}

func (a *User) Render(w http.ResponseWriter, r *http.Request) error {
	a.SelfURL = fmt.Sprintf("http://localhost:3000/v1/users/%v", a.ID)
	return nil
}

func UserResponse(user *models.User) *User {
	return &User{User: user}
}

func UserListResponse(users []*models.User) []render.Renderer {
	list := []render.Renderer{}

	for _, user := range users {
		list = append(list, UserResponse(user))
	}

	return list
}