package users

import (
	gorm "github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// UserGorm is the GORM database connector of the User model
// This connector can be used for sqlite3, mysql and postgresql databases
type UserGorm struct {
	DB *gorm.DB
}


//GetByID queries the database for a User instance by the user ID
// Returns *User and and error
func (o *UserGorm) GetByID(id uint) (*User, error) {
	user := &User{}
	err := o.DB.Where("id = ?", id).First(user).Error

	return user, err
}

//GetByEmail queries the database for a User instance by email
// Returns *User and and error
func (o *UserGorm) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := o.DB.Where("email = ?", email).First(user).Error

	return user, err
}

//Get queries the database via the provided query
// Returns *User and and error
func (o *UserGorm) Get(query interface{}) (*User, error) {
	user := &User{}
	err := o.DB.Where(query).First(user).Error

	return user, err
}

//Find queries the database via the provided query
// Returns a slice of User and and error
func (o *UserGorm) Find(query interface{}) (*[]Users, error) {
	users := &[]Users{}
	err := o.DB.Where(query).Find(users).Error

	return users, err
}

func (o *UserGorm) compareHashAndPassword(password1 string, password2 string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if err != nil {
		return err
	}

	return nil
}

// Authenticate provides functionality to authenticate a user with the provided password
// Authenticate first queries for a user based on the value of authBy, valid search criteria are AuthByUsername, AuthByEmail, AuthByUsernameAndEmail
// After a successful user query the given password will compared to the password hash in the User record.
// Returns *User and error
func (o *UserGorm) Authenticate(user string, password string, authBy uint) (*User, error) {
	u := &User{}
	var err error

	if authBy == AuthByUsername {
		err = o.DB.Where("username = ?", user).First(u).Error
	} else if authBy == AuthByEmail {
		err = o.DB.Where("email = ?", user).First(u).Error
	} else if authBy == AuthByUsernameOrEmail {
		err = o.DB.Where("email = ? OR username = ?", user, user).First(u).Error
	} else {
		err = errors.New("Invalid authBy value")
	}

	if err != nil {
		return u, err
	}

	err = o.compareHashAndPassword(u.Password, password)

	return u, nil
}

// NewPageManagerGorm returns a database connector for Gorm databases
func NewPageManagerGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{DB: db}
}
