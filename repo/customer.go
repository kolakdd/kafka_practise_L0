package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CustomerRepo interface {
	CreateCustomer(name string) error
	GetRandomCustomerUID() (int, error)
	CheckExist(uid int) error
}

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) CustomerRepo {
	return &customerRepo{db}
}

func (r *customerRepo) CheckExist(uid int) error {
	query := `
		SELECT count(*)
		FROM public."Customer" c
		where c.uid = $1
		`
	var count int
	err := r.db.Get(&count, query, uid)
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("count != 1")
	}
	return nil
}

func (r *customerRepo) CreateCustomer(name string) error {
	query := `INSERT INTO public."Customer" (name) VALUES ($1)`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (r *customerRepo) GetRandomCustomerUID() (int, error) {
	query := `
		SELECT uid
		FROM public."Customer"
		ORDER BY RANDOM()
		LIMIT 1
		`
	var uid int
	err := r.db.Get(&uid, query)
	if err != nil {
		return -1, err
	}
	return uid, nil
}
