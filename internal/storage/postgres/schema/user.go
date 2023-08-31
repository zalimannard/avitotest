package schema

func (s *Storage) GetTotalUsersCount() (totalUsers int, err error) {
	err = s.Db.QueryRow(`
		SELECT COUNT(*) FROM users
	`).Scan(&totalUsers)
	if err != nil {
		return 0, err
	}
	return totalUsers, nil
}

func (s *Storage) SelectRandomUsers(limit int) (userIds []int, err error) {
	rows, err := s.Db.Query(`
		SELECT id FROM users ORDER BY RANDOM() LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIds = append(userIds, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIds, nil
}

func (s *Storage) UserExists(userId int) (bool, error) {
	var exists bool
	err := s.Db.QueryRow(`
		SELECT
		    EXISTS(
		    	SELECT
		    	    1
		    	FROM 
		    	    users 
		    	WHERE 
		    	    id = $1
		    )
	`, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
