package configuration

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// DBConfiguration provides configuration options for connecting to
// the database.
type DBConfiguration struct {
	Type     DBType `default:"postgres"`
	Postgres PostgresConfiguration
}

// ConnectionString generates a connection string for use with the SQL
// package.
func (dbConfig DBConfiguration) ConnectionString() (string, error) {
	switch dbConfig.Type {
	case DBTypePostgres:
		return dbConfig.Postgres.ConnectionString()
	default:
		return "", fmt.Errorf("Unsupported database type '%s'", dbConfig.Type)
	}
}

// DBType gives the type of database we are connecting to.
type DBType string

// Decode provides a hook to parse the database type from an environment
// variables.
func (dbType *DBType) Decode(value string) error {
	if _, ok := supportedDBTypes[value]; !ok {
		return fmt.Errorf("Unsupported database type '%s'", value)
	}

	*dbType = DBType(value)

	return nil
}

var _ envconfig.Decoder = (*DBType)(nil)

const (
	// DBTypeInvalid gives an invalid database type.
	DBTypeInvalid DBType = ""

	// DBTypePostgres tells the service to use SQLite3 as it's backing store.
	DBTypePostgres DBType = "postgres"
)

var (
	supportedDBTypes = map[string]any{
		string(DBTypePostgres): nil,
	}
)

// DBTypeConfiguration provides useful methods that can be called on
// database configuration types.
type DBTypeConfiguration interface {
	ConnectionString() (string, error)
}

// PostgresConfiguration holds any configuration that is specific to
// using a Postgres DB.
type PostgresConfiguration struct {
	Host     string `required:"true"`
	Port     int    `default:"5432"`
	Role     string `default:"ley"`
	Password string `required:"true"`
	DBName   string `default:"ley"`
	SSLMode  string `default:"enable"`
}

// ConnectionString generates a connection string for use with the SQL
// package.
func (dbConfig PostgresConfiguration) ConnectionString() (string, error) {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Role,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	), nil
}

var _ DBTypeConfiguration = (*PostgresConfiguration)(nil)
