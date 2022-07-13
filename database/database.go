package database

import (
	"context"
	"database/sql"
	"fmt"
	"webserver/entity"

	_ "github.com/denisenkom/go-mssqldb"
)

type DatabaseIface interface {
	CloseConnection()
	Login(ctx context.Context, userName string) (*entity.User, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, userid int) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (string, error)
	UpdateUser(ctx context.Context, userId int, user entity.User) (string, error)
	DeleteUser(ctx context.Context, userId int) (string, error)

	GetOrders(ctx context.Context) ([]entity.OrderWithItems, error)
	GetOrderByID(ctx context.Context, orderid int) (*entity.OrderWithItems, error)
	CreateOrder(ctx context.Context, user entity.OrderWithItems) (string, error)
	UpdateOrder(ctx context.Context, orderId int, order entity.OrderWithItems) (string, error)
	DeleteOrder(ctx context.Context, orderId int) (string, error)
}

type Database struct {
	ConnectionString string
	SqlDb            *sql.DB
}

func NewSqlConnection(connectionString string) DatabaseIface {
	s := Database{
		ConnectionString: connectionString,
	}

	db, err := sql.Open("sqlserver", s.ConnectionString)
	if err != nil {
		fmt.Printf("[mssql] Error connecting to SQL Server: %v", err)
	}

	s.SqlDb = db
	s.SqlDb.SetMaxIdleConns(25)
	s.SqlDb.SetMaxOpenConns(25)

	return &s
}

func (d *Database) CloseConnection() {
	fmt.Println("connection closed!")
	d.SqlDb.Close()
}
