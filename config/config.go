package config

import (
	"base-site-api/log"
	"io/ioutil"
	"os"

	"base-site-api/models"

	"github.com/jinzhu/gorm"
	// sqlite driver needed for gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Constants for whole application setup
type Constants struct {
	ADDRESS      string
	ENV          string
	TemplatePath string
}

// Config application with all constants and database
type Config struct {
	Constants
	Database   *gorm.DB
	SigningKey []byte
}

func initDB(env string) (*gorm.DB, error) {
	conn, err := gorm.Open(os.Getenv("DB_DIALECT"), os.Getenv("DB_URI"))
	if err != nil {
		return nil, err
	}

	conn.LogMode(env == "development")
	conn.Debug().AutoMigrate(
		&models.Article{},
		&models.User{},
		&models.ForgotPasswordToken{},
		&models.PageCategory{},
		&models.Page{},
		&models.UploadType{},
		&models.UploadCategory{},
		&models.Upload{},
	)

	return conn, nil
}

// New return application Config
func New() (*Config, error) {
	var err error

	constants := Constants{}
	constants.ENV = os.Getenv("GO_ENV")
	constants.ADDRESS = os.Getenv("ADDRESS")
	constants.TemplatePath = os.Getenv("TEMPLATE_PATH")

	signingKey, err := ioutil.ReadFile(os.Getenv("JWT_KEY_PATH"))

	if err != nil {
		log.Fatalf("Init jwt open err: %s", err)
	}

	db, err := initDB(constants.ENV)

	if err != nil {
		return nil, err
	}

	return &Config{
		Constants:  constants,
		Database:   db,
		SigningKey: signingKey,
	}, nil
}
