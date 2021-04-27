# Extension for go-pg/pg

##Methods:
* `ConnOptsFromDsn(dsn string) *pg.Options` - Configuration from DSN [host=localhost port=5432 user=postgres dbname=postgres password=]
* `GetCurrentSchema(db *pg.DB) (schema string, err error)` - Get current schema
* `InitMigrationTableIfNeeded(db *pg.DB, log *logrus.Entry)` - Create default schema and migration table

##Usage
```go
package main
import  (
	"gitlab.mobio.ru/go-packages/pg-ext"
    "github.com/go-pg/pg/v10"
    "github.com/sirupsen/logrus"
)
dsn := "host=localhost port=5432 user=postgres dbname=postgres password="
pgOpts := pg_ext.ConnOptsFromDsn(dsn)
pgClient = pg.Connect(pgOpts)

logger := logrus.New()
pg_ext.InitMigrationTableIfNeeded(pgClient, logger)
```

##Dependencies
* `github.com/go-pg/pg` - PostgresSQL ORM
* `github.com/go-pg/migrations` - PostgresSQL Migrations
* `github.com/sirupsen/logrus` - Logger with hooks