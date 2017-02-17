package users

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
)

const (
	AuthByUsername        uint = 0
	AuthByEmail           uint = 1
	AuthByUsernameOrEmail uint = 2
)

// User implements a user model
type User struct {
	GID uint          `gorm:"column:id; primary_key" bson:"-"`
	MID bson.ObjectId `bson:"_id" gorm:"-"`

	Username  string
	Email     string
	Password  string
	ChangePassword bool
	IsAdmin bool
	Profile Profile
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u *User) ID() string {
	if u.GID != 0 {
		return strconv.FormatUint(uint64(u.GID), 10)
	} else if u.MID.Valid() {
		return u.MID.String()
	}

	return ""
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
	Init(drop bool) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Get(query interface{}, values ...interface{}) (*User, error)
	Find(query interface{}, values ...interface{}) (*[]User, error)
	Authenticate(user string, password string, authBy uint) (*User, error)
	Create(*User) error
	Update(*User) error
	Delete(*User) error
}

// Users implements the controller functions.
type Users struct {
	UserManager
}

// UserController takes a UserManager and returns a Users controller struct
func UserController(m UserManager) *Users {
	return &Users{UserManager: m}
}
