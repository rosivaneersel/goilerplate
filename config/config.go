package config

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
)

// Database implements database connection configuration details
type Database struct {
	Type string // Options: mysql, postgres, sqlite3, mongodb
	DB string // Name of the database (which is a filepath for sqlite3
	Host string // For mysql, postgres, mongodb
	User string // For mysql, postgres, mongodb
	Password string // For mysql, postgres, mongodb

	SSL bool // For postgres only

	Charset string // For mysql only
	ParseTime bool // For mysql only
	Local string // For mysql only
}

func (d *Database) getSqLite3ConnectionString() (string, error) {
	if d.GetType() != "sqlite3" {
		return "", &configError{"Database:Type", "Field not or incorrectly set"}
	}

	if d.DB == "" {
		return "", &configError{"Database:DB", "Field not set"}
	}

	return d.DB, nil
}

func (d *Database) getPostgresConnectionString() (string, error) {
	if d.GetType() != "postgres" {
		return "", &configError{"Database:Type", "Field not or incorrectly set"}
	}

	if d.Host == "" {
		return "", &configError{"Database:Host", "Field not set"}
	}

	if d.DB == "" {
		return "", &configError{"Database:DB", "Field not set"}
	}

	ssl := "disable"
	if d.SSL {
		ssl = "enable"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.User, d.Password, d.DB, ssl), nil
}

func (d *Database) getMySQLConnectionString() (string, error) {
	if d.GetType() != "mysql" {
		return "", &configError{"Database:Type", "Field not or incorrectly set"}
	}

	if d.DB == "" {
		return "", &configError{"Database:DB", "Field not set"}
	}

	if d.Local == "" {
		d.Local = "Local"
	}

	if d.Charset == "" {
		d.Charset = "utf8"
	}

	if strings.ToLower(d.Host) == "localhost" || d.Host == "127.0.0.1" {
		d.Host = ""
	}

	return fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=%t&loc=%s", d.User, d.Password, d.Host, d.DB, d.Charset, d.ParseTime, d.Local), nil
}

// GetDBConnectionString creates and returns a formatted database connection string
func (d *Database) GetDBConnectionString() (string, error) {
	switch d.GetType() {
		case "sqlite3":
			return d.getSqLite3ConnectionString()
		case "postgres":
			return d.getPostgresConnectionString()
		case "mysql":
			return d.getMySQLConnectionString()
		default:
			return "", &configError{"Database:Type", "Unsported database type"}
	}
}

// GetType returns the database type in lowercase
func (d *Database) GetType() string {
	return strings.ToLower(d.Type)
}

// Config implements configuration details.
// File: Contains a string to the configuration file
// Database: Contains database configuration details
type Config struct {
	File string `json:"-"`
	Database Database `json:"database"`
}

// Load will load the database file into the Config instance
func (c *Config) Load() error {
	if c.File == "" {
		return ErrFileNotSet
	}

	data, err := ioutil.ReadFile(c.File)
	if err !=  nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

// Save will save the Config instance to the configuration file
func (c *Config) Save() error {
	if c.File == "" {
		return ErrFileNotSet
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.File, data, 0644)
	if err != nil{
		return err
	}

	return nil
}

//ToDo: Implement MongoDB