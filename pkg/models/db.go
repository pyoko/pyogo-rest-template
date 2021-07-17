package models

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/pyoko/gorest/pkg/settings"
)

type DB struct {
	*gorm.DB
}

func DbConnect() (*DB, error) {
	host := settings.ReadConfig("DB_HOST")
	port := settings.ReadConfig("DB_PORT")
	database := settings.ReadConfig("DB_DATABASE")
	user := settings.ReadConfig("DB_USERNAME")
	password := settings.ReadConfig("DB_PASSWORD")

	// set up dns string
	var dns strings.Builder
	dns.WriteString(user)
	if password != "" {
		dns.WriteString(":")
		dns.WriteString(password)
	}
	dns.WriteString("@tcp(%s:%s)/%s?parseTime=true")

	// connect
	dbConnection, err := gorm.Open(mysql.Open(fmt.Sprintf(dns.String(), host, port, database)), &gorm.Config{})

	return &DB{dbConnection}, err
}