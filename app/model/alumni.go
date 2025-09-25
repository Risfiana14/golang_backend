package model

import "time"

type Alumni struct {
	ID         int       `json:"id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  *string   `json:"no_telepon,omitempty"`
	Alamat     *string   `json:"alamat,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateAlumniRequest struct {
	NIM        string  `json:"nim"`
	Nama       string  `json:"nama"`
	Jurusan    string  `json:"jurusan"`
	Angkatan   int     `json:"angkatan"`
	TahunLulus int     `json:"tahun_lulus"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
}

type UpdateAlumniRequest struct {
	Nama       string  `json:"nama"`
	Jurusan    string  `json:"jurusan"`
	Angkatan   int     `json:"angkatan"`
	TahunLulus int     `json:"tahun_lulus"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
}
