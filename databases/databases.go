package databases

import (
	"errors"
	"github.com/BalkanTech/goilerplate/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/mgo.v2"
)

// ErrNotGorm is used in case when the database type in the config file isn't a Gorm type of database
var ErrNotGorm = errors.New("Not a Gorm database")
var ErrNotMongoDB = errors.New("Not a Mongo database")

//NewGormConnection reads the provided config and returns an active Gorm database connection or an error
func NewGormConnection(c *config.Config) (*gorm.DB, error) {
	if c.Database.GetType() != "sqlite3" &&
		c.Database.GetType() != "postgres" &&
		c.Database.GetType() != "mysql" {
		return nil, ErrNotGorm
	}

	s, err := c.Database.GetDBConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(c.Database.GetType(), s)
	if c.Debug {
		db.LogMode(true)
	}
	return db, err

}

type MongoDBConnection struct {
	Session *mgo.Session
	DB      *mgo.Database
}

//NewMgoConnection reads the provided config and returns an active MGO session or an error
func NewMgoConnection(c *config.Config) (*MongoDBConnection, error) {
	if c.Database.GetType() != "mongodb" {
		return nil, ErrNotMongoDB
	}

	s, err := c.Database.GetDBConnectionString()
	if err != nil {
		return nil, err
	}

	session, err := mgo.Dial(s)
	if err != nil {
		return nil, err
	}

	mode, err := c.Database.GetMongoMode()
	if err != nil {
		return nil, err
	}
	session.SetMode(mode, true)

	db := session.DB("") // The DB name has been provided via the dial string

	return &MongoDBConnection{Session: session, DB: db}, nil
}
