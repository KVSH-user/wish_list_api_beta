package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"wish_list/internal/entity"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, password, dbName string) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{db: db}

	cwd, _ := os.Getwd()
	log.Println("Current working directory:", cwd)

	err = goose.Up(storage.db, "db/migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
}

func (s *Storage) CreateList(name, alias string, uid int) (int, error) {
	const op = "storage.postgres.CreateList"

	var id int

	query := `
		INSERT INTO wishlist (name, uid, alias) VALUES ($1, $2, $3) RETURNING wishlist_id;
		`

	err := s.db.QueryRow(query, name, uid, alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) CreateItem(wishlistId int, giftName, url string) (int, error) {
	const op = "storage.postgres.CreateItem"

	var id int

	query := `
		INSERT INTO items (wishlist_id, name, url) VALUES ($1, $2, $3) RETURNING gift_id;
		`

	err := s.db.QueryRow(query, wishlistId, giftName, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetList(alias string) ([]entity.GiftList, error) {
	const op = "storage.postgres.GetList"

	query := `
	SELECT 
	    items.gift_id, 
	    items.name, 
	    items.url, 
	    wishlist.wishlist_id, 
	    wishlist.name
	FROM items
	JOIN wishlist ON items.wishlist_id = wishlist.wishlist_id
	WHERE wishlist.alias = $1;
	`
	rows, err := s.db.Query(query, alias)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []entity.GiftList

	for rows.Next() {
		var l entity.GiftList

		if err := rows.Scan(&l.GiftId, &l.Name, &l.Url, &l.WishListId, &l.WishListName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (s *Storage) GetAllLists(uid int) ([]entity.WishList, error) {
	const op = "storage.postgres.GetAllLists"

	query := `
	SELECT 
	    wishlist.wishlist_id, 
	    wishlist.name,
		wishlist.uid,
		wishlist.alias
	FROM wishlist
	WHERE wishlist.uid = $1;
	`
	rows, err := s.db.Query(query, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []entity.WishList

	for rows.Next() {
		var l entity.WishList

		if err := rows.Scan(&l.WishListId, &l.Name, &l.UID, &l.Alias); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (s *Storage) WishListDel(wishListId int) error {
	const op = "storage.postgres.WishListDel"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = tx.Exec("DELETE FROM items WHERE wishlist_id = $1", wishListId); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = tx.Exec("DELETE FROM wishlist WHERE wishlist_id = $1", wishListId); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetByWishId(wishListId int) ([]entity.GiftList, error) {
	const op = "storage.postgres.GetByWishId"

	query := `
	SELECT 
	    gift_id,
	    wishlist_id,
	    name,
	    url
	FROM items
	WHERE wishlist_id = $1;
	`
	rows, err := s.db.Query(query, wishListId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []entity.GiftList

	for rows.Next() {
		var l entity.GiftList

		if err := rows.Scan(&l.GiftId, &l.WishListId, &l.Name, &l.Url); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (s *Storage) DelItemById(itemId int) error {
	const op = "storage.postgres.DelItemById"

	_, err := s.db.Exec(`DELETE FROM items WHERE gift_id = $1`, itemId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
