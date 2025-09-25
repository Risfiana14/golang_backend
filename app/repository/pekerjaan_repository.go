package repository

import (
	"database/sql"
	"time"
	"tugas5/app/model"
)

type PekerjaanRepository interface {
	GetAll() ([]model.Pekerjaan, error)
	GetByID(id int) (*model.Pekerjaan, error)
	GetByAlumniID(alumniID int) ([]model.Pekerjaan, error)
	Create(req model.CreatePekerjaanRequest) (*model.Pekerjaan, error)
	Update(id int, req model.UpdatePekerjaanRequest) (*model.Pekerjaan, error)
	Delete(id int) error
}

type pekerjaanRepository struct {
	db *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) PekerjaanRepository {
	return &pekerjaanRepository{db: db}
}

func (r *pekerjaanRepository) GetAll() ([]model.Pekerjaan, error) {
	rows, err := r.db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *pekerjaanRepository) GetByID(id int) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	row := r.db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan
		WHERE id = $1
	`, id)

	err := row.Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pekerjaanRepository) GetByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	rows, err := r.db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan
		WHERE alumni_id = $1
		ORDER BY created_at DESC
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *pekerjaanRepository) Create(req model.CreatePekerjaanRequest) (*model.Pekerjaan, error) {
	var id int
	var tanggalMulai, tanggalSelesai *time.Time

	// Parse tanggal mulai
	if req.TanggalMulaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
		if err != nil {
			return nil, err
		}
		tanggalMulai = &t
	}

	// Parse tanggal selesai jika ada
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return nil, err
		}
		tanggalSelesai = &t
	}

	err := r.db.QueryRow(`
		INSERT INTO pekerjaan (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		                             lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		                             status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, tanggalMulai, tanggalSelesai, req.StatusPekerjaan, req.DeskripsiPekerjaan,
		time.Now(), time.Now()).Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *pekerjaanRepository) Update(id int, req model.UpdatePekerjaanRequest) (*model.Pekerjaan, error) {
	var tanggalMulai, tanggalSelesai *time.Time

	// Parse tanggal mulai
	if req.TanggalMulaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
		if err != nil {
			return nil, err
		}
		tanggalMulai = &t
	}

	// Parse tanggal selesai jika ada
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return nil, err
		}
		tanggalSelesai = &t
	}

	result, err := r.db.Exec(`
		UPDATE pekerjaan
		SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3, lokasi_kerja = $4,
		    gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7, 
		    status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = $10
		WHERE id = $11
	`, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, tanggalMulai, tanggalSelesai, req.StatusPekerjaan, req.DeskripsiPekerjaan,
		time.Now(), id)

	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.GetByID(id)
}

func (r *pekerjaanRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM pekerjaan WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
