package users

import (
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"log"
)

// UserGorm is the GORM database connector of the User model
// This connector can be used for sqlite3, mysql and postgresql databases
type UserGorm struct {
	DB *gorm.DB
}

// Init performs Gorm's automigrate function to initialize the table and drops the table if drop is set to true
func (o *UserGorm) Init(drop bool) error {
	if drop {
		err := o.DB.DropTable(&User{}).Error
		if err != nil {
			return err
		}
	}
	err := o.DB.AutoMigrate(&User{}, &Profile{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetByID queries the database for a User instance by the user ID
// Returns *User and and error
func (o *UserGorm) GetByID(id string) (*User, error) {
	user := &User{}

	if id == "" {
		return user, errors.New("Empty ID")
	}

	gid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return user, err
	}

	err = o.DB.Where("id = ?", gid).First(user).Error

	o.DB.Model(user).Related(&user.Profile)

	return user, err
}

//GetByEmail queries the database for a User instance by email
// Returns *User and and error
func (o *UserGorm) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := o.DB.Where("email = ?", email).First(user).Error

	o.DB.Model(user).Related(&user.Profile)

	return user, err
}

//Get queries the database via the provided query
// Returns *User and and error
func (o *UserGorm) Get(query interface{}, values ...interface{}) (*User, error) {
	user := &User{}
	err := o.DB.Where(query, values...).First(user).Error

	o.DB.Model(user).Related(&user.Profile)

	return user, err
}

// Find queries the database via the provided query
// Returns a slice of User and and error
func (o *UserGorm) Find(query interface{}, values ...interface{}) (*[]User, error) {
	users := &[]User{}
	err := o.DB.Where(query, values...).Find(users).Error

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

	o.DB.Model(u).Related(&u.Profile)

	if err != nil {
		return nil, err
	}

	err = o.compareHashAndPassword(u.Password, password)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Create will create a new record in the database for the provided user or return an error
func (o *UserGorm) Create(u *User) error {
	err := o.DB.Create(u).Error
	if err != nil {
		return err
	}

	p := &Profile{UserID: u.GID}
	err = o.DB.Create(p).Error
	if err != nil {
		return err
	}

	return nil
}

// Update will update the provided record in the database or return an error
func (o *UserGorm) Update(u *User) error {
	//u.UpdatedAt = time.Now()
	err := o.DB.Save(u).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete will delete the provided record in the database or return an error
func (o *UserGorm) Delete(u *User) error {
	err := o.DB.Delete(u).Error
	if err != nil {
		return err
	}
	return nil
}

// NewPageManagerGorm returns a database connector for Gorm databases
func NewUserManagerGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{DB: db}
}