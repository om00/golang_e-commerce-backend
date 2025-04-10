package Database

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/om00/golang-ecommerce/Models"
)

func (db *DB) GetProuducts(ctx context.Context, filter Models.ProductQuery) ([]Models.Product, error) {

	builder := squirrel.Select("*").
		From("Product")

	if filter.Name != "" {
		builder = builder.Where("MATCH(produtName) AGAINST(?)", filter.Name)
	}
	if filter.Price > 0.0 {
		builder = builder.Where(squirrel.Eq{"price": filter.Price})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.mainDB.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Models.Product

	// Use StructScan to map each row to struct
	for rows.Next() {
		var p Models.Product
		err := rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
