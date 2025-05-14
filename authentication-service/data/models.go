package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const QueryTimeoutDuration = time.Second * 3

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}


type User struct{
	ID int `json:"id"`
	Email string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Password Password `json:"-"`
	Active int `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Password struct {
	text *string
	hash []byte
}

func (pass *Password) SetPassword(text string) error{
	
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	pass.text = &text
	pass.hash = hash

	return nil
}

func (pass *User) PasswordMatches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(pass.Password.hash), []byte(text))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *UserStore) GetUsers(ctx context.Context) ([]*User, error){
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, active, created_at, updated_at
			FROM users ORDER BY last_name
	`
	rows, err := u.db.QueryContext(ctx, query)

	if err != nil{
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Active, &user.CreatedAt, &user.UpdatedAt);
		err != nil{
			return nil, err
		}
		users = append(users, &user)		
	}
	
	return users, nil
} 

func(u *UserStore) GetUserByID(ctx context.Context, UserID int) (*User, error){

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, active, created_at, updated_at
		FROM users WHERE id = $1`

	var user User
	
	err := u.db.QueryRowContext(ctx, query, UserID).
	Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Active, &user.CreatedAt, &user.UpdatedAt)

	if err != nil{
		return nil, err
	}

	
	return &user, nil

}



func(u *UserStore) GetByEmail(ctx context.Context, email string) (*User, error){

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, active, created_at, updated_at
		FROM users WHERE email = $1`

	var user User
	
	err := u.db.QueryRowContext(ctx, query, email).
	Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password.hash, &user.Active, &user.CreatedAt, &user.UpdatedAt)

	if err != nil{
		return nil, err
	}

	
	return &user, nil

}



func(u *UserStore) UpdateUser(ctx context.Context, user *User) error{

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		UPDATE users SET
		email = $1,
		first_name = $2,
		last_name = $3,
		active = $4,
		updated_at = $5
		WHERE id = $6
	`
	
	_, err := u.db.ExecContext(
		ctx, query, user.Email, user.FirstName, user.LastName, user.Active, time.Now(),user.ID)

	
	if err != nil{
		return err
	}
	
	return nil

}


func(u *UserStore) DeleteUser(ctx context.Context, ID int) error{

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		DELETE FROM users WHERE id = $1`
	
	_, err := u.db.ExecContext(ctx, query, ID)

	
	if err != nil{
		return err
	}
	
	return nil
}


func(u *UserStore) CreateUser(ctx context.Context, user User) (int,error){

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		INSERT INTO users (email, first_name, last_name, password, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	var newID int
	err := u.db.QueryRowContext(
		ctx, query, user.Email, user.FirstName, user.LastName, user.Password.hash, user.Active, time.Now(), time.Now()).
	Scan(newID)
	
	
	if err != nil{
		return 0, err
	}
	
	return newID, nil

}



