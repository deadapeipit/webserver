package service

import (
	"fmt"
	"webserver/entity"
)

type UserServiceIface interface {
	Register(user *entity.User) (users *entity.User, err error)
}

type UserSvc struct {
	User entity.User
}

func NewUserService() UserServiceIface {
	return &UserSvc{}
}

func (u *UserSvc) Register(user *entity.User) (users *entity.User, err error) {
	u.User = *user
	fmt.Println("Data berhasil dimasukan.")
	return user, nil
}
