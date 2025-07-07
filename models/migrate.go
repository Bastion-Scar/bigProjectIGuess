package models

import (
	"awesomeProject10/zapLogger"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type okLogs struct {
	id       int64 `gorm:"primaryKey;autoIncrement"`
	path     string
	raw      string
	ip       string
	duration string
}

func InitDb() {
	zapLogger.Init()
	logger := zapLogger.Log

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Ошибка загрузки .env файла")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных:" + err.Error())
	}

	err = db.AutoMigrate(&okLogs{})
	if err != nil {
		panic("Не удалось мигрировать таблицу:" + err.Error())
	}
}
