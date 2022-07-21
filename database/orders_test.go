package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
	"webserver/entity"

	"github.com/DATA-DOG/go-sqlmock"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockTvpConverter struct{}

func (converter *mockTvpConverter) ConvertValue(raw interface{}) (driver.Value, error) {

	// Since this function will take the place of every call of ConvertValue, we will inevitably
	// the fake string we return from this function so we need to check whether we've recieved
	// that or a TVP. More extensive logic may be required
	switch inner := raw.(type) {
	case string:
		return raw.(string), nil
	case mssql.TVP:

		// First, verify the type name
		if !strings.EqualFold(inner.TypeName, "ut_OrderItems") {
			return nil, fmt.Errorf("Invalid type")
		}

		// VERIFICATION LOGIC HERE

		// Finally, return a fake value that we can use when verifying the arguments
		return "PASSED", nil
	}

	// We had an invalid type; return an error
	return nil, fmt.Errorf("Invalid type")
}

var Items = []entity.Item{
	{
		ItemCode:    "ITEM_001",
		Description: "Iphone 10X",
		Quantity:    1,
		//OrderId:     1,
	},
	{
		ItemCode:    "ITEM_002",
		Description: "Samsung S21",
		Quantity:    1,
		//OrderId:     1,
	},
	{
		ItemCode:    "ITEM_001",
		Description: "Iphone 10X",
		Quantity:    2,
		//OrderId:     2,
	},
	{
		ItemCode:    "ITEM_002",
		Description: "Samsung S21",
		Quantity:    2,
		//OrderId:     2,
	},
}

var Tvp = mssql.TVP{
	TypeName: "ut_OrderItems",
	Value:    Items,
}

var OrderWithItems = entity.OrderWithItems{
	Order: entity.Order{
		OrderId:      1,
		CustomerName: "Blacky",
		OrderedAt:    time.Now(),
	},
	Items: []entity.Item{
		{
			ItemCode:    "ITEM_001",
			Description: "Iphone 10X",
			Quantity:    1,
			//OrderId:     1,
		},
		{
			ItemCode:    "ITEM_002",
			Description: "Samsung S21",
			Quantity:    1,
			//OrderId:     1,
		},
	},
}

var Order = entity.Order{
	OrderId:      1,
	CustomerName: "Blacky",
	OrderedAt:    time.Now(),
}

var Orders = []entity.Order{
	{
		OrderId:      1,
		CustomerName: "Blacky",
		OrderedAt:    time.Now(),
	},
	{
		OrderId:      2,
		CustomerName: "Bone",
		OrderedAt:    time.Now(),
	},
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestDatabase_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tesName := "Test get orders"
	ctx := context.Background()
	db, mock := NewMock()
	dbtes := Database{
		SqlDb: db,
	}
	t.Run(tesName, func(t *testing.T) {
		rows := mock.NewRows([]string{"OrderId", "CustomerName", "OrderedAt"}).
			AddRow(1, "Blacky", time.Now()).
			AddRow(2, "Bone", time.Now())
		rowsitem := mock.NewRows([]string{"ItemCode", "Description", "Quantity", "OrderId"}).
			AddRow("ITEM_001", "Iphone 10X", "1", "1").
			AddRow("ITEM_002", "Samsung S21", "1", "1").
			AddRow("ITEM_001", "Iphone 10X", "2", "2").
			AddRow("ITEM_002", "Samsung S21", "2", "2")
		mock.ExpectQuery("sp_getOrders").WillReturnRows(rows, rowsitem)
		got, err := dbtes.GetOrders(ctx)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tesName := "Test get order by ID"
	ctx := context.Background()
	db, mock := NewMock()
	dbtes := Database{
		SqlDb: db,
	}
	t.Run(tesName, func(t *testing.T) {
		rows := mock.NewRows([]string{"OrderId", "CustomerName", "OrderedAt"}).
			AddRow(1, "Blacky", time.Now())
		rowsitem := mock.NewRows([]string{"ItemCode", "Description", "Quantity"}).
			AddRow("ITEM_001", "Iphone 10X", "1").
			AddRow("ITEM_002", "Samsung S21", "1")
		mock.ExpectQuery("sp_getOrderByID").
			WithArgs(Order.OrderId).
			WillReturnRows(rows, rowsitem)
		got, err := dbtes.GetOrderByID(ctx, Order.OrderId)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func TestDatabase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tesName := "Test create orders"
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.ValueConverterOption(&mockTvpConverter{}))
	dbtes := Database{
		SqlDb: db,
	}

	t.Run(tesName, func(t *testing.T) {
		mock.ExpectExec("sp_createOrder").
			WithArgs(Order.CustomerName, "PASSED").
			WillReturnResult(sqlmock.NewResult(1, 1))
		got, err := dbtes.CreateOrder(ctx, OrderWithItems)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func TestDatabase_UpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tesName := "Test update orders"
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.ValueConverterOption(&mockTvpConverter{}))
	dbtes := Database{
		SqlDb: db,
	}

	t.Run(tesName, func(t *testing.T) {
		//prep := mock.ExpectPrepare("sp_updateOrder")
		mock.ExpectExec("sp_updateOrder").
			WithArgs(1, Order.CustomerName, "PASSED").
			WillReturnResult(sqlmock.NewResult(1, 1))
		got, err := dbtes.UpdateOrder(ctx, 1, OrderWithItems)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func TestDatabase_DeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tesName := "Test delete order"
	ctx := context.Background()
	db, mock := NewMock()
	dbtes := Database{
		SqlDb: db,
	}

	t.Run(tesName, func(t *testing.T) {
		mock.ExpectExec("sp_deleteOrder").
			WithArgs(Order.OrderId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		got, err := dbtes.DeleteOrder(ctx, Order.OrderId)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}
