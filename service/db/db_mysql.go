package db

import (
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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
	dbDriverName := baseconfig.BasicConfigs.DbDriverName

	path := ""
	if dbDriverName == "mysql"{
		path = strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	}else if dbDriverName == "postgres"{
		path = strings.Join([]string{"host=", ip, " port=" , port, " user=", userName, " password=", password, " dbname=", dbName, " sslmode=disable"}, "")
	}
	fmt.Println(path)
	var err error = nil
	//DB, err = sql.Open("mysql", path)
	DB, err = sql.Open(dbDriverName, path)
	if err != nil{
		panic(err)
	}
	DB.SetMaxIdleConns(baseconfig.BasicConfigs.MaxIdleConns)
	DB.SetMaxOpenConns(baseconfig.BasicConfigs.MaxOpenConns)

	if err := DB.Ping(); err != nil{
		level.Info(logger).Log("msg", "Database Connections init error")
		panic(err)
	}
	level.Info(logger).Log("msg", "Database Connections init success")
}
