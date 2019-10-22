// Package dburl allows you to configure your database configuration based
// on DATABASE_URL environment variable.
//
// The package provides Config(string) (*DBConfig, error) function
// to generate configuration based on environment variable passed
// as an argument.
// Function returns a database config structure, populated with
// all the data specified in your environment variable.
// Package helps to create config for docker based applications,
// which depends on external database.
//
// If you'd rather not use an environment variable,
// you can pass a URL in directly instead to Parse(string) (error) method.
//
package dburl
