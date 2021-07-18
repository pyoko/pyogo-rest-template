package presenters

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pyoko/pyogo-rest-template/pkg/models"
)

type Post struct {
	*models.Post

	// Additional fields.
	SelfURL string `json:"SelfURL" xml:"SelfURL"`

	// Omitted fields.
	// URL interface{} `json:"url,omitempty" xml:"url,omitempty"`
}

func (a *Post) Render(w http.ResponseWriter, r *http.Request) error {
	a.SelfURL = fmt.Sprintf("http://localhost:3000/v1/posts/%v", a.ID)
	return nil
}

func PostResponse(post *models.Post) *Post {
	return &Post{Post: post}
}

func PostListResponse(posts []*models.Post) []render.Renderer {
	list := []render.Renderer{}

	for _, post := range posts {
		list = append(list, PostResponse(post))
	}

	return list
}