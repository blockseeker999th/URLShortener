package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	render.JSON(w, r, v)
}
