package users

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	AuthByUsername uint = 0
	AuthByEmail uint = 1
	AuthByUsernameOrEmail uint = 2
)

// User implements a user model
type User struct {
	ID uint
	Username string
	Email  string
	Password string
}

//SetPassword hashes the given password and stores it into the User instance
func (u *User) SetPassword(password string) error {
	p := []byte(password)

	h, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(h)
	return nil
}

// UserManager is an interface for database connectors which provide the interaction between the model and the database
type UserManager interface {
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Get(query interface{}) (*User, error)
	//Find(query interface{}) (*Users, error)
	Authenticate(user string, password string, authBy uint) (*User, error)
}

// Users implements the controller functions.
type Users struct {
	UserManager
}

// UserController takes a UserManager and returns a Users controller struct
func UserController(m UserManager) *Users {
	return &Users{UserManager: m}
}
