package dao

import (
	"fmt"
	"github.com/sdblg/hoa-auth/src/models"
	log "github.com/sirupsen/logrus"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBArguments struct {
	dbHost string
	port   string
	name   string
	user   string
	pass   string
}

func (d *DBArguments) Connect() (*gorm.DB, error) {

	if err := d.InitFromEnv(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(d.connectStr()), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal("Connecting to database ::", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Connecting to database :: '%v'", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	log.Infof("Connected to database :: %v", d.shortConnectStr())

	db.Logger.LogMode(4)
	err = db.AutoMigrate(&models.User{}, &models.Address{})
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

func (d *DBArguments) InitFromEnv() error {
	if d == nil {
		return fmt.Errorf("*DBArguments struct is nil, please instantiate first")
	}

	var ok bool

	if d.dbHost, ok = os.LookupEnv("DB_HOST"); !ok {
		return fmt.Errorf("check DBArguments config")
	}
	if d.port, ok = os.LookupEnv("DB_PORT"); !ok {
		return fmt.Errorf("check DBArguments config")
	}
	if d.name, ok = os.LookupEnv("DB_SCHEMA"); !ok {
		return fmt.Errorf("check DBArguments config")
	}
	if d.user, ok = os.LookupEnv("DB_USERNAME"); !ok {
		return fmt.Errorf("check DBArguments config")
	}
	if d.pass, ok = os.LookupEnv("DB_PASSWORD"); !ok {
		return fmt.Errorf("check DBArguments config")
	}

	return nil
}

func (d *DBArguments) shortConnectStr() string {
	return fmt.Sprintf("host=%v port=%v dbname=%v user=%v", d.dbHost, d.port, d.name, d.user)
}

func (d *DBArguments) connectStr() string {
	return fmt.Sprintf("%v password=%v sslmode=disable", d.shortConnectStr(), d.pass)
}
