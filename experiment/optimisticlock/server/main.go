package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lab46/monorepo/gopkg/sql/sqldb"
	"github.com/lab46/monorepo/gopkg/testutil/testenv"
)

var db *sqlx.DB

func init() {
	database, err := sqldb.Connect("postgres", testenv.EnvConfig.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	db = database
}

// Coupon struct
type Coupon struct {
	ID     int64  `db:"id"`
	Code   string `db:"code"`
	Status int    `db:"status"`
}

// CouponQuantity struct
type CouponQuantity struct {
	ID       int64 `db:"id"`
	CouponID int64 `db:"coupon_id"`
	Quantity int   `db:"quantity"`
}

// QuantityHistory struct
type QuantityHistory struct {
	ID               int64 `db:"id"`
	CouponID         int64 `db:"coupon_id"`
	ModifiedQuantity int   `db:"modified_quantity"`
}

func main() {
	http.HandleFunc("/coupon/validate", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		log.Printf("voucher code: %s", code)
		c, err := checkAvailability(code)
		if err != nil {
			w.Write([]byte("voucher is not available"))
			return
		}

		q, err := getCouponQuantity(c.ID)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if q.Quantity > 4999 {
			w.Write([]byte("maximum quantity reached"))
			return
		}

		err = updateQuantity(q.Quantity+1, q.Quantity, q.CouponID, 1)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("validate coupon success"))
	})
	log.Println("setup webserver")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func checkAvailability(code string) (Coupon, error) {
	c := Coupon{}
	query := "SELECT id FROM coupon.coupon WHERE code = $1 AND status <> 0"
	err := db.Get(&c, query, code)
	return c, err
}

func getCouponQuantity(couponID int64) (CouponQuantity, error) {
	cq := CouponQuantity{}
	query := "SELECT id, coupon_id, quantity FROM coupon.coupon_quantity WHERE coupon_id = $1"
	err := db.Get(&cq, query, couponID)
	return cq, err
}

func updateQuantity(newQuantity, oldQuantity int, couponID int64, modifiedQuantity int) error {
	tx, err := db.Beginx()
	defer tx.Rollback()
	updateQuery := "UPDATE coupon.coupon_quantity SET quantity = $1 WHERE coupon_id = $2 AND quantity = $3"
	insertQuery := "INSERT INTO coupon.coupon_history(coupon_id, modified_quantity) VALUES($1, $2)"

	res, err := tx.Exec(updateQuery, newQuantity, couponID, oldQuantity)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("quantity have changed before update")
	}

	_, err = tx.Exec(insertQuery, couponID, modifiedQuantity)
	if err != nil {
		return err
	}
	return tx.Commit()
}
