package mysqlsrv

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lab259/http"
)

// MySQLServiceConfiguration describes the `MySqlService` configuration.
type MySQLServiceConfiguration struct {
	Host        string `yaml:"host"`
	User        string `yaml:"user"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Password    string `yaml:"password"`
	MaxPoolSize int    `yaml:"max_pool_size"`
}

// MySQLService implements the MySQL service itself.
type MySQLService struct {
	running       bool
	db            *sql.DB
	Configuration MySQLServiceConfiguration
}

// MySQLServiceConnHandler is the handler description for being used on
// `MySQLService.RunWithConn` method.
type MySQLServiceConnHandler func(conn *sql.Conn) error

// LoadConfiguration is an abstract method that should be overwritten on the
// actual usage of this service.
func (service *MySQLService) LoadConfiguration() (interface{}, error) {
	return nil, errors.New("not implemented")
}

// ApplyConfiguration implements the type verification of the given
// `configuration` and applies it to the service.
func (service *MySQLService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case MySQLServiceConfiguration:
		service.Configuration = c
	case *MySQLServiceConfiguration:
		service.Configuration = *c
	default:
		return http.ErrWrongConfigurationInformed
	}

	return nil
}

// Restart restarts the service.
func (service *MySQLService) Restart() error {
	if service.running {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	return service.Start()
}

// Start initialize the mysql connection and saves the db.
func (service *MySQLService) Start() error {
	if !service.running {
		var err error

		conf := service.Configuration

		// Create the connection string
		connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&collation=utf8mb4_unicode_ci",
			conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

		// Open database
		service.db, err = sql.Open("mysql", connString)

		if err != nil {
			return err
		}

		// Sets the max pool size
		if conf.MaxPoolSize > 0 {
			service.db.SetMaxOpenConns(conf.MaxPoolSize)
		}

		// Pings the session to ensure it is working
		err = service.db.Ping()

		if err != nil {
			return err
		}

		service.running = true
	}
	return nil
}

// Stop stops the service.
func (service *MySQLService) Stop() error {
	if service.running {
		err := service.db.Close()
		if err != nil {
			return err
		}
		service.running = false
	}
	return nil
}

// RunWithConn runs a handler passing a new instance of the a conn.
func (service *MySQLService) RunWithConn(handler MySQLServiceConnHandler) error {
	if !service.running {
		return http.ErrServiceNotRunning
	}

	conn, err := service.db.Conn(context.Background())

	if err != nil {
		return err
	}

	defer conn.Close()
	return handler(conn)
}

// RunWithTx runs a handler passing a new instance of the a transaction.
func (service *MySQLService) RunWithTx(handler func(conn *sql.Tx) error) (globalErr error) {
	if !service.running {
		return http.ErrServiceNotRunning
	}

	tx, err := service.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		return
	}()

	err = handler(tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
