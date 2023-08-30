package schema

import (
	"avitotest/internal/storage/schema"
	"github.com/lib/pq"
)

func (s *Storage) InsertSegment(request schema.Segment) (createdEntityId int, err error) {
	err = s.Db.QueryRow(`
		INSERT INTO
			segments (
			    slug
			)
		VALUES (
		    $1
		)
		RETURNING ( 
		    id
		)
	`, request.Slug).Scan(&createdEntityId)
	if err != nil {
		return 0, err
	}
	return createdEntityId, nil
}

func (s *Storage) DeleteSegmentBySlug(slug string) (err error) {
	_, err = s.Db.Exec(`
		DELETE FROM
		    segments
		WHERE
		    slug = $1
	`, slug)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SelectSegmentIdsBySlugs(slugs []string) ([]int, error) {
	rows, err := s.Db.Query(`
		SELECT 
		    id 
		FROM 
		    segments 
		WHERE 
		    slug = ANY($1)
	`, pq.Array(slugs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Storage) SelectSegmentsByUserId(userId int) (segments []schema.Segment, err error) {
	rows, err := s.Db.Query(`
		SELECT
			s.id,
			s.slug
		FROM
			users_segments AS us
		JOIN
			segments AS s ON us.id_segment = s.id
		WHERE
			us.id_user = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var segment schema.Segment
		err := rows.Scan(&segment.Id, &segment.Slug)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}
