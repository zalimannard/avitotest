package segment

import (
	"avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/http-server/validators"
	"avitotest/internal/lib/api/response"
	"avitotest/internal/lib/logger/sl"
	"avitotest/internal/storage/schema"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type CreateSegmentRequest struct {
	Slug string `json:"slug" validate:"required,validateslug"`
}

type CreateSegmentResponse struct {
	response.Response
	Id int `json:"id"`
}

type CreateSegmentHandler interface {
	InsertSegment(request schema.Segment) (createdEntityId int, err error)
}

func parseAndValidate(r *http.Request, log slog.Logger) (*CreateSegmentRequest, error) {
	var req CreateSegmentRequest

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

func Insert(log slog.Logger, handler CreateSegmentHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req, err := parseAndValidate(r, log)
		if err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("Request body decoded", slog.Any("request", *req))

		id, err := handler.InsertSegment(schema.Segment{Slug: req.Slug})
		if err != nil {
			log.Error("Failed to save segment", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusConflict, "Failed to save segment")
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, CreateSegmentResponse{
			Response: response.OkMessage(),
			Id:       id,
		})
	}
}
