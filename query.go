package query

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a generics facilitator for executing generic queries
// against the underlying pgx connection object.
type Querier[T any] struct {
	query      string
	conn       *pgx.Conn
	scanValues []pgtype.Value
	mapper     func(a ...pgtype.Value) T
}

// NewQuerier builds a new Querier instance. Provided a pgx.Conn, a
// a query, and a list of []pgtype.Value's
func NewQuerier[T any](c *pgx.Conn, q string, r []pgtype.Value, mapper func(a ...pgtype.Value) T) *Querier[T] {
	return &Querier[T]{
		query:      q,
		conn:       c,
		scanValues: r,
		mapper:     mapper,
	}
}

func (q *Querier[T]) QueryRow(ctx context.Context, args ...interface{}) (T, error) {
	var result T
	row := q.conn.QueryRow(ctx, q.query, args...)

	// propose to pgx a Scan() method that takes in a pgtype.Value
	// variadic array to avoid this.
	scannables := make([]interface{}, len(q.scanValues))
	for i, c := range q.scanValues {
		scannables[i] = c
	}

	err := row.Scan(scannables...)

	if err != nil {
		return result, err
	}

	// This interface upgrading and downgrading is clearly not performant
	// and most definitely not recommended for production usage.
	mappables := make([]pgtype.Value, len(scannables))
	for i := range mappables {
		mappables[i] = scannables[i].(pgtype.Value)

	}
	result = q.mapper(mappables...)

	return result, nil
}

func (q *Querier[T]) Query(ctx context.Context, args ...interface{}) ([]T, error) {
	rows, err := q.conn.Query(ctx, q.query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// propose to pgx a Scan() method that takes in a pgtype.Value
	// variadic array to avoid this.
	scannables := make([]interface{}, len(q.scanValues))
	for i, c := range q.scanValues {
		scannables[i] = c
	}

	var results []T
	for rows.Next() {
		err = rows.Scan(scannables...)
		if err != nil {
			return nil, err
		}

		// This interface upgrading and downgrading is clearly not performant
		// and most definitely not recommended for production usage.
		mappables := make([]pgtype.Value, len(scannables))
		for i := range mappables {
			mappables[i] = scannables[i].(pgtype.Value)

		}

		results = append(results, q.mapper(mappables...))
	}

	if err := rows.Err(); err != nil {
		return results, err
	}

	return results, nil

}
