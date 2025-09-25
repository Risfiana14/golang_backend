package repository

import (
	"database/sql"
	"fmt"
	"log"
	"tugas5/app/model"
	"tugas5/database"
)

// GetUsersRepo -> ambil data users dari DB
func GetUsersRepo(search, sortBy, order string, limit, offset int) ([]model.User, error) {
	query := fmt.Sprintf(`
	SELECT id, name, email, created_at
	FROM users
	WHERE name ILIKE $1 OR email ILIKE $1
	ORDER BY %s %s
	LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)

	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}

	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// CountUsersRepo -> hitung total data untuk pagination
func CountUsersRepo(search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE name ILIKE $1 OR
email ILIKE $1`
	err := database.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
