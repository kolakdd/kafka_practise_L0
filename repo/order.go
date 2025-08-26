package repo

import (
	"fmt"
	"kafkapractisel0/models"

	"github.com/jmoiron/sqlx"
)

type OrderRepo interface {
	SelectOrderById(uid int64) (*models.Order, error)
	CreateOrderTx(tx *sqlx.Tx, payment models.OrderMessage) (int, error)
	CreateOrderXItemsTx(tx *sqlx.Tx, order_uid int, items_uid []int) error
	SelectNewestWithOffset(offset int) (*models.Order, error)
}

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) OrderRepo {
	return &orderRepo{db}
}

func (r *orderRepo) SelectOrderById(uid int64) (*models.Order, error) {
	var order models.Order

	err := r.db.Get(&order, `
		select 
			o.*,
			d.uid as "delivery.uid", d."name" as "delivery.name", d.phone as "delivery.phone",
			d.zip as "delivery.zip", d.city as "delivery.city", d.address as "delivery.address",
			d.region as "delivery.region", d.email as "delivery.email", d.created_at as "delivery.created_at",
			p."transaction" as "payment.transaction", p.request_id as "payment.request_id",
			p.currency as "payment.currency", p.provider as "payment.provider",
			p.amount as "payment.amount", p.payment_dt as "payment.payment_dt",
			p.bank as "payment.bank", p.delivery_cost as "payment.delivery_cost",
			p.goods_total as "payment.goods_total", p.custom_fee as "payment.custom_fee",
			p.created_at as "payment.created_at",
			c."uid" as "customer.uid", c."name" as "customer.name", c.created_at as "customer.created_at"
		from "Order" o 
		join "Delivery" d on o.delivery = d.uid 
		join "Payment" p on o.payment = p.transaction 
		join "Customer" c on o.customer_id = c.uid 
		where o.uid = $1
	`, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	itemsQuery := `
        select i.*
        from "Items" i
        join "OrderXItems" oxi ON i.chrt_id = oxi.item_id
        where oxi.order_uid = $1
    `
	var items []models.Item
	err = r.db.Select(&items, itemsQuery, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	order.Items = items
	return &order, nil

}

func (r *orderRepo) SelectNewestWithOffset(offset int) (*models.Order, error) {
	var order models.Order
	err := r.db.Get(&order, `
		select 
			o.*,
			d.uid as "delivery.uid", d."name" as "delivery.name", d.phone as "delivery.phone",
			d.zip as "delivery.zip", d.city as "delivery.city", d.address as "delivery.address",
			d.region as "delivery.region", d.email as "delivery.email", d.created_at as "delivery.created_at",
			p."transaction" as "payment.transaction", p.request_id as "payment.request_id",
			p.currency as "payment.currency", p.provider as "payment.provider",
			p.amount as "payment.amount", p.payment_dt as "payment.payment_dt",
			p.bank as "payment.bank", p.delivery_cost as "payment.delivery_cost",
			p.goods_total as "payment.goods_total", p.custom_fee as "payment.custom_fee",
			p.created_at as "payment.created_at",
			c."uid" as "customer.uid", c."name" as "customer.name", c.created_at as "customer.created_at"
		from "Order" o 
		join "Delivery" d on o.delivery = d.uid 
		join "Payment" p on o.payment = p.transaction 
		join "Customer" c on o.customer_id = c.uid 
		order by d.created_at desc
		offset $1 limit 1
	`, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	itemsQuery := `
        select i.*
        from "Items" i
        join "OrderXItems" oxi ON i.chrt_id = oxi.item_id
        where oxi.order_uid = $1
    `
	var items []models.Item
	err = r.db.Select(&items, itemsQuery, order.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	order.Items = items
	return &order, nil

}

func (r *orderRepo) CreateOrderTx(tx *sqlx.Tx, order models.OrderMessage) (int, error) {
	var uid int
	query := `
		INSERT INTO public."Order"
		(customer_id, track_number, entry, delivery, payment, locale, 
		internal_signature, delivery_service, shardkey, 
		sm_id, oof_shard)
		VALUES
		(:customer_id, :track_number, :entry,:delivery, :payment.transaction,
		:locale, :internal_signature, :delivery_service, :shardkey,
		:sm_id, :oof_shard)
		RETURNING uid
		`
	query, args, err := sqlx.Named(query, order)
	if err != nil {
		return uid, fmt.Errorf("sqlx named error: %w", err)
	}
	query = tx.Rebind(query)
	err = tx.Get(&uid, query, args...)
	if err != nil {
		return uid, fmt.Errorf("order select error: %w", err)
	}
	return uid, nil
}

func (r *orderRepo) CreateOrderXItemsTx(tx *sqlx.Tx, order_uid int, items_uid []int) error {
	for i := 0; i < len(items_uid); i++ {
		query := `
		INSERT INTO public."OrderXItems"
		(order_uid, item_id)
		VALUES ($1, $2)`
		_, err := tx.Exec(query, order_uid, items_uid[i])
		if err != nil {
			return fmt.Errorf("order failed after exec: %w", err)
		}
	}
	return nil
}
