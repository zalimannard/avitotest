package schema

import "github.com/lib/pq"

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
