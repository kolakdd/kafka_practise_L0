package repo

import (
	"fmt"
	"kafkapractisel0/models"

	"github.com/jmoiron/sqlx"
)

type ItemsRepo interface {
	CreateItem(item models.Item) error
	CheckExistMulti(uids []int) error
}

type itemsRepo struct {
	db *sqlx.DB
}

func NewItemsRepo(db *sqlx.DB) ItemsRepo {
	return &itemsRepo{db}
}

func (r *itemsRepo) CheckExistMulti(uids []int) error {
	query, args, err := sqlx.In(`
			SELECT count(*)
			FROM public."Items" i
			WHERE i.chrt_id IN (?);
			`,
		uids)
	if err != nil {
		return err
	}
	var count int
	query = r.db.Rebind(query)
	err = r.db.Get(&count, query, args...)
	if err != nil {
		return err
	}
	if count != len(uids) {
		return fmt.Errorf("kek count != len(uids), count = %v, uids = %v", count, uids)
	}
	return nil
}

func (r *itemsRepo) CreateItem(item models.Item) error {
	query := `
	INSERT INTO public."Items"
	(track_number, price, rid, "name", sale, "size", currency, total_price, nm_id, brand, status)
	VALUES(:track_number, :price, :rid, :name, :sale, :size, :currency, :total_price, :nm_id, :brand, :status);
	`
	_, err := r.db.NamedExec(query, item)
	if err != nil {
		return err
	}
	return nil
}
