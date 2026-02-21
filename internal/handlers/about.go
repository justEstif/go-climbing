package handlers

import (
	"net/http"

	"github.com/justestif/go-climbing/components"
)

func About(w http.ResponseWriter, r *http.Request) {
	components.About().Render(r.Context(), w)
}
