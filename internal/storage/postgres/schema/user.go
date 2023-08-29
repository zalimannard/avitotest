package schema

import "database/sql"

type User struct {
	Id   int
	Slug string
}

func InsertUser(db *sql.DB, request User) (createdEntityId int, err error) {
	err = db.QueryRow(`
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
	`, request.Slug).Scan(&createdEntityId)
	if err != nil {
		return 0, err
	}
	return createdEntityId, nil
}

func SelectUserById(db *sql.DB, id int) (user User, err error) {
	err = db.QueryRow(`
		SELECT
		    id,
		    name
		FROM
		    users
		WHERE
		    id = $1
	`, id).Scan(
		&user.Id,
		&user.Slug)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func DeleteUserById(db *sql.DB, id int) (err error) {
	_, err = db.Exec(`
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
