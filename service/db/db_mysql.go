package db

import (
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/prometheus/service/baseconfig"
	"strings"
)

var (
	DB *sql.DB
)

// InitDB 数据库链接初始化
func InitDB(logger log.Logger){

	level.Info(logger).Log("msg", "Init Database Connections")

	userName := baseconfig.BasicConfigs.Username
	password := baseconfig.BasicConfigs.Password
	ip       := baseconfig.BasicConfigs.Ip
	port 	 := baseconfig.BasicConfigs.Port
	dbName   := baseconfig.BasicConfigs.Dbname

	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	var err error = nil
	DB, err = sql.Open("mysql", path)
	if err != nil{
		panic(err)
	}
	DB.SetMaxIdleConns(baseconfig.BasicConfigs.MaxIdleConns)
	DB.SetMaxOpenConns(baseconfig.BasicConfigs.MaxOpenConns)

	if err := DB.Ping(); err != nil{
		panic(err)
	}
	level.Info(logger).Log("msg", "Database Connections init success")
}
