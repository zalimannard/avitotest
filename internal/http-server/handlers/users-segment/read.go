package users_segment

import (
	"avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/lib/api/response"
	"avitotest/internal/lib/logger/sl"
	"avitotest/internal/storage/schema"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type ReadUsersSegmentsResponse struct {
	response.Response
	UserId int      `json:"userId"`
	Slugs  []string `json:"slugs"`
}

type ReadUsersSegmentsHandler interface {
	SelectSegmentsByUserId(userId int) (segments []schema.Segment, err error)
}

func parseAndValidateReadRequest(r *http.Request, log slog.Logger) (int, error) {
	userIdStr := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Error("Failed to convert userId to int", sl.Err(err))
		return 0, errors.New("Failed to convert userId to int")
	}
	return userId, nil
}

func ReadSlugs(log slog.Logger, handler ReadUsersSegmentsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userId, err := parseAndValidateReadRequest(r, log)
		if err != nil {
			log.Error("Error in parsing and validating request", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusBadRequest, "Error in request parsing and validation")
			return
		}

		segments, err := handler.SelectSegmentsByUserId(userId)
		if err != nil {
			log.Error("Failed to read segments", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusConflict, "Failed to retrieve segments")
			return
		}

		slugs := make([]string, 0)
		for _, segment := range segments {
			slugs = append(slugs, segment.Slug)
		}

		render.JSON(w, r, ReadUsersSegmentsResponse{
			Response: response.OkMessage(),
			UserId:   userId,
			Slugs:    slugs,
		})
	}
}
