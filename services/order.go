package services

import (
	"kafkapractisel0/models"
	"kafkapractisel0/repo"

	"github.com/jmoiron/sqlx"
)

type OrderService interface {
	GetOrderById(uid int64) (*models.Order, error)
	CreateOrder(order models.OrderMessage) (*models.Order, error)
}

type orderService struct {
	db *sqlx.DB
	c  repo.CustomerRepo
	d  repo.DeliveryRepo
	i  repo.ItemsRepo
	o  repo.OrderRepo
	p  repo.PaymentRepo
}

func NewOrderService(db *sqlx.DB, c repo.CustomerRepo, d repo.DeliveryRepo, i repo.ItemsRepo, o repo.OrderRepo, p repo.PaymentRepo) OrderService {
	return &orderService{db, c, d, i, o, p}
}

func (s *orderService) GetOrderById(uid int64) (*models.Order, error) {
	return s.o.SelectOrderById(uid)
}

func (s *orderService) CreateOrder(order models.OrderMessage) (*models.Order, error) {
	var err error
	if err = s.c.CheckExist(order.CustomerID); err != nil {
		return nil, err
	}
	if err = s.d.CheckExist(order.DeliveryUID); err != nil {
		return nil, err
	}
	if err = s.i.CheckExistMulti(order.ItemsID); err != nil {
		return nil, err
	}

	tx, err := s.db.Beginx()

	if err != nil {
		return nil, err
	}
	if err = s.p.CreatePaymentTx(tx, order.Payment); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	order_uid, err := s.o.CreateOrderTx(tx, order)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = s.o.CreateOrderXItemsTx(tx, order_uid, order.ItemsID); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()

	orderDB, err := s.GetOrderById(order_uid)
	if err != nil {
		return nil, err
	}
	return orderDB, nil
}
