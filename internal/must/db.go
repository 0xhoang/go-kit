package must

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDb(conStr string) *gorm.DB {
	sqlDB, err := sql.Open("mysql", conStr)

	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn: sqlDB,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	if err != nil {
		panic(err)
	}

	return gormDB
}
