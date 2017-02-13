package config

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"gopkg.in/mgo.v2"
	"time"
)

type Static struct {
	Path string
	URL string
}

func (s *Static) GetPath() string {
	if s.Path != "" {
		if string(s.Path[len(s.Path)-1]) != "/" {
			return s.Path + "/"
		}
		return s.Path
	}
	return "./static"
}

func (s *Static) GetURL() string {
	if s.URL != "" { return s.URL }
	return "/static/"
}

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

	Mode string // MongoDB only
}

func (d *Database) getMongoDBConnectionString() (string, error)  {
	if d.GetType() != "mongodb" {
		return "", &configError{"Database:Type", "Field not or incorrectly set"}
	}

	if d.Host == "" {
		return "", &configError{"Database:Host", "Field not set"}
	}

	if d.DB == "" {
		return "", &configError{"Database:DB", "Field not set"}
	}

	str := "mongodb://"
	if d.User != "" {
		str += fmt.Sprintf("%s:%s@%s/%s", d.User, d.Password, d.Host, d.DB)
	} else {
		str += fmt.Sprintf("%s/%s", d.Host, d.DB)
	}

	return str, nil
}

// GetMongoMode returns a mgo.Mode based upon the settings of the configuration file. The default mode is mgo.Strong
func (d *Database) GetMongoMode() (mgo.Mode, error)  {
	if d.GetType() != "mongodb" {
		return -1, &configError{"Database:Type", "Field not or incorrectly set"}
	}

	switch(strings.ToLower(d.Mode)) {
		case "primary":
			return mgo.Primary, nil
		case "primary_preferred":
			return mgo.PrimaryPreferred, nil
		case "secondary":
			return mgo.Secondary, nil
		case "secondary_preferred":
			return mgo.SecondaryPreferred, nil
		case "nearest":
			return mgo.Nearest, nil
		case "eventual":
			return mgo.Eventual, nil
		case "monotonic":
			return mgo.Monotonic, nil
		case "strong":
			return mgo.Strong, nil
		default:
			return mgo.Strong, nil
	}
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
		case "mongodb":
			return d.getMongoDBConnectionString()
		default:
			return "", &configError{"Database:Type", "Unsported database type"}
	}
}

// GetType returns the database type in lowercase
func (d *Database) GetType() string {
	return strings.ToLower(d.Type)
}

// IsGorm returns true or false whether the database type is set to use a Gorm type of database
func (d *Database) IsGorm() bool {
	dbtype := d.GetType()

	return dbtype == "sqlite3" || dbtype == "postgres" || dbtype == "mysql"
}

// IsMGO returns true or false whether the database type is set to use a MongoDB type of database
func (d *Database) IsMGO() bool {
	return d.GetType() == "mongodb"
}

// IsValidType returns true if the Type is set to a valid value, false if set to a false value
func (d *Database) IsValidType() bool {
	return d.IsGorm() || d.IsMGO()
}

type Server struct {
	Host string
	Port uint64
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	MaxHeaderBytes int
	//ToDo: TLS Config
}

func (s *Server) EnsureDefaults() {
	if s.Port == 0 {
		s.Port = 8000
	}

	if s.ReadTimeout == 0 {
		s.ReadTimeout = 10
	}

	if s.WriteTimeout == 0 {
		s.WriteTimeout = 10
	}

	if s.MaxHeaderBytes == 0 {
		s.MaxHeaderBytes = 1 << 20
	}
}

func (s *Server) Addr() string{
	s.EnsureDefaults()

	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Config implements configuration details.
// File: Contains a string to the configuration file
// Database: Contains database configuration details
type Config struct {
	File string `json:"-"`
	Database Database
	Server Server
	Static Static
	TemplatesPath string
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

func (c *Config) GetTemplatesPath() string {
	if c.TemplatesPath != "" {
		if string(c.TemplatesPath[len(c.TemplatesPath)-1]) != "/" {
			return c.TemplatesPath + "/"
		}
		return c.TemplatesPath
	}
	return "templates/"
}