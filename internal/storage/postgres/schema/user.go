package schema

import (
	"avitotest/internal/storage/schema"
)

func (s *Storage) InsertUser(request schema.User) (createdEntityId int, err error) {
	err = s.Db.QueryRow(`
		INSERT INTO
			users (
			    name
			)
		VALUES (
		    $1
		)
		RETURNING ( 
		    id
		)
	`, request.Name).Scan(&createdEntityId)
	if err != nil {
		return 0, err
	}
	return createdEntityId, nil
}

func (s *Storage) SelectUserById(id int) (user schema.User, err error) {
	err = s.Db.QueryRow(`
		SELECT
		    id,
		    name
		FROM
		    users
		WHERE
		    id = $1
	`, id).Scan(
		&user.Id,
		&user.Name)
	if err != nil {
		return schema.User{}, err
	}
	return user, err
}

func (s *Storage) DeleteUserById(id int) (err error) {
	_, err = s.Db.Exec(`
		DELETE FROM
		    users
		WHERE
		    id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
