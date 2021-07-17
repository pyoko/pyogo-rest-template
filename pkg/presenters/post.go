package presenters

import (
	"fmt"
	"net/http"

	"github.com/pyoko/gorest/pkg/models"
)

type Post struct {
	*models.Post

	// Additional fields.
	SelfURL string `json:"self_url" xml:"self_url"`

	// Omitted fields.
	URL interface{} `json:"url,omitempty" xml:"url,omitempty"`
}

func (a *Post) Render(w http.ResponseWriter, r *http.Request) error {
	a.SelfURL = fmt.Sprintf("http://localhost:3000/posts/%v", a.ID)
	return nil
}

func PostResponse(post *models.Post) *Post {
	return &Post{Post: post}
}