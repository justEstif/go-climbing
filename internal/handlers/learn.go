package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
)

func LearnListPage(w http.ResponseWriter, r *http.Request) {
	items, err := database.DB.ListAllLearnContent(r.Context())
	if err != nil {
		items = nil
	}

	// Group by category (DB already returns sorted by category)
	var groups []components.LearnCategory
	for _, item := range items {
		if len(groups) == 0 || groups[len(groups)-1].Name != item.Category {
			groups = append(groups, components.LearnCategory{Name: item.Category})
		}
		groups[len(groups)-1].Items = append(groups[len(groups)-1].Items, item)
	}

	components.LearnList(groups).Render(r.Context(), w)
}

func LearnDetailPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid learn content ID", http.StatusBadRequest)
		return
	}

	item, err := database.DB.GetLearnContent(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	components.LearnDetail(item).Render(r.Context(), w)
}
