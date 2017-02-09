package users

import (
	"github.com/BalkanTech/goilerplate/config"
	db "github.com/BalkanTech/goilerplate/databases"
	"testing"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

const testOK = "\u2714"
const testFailed = "\u2718"

var GormConfig = &config.Config{
		Database: config.Database{
		Type: "sqlite3",
		DB:   "/tmp/test.db",
	},
}

var MgoConfig = &config.Config{
	Database: config.Database{
		Type: "mongodb",
		Host: "localhost",
		DB:   "test",
		User: "test",
		Password: "test",
	},
}

func getTestUser() *User {
	user := &User{GID: 1, MID: bson.ObjectIdHex("588c89a388c8f1120a0828f6"), Username: "test", Email: "test@example.com"}
	user.SetPassword("test")

	return user
}

func getUserControllerGorm() (*Users, *gorm.DB, error) {
	dbase, err := db.NewGormConnection(GormConfig)
	if err != nil {
		return nil, nil, err
	}

	m := NewUserManagerGorm(dbase)

	return UserController(m),dbase, nil
}

func getUserControllerMgo() (*Users, *mgo.Session, error) {
	dbase, err := db.NewMgoConnection(MgoConfig)
	if err != nil {
		return nil, nil, err
	}

	m := NewUserManagerMgo(dbase, "users")

	return UserController(m), dbase.Session, nil
}

// TestInit will test the initialization processes with dropping the table
func TestInit(t *testing.T) {
	t.Log("Testing Init for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()

		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Init(true)
		if err != nil {
			t.Fatalf("\tExpected to be able to drop an initialize the table, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to drop and initialize the table. ", testOK)
	}

	t.Log("Testing Init for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()

		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Init(true)
		if err != nil {
			t.Fatalf("\tExpected to be able to drop an initialize the table, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to drop and initialize the table. ", testOK)
	}
}


// TestCreate tries to create a new user in the database
func TestCreate(t *testing.T) {
	testUser := getTestUser()

	t.Log("Testing Create for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Create(testUser)
		if err != nil {
			t.Fatalf("\tExpected to be able to create, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to create. ", testOK)
	}

	t.Log("Testing Create for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		t.Log(testUser)
		err = users.Create(testUser)
		if err != nil {
			t.Fatalf("\tExpected to be able to create, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to create. ", testOK)
	}
}

// TestGetByID tries to get a user from the database by ID
func TestGetByID(t *testing.T) {
	user1 := getTestUser()

	t.Log("Testing GetByID for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		id := strconv.FormatUint(uint64(user1.GID), 10)

		user2, err := users.GetByID(id)
		if err != nil {
			t.Fatalf("\tExpected to be able to read by ID, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read by ID. ", testOK)

		if user2.GID != user1.GID || user2.Username != user1.Username || user2.Email != user1.Email {
			t.Fatalf("\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\""+
			", but got ID: %d, Username: \"%s\", Email: \"%s\" instead. %v",
				user1.GID, user1.Username, user1.Email,
				user2.GID, user2.Username, user2.Email, testFailed)
		}
		t.Logf("\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\". %v",
			user1.GID, user1.Username, user1.Email, testOK)
	}

	t.Log("Testing GetByID for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		id := user1.MID.Hex()
		t.Log(id)

		user2, err := users.GetByID(id)
		if err != nil {
			t.Fatalf("\tExpected to be able to read by ID, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read by ID. ", testOK)

		if user2.MID.Hex() != user1.MID.Hex() || user2.Username != user1.Username || user2.Email != user1.Email {
			t.Fatalf("\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\""+
				", but got ID: %s, Username: \"%s\", Email: \"%s\" instead. %v",
				user1.MID.Hex(), user1.Username, user1.Email,
				user2.MID.Hex(), user2.Username, user2.Email, testFailed)
		}
		t.Logf("\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\". %v",
			user1.MID.Hex(), user1.Username, user1.Email, testOK)
	}
}

// TestGetByID tries to get a user from the database by email
func TestGetByEmail(t *testing.T) {
	user1 := getTestUser()

	t.Log("Testing GetByEmail for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		user2, err := users.GetByEmail(user1.Email)
		if err != nil {
			t.Fatalf("\tExpected to be able to read by Email, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read by Email. ", testOK)

		if user2.GID != user1.GID || user2.Username != user1.Username || user2.Email != user1.Email {
			t.Fatalf("\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\""+
			", but got ID: %d, Username: \"%s\", Email: \"%s\" instead. %v",
					user1.GID, user1.Username, user1.Email,
					user2.GID, user2.Username, user2.Email, testFailed)
			}
			t.Logf("\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\". %v",
				user1.GID, user1.Username, user1.Email, testOK)
	}

	t.Log("Testing GetByEmail for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		user2, err := users.GetByEmail(user1.Email)
		if err != nil {
			t.Fatalf("\tExpected to be able to read by Email, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read by Email. ", testOK)

		if user2.MID.Hex() != user1.MID.Hex() || user2.Username != user1.Username || user2.Email != user1.Email {
			t.Fatalf("\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\""+
				", but got ID: %s, Username: \"%s\", Email: \"%s\" instead. %v",
				user1.MID.Hex(), user1.Username, user1.Email,
				user2.MID.Hex(), user2.Username, user2.Email, testFailed)
		}
		t.Logf("\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\". %v",
			user1.MID.Hex(), user1.Username, user1.Email, testOK)
	}
}

// TestGet tries to get a user from the database by a custom query
func TestGet(t *testing.T) {
	user1 := getTestUser()

	t.Log("Testing Get for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		user2, err := users.Get("email = ? AND username = ?", user1.Email, user1.Username)
		if err != nil {
			t.Fatalf("\tExpected to be able to read by custom query, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read by custom query. ", testOK)

		if user2.GID != user1.GID || user2.Username != user1.Username || user2.Email != user1.Email {
			t.Fatalf("\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\""+", but got ID: %d, Username: \"%s\", Email: \"%s\" instead. %v",
				user1.GID, user1.Username, user1.Email,
				user2.GID, user2.Username, user2.Email, testFailed)
		}
		t.Logf("\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\". %v",
			user1.GID, user1.Username, user1.Email, testOK)
	}

	t.Log("Testing Get for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		user2, err := users.Get(bson.M{"email": user1.Email, "username": user1.Username})
		if err != nil {
			t.Fatalf("\tExpected to be able to read by custom query, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read by custom query. ", testOK)

		if user2.MID.Hex() != user1.MID.Hex() || user2.Username != user1.Username || user2.Email != user1.Email {
			t.Fatalf("\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\""+", but got ID: %s, Username: \"%s\", Email: \"%s\" instead. %v",
				user1.MID.Hex(), user1.Username, user1.Email,
				user2.MID.Hex(), user2.Username, user2.Email, testFailed)
		}
		t.Logf("\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\". %v",
			user1.MID.Hex(), user1.Username, user1.Email, testOK)
	}
}

// TestFind tries to get a user array from the database by a custom query
func TestFind(t *testing.T) {
	user1 := getTestUser()

	t.Log("Testing Find for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		usr, err := users.Find("id = ?", user1.GID)
		if err != nil {
			t.Fatalf("\tExpected to be able to read many by custom query, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read many by custom query. ", testOK)

		if len(*usr) != 1 {
			t.Fatalf("\t\tExpected array length to be 1, but got %d instead. %v", len(*usr), testFailed)
		}
		t.Log("\t\tExpected array length to be 1. ", testOK)
	}

	t.Log("Testing Find for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		usr, err := users.Find(bson.M{"_id": user1.MID})
		if err != nil {
			t.Fatalf("\tExpected to be able to read many by custom query, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to read many by custom query. ", testOK)

		if len(*usr) != 1 {
			t.Fatalf("\t\tExpected array length to be 1, but got %d instead. %v", len(*usr), testFailed)
		}
		t.Log("\t\tExpected array length to be 1. ", testOK)
	}
}


// TestAuthentication tries to perform authentication by Username, Email, Username or email
func TestAuthentication(t *testing.T) {
	user1 := getTestUser()

	tests := []struct{
		Arg string
		Method uint
		Description string
	}{
		{"test", AuthByUsername, "username"},
		{"test@example.com", AuthByEmail, "email"},
		{"test", AuthByUsernameOrEmail, "username or email"},
		{"test@example.com", AuthByUsernameOrEmail, "username or email"},
	}

	t.Log("Testing Authentication for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		for _,test := range tests {
			t.Logf("\tAuthenticate by %s with \"%s\"", test.Description, test.Arg)
			{
				user2, err := users.Authenticate(test.Arg, "test", test.Method)
				if err != nil {
					t.Fatalf("\t\tExpected to be able to authenticate, but got error: %s. %v", err, testFailed)
				}
				t.Log("\t\tExpected to be able to authenticate. ", testOK)

				if user2.GID != user1.GID || user2.Username != user1.Username || user2.Email != user1.Email {
					t.Fatalf("\t\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\""+
						", but got ID: %d, Username: \"%s\", Email: \"%s\" instead. %v",
						user1.GID, user1.Username, user1.Email,
						user2.GID, user2.Username, user2.Email, testFailed)
				}
				t.Logf("\t\t\tExpected ID: %d, Username: \"%s\", Email: \"%s\". %v",
					user1.GID, user1.Username, user1.Email, testOK)
			}
		}
	}

	t.Log("Testing Authentication for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		for _,test := range tests {
			t.Logf("\tAuthenticate by %s with \"%s\"", test.Description, test.Arg)
			{
				user2, err := users.Authenticate(test.Arg, "test", test.Method)
				if err != nil {
					t.Fatalf("\t\tExpected to be able to authenticate, but got error: %s. %v", err, testFailed)
				}
				t.Log("\t\tExpected to be able to authenticate. ", testOK)

				if user2.MID.Hex() != user1.MID.Hex() || user2.Username != user1.Username || user2.Email != user1.Email {
					t.Fatalf("\t\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\""+
						", but got ID: %s, Username: \"%s\", Email: \"%s\" instead. %v",
						user1.MID.Hex(), user1.Username, user1.Email,
						user2.MID.Hex(), user2.Username, user2.Email, testFailed)
				}
				t.Logf("\t\t\tExpected ID: %s, Username: \"%s\", Email: \"%s\". %v",
					user1.MID.Hex(), user1.Username, user1.Email, testOK)
			}
		}
	}
}


// TestUpdate tries to save an updated user to the database
func TestUpdate(t *testing.T) {
	user1 := getTestUser()
	user1.Email = "testuser@example.com"

	t.Log("Testing Update for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Update(user1)
		if err != nil {
			t.Fatalf("\tExpected to be able to update, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to update. ", testOK)
	}

	t.Log("Testing Update for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Update(user1)
		if err != nil {
			t.Fatalf("\tExpected to be able to update, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to update. ", testOK)
	}
}

// TestDelete tries to delete the test user in the database
func TestDelete(t *testing.T) {
	user1 := getTestUser()

	t.Log("Testing Delete for Gorm")
	{
		users, d, err := getUserControllerGorm()
		defer d.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Delete(user1)
		if err != nil {
			t.Fatalf("\tExpected to be able to delete, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to delete. ", testOK)
	}

	t.Log("Testing Delete for MGO")
	{
		users, s, err := getUserControllerMgo()
		defer s.Close()
		if err != nil {
			t.Fatalf("\tExpected to get a user controller, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to get a user controller", testOK)

		err = users.Delete(user1)
		if err != nil {
			t.Fatalf("\tExpected to be able to delete, but got error: %s. %v", err, testFailed)
		}
		t.Log("\tExpected to be able to delete. ", testOK)
	}
}