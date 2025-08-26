package repo

import (
	"fmt"
	"kafkapractisel0/models"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepo interface {
	CreateDelivery(dto models.Delivery) error
	GetRandomDelivery() (int, error)
	CheckExist(uid int) error
}

type deliveryRepo struct {
	db *sqlx.DB
}

func NewDeliveryRepo(db *sqlx.DB) DeliveryRepo {
	return &deliveryRepo{db}
}

func (r *deliveryRepo) CheckExist(uid int) error {
	query := `
		SELECT count(*)
		FROM public."Delivery" d
		where d.uid = $1
		`
	var count int
	err := r.db.Get(&count, query, uid)
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("count != 1, count = %v", count)
	}
	return nil
}

func (r *deliveryRepo) CreateDelivery(dto models.Delivery) error {
	query := `
		INSERT INTO public."Delivery"
			("name", phone, zip, city, address, region, email)
		VALUES
			(:name, :phone, :zip, :city, :address, :region, :email);
			`
	_, err := r.db.NamedExec(query, dto)
	if err != nil {
		return err
	}
	return nil
}

func (r *deliveryRepo) GetRandomDelivery() (int, error) {
	query := `
		SELECT uid
		FROM public."Delivery"
		ORDER BY RANDOM()
		LIMIT 1
		`
	var uid int
	err := r.db.Get(&uid, query)
	if err != nil {
		return uid, err
	}
	return uid, nil
}
