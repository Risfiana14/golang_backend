package repository

import (
	"tugas5/app/model"
	"database/sql"
	"time"
	"fmt"
)

type AlumniRepository interface {
	GetAll(search, sortBy, order string, limit, offset int) ([]model.Alumni, error)
	GetByID(id int) (*model.Alumni, error)
	Create(req model.CreateAlumniRequest) (*model.Alumni, error)
	Update(id int, req model.UpdateAlumniRequest) (*model.Alumni, error)
	Delete(id int) error
}

type alumniRepository struct {
	db *sql.DB
}

func NewAlumniRepository(db *sql.DB) AlumniRepository {
	return &alumniRepository{db: db}
}

func (r *alumniRepository) GetAll(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := fmt.Sprintf(`
	SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
	FROM alumni
	WHERE nama ILIKE $1 OR email ILIKE $1
	ORDER BY %s %s
	LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

func (r *alumniRepository) GetByID(id int) (*model.Alumni, error) {
	var a model.Alumni
	row := r.db.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE id = $1
	`, id)

	err := row.Scan(
		&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
		&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *alumniRepository) Create(req model.CreateAlumniRequest) (*model.Alumni, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), time.Now()).Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *alumniRepository) Update(id int, req model.UpdateAlumniRequest) (*model.Alumni, error) {
	result, err := r.db.Exec(`
		UPDATE alumni
		SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, email = $5, no_telepon = $6, alamat = $7, updated_at = $8
		WHERE id = $9
	`, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), id)

	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.GetByID(id)
}

func (r *alumniRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM alumni WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
