package query

import (
	"time"

	"github.com/jackc/pgtype"
)

type Product struct {
	Number int
	Name   string
	Price  int
}

var ProductColumns = []pgtype.Value{&pgtype.Int4{}, &pgtype.Text{}, &pgtype.Int4{}}

func ProductMapper(a ...pgtype.Value) Product {
	p := Product{}

	num, _ := a[0].(*pgtype.Int4)
	p.Number = int(num.Int)

	str, _ := a[1].(*pgtype.Text)
	p.Name = str.String

	price, _ := a[2].(*pgtype.Int4)
	p.Price = int(price.Int)

	return p
}

type Purchase struct {
	Number        int
	ProductNumber int
	PricePaid     int
	CustomerID    string
	PurchasedAt   time.Time
}

var PurchaseColumns = []pgtype.Value{&pgtype.Int4{}, &pgtype.Int4{}, &pgtype.Int4{}, &pgtype.Text{}, &pgtype.Timestamptz{}}

func PurchaseMapper(a ...pgtype.Value) Purchase {
	p := Purchase{}

	num, _ := a[0].(*pgtype.Int4)
	p.Number = int(num.Int)

	pNum, _ := a[1].(*pgtype.Int4)
	p.ProductNumber = int(pNum.Int)

	pricePaid, _ := a[2].(*pgtype.Int4)
	p.PricePaid = int(pricePaid.Int)

	customer := a[3].(*pgtype.Text)
	p.CustomerID = customer.String

	purchasedAt := a[4].(*pgtype.Timestamptz)
	p.PurchasedAt = purchasedAt.Time

	return p
}
