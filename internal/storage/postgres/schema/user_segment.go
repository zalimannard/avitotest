package schema

import (
	"github.com/lib/pq"
	"strconv"
)

func (s *Storage) AssignSegmentsToUser(userId int, segmentIds []int) error {
	_, err := s.Db.Exec(`
		INSERT INTO 
		    users_segments (
		        id_user, 
		        id_segment
		    )
		VALUES (
		    $1, 
		    unnest($2::integer[])
		)
	`, userId, pq.Array(segmentIds))
	return err
}
func (s *Storage) DischargeSegmentsToUser(userId int, segmentIds []int) error {
	_, err := s.Db.Exec(`
		DELETE FROM
		    users_segments
		WHERE
		    id_user = $1
		  AND 
		    id_segment IN (SELECT unnest($2::integer[]))
	`, userId, pq.Array(segmentIds))
	return err
}

func (s *Storage) SelectSegmentIdBySlug(slug string) (int, error) {
	var segmentId int
	err := s.Db.QueryRow(`
		SELECT
			id
		FROM
			segments
		WHERE
			slug = $1
	`, slug).Scan(&segmentId)
	if err != nil {
		return 0, err
	}
	return segmentId, nil
}

func (s *Storage) AssignSegmentsToUsers(userIds []int, segmentId int) error {
	query := `
		INSERT INTO 
		    users_segments (
		        id_user, 
		        id_segment
		    )
		VALUES 
	`

	values := []interface{}{}
	for i, userID := range userIds {
		if i > 0 {
			query += ","
		}
		query += "($" + strconv.Itoa(i*2+1) + ", $" + strconv.Itoa(i*2+2) + ")"
		values = append(values, userID, segmentId)
	}

	_, err := s.Db.Exec(query, values...)
	return err
}
