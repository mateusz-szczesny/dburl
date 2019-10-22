package dburl_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/mateusz-szczesny/dburl"
)

func setUpEnvVariable(key string, value string) error {
	return os.Setenv(key, value)
}

func TestConfig(t *testing.T) {
	type envVariable struct {
		key   string
		value string
	}
	type args struct {
		env string
	}
	tests := []struct {
		url     string
		name    string
		args    args
		want    *dburl.DBConfig
		envVars []envVariable
		wantErr bool
	}{
		{
			name: "postgres - valid",
			args: args{
				env: dburl.DefaultEnv,
			},
			envVars: []envVariable{
				{
					key:   dburl.DefaultEnv,
					value: "postgres://USER:PASSWORD@HOST:5432/NAME",
				},
			},
			want: &dburl.DBConfig{
				Dialect:  "postgres",
				User:     "USER",
				Password: "PASSWORD",
				Host:     "HOST",
				Port:     5432,
				DBName:   "NAME",
			},
			wantErr: false,
		},
		{
			name: "mssql - valid",
			args: args{
				env: "DB_URL",
			},
			envVars: []envVariable{
				{
					key:   "DB_URL",
					value: "mssql://USER:PASSWORD@HOST:1433/NAME",
				},
			},
			want: &dburl.DBConfig{
				Dialect:  "mssql",
				User:     "USER",
				Password: "PASSWORD",
				Host:     "HOST",
				Port:     1433,
				DBName:   "NAME",
			},
			wantErr: false,
		},
		{
			name: "sqlite3 - valid",
			args: args{
				env: dburl.DefaultEnv,
			},
			envVars: []envVariable{
				{
					key:   dburl.DefaultEnv,
					value: "sqlite3:///some/path/to/database.db",
				},
			},
			want: &dburl.DBConfig{
				Dialect: "sqlite3",
				Path:    "/some/path/to/database.db",
			},
			wantErr: false,
		},
		{
			name: "sqlite3 in memory - valid",
			args: args{
				env: dburl.DefaultEnv,
			},
			envVars: []envVariable{
				{
					key:   dburl.DefaultEnv,
					value: "sqlite3://:memory:",
				},
			},
			want: &dburl.DBConfig{
				Dialect: "sqlite3",
				Path:    ":memory:",
			},
			wantErr: false,
		},
		{
			name: "mssql - invalid",
			args: args{
				env: "DB_URL",
			},
			envVars: []envVariable{
				{
					key:   "DB_URL",
					value: "mssql://USER:PASS/WORD@HOST:1433/NAME",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unsupported engine",
			args: args{
				env: "DB_URL",
			},
			envVars: []envVariable{
				{
					key:   "DB_URL",
					value: "unsupported_engine://<some_credentials>",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "env variable not found",
			args: args{
				env: "DB_URL_XYZ",
			},
			envVars: []envVariable{
				{
					key:   "DB_URL",
					value: "unsupported_engine://<some_credentials>",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, envVar := range tt.envVars {
				setUpEnvVariable(envVar.key, envVar.value)
			}
			got, err := dburl.Config(tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBConfig_GetConnectionString(t *testing.T) {
	tests := []struct {
		name string
		d    *dburl.DBConfig
		want string
	}{
		{
			name: "postgresql",
			d: &dburl.DBConfig{
				Dialect:  "postgres",
				User:     "USER",
				Password: "PASSWORD",
				Host:     "HOST",
				Port:     5432,
				DBName:   "NAME",
			},
			want: "host=HOST port=5432 user=USER password=PASSWORD dbname=NAME sslmode=disable",
		},
		{
			name: "mssql",
			d: &dburl.DBConfig{
				Dialect:  "mssql",
				User:     "USER",
				Password: "PASSWORD",
				Host:     "HOST",
				Port:     5432,
				DBName:   "NAME",
			},
			want: "sqlserver://USER:PASSWORD@HOST:5432?database=NAME",
		},
		{
			name: "sqlite3",
			d: &dburl.DBConfig{
				Dialect: "sqlite3",
				Path:    "/some/path/to/database.db",
			},
			want: "/some/path/to/database.db",
		},
		{
			name: "sqlite3 in memory",
			d: &dburl.DBConfig{
				Dialect: "sqlite3",
				Path:    ":memory:",
			},
			want: ":memory:",
		},
		{
			name: "empty",
			d:    &dburl.DBConfig{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetConnectionString(); got != tt.want {
				t.Errorf("DBConfig.GetConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
