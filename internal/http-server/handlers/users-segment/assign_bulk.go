package users_segment

import (
	error_handler "avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/http-server/validators"
	"avitotest/internal/lib/api/response"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type AssignPercentRequest struct {
	Slug    string  `json:"slug" validate:"required,validateslug"`
	Percent float64 `json:"percent" validate:"required,gte=0,lte=100"`
}

type AssignPercentHandler interface {
	GetTotalUsersCount() (int, error)
	SelectRandomUsers(limit int) ([]int, error)
	SelectSegmentIdBySlug(slug string) (int, error)
	AssignSegmentsToUsers(userIds []int, segmentId int) error
	InsertIntoHistoryBulk(userIds []int, slug string, actionType string) error
}

func AssignPercent(log slog.Logger, handler AssignPercentHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req AssignPercentRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request JSON", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Failed to decode request JSON")
			return
		}

		if err := validators.Instance.Struct(req); err != nil {
			log.Error("Invalid request format", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Invalid request format")
			return
		}

		totalUsers, err := handler.GetTotalUsersCount()
		if err != nil {
			log.Error("Failed to get total users count", err)
			error_handler.HandleError(w, r, http.StatusInternalServerError, "Failed to get total users count")
			return
		}

		usersToAssignCount := int(float64(totalUsers) * req.Percent / 100)
		userIds, err := handler.SelectRandomUsers(usersToAssignCount)
		if err != nil {
			log.Error("Failed to fetch random users", err)
			error_handler.HandleError(w, r, http.StatusInternalServerError, "Failed to fetch random users")
			return
		}

		segmentId, err := handler.SelectSegmentIdBySlug(req.Slug)
		if err != nil {
			log.Error("Failed to fetch segment ID by slug", err)
			error_handler.HandleError(w, r, http.StatusNotFound, "Failed to fetch segment ID by slug")
			return
		}

		if err := handler.AssignSegmentsToUsers(userIds, segmentId); err != nil {
			log.Error("Failed to assign segment to users", err)
			error_handler.HandleError(w, r, http.StatusBadRequest, "Failed to assign segment to users")
			return
		}

		if err := handler.InsertIntoHistoryBulk(userIds, req.Slug, "added"); err != nil {
			log.Error("Failed to insert segment assignment into history", err)
			error_handler.HandleError(w, r, http.StatusInternalServerError, "Failed to insert segment assignment into history")
			return
		}

		render.JSON(w, r, response.OkMessage())
	}
}
