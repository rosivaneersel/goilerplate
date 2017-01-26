package databases_test

import (
	"testing"
	"github.com/BalkanTech/goilerplate/config"
	db "github.com/BalkanTech/goilerplate/databases"
	"reflect"
	"github.com/jinzhu/gorm"
)

const testOK = "\u2714"
const testFailed = "\u2718"

func TestNewGormConnector(t *testing.T) {
	Tests := []config.Config{
		config.Config{
			Database: config.Database{
				Type: "sqlite3",
				DB:   "/tmp/test.db",
			},
		},
		config.Config{
			Database: config.Database{
				Type:     "postgres",
				Host:     "localhost",
				User:     "test",
				Password: "test",
				DB:       "test",
			},
		},
		config.Config{
			Database: config.Database{
				Type:     "mysql",
				Host:     "",
				User:     "test",
				Password: "test",
				DB:       "test",
			},
		},
	}
	expectedType := reflect.TypeOf(&gorm.DB{})

	t.Log("Testing NewGormConnector")
	{
		for _, test := range Tests {
			t.Logf("\t%s", test.Database.GetType())
			{
				dbase, err := db.NewGormConnection(&test)
				if err != nil {
					t.Fatalf("\t\tExpected to get a database instance, but got an error: %s. %v", err, testFailed)
				}
				t.Log("\t\tExpected to get a database instance. ", testOK)
				defer dbase.Close()

				receivedType := reflect.TypeOf(dbase)
				if receivedType != expectedType {
					t.Fatalf("\t\tExpeted a \"%s\" type, but got \"%s\" instead. %v", expectedType, receivedType, testFailed)
				}
				t.Logf("\t\tExpected a \"%s\" type. %v", expectedType, testOK)
			}
		}
	}
}