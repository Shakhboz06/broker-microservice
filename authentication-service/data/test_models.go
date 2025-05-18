package data

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PostGresTest struct {
	Conn *sql.DB
}

var testPass = "verysecret"
var testHash, _ = bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)

func NewPostGresTest(db *sql.DB) *PostGresTest {
	return &PostGresTest{
		Conn: db,
	}
}

func (u *PostGresTest) GetUsers(ctx context.Context) ([]*User, error) {
	users := []*User{}

	return users, nil
}

func (u *PostGresTest) GetUserByID(ctx context.Context, UserID int) (*User, error) {

	user := User{
		ID:        1,
		FirstName: "first",
		LastName:  "last",
		Email:     "me@here.come",
		Password:  Password{hash: testHash},
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil

}

func (u *PostGresTest) GetByEmail(ctx context.Context, email string) (*User, error) {


	user := User{
		ID:        1,
		FirstName: "first",
		LastName:  "last",
		Email:     "me@here.come",
		Password:  Password{hash: testHash},
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil

}

func (u *PostGresTest) UpdateUser(ctx context.Context, user *User) error {

	return nil
}

func (u *PostGresTest) DeleteUser(ctx context.Context, ID int) error {

	return nil
}

func (u *PostGresTest) CreateUser(ctx context.Context, user User) (int, error) {

	return 1, nil

}


func (u *PostGresTest) SetPassword(text string) error {
	return nil
}