package repo

import (
	"kafkapractisel0/models"

	"github.com/jmoiron/sqlx"
)

type PaymentRepo interface {
	CreatePaymentTx(tx *sqlx.Tx, payment models.Payment) error
}

type paymentRepo struct {
	db *sqlx.DB
}

func NewPaymentRepo(db *sqlx.DB) PaymentRepo {
	return &paymentRepo{db}
}

func (r *paymentRepo) CreatePaymentTx(tx *sqlx.Tx, payment models.Payment) error {
	query := `
		INSERT INTO public."Payment"
		("transaction", request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES(:transaction, :request_id, :currency, :provider, :amount, 
		:payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
		`
	_, err := tx.NamedExec(query, payment)
	if err != nil {
		return err
	}
	return nil
}
