package users

import (
	"testing"
	"reflect"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

const testOK = "\u2714"
const testFailed = "\u2718"

var users Users

// TestGormUserManagerAndController
func TestGormUserManagerAndController(t *testing.T) {
	var dbase *gorm.DB

	t.Log("While testing with Gorm")
	{
		m := NewUserManagerGorm(dbase)
		expected_type := reflect.TypeOf(UserGorm{})
		if reflect.TypeOf(m) == expected_type {
			t.Fatalf("\tShould return a type: %s, but got: %s. %v", expected_type, reflect.TypeOf(m), testFailed)
		}
		t.Logf("\tShould return a type: %s. %v", expected_type, testOK)

		users := UserController(m)
		t.Logf("%+v", *users)
		/*u, err := users.Get(0)
		if err != nil {
			t.Fatalf("\tShould get a user. %s. %v", err, testFailed)
		}
		t.Log("\tShould be able to get a user. ", testOK)

		var expectedID uint = 0
		var expectedUsername string = "User0"
		if u.ID != expectedID || u.Username != expectedUsername {
			t.Fatalf("\tExpected a user with ID %d and Username \"%s\", got ID %d and Username %s instead. %v", expectedID, expectedUsername, u.ID, u.Username, testFailed)
		}
		t.Logf("\tExpected a user with ID %d and Username \"%s\". %v", expectedID, expectedUsername, testOK)*/
	}
}

// TestGormUserManagerAndController
func TestMgoUserManagerAndController(t *testing.T) {
	var dbase *mgo.Database
	t.Log("While testing with MGO")
	{
		m := NewUserManagerMgo(dbase)
		expected_type := reflect.TypeOf(UserMgo{})
		if reflect.TypeOf(m) == expected_type {
			t.Fatalf("\tShould return a type: %s, but got: %s. %v", expected_type, reflect.TypeOf(m), testFailed)
		}
		t.Logf("\tShould return a type: %s. %v", expected_type, testOK)

		users := UserController(m)
		t.Logf("%+v", *users)
		/*u, err := users.Get(0)
		if err != nil {
			t.Fatalf("\tShould get a user. %s. %v", err, testFailed)
		}
		t.Log("\tShould be able to get a user. ", testOK)

		var expectedID uint = 0
		var expectedUsername string = "User0"
		if u.ID != expectedID || u.Username != expectedUsername {
			t.Fatalf("\tExpected a user with ID %d and Username \"%s\", got ID %d and Username %s instead. %v", expectedID, expectedUsername, u.ID, u.Username, testFailed)
		}
		t.Logf("\tExpected a user with ID %d and Username \"%s\". %v", expectedID, expectedUsername, testOK)*/
	}
}


