package repository

import (
	"tugas5/app/model"
	"database/sql"
)

func Login(db *sql.DB, username string, password string) (model.User, error) {
	var user model.User
	row := db.QueryRow("SELECT id, username, email, password_hash, role FROM users WHERE username = $1", username)
	var hashedPassword string
	err := row.Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}	
		return user, err
	}

	// Bandingkan password yang diberikan dengan hash yang disimpan
	if hashedPassword != password {
		return user, err
	}
	return user, nil
}