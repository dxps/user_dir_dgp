package handlers

import (
	"net/http"

	"github.com/dxps/user_dir_dgp/internal/ui"
)

type HttpHandlers struct{}

func (h *HttpHandlers) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	ui.Page("Home").Render(r.Context(), w)
}
