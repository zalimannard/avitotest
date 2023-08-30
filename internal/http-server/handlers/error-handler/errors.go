package error_handler

import (
	"avitotest/internal/lib/api/response"
	"github.com/go-chi/render"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	render.JSON(w, r, response.ErrorMessage(message))
}
