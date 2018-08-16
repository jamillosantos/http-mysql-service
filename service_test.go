package mysqlsrv

import (
	"log"
	"testing"

	"database/sql"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

func TestService(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "MySQL Test Suite")
}

func pingConn(conn *sql.Conn) error {
	return conn.PingContext(context.Background())
}

var _ = Describe("MySQLService", func() {
	It("should fail loading a configuration", func() {
		var service MySQLService
		configuration, err := service.LoadConfiguration()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("not implemented"))
		Expect(configuration).To(BeNil())
	})

	It("should fail applying configuration", func() {
		var service MySQLService
		err := service.ApplyConfiguration(map[string]interface{}{
			"address": "localhost",
		})
		Expect(err).To(Equal(http.ErrWrongConfigurationInformed))
	})

	It("should apply the configuration using a pointer", func() {
		var service MySQLService
		err := service.ApplyConfiguration(&MySQLServiceConfiguration{
			Host:        "host",
			User:        "user",
			Password:    "password",
			Database:    "database",
			Port:        3306,
			MaxPoolSize: 1,
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Host).To(Equal("host"))
		Expect(service.Configuration.User).To(Equal("user"))
		Expect(service.Configuration.Password).To(Equal("password"))
		Expect(service.Configuration.Database).To(Equal("database"))
		Expect(service.Configuration.Port).To(Equal(3306))
		Expect(service.Configuration.MaxPoolSize).To(Equal(1))
	})

	It("should apply the configuration using a copy", func() {
		var service MySQLService
		err := service.ApplyConfiguration(MySQLServiceConfiguration{
			Host:        "host",
			User:        "user",
			Password:    "password",
			Database:    "database",
			Port:        3306,
			MaxPoolSize: 1,
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Host).To(Equal("host"))
		Expect(service.Configuration.User).To(Equal("user"))
		Expect(service.Configuration.Password).To(Equal("password"))
		Expect(service.Configuration.Database).To(Equal("database"))
		Expect(service.Configuration.Port).To(Equal(3306))
		Expect(service.Configuration.MaxPoolSize).To(Equal(1))
	})

	It("should start the service", func() {
		var service MySQLService
		Expect(service.ApplyConfiguration(MySQLServiceConfiguration{
			Host:        "localhost",
			User:        "root",
			Password:    "",
			Database:    "",
			Port:        3306,
			MaxPoolSize: 1,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		defer service.Stop()
		Expect(service.RunWithConn(pingConn)).To(BeNil())
	})

	It("should stop the service", func() {
		var service MySQLService
		Expect(service.ApplyConfiguration(MySQLServiceConfiguration{
			Host:        "localhost",
			User:        "root",
			Password:    "",
			Database:    "",
			Port:        3306,
			MaxPoolSize: 1,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Stop()).To(BeNil())
		Expect(service.RunWithConn(func(conn *sql.Conn) error {
			return nil
		})).To(Equal(http.ErrServiceNotRunning))
	})

	It("should restart the service", func() {
		var service MySQLService
		Expect(service.ApplyConfiguration(MySQLServiceConfiguration{
			Host:        "localhost",
			User:        "root",
			Password:    "",
			Database:    "",
			Port:        3306,
			MaxPoolSize: 1,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Restart()).To(BeNil())
		Expect(service.RunWithConn(pingConn)).To(BeNil())
	})

	It("should initialize a transaction", func() {
		var service MySQLService
		Expect(service.ApplyConfiguration(MySQLServiceConfiguration{
			Host:        "localhost",
			User:        "root",
			Password:    "",
			Database:    "",
			Port:        3306,
			MaxPoolSize: 1,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		defer service.Stop()
		Expect(service.RunWithTx(func(tx *sql.Tx) error {
			Expect(tx).NotTo(BeNil())
			return nil
		})).To(BeNil())
	})

	PIt("should run a transaction committing changes")

	PIt("should run a transaction rolling changes due to return error")

	PIt("should run a transaction rolling changes due to panic")
})
