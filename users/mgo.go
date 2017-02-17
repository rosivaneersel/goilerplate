package users

import (
	db "github.com/BalkanTech/goilerplate/databases"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"errors"
)

// UserMgo is the MGO database connector of the User model
// This connector can be used for mongodb
type UserMgo struct {
	DB *db.MongoDBConnection
	C  *mgo.Collection
}

// Init drops the collection if drop is set to true
func (o *UserMgo) Init(drop bool) error {
	if drop {
		err := o.C.DropCollection()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetByID queries the database for a User instance by the user ID
// Returns *User and and error
func (o *UserMgo) GetByID(id string) (*User, error) {
	user := &User{}
	err := o.C.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)

	return user, err
}

// GetByEmail queries the database for a User instance by email
// Returns *User and and error
func (o *UserMgo) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := o.C.Find(bson.M{"email": email}).One(&user)

	return user, err
}

// Get queries the database via the provided query
// Returns *User and and error
func (o *UserMgo) Get(query interface{}, values ...interface{}) (*User, error) {
	user := &User{}
	err := o.C.Find(query).One(user)

	return user, err
}

// Find queries the database via the provided query
// Returns a slice of User and and error
func (o *UserMgo) Find(query interface{}, values ...interface{}) (*[]User, error) {
	users := &[]User{}
	err := o.C.Find(query).All(users)

	return users, err
}

func (o *UserMgo) compareHashAndPassword(password1 string, password2 string) error {
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
func (o *UserMgo) Authenticate(user string, password string, authBy uint) (*User, error) {
	u := &User{}
	var err error

	if authBy == AuthByUsername {
		err = o.C.Find(bson.M{"username": user}).One(u)
	} else if authBy == AuthByEmail {
		err = o.C.Find(bson.M{"email": user}).One(u)
	} else if authBy == AuthByUsernameOrEmail {
		err = o.C.Find(bson.M{"$or": []interface{}{bson.M{"username": user}, bson.M{"email": user}}}).One(u)
	} else {
		err = errors.New("Invalid authBy value")
	}

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
func (o *UserMgo) Create(u *User) error {
	time := time.Now()

	u.CreatedAt = time
	u.UpdatedAt = time
	if !u.MID.Valid() {
		u.MID = bson.NewObjectId()
	}
	err := o.C.Insert(u)
	if err != nil {
		return err
	}

	return nil
}

// Update will update the provided record in the database or return an error
func (o *UserMgo) Update(u *User) error {
	u.UpdatedAt = time.Now()
	err := o.C.Update(bson.M{"_id": u.MID}, u)
	if err != nil {
		return err
	}

	return nil
}

// Delete will delete the provided record in the database or return an error
func (o *UserMgo) Delete(u *User) error {
	err := o.C.Remove(bson.M{"_id": u.MID})
	if err != nil {
		return err
	}
	return nil
}


// NewPageManagerMgo returns a database connector for Mgo
func NewUserManagerMgo(db *db.MongoDBConnection, collection string) *UserMgo {
	return &UserMgo{DB: db, C: db.DB.C(collection)}
}