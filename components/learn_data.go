package components

import (
	"strings"

	"github.com/justestif/go-climbing/internal/database"
)

type LearnCategory struct {
	Name  string
	Items []database.LearnContent
}

// toEmbedURL converts a YouTube watch URL to an embed URL.
func toEmbedURL(watchURL string) string {
	return strings.Replace(watchURL, "watch?v=", "embed/", 1)
}
