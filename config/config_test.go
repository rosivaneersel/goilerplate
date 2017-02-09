package config

import (
	"testing"
	"reflect"
)

const testOK = "\u2714"
const testFailed = "\u2718"

func TestSaveAndLoadConfig(t *testing.T) {
	file := "/tmp/config.json"
	expectedConfig := &Config{
		File: file,
		Database: Database{
			Type: "sqlite",
			DB: "/tmp/test.db",
		},
	}

	t.Log("While running tests for the config package")
	{
		t.Log("\tTesting config without a file set.")
		{
			expectedErr := ErrFileNotSet
			c := &Config{}

			t.Log("\t\tTesting Save")
			{
				err := c.Save()

				if err != expectedErr {
					t.Fatalf("\t\t\tExpected to fail with error \"%s\", but got \"%s\" instead. %v", expectedErr, err, testFailed)
				}
				t.Logf("\t\t\tExpected to fail with error \"%s\". %v", expectedErr, testOK)
			}

			t.Log("\t\tTesting Load")
			{
				err := c.Load()

				if err != expectedErr {
					t.Fatalf("\t\t\tExpected to fail with error \"%s\", but got \"%s\" instead. %v", expectedErr, err, testFailed)
				}
				t.Logf("\t\t\tExpected to fail with error \"%s\". %v", expectedErr, testOK)
			}
		}

		t.Log("\tTesting Save")
		{
			err := expectedConfig.Save()
			if err != nil {
				t.Fatalf("\t\tExpected to be able to save config, but got error: %s. %v", err, testFailed)
			}
			t.Log("\t\tExpected to be able to save config. ", testOK)
		}

		t.Log("\tTesting Load")
		{
			c := &Config{File: file}
			err := c.Load()
			if err != nil {
				t.Fatalf("\t\tExpected to be able to load config, but got error: %s. %v", err, testFailed)
			}
			t.Log("\t\tExpected to be able to read config. ", testOK)

			if !reflect.DeepEqual(expectedConfig, c) {
				t.Fatalf("\t\tExpected content to be %+v, but got %+v instead. %v", *expectedConfig, *c, testFailed)
			}
			t.Log("\t\tExpected content. ", testOK)
		}
	}
}

func TestDBConnectionStringCreators(t *testing.T) {
	type Test struct {
		c Config
		expected string
	}

	Tests := []Test{
		{
			Config{
				Database: Database{
					Type: "sqlite3",
					DB:   "/tmp/test.db",
				},
			},
			"/tmp/test.db",
		},
		{
			Config{
				Database: Database{
					Type: "postgres",
					DB: "mydb",
					Host: "myhost",
					User: "myuser",
					Password: "mypassword",

					//SSL: false

					//Charset: "utf-8",
					//ParseTime: true,
					//Local: "Local",
				},
			},
			"host=myhost user=myuser password=mypassword dbname=mydb sslmode=disable",
		},
		{
			Config{
				Database: Database{
					Type: "mysql",
					DB: "mydb",
					Host: "myhost",
					User: "myuser",
					Password: "mypassword",

					//SSL: false

					//Charset: "utf-8",
					//ParseTime: true,
					//Local: "Local",
				},
			},
			"myuser:mypassword@myhost/mydb?charset=utf8&parseTime=false&loc=Local",
		},
		{
			Config{
				Database: Database{
					Type: "mysql",
					DB: "mydb",
					Host: "localhost",
					User: "myuser",
					Password: "mypassword",

					//SSL: false

					//Charset: "utf-8",
					//ParseTime: true,
					//Local: "Local",
				},
			},
			"myuser:mypassword@/mydb?charset=utf8&parseTime=false&loc=Local",
		},
		{
			Config{
				Database: Database{
					Type: "mongodb",
					DB: "test",
					Host: "localhost",
					Mode: "Monotonic",
				},
			},
			"mongodb://localhost/test",
		},
		{
			Config{
				Database: Database{
					Type: "mongodb",
					DB: "test",
					User: "test",
					Password: "test",
					Host: "127.0.0.1",
					Mode: "Monotonic",
				},
			},
			"mongodb://test:test@127.0.0.1/test",
		},
	}

	t.Log("Testing Connection Strings")
	{
		for _, test := range Tests {
			t.Logf("\t%s (%s)", test.c.Database.Type, test.c.Database.Host)
			{
				s, err := test.c.Database.GetDBConnectionString()
				if err != nil {
					t.Fatalf("\t\tExpected to get a string, but got an error: %s. %v", err, testFailed)
				}
				t.Log("\t\tExpected to get a string", testOK)

				if s != test.expected {
					t.Fatalf("\t\tExpected \"%s\", received \"%s\". %v", test.expected, s, testFailed)
				}
				t.Logf("\t\tExpected \"%s\". %v", test.expected, testOK)

			}
		}
	}
}

func TestServerAddr(t *testing.T) {
	tests := []struct{
		Cfg Config
		Expected string
		Description string
	}{
		{Config{Server:Server{}}, ":8000", "Empty host and port"},
		{Config{Server:Server{Host:"localhost"}}, "localhost:8000", "Localhost with empty port"},
		{Config{Server:Server{Host:"localhost", Port: 8080}}, "localhost:8080", "Localhost on port 80"},
	}

	for _, test := range tests{
		t.Logf("Testing creation of webserver address with host: %s and port: %d", test.Cfg.Server.Host, test.Cfg.Server.Port)
		{
			addr := test.Cfg.Server.Addr()
			if addr != test.Expected {
				t.Fatalf("\tExpected: \"%s\", but got \"%s\" instead. %v", test.Expected, addr, testFailed)
			}
			t.Logf("\tExpected: \"%s\". %v", test.Expected, testOK)
		}
	}
}
