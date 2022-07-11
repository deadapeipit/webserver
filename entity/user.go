package entity

import "time"

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	Id                    int          `json:"id"`
	Uid                   string       `json:"uid"`
	Password              string       `json:"password"`
	FirstName             string       `json:"first_name"`
	LastName              string       `json:"last_name"`
	Username              string       `json:"user_name"`
	Email                 string       `json:"email"`
	Avatar                string       `json:"avatar"`
	Gender                string       `json:"gender"`
	PhoneNumber           string       `json:"phone_number"`
	SocialInsuranceNumber string       `json:"social_insurance_number"`
	DateOfBirth           string       `json:"date_of_birth"`
	Employment            Employment   `json:"employment"`
	Address               Address      `json:"address"`
	CreditCard            CreditCard   `json:"credit_card"`
	Subscription          Subscription `json:"subscription"`
}

type UserResponse struct {
	Id        int     `json:"id"`
	Uid       string  `json:"uid"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"user_name"`
	Address   Address `json:"address"`
}

func (u *UserRequest) ToUserResponse() *UserResponse {
	retval := UserResponse{
		Id:        u.Id,
		Uid:       u.Uid,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Address:   u.Address,
	}
	return &retval
}

func ToArrayUserResponse(i []UserRequest) *[]UserResponse {
	var retval []UserResponse
	for _, row := range i {
		tempUser := *row.ToUserResponse()
		retval = append(retval, tempUser)
	}
	return &retval
}
