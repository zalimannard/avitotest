package schema

import "database/sql"

type UsersSegment struct {
	Id        int
	UserId    int
	SegmentId int
}

func InsertUsersSegment(db *sql.DB, request UsersSegment) (createdEntityId int, err error) {
	err = db.QueryRow(`
		INSERT INTO
			users_segments (
			    id_user,
			    id_segment
			)
		VALUES (
		    $1,
		    $1
		)
		RETURNING ( 
		    id
		)
	`, request.Id, request.UserId, request.SegmentId).Scan(&createdEntityId)
	if err != nil {
		return 0, err
	}
	return createdEntityId, nil
}

func SelectUsersSegmentsByUserId(db *sql.DB, userId int) (usersSegments []UsersSegment, err error) {
	rows, err := db.Query(`
		SELECT
		    id,
		    id_user,
		    id_segment
		FROM
		    users_segments
		WHERE
		    id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var usersSegment UsersSegment
		err = rows.Scan(
			&usersSegment.Id,
			&usersSegment.UserId,
			&usersSegment.SegmentId)
		if err != nil {
			return nil, err
		}
		usersSegments = append(usersSegments, usersSegment)
	}
	return usersSegments, err
}

func DeleteUsersSegmentById(db *sql.DB, id int) (err error) {
	_, err = db.Exec(`
		DELETE FROM
		    users_segments
		WHERE
		    id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
