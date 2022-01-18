# go-generic-pgx

Playing around with Go 1.18's new Generics feature. (... requires go1.18+ to run )

Given these tables:
```sql
create table products(product_no integer, product_name text, price integer)
create table purchases(purchase_no integer, product_no integer, price_paid integer, customer_id text)
```

This works by building a generic querier object. This leverages the faciliator pattern to support that Golang does not have [parameterized methods](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods).

First, we can define the columns we are scanning with an ordered array of `pgtype.Value` interface values.
```go
var ProductColumns = []pgtype.Value{&pgtype.Int4{}, &pgtype.Text{}, &pgtype.Int4{}}
```

With this, we can build a callback function that operates on this array to map to a specific struct.
```go
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
```

Not that this is a huge improvement in typing code, but one thing to note is that this is much easier to test than the current interface provided by `database/sql` package.

With the columns we are scanning defined, and a mapping callback to map the columns into a Go struct, we can now execute a query.

```go
// instantiate a new querier
query := `select product_no, product_name, price from products where product_no = $1`
q := NewQuerier(conn, query, ProductColumns, ProductMapper)

// now we can use this querier to get our product directly back
product, err := q.QueryRow(context.Background(), 0)

// Prints {0 a 100}
fmt.Println(product)

// Now lets reuse that logic for our purchase query
query = `select purchase_no, product_no, price_paid, customer_id, purchased_at from purchases where purchase_no = $1`
q1 := NewQuerier(conn, query, PurchaseColumns, PurchaseMapper)

purchase, err := q1.QueryRow(context.Background(), 0)

// Prints {0 1 1000 alice 2022-01-18 07:34:28.796139 -0700 MST}
fmt.Println(product)
```
