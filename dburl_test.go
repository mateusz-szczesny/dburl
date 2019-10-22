package dburl_test

import (
	"os"
	"testing"

	"github.com/mateusz-szczesny/dburl"
)

func setUpEnvVariable(key string, value string) error {
	return os.Setenv(key, value)
}

func TestPostgreSQL(t *testing.T) {
	setUpEnvVariable(dburl.DefaultEnv, "postgres://USER:PASSWORD@HOST:5432/NAME")
	got, err := dburl.Config(dburl.DefaultEnv)
	if err != nil {
		t.Errorf("url cannot be parsed | error: %s", err)
	}
	if got.Dialect != "postgres" {
		t.Errorf("Dialect does not match | want: %s, get: %s", "postgres", got.Dialect)
	}
	if got.User != "USER" {
		t.Errorf("User does not match | want: %s, get: %s", "USER", got.User)
	}
	if got.Password != "PASSWORD" {
		t.Errorf("Password does not match | want: %s, get: %s", "PASSWORD", got.Password)
	}
	if got.Host != "HOST" {
		t.Errorf("Host does not match | want: %s, get: %s", "HOST", got.Host)
	}
	if got.Port != 5432 {
		t.Errorf("Port does not match | want: %d, get: %d", 5432, got.Port)
	}
	if got.DBName != "NAME" {
		t.Errorf("DBName does not match | want: %s, get: %s", "NAME", got.DBName)
	}
}
func TestMSSQL(t *testing.T) {
	setUpEnvVariable(dburl.DefaultEnv, "mssql://USER:PASSWORD@HOST:1234/NAME")
	got, err := dburl.Config(dburl.DefaultEnv)
	if err != nil {
		t.Errorf("url cannot be parsed | error: %s", err)
	}
	if got.Dialect != "mssql" {
		t.Errorf("Dialect does not match | want: %s, get: %s", "mssql", got.Dialect)
	}
	if got.User != "USER" {
		t.Errorf("User does not match | want: %s, get: %s", "USER", got.User)
	}
	if got.Password != "PASSWORD" {
		t.Errorf("Password does not match | want: %s, get: %s", "PASSWORD", got.Password)
	}
	if got.Host != "HOST" {
		t.Errorf("Host does not match | want: %s, get: %s", "HOST", got.Host)
	}
	if got.Port != 1234 {
		t.Errorf("Port does not match | want: %d, get: %d", 1234, got.Port)
	}
	if got.DBName != "NAME" {
		t.Errorf("DBName does not match | want: %s, get: %s", "NAME", got.DBName)
	}
}
func TestSQLITE(t *testing.T) {
	setUpEnvVariable(dburl.DefaultEnv, "sqlite:///PATH")
	got, err := dburl.Config(dburl.DefaultEnv)
	if err != nil {
		t.Errorf("url cannot be parsed | error: %s", err)
	}
	if got.Dialect != "sqlite" {
		t.Errorf("Dialect does not match | want: %s, get: %s", "sqlite", got.DBName)
	}
	if got.Path != "/PATH" {
		t.Errorf("DBName does not match | want: %s, get: %s", "/PATH", got.DBName)
	}
}

func TestSQLITEInMemory(t *testing.T) {
	setUpEnvVariable(dburl.DefaultEnv, "sqlite://:memory:")
	got, err := dburl.Config(dburl.DefaultEnv)
	if err != nil {
		t.Errorf("url cannot be parsed | error: %s", err)
	}
	if got.Dialect != "sqlite" {
		t.Errorf("Dialect does not match | want: %s, get: %s", "sqlite", got.DBName)
	}
	if got.Path != ":memory:" {
		t.Errorf("DBName does not match | want: %s, get: %s", ":memory:", got.DBName)
	}
}

func TestAccessDataInvalid(t *testing.T) {
	setUpEnvVariable(dburl.DefaultEnv, "mssql://USER:PAS/SWORD@HOST:1234/NAME")
	got, err := dburl.Config(dburl.DefaultEnv)
	if err != dburl.ErrDatabaseURLCannotBeParsed {
		t.Errorf("url should not be parsed | error: %s", err)
	}
	if got != nil {
		t.Errorf("Dialect does not match | want: %s, get: %s", "mssql", got.Dialect)
	}
}

func TestUnsupportedEngine(t *testing.T) {
	setUpEnvVariable(dburl.DefaultEnv, "unsupported_engine://<some_credentials>")
	got, err := dburl.Config(dburl.DefaultEnv)
	if err != dburl.ErrDatabaseEngineNotSupported {
		t.Errorf("engine is not suported | error: %s", err)
	}
	if got != nil {
		t.Errorf("Dialect is not supported | want: %s, get: %s", "mssql", got.Dialect)
	}
}

func TestNotFoundVariable(t *testing.T) {
	setUpEnvVariable("RANDOM_NAME", "VALUE")
	_, err := dburl.Config("NAME_RANDOM")
	if err != dburl.ErrDatabaseURLNotFound {
		t.Error("Env variable should not match!")
	}
}

func TestDefaultVariable(t *testing.T) {
	setUpEnvVariable("DATABASE_URL", "VALUE")
	_, err := dburl.Config(dburl.DefaultEnv)
	if err == dburl.ErrDatabaseURLNotFound {
		t.Error("Env variable should be found!")
	}
}
