package database

import (
	"time"

	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"

	// blank import migrations
	_ "seater/database/migrations"
)

// migration status
const (
	MigrationsStatusUpdate   = "update"
	MigrationsStatusRollback = "rollback"
)

// migrationsSQL is the DDL for migration table
const migrationsSQL = `
CREATE TABLE IF NOT EXISTS migrations (
	id_migration SERIAL PRIMARY KEY,
	name varchar(255) DEFAULT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	statements text,
	rollback_statements text,
	status varchar(32)
)
`

// Migrations defines migrations model
type Migrations struct {
	IDMigration        int64     `orm:"auto;pk;column(id_migration)"`
	Name               string    `orm:"null"`
	CreatedAt          time.Time `orm:"auto_now;type(datetime)"`
	Statements         string    `orm:"null;type(text)"`
	RollbackStatements string    `orm:"null;type(text)"`
	Status             string    `orm:"size(32);null"`
}

// ListMigrations lists all migrations
func ListMigrations() (migrationsList []*Migrations, err error) {
	o := NewOrm()
	qs := o.QueryTable(new(Migrations))
	qs = qs.OrderBy("CreatedAt")
	_, err = qs.All(&migrationsList)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// GetLastestMigrations returns the latest migrations
func GetLastestMigrations(onlyUpdate ...bool) (migrations *Migrations, err error) {
	o := NewOrm()
	migrations = new(Migrations)
	qs := o.QueryTable(migrations)

	if !(len(onlyUpdate) > 0 && !onlyUpdate[0]) {
		qs = qs.Filter("Status", MigrationsStatusUpdate)
	}

	qs = qs.OrderBy("-CreatedAt").Limit(1)
	err = qs.One(migrations)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

func getCreateTimeByName(name string) (createdAt int64, err error) {
	createdAtStr := name[len(name)-15:]
	if t, err := time.Parse(migration.DateFormat, createdAtStr); err != nil {
		err = errors.Trace(err)
	} else {
		createdAt = t.Unix()
	}
	return
}

// CheckCreateMigrationsTable check and create migrations table
func CheckCreateMigrationsTable() (err error) {
	o := NewOrm()
	_, err = o.Raw(migrationsSQL).Exec()
	if err != nil {
		return errors.Trace(err)
	}
	return
}

// UpgradeDB upgrade the migration from lasttime
func UpgradeDB() (err error) {
	migrations, err := GetLastestMigrations()
	if err != nil {
		return
	}
	var lastTime int64
	if migrations != nil {
		lastTime, err = getCreateTimeByName(migrations.Name)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return migration.Upgrade(lastTime)
}

// DowngradeDB rollback the latest migration
func DowngradeDB() (err error) {
	migrations, err := GetLastestMigrations()
	if err != nil {
		return errors.Trace(err)
	}
	if migrations == nil {
		return
	}
	err = migration.Rollback(migrations.Name)
	if err != nil {
		return errors.Trace(err)
	}
	return
}

// ResetDB rollback all migration
func ResetDB() (err error) {
	return migration.Reset()
}

func init() {
	orm.RegisterModel(new(Migrations))
}
