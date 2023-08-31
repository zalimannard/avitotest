package users_segment

import (
	"avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/http-server/validators"
	"avitotest/internal/lib/api/response"
	"avitotest/internal/storage/schema"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type AssignSlugsRequest struct {
	Slugs []string `json:"slugs" validate:"required,dive,validateslug"`
}

type AssignSlugsHandler interface {
	SelectSegmentsByUserId(userId int) (segments []schema.Segment, err error)
	SelectSegmentIdsBySlugs(slugs []string) ([]int, error)
	AssignSegmentsToUser(userId int, segmentIds []int) error
	InsertIntoHistory(userId int, slug string, actionType string) error
}

func AssignSlugs(log slog.Logger, handler AssignSlugsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userIdStr := chi.URLParam(r, "userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Error("Failed to convert userId to int", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Failed to convert userId to int")
			return
		}

		var req AssignSlugsRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request JSON", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Failed to decode request JSON")
			return
		}

		if err := validators.Instance.Struct(req); err != nil {
			log.Error("Invalid slug format", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Invalid slug format")
			return
		}

		currentSegments, err := handler.SelectSegmentsByUserId(userId)
		if err != nil {
			log.Error("Failed to fetch current slugs for user", err)
			error_handler.HandleError(w, r, http.StatusInternalServerError, "Error fetching current slugs for user")
			return
		}
		currentSlugs := extractSlugsFromSegments(currentSegments)

		var newSlugs []string
		for _, slug := range req.Slugs {
			if !contains(currentSlugs, slug) {
				newSlugs = append(newSlugs, slug)
			}
		}

		segmentIds, err := handler.SelectSegmentIdsBySlugs(newSlugs)
		if err != nil {
			log.Error("Failed to fetch segment IDs by slugs", err)
			error_handler.HandleError(w, r, http.StatusNotFound, "Failed to fetch segment IDs by slugs")
			return
		}

		if err := handler.AssignSegmentsToUser(userId, segmentIds); err != nil {
			log.Error("Failed to assign slugs to user", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Failed to assign slugs to user")
			return
		}

		for _, newSlug := range newSlugs {
			if err := handler.InsertIntoHistory(userId, newSlug, "added"); err != nil {
				log.Error("Failed to insert slug assignment into history", err)
				error_handler.HandleError(w, r, http.StatusInternalServerError, "Failed to insert slug assignment into history")
				return
			}
		}

		render.JSON(w, r, response.OkMessage())
	}
}

func extractSlugsFromSegments(segments []schema.Segment) []string {
	var slugs []string
	for _, segment := range segments {
		slugs = append(slugs, segment.Slug)
	}
	return slugs
}

func contains(slugs []string, slug string) bool {
	for _, s := range slugs {
		if s == slug {
			return true
		}
	}
	return false
}
