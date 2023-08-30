package segment

import (
	"avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/http-server/validators"
	"avitotest/internal/lib/api/response"
	"avitotest/internal/lib/logger/sl"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
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

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("Failed to decode request body", sl.Err(err))
		if errors.Is(err, io.EOF) {
			return nil, errors.New("Empty request body")
		}
		return nil, errors.New("Failed to decode request")
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
