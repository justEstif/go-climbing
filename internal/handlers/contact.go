package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/justestif/go-climbing/components"
)

func ContactForm(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	components.Contact(csrfToken).Render(r.Context(), w)
}
