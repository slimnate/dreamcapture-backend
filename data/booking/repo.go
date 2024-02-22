package booking

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/slimnate/dreamcapture-backend/data"
)

type Repository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *Repository {
	db.Begin()
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS bookings(
	id SERIAL UNIQUE PRIMARY KEY,
	name VARCHAR(250) NOT NULL,
	phone VARCHAR(20) NOT NULL,
	email VARCHAR(100) NOT NULL,
	package_type VARCHAR(50) NOT NULL,
	session_type VARCHAR(50) NOT NULL,
	date DATE NOT NULL,
	time TIME NOT NULL,
	subjects VARCHAR(500),
	additional_info VARCHAR(1000),
	referral VARCHAR(100)
	)
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Create(obj Booking, id int64) (*Booking, error) {
	var lastInsertId int64

	query := `
	INSERT INTO
	bookings(name, phone, email, package_type, session_type, date, time, subjects, additional_info, referral)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	err := r.db.QueryRow(query, obj.Name, obj.Phone, obj.Email, obj.PackageType, obj.SessionType, obj.Date, obj.Time, obj.Subjects, obj.AdditionalInfo, obj.Referral).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	obj.Id = lastInsertId

	return &obj, nil
}

func (r *Repository) All() ([]Booking, error) {
	query := `
	SELECT id, name, phone, email, package_type, session_type, date, time, subjects, additional_info, referral
	FROM bookings
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.Id, &b.Name, &b.Phone, &b.Email, &b.PackageType, &b.SessionType, &b.Date, &b.Time, &b.Subjects, &b.AdditionalInfo, &b.Referral); err != nil {
			return nil, err
		}
		all = append(all, b)
	}
	return all, nil
}

func (r *Repository) GetByID(id int64) (*Booking, error) {
	query := `
	SELECT id, name, phone, email, package_type, session_type, date, time, subjects, additional_info, referral
	FROM bookings
	WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var b Booking
	if err := row.Scan(&b.Id, &b.Name, &b.Phone, &b.Email, &b.PackageType, &b.SessionType, &b.Date, &b.Time, &b.Subjects, &b.AdditionalInfo, &b.Referral); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, data.ErrNotExists
		}
		return nil, err
	}
	return &b, nil

}

func (r *Repository) Update(id int64, b Booking) (*Booking, error) {
	if id == 0 {
		return nil, errors.New(fmt.Sprintf("Invalid ID to update: %d", id))
	}

	query := `
		UPDATE bookings
		SET name = $1, phone = $2, email = $3, package_type = $4, session_type = $5, date = $6, time = $7, subjects = $8, additional_info = $9, referral = $10
		WHERE id = $11
	`
	res, err := r.db.Exec(query, b.Name, b.Phone, b.Email, b.PackageType, b.SessionType, b.Date, b.Time, b.Subjects, b.AdditionalInfo, b.Referral, id)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, data.ErrUpdateFailed
	}

	updated, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return data.ErrDeleteFailed
	}

	return err
}
