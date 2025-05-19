package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *HttpHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var creds LoginRequest
	if err := render.Decode(r, &creds); err != nil {
		slog.Error("Failed to decode login request", "error", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	slog.Info("Login request received", "username", creds.Username, "password", creds.Password)

	// TODO: validate credentials
	// TODO: create session
	// TODO: redirect to home page
}
