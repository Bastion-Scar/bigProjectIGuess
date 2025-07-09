package models

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type OkLogs struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	Path     string
	Raw      string
	IP       string
	Duration string
}

func InitDb() *gorm.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", zap.Error(err))
	}
	err = db.AutoMigrate(&OkLogs{})
	if err != nil {
		log.Fatal("Не удалось мигрировать БД:", zap.Error(err))
	}
	return db
}
