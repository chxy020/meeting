package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"

	// import mysql
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbRegisterLock sync.Mutex
	currentDb      string
)

const (
	defaultAlias = "default"
)

// CurrentDb returns current db
func CurrentDb() string {
	return currentDb
}

// NewOrm returns orm with current active primary db
func NewOrm() (ormer orm.Ormer) {
	return orm.NewOrm()
}

// RegisterDb register db with alias
func RegisterDb(dbHost string, retry ...bool) (err error) {
	if len(retry) > 0 && !retry[0] {
		return registerDatabase(dbHost)
	}

	timeoutC := time.After(120 * time.Second)
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-timeoutC:
			err = errors.Errorf("db register timeout %s", 120*time.Second)
			return
		case <-ticker.C:
			err = registerDatabase(dbHost)
			if err == nil {
				return
			}
		}
	}
}

func registerDatabase(dbHost string) (err error) {
	dbRegisterLock.Lock()
	defer dbRegisterLock.Unlock()

	dbUser := beego.AppConfig.String("dbuser")
	dbPassword := beego.AppConfig.String("dbpassword")
	dbName := beego.AppConfig.String("dbname")
	err = orm.RegisterDataBase(defaultAlias, "mysql",
		dbConnStr(dbHost, 3306, dbUser,
			dbPassword, dbName), -1, 6)
	if err == nil {
		currentDb = dbHost
		return nil
	}
	return errors.Trace(err)
}

func dbConnStr(host string, port int64, user, password, dbname string) string {
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, dbname)
}
