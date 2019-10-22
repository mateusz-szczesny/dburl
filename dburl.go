package dburl

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DBConfig is structure representing database connection config
type DBConfig struct {
	Dialect  string
	Path     string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

const (
	// DefaultEnv is default env varaible name
	DefaultEnv string = "DATABASE_URL"
	// SQLITE is dialect name for sqlite3 database engine
	SQLITE string = "sqlite3"
	// MSSQL is dialect name for mssql database engine
	MSSQL string = "mssql"
	// POSTGRES is dialect name for postgres database engine
	POSTGRES = "postgres"
)

var (
	// ErrDatabaseURLNotFound is returned when DATABASE_URL env variable is not set
	ErrDatabaseURLNotFound = errors.New("env varialbe not found")
	// ErrDatabaseURLInvalidSyntac is returned when DATABSE_URL env variable has invalid syntax
	ErrDatabaseURLInvalidSyntac = errors.New("url's syntax is not valid to be parsed")
	// ErrDatabaseEngineNotSupported is returned when database engine is incorrect or unsupported
	ErrDatabaseEngineNotSupported = errors.New("engine incorrect or unsupported")
	// ErrDatabaseURLCannotBeParsed is returned when data are not valid (e.g. password contains any of (':','@','/'))
	ErrDatabaseURLCannotBeParsed = errors.New("usr data cannot be parsed")
)

// Config creates DBConfig based on environment
func Config(env string) (*DBConfig, error) {
	config := &DBConfig{}

	url, ok := os.LookupEnv(env)
	if ok {
		err := config.Parse(url)
		if err != nil {
			return nil, err
		}

		return config, nil
	}

	return nil, ErrDatabaseURLNotFound
}

// Parse allows to fill configuration structure with details from parameter
func (d *DBConfig) Parse(url string) error {
	if url == "sqlite3://:memory:" {
		d.Dialect = SQLITE
		d.Path = ":memory:"
		return nil
	}

	s := strings.Split(url, "://")
	d.Dialect = s[0]
	switch d.Dialect {
	case MSSQL, POSTGRES:
		accessData := strings.FieldsFunc(s[1], func(r rune) bool {
			return r == ':' || r == '@' || r == '/'
		})

		if len(accessData) != 5 {
			return ErrDatabaseURLCannotBeParsed
		}

		d.User = accessData[0]
		d.Password = accessData[1]
		d.Host = accessData[2]
		d.Port, _ = strconv.Atoi(accessData[3])
		d.DBName = accessData[4]

		return nil
	case SQLITE:
		d.Path = s[1]
		return nil
	default:
		return ErrDatabaseEngineNotSupported
	}
}

// GetConnectionString returns connection string for specified database config
func (d *DBConfig) GetConnectionString() string {
	var connection string
	switch d.Dialect {
	case MSSQL:
		connection = getMSSQLConnectionString(d)
	case SQLITE:
		connection = getSQLITEConnectionString(d)
	case POSTGRES:
		connection = getPostgresConnectionString(d)
	default:
		connection = ""
	}
	return connection
}

func getMSSQLConnectionString(d *DBConfig) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

func getSQLITEConnectionString(d *DBConfig) string {
	return fmt.Sprintf("%s", d.Path)
}

func getPostgresConnectionString(d *DBConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.DBName)
}
