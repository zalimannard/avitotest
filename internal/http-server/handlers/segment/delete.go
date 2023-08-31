package segment

import (
	"avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/http-server/validators"
	"avitotest/internal/lib/api/response"
	"avitotest/internal/lib/logger/sl"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type DeleteSegmentRequest struct {
	Slug string `json:"slug" validate:"required,validateslug"`
}

type DeleteSegmentResponse struct {
	response.Response
}

type DeleteSegmentHandler interface {
	DeleteSegmentBySlug(slug string) error
}

func parseAndValidateDeleteRequest(r *http.Request, log slog.Logger) (*DeleteSegmentRequest, error) {
	var req DeleteSegmentRequest

	req.Slug = r.URL.Query().Get("slug")
	if req.Slug == "" {
		log.Error("Slug parameter missing")
		return nil, errors.New("Slug parameter is required")
	}

	if err := validators.Instance.Struct(req); err != nil {
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			return nil, response.ValidationError(validateErr)
		}
		return nil, err
	}

	return &req, nil
}

func Delete(log slog.Logger, handler DeleteSegmentHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req, err := parseAndValidateDeleteRequest(r, log)
		if err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("Request body decoded", slog.Any("request", *req))

		if err := handler.DeleteSegmentBySlug(req.Slug); err != nil {
			log.Error("Failed to delete segment", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusConflict, "Failed to delete segment")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		render.JSON(w, r, DeleteSegmentResponse{
			Response: response.OkMessage(),
		})
	}
}
