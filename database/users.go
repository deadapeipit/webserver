package database

import (
	"context"
	"database/sql"
	"log"
	"time"
	"webserver/entity"
)

func (s *Database) GetUsers(ctx context.Context) ([]entity.User, error) {
	var result []entity.User

	rows, err := s.SqlDb.QueryContext(ctx, "select id, username, email, password, age, createdat, updatedat from users")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.User
		err := rows.Scan(
			&row.Id,
			&row.Username,
			&row.Email,
			&row.Password,
			&row.Age,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s *Database) GetUserByID(ctx context.Context, i int) (*entity.User, error) {
	result := &entity.User{}

	rows, err := s.SqlDb.QueryContext(ctx, "select id, username, email, password, age, createdat, updatedat from users where id = @ID",
		sql.Named("ID", i))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.Id,
			&result.Username,
			&result.Email,
			&result.Password,
			&result.Age,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) CreateUser(ctx context.Context, i entity.User) (string, error) {
	var result string
	now := time.Now()
	_, err := s.SqlDb.ExecContext(ctx, "insert into users (id, username, email, password, age, createdat, updatedat) values (@id, @username, @email, @password, @age, @createdat, @updatedat)",
		sql.Named("id", i.Id),
		sql.Named("username", i.Username),
		sql.Named("email", i.Email),
		sql.Named("password", i.Password),
		sql.Named("age", i.Age),
		sql.Named("createdat", now),
		sql.Named("updatedat", now))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result = "Inserted"

	return result, nil
}

func (s *Database) UpdateUser(ctx context.Context, id int, i entity.User) (string, error) {
	var result string

	_, err := s.SqlDb.ExecContext(ctx, "update users set username = @username,email = @email, password = @password, age = @age, updatedat = @updatedat where id = @id",
		sql.Named("id", id),
		sql.Named("username", i.Username),
		sql.Named("email", i.Email),
		sql.Named("password", i.Password),
		sql.Named("age", i.Age),
		sql.Named("updatedat", time.Now()))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result = "Updated"

	return result, nil
}

func (s *Database) DeleteUser(ctx context.Context, i int) (string, error) {
	var result string

	_, err := s.SqlDb.ExecContext(ctx, "delete from users where id=@id",
		sql.Named("id", i))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result = "Deleted"

	return result, nil
}
