package dburl

import (
	"errors"
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
	// SQLITE is dialect name for sqlite database engine
	SQLITE string = "sqlite"
	// MSSQL is dialect name for mssql database engine
	MSSQL string = "mssql"
	// POSTGRESQL is dialect name for postgresql database engine
	POSTGRESQL = "postgres"
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
		err := config.parse(url)
		if err != nil {
			return nil, err
		}

		return config, nil
	}

	return nil, ErrDatabaseURLNotFound
}

func (d *DBConfig) parse(url string) error {
	if url == "sqlite://:memory:" {
		d.Dialect = SQLITE
		d.Path = ":memory:"
		return nil
	}

	s := strings.Split(url, "://")
	d.Dialect = s[0]
	switch d.Dialect {
	case MSSQL, POSTGRESQL:
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
