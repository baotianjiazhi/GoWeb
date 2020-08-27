package mysql

import (
	"fmt"
	"bluebell/settings"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() (err error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=%s",
		settings.DataBaseSetting.User,
		settings.DataBaseSetting.Password,
		settings.DataBaseSetting.Host,
		settings.DataBaseSetting.Port,
		settings.DataBaseSetting.DbName,
		settings.DataBaseSetting.CharSet,
		settings.DataBaseSetting.ParseTime)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("mysql connect err", zap.Error(err))
	}

	//db.SetMaxOpenConns(settings.DataBaseSetting.MaxOpen)
	//db.SetMaxIdleConns(settings.DataBaseSetting.MaxIdle)

	return
}

func Close() {
	_ = db.Close()
}
