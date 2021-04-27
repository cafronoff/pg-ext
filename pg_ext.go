package pg_ext

import (
	"context"
	"strings"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

const MigrationUsageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

const MigrationTable = "gopg_migrations"

type DbLogger struct {
	LogFunc func(query string, params []interface{})
	ErrFunc func(err error)
}

func (d DbLogger) BeforeQuery(_ context.Context, _ *pg.QueryEvent) {}

func (d DbLogger) AfterQuery(q *pg.QueryEvent) {
	query, err := q.UnformattedQuery()
	if err != nil {
		d.ErrFunc(err)
	}
	d.LogFunc(string(query), q.Params)
}

//GetCurrentSchema Get current schema
func GetCurrentSchema(db *pg.DB) (schema string, err error) {
	_, err = db.Query(pg.Scan(&schema), "show search_path")
	return schema, err
}

//InitMigrationTableIfNeeded Create default schema and migration table
func InitMigrationTableIfNeeded(db *pg.DB, log *logrus.Entry) {

	schema, err := GetCurrentSchema(db)
	if err != nil {
		log.Panic(err)
	}

	setSchemaForMigration(schema)

	if err != nil {

	}
	exist, err := isExistMigrationTable(db, schema)
	if err != nil {
		log.Panic(err)
	}

	if !exist {
		_, _, err = migrations.Run(db, "init")
		if err != nil {
			log.Panic(err)
		}
	}
}

func isExistMigrationTable(db *pg.DB, schema string) (exist bool, err error) {
	_, err = db.Query(pg.Scan(&exist), "SELECT EXISTS ("+
		"SELECT 1 "+
		"FROM   information_schema.tables "+
		"WHERE  table_schema = ?0 "+
		"AND    table_name = ?1"+
		");", schema, MigrationTable)
	return exist, err
}

func setSchemaForMigration(schema string) {
	migrations.SetTableName(schema + "." + MigrationTable)
}

//ConnOptsFromDsn Configuration from DSN [host=localhost port=5432 user=postgres dbname=postgres password=]
func ConnOptsFromDsn(dsn string) *pg.Options {
	host := "localhost"
	port := "5432"
	opts := &pg.Options{}

	dsn = strings.Trim(dsn, " ")
	for _, opt := range strings.Split(dsn, " ") {
		kvArr := strings.Split(opt, "=")
		//vs: skip double space and spaces at the beginning and at the end, and incorrect settings
		if len(kvArr) != 2 {
			continue
		}
		key, val := kvArr[0], kvArr[1]
		switch key {
		case "host":
			host = val
			opts.Addr = host + ":" + port
			break
		case "port":
			port = val
			opts.Addr = host + ":" + port
			break
		case "user":
			opts.User = val
			break
		case "password":
			opts.Password = val
			break
		case "dbname":
			opts.Database = val
			break
		}
	}
	return opts
}
