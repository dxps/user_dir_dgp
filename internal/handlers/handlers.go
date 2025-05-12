package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dxps/user_dir_dgp/internal/ui/pages"
)

type HttpHandlers struct{}

func (h *HttpHandlers) HomePageHandler(w http.ResponseWriter, r *http.Request) {

	if err := pages.Page("Home").Render(r.Context(), w); err != nil {
		slog.Error("Failed to render homepage", "error", err)
	}
}
