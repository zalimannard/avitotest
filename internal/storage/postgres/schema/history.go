package schema

import (
	"avitotest/internal/storage/schema"
	"fmt"
	"time"
)

func (s *Storage) InsertIntoHistory(userId int, slug string, actionType string) error {
	_, err := s.Db.Exec(`
		INSERT INTO history (id_user, 
		                     segment_name, 
		                     action_type, 
		                     action_date)
		VALUES ($1, 
		        $2, 
		        $3, 
		        $4)
	`, userId, slug, actionType, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ReadHistoryRecordsForMonth(year, month int) ([]schema.HistoryRecord, error) {
	var records []schema.HistoryRecord
	startDate := fmt.Sprintf("%d-%02d-01", year, month)
	endDate := fmt.Sprintf("%d-%02d-01", year, month+1)
	if month == 12 {
		endDate = fmt.Sprintf("%d-01-01", year+1)
	}

	rows, err := s.Db.Query(`
		SELECT 
		    id_user, 
		    segment_name, 
		    action_type, 
		    action_date
		FROM 
		    history
		WHERE 
		    action_date >= $1 
		  AND 
		    action_date < $2
		ORDER BY 
		    action_date ASC
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record schema.HistoryRecord
		err := rows.Scan(&record.UserId, &record.Segment, &record.Action, &record.Timestamp)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
