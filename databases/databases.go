package databases

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/BalkanTech/goilerplate/config"
	"strings"
	"errors"
	"gopkg.in/mgo.v2"
)

// ErrNotGorm is used in case when the database type in the config file isn't a Gorm type of database
var ErrNotGorm = errors.New("Not a Gorm database")

//NewGormConnection reads the provided config and returns an active Gorm database connection or an error
func NewGormConnection(c *config.Config) (*gorm.DB, error) {
	if strings.ToLower(c.Database.Type) != "sqlite3" &&
	  strings.ToLower(c.Database.Type) != "postgres" &&
	  strings.ToLower(c.Database.Type) != "mysql" {
		return nil, ErrNotGorm
	}

	s, err := c.Database.GetDBConnectionString()
	if err != nil {
		return nil, err
	}

	return gorm.Open(c.Database.GetType(), s)

}

func NewMgoConnection(c *config.Config) (*mgo.Session, error) {
	return mgo.Dial("server1.example.com,server2.example.com")
}

//ToDO: Implement MongoDB
//ToDo: Add comments
