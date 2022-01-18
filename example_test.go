package query

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
)

func TestQueryRow(t *testing.T) {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://localhost:5432/hingedb?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	query := `select product_no, product_name, price from products where product_no = $1`

	q := NewQuerier(conn, query, ProductColumns, ProductMapper)

	product, err := q.QueryRow(context.Background(), 0)
	if err != nil {
		t.Error(err)
	}

	t.Log(product)

	query = `select purchase_no, product_no, price_paid, customer_id, purchased_at from purchases where purchase_no = $1`
	q1 := NewQuerier(conn, query, PurchaseColumns, PurchaseMapper)

	purchase, err := q1.QueryRow(context.Background(), 0)
	if err != nil {
		t.Error(err)
	}

	t.Log(purchase)

}

func TestQuery(t *testing.T) {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://localhost:5432/hingedb?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	query := `select product_no, product_name, price from products where price = $1`

	q := NewQuerier(conn, query, ProductColumns, ProductMapper)

	p, err := q.Query(context.Background(), 42069)
	if err != nil {
		t.Error(err)
	}
	t.Log(p)

}
