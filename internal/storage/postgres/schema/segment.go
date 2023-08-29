package schema

import "database/sql"

type Segment struct {
	Id   int
	Slug string
}

func InsertSegment(db *sql.DB, request Segment) (createdEntityId int, err error) {
	err = db.QueryRow(`
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

func SelectSegmentById(db *sql.DB, id int) (segment Segment, err error) {
	err = db.QueryRow(`
		SELECT
		    id,
		    slug
		FROM
		    segments
		WHERE
		    id = $1
	`, id).Scan(
		&segment.Id,
		&segment.Slug)
	if err != nil {
		return Segment{}, err
	}
	return segment, err
}

func DeleteSegmentById(db *sql.DB, id int) (err error) {
	_, err = db.Exec(`
		DELETE FROM
		    segments
		WHERE
		    id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
