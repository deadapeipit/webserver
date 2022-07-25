package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"webserver/entity"

	mssql "github.com/denisenkom/go-mssqldb"
)

//Orders using stored procedure, script at database_script/db_script.txt

func (s *Database) GetOrders(ctx context.Context) ([]entity.OrderWithItems, error) {
	var result []entity.OrderWithItems
	var orders []entity.Order
	var items []entity.ItemWithOrderID

	rows, err := s.SqlDb.QueryContext(ctx, "sp_getOrders")
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.Order
		err := rows.Scan(
			&row.OrderId,
			&row.CustomerName,
			&row.OrderedAt,
		)
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}
		orders = append(orders, row)
	}
	if !rows.NextResultSet() {
		log.Fatal("[mssql] Expected more resultset")
		return nil, errors.New("[mssql] Expected more resultset")
	}
	for rows.Next() {
		var item entity.ItemWithOrderID
		err := rows.Scan(
			&item.ItemCode,
			&item.Description,
			&item.Quantity,
			&item.OrderId,
		)
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}
		items = append(items, item)
	}
	for _, o := range orders {
		var tempOrder entity.OrderWithItems
		var tempItems []entity.Item
		tempOrder.Order = o
		for _, i := range items {
			if tempOrder.OrderId == i.OrderId {
				tempItems = append(tempItems, *i.ToWithoutOrderID())
			}
		}
		tempOrder.Items = tempItems
		result = append(result, tempOrder)
	}

	return result, nil
}

func (s *Database) GetOrderByID(ctx context.Context, i int) (*entity.OrderWithItems, error) {
	result := &entity.OrderWithItems{}
	order := &entity.Order{}
	var items []entity.Item

	rows, err := s.SqlDb.QueryContext(ctx, "sp_getOrderByID",
		sql.Named("pOrderId", i))
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&order.OrderId,
			&order.CustomerName,
			&order.OrderedAt,
		)
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}
	}
	if !rows.NextResultSet() {
		//log.Fatal("[mssql] Expected more resultset")
		return nil, errors.New("[mssql] Expected more resultset")
	}
	if rows.Next() {
		var item entity.Item
		err := rows.Scan(
			&item.ItemCode,
			&item.Description,
			&item.Quantity,
		)
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}
		items = append(items, item)
	}

	result.Order = *order
	result.Items = items

	return result, nil
}

func (s *Database) CreateOrder(ctx context.Context, i entity.OrderWithItems) (string, error) {
	var result string

	tvp := mssql.TVP{
		TypeName: "ut_OrderItems",
		Value:    i.Items,
	}

	_, err := s.SqlDb.ExecContext(ctx, "sp_createOrder",
		sql.Named("pCustomerName", i.CustomerName),
		sql.Named("pItems", tvp))
	if err != nil {
		//log.Fatal(err)
		return "", err
	}

	result = "Inserted"

	return result, nil
}

func (s *Database) UpdateOrder(ctx context.Context, id int, i entity.OrderWithItems) (string, error) {
	var result string

	tvp := mssql.TVP{
		TypeName: "ut_OrderItems",
		Value:    i.Items,
	}
	_, err := s.SqlDb.ExecContext(ctx, "sp_updateOrder",
		sql.Named("pOrderId", id),
		sql.Named("pCustomerName", i.CustomerName),
		sql.Named("pItems", tvp))
	if err != nil {
		//log.Fatal(err)
		return "", err
	}

	result = "Updated"

	return result, nil
}

func (s *Database) DeleteOrder(ctx context.Context, i int) (string, error) {
	var result string

	_, err := s.SqlDb.ExecContext(ctx, "sp_deleteOrder",
		sql.Named("pOrderId", i))
	if err != nil {
		//log.Fatal(err)
		return "", err
	}

	result = "Deleted"

	return result, nil
}
