package history

import (
	"avitotest/internal/http-server/handlers/error-handler"
	"avitotest/internal/lib/logger/sl"
	"avitotest/internal/storage/schema"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type HistoryReportForMonthRequest struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

type HistoryReportForMonthHandler interface {
	ReadHistoryRecordsForMonth(year, month int) ([]schema.HistoryRecord, error)
}

func parseAndValidateHistoryRequest(r *http.Request, log slog.Logger) (HistoryReportForMonthRequest, error) {
	var req HistoryReportForMonthRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("Error decoding request body", sl.Err(err))
		return req, errors.New("Error decoding request body")
	}

	if req.Year <= 0 || req.Month <= 0 || req.Month > 12 {
		log.Error("Invalid year or month", sl.Err(errors.New("invalid year or month")))
		return req, errors.New("invalid year or month")
	}

	return req, nil
}

func ReadHistoryRecordsForMonth(log slog.Logger, handler HistoryReportForMonthHandler, serverAddress string, reportDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseAndValidateHistoryRequest(r, log)
		if err != nil {
			log.Error("Error in parsing and validating request", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusBadRequest, "Error in request parsing and validation")
			return
		}

		records, err := handler.ReadHistoryRecordsForMonth(req.Year, req.Month)
		if err != nil {
			log.Error("Failed to fetch history records", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusConflict, "Failed to retrieve history records")
			return
		}

		fileName, err := generateCSV(records, reportDir)
		if err != nil {
			log.Error("Failed to create csv", sl.Err(err))
			error_handler.HandleError(w, r, http.StatusConflict, "Failed to create csv")
			return
		}

		url := fmt.Sprintf("http://%s/%s", serverAddress, fileName)
		render.JSON(w, r, map[string]string{"url": url})
	}
}

func generateCSV(records []schema.HistoryRecord, reportDir string) (string, error) {
	if err := os.MkdirAll("reports", 0755); err != nil {
		return "", err
	}
	fileName := path.Join(reportDir, fmt.Sprintf("report_%d.csv", time.Now().UnixNano()))
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"User ID", "Segment", "Action", "Date and Time"})

	for _, record := range records {
		writer.Write([]string{
			strconv.Itoa(record.UserId),
			record.Segment,
			record.Action,
			record.Timestamp.Format(time.RFC3339),
		})
	}

	return fileName, nil
}
