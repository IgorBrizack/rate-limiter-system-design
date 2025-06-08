package database

import (
	"fmt"
	"log"
	"os"

	"github.com/IgorBrizack/rate-limiter-system-design/internal/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	dbUser string
	dbPass string
	dbName string
	dbHost string
}

func NewDatabase() *Database {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env.")
	}

	db := &Database{
		dbUser: os.Getenv("MYSQL_USER"),
		dbPass: os.Getenv("MYSQL_PASSWORD"),
		dbName: os.Getenv("MYSQL_DATABASE"),
		dbHost: os.Getenv("DB_HOST"),
	}

	return db
}

func (d *Database) DB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(d.GetDbDSN()), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{})

	if err != nil {
		log.Fatal("Failed to connect to the  database:", err)
	}

	return db
}

func (d *Database) GetDbDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		d.dbUser, d.dbPass, d.dbHost, d.dbName)
}
