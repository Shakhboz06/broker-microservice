package data

import "context"

type Repository interface {
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, UserID int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, ID int) error
	CreateUser(ctx context.Context, user User) (int,error)
	SetPassword(text string) error
}