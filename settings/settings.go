package settings

import (
	"fmt"
	"go.uber.org/zap"
	"log"

	"github.com/go-ini/ini"
)

type DataBase struct {
	Host      string
	Port      string
	User      string
	Password  string
	DbName    string
	CharSet   string
	ParseTime string
	MaxIdle   int
	MaxOpen   int
}

type App struct {
	Name string
	Mode string
	Port string
	TokenExpireDuration int
}

type Log struct {
	Level      string
	FileName   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type Redis struct {
	Host     string
	Port     string
	DB       int
	PoolSize int
}

type SnowFlake struct {
	StartTime string
	MachineID int64
}

var DataBaseSetting = &DataBase{}
var AppSetting = &App{}
var LogSetting = &Log{}
var RedisSetting = &Redis{}
var SnowFlakeSetting = &SnowFlake{}


var cfg *ini.File

func Setup(filename string) {
	var err error
	cfg, err = ini.Load(filename)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf.ini':%v", err)
	}
	mapTo("mysql", DataBaseSetting)
	mapTo("app", AppSetting)
	mapTo("log", LogSetting)
	mapTo("redis", RedisSetting)
	fmt.Println(RedisSetting)
	mapTo("snowflake", SnowFlakeSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		zap.L().Error("Cfg.MapTo %s err", zap.Error(err))
	}
}
