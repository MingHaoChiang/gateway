package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"com.hyweb/gateway/common"
	"com.hyweb/gateway/web"
	"github.com/facebookgo/grace/gracehttp"
	_ "github.com/go-sql-driver/mysql" // WARNING!
	"github.com/gorilla/handlers"
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
)

func main() {
	serverName := "GATEWAY"
	log.Printf("Gateway Server started")

	configFile := os.Getenv(serverName + "_CONFIG")
	if configFile == "" {
		log.Fatalln("env param " + serverName + "_CONFIG not set.")
	}
	common.ReadConfig(configFile)

	db, err := sql.Open("mysql", common.Config.DATABASE)
	common.Conn = &common.MySQLConn{DB: db}
	if err != nil {
		log.Println("create mysql connection error", err.Error())
		return
	}
	common.Conn.DB.SetMaxIdleConns(30)
	common.Conn.DB.SetConnMaxLifetime(600 * time.Second)
	defer common.Conn.DB.Close()

	accessLogFilePath := os.Getenv(serverName + "_ACCESS_LOG")
	if accessLogFilePath == "" {
		log.Fatalln("env param " + serverName + "_ACCESS_LOG not set.")
	}
	accessLogFile, err := os.OpenFile(accessLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer accessLogFile.Close()
	if err != nil {
		log.Println("create access log error", err.Error())
	}
	webLogFilePath := os.Getenv(serverName + "_WEB_LOG")
	if webLogFilePath == "" {
		log.Fatalln("env param " + serverName + "_WEB_LOG not set.")
	}
	common.WebLog, err = os.OpenFile(webLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer common.WebLog.Close()
	if err != nil {
		log.Println("create web log error", err.Error())
	}
	router := web.NewRouter()
	loggedRouter := handlers.CombinedLoggingHandler(accessLogFile, router)
	pidFilePath := os.Getenv(serverName + "_PID")
	if pidFilePath == "" {
		log.Fatalln("env param " + serverName + "_PID not set.")
	}
	pidFile, err := os.Create(pidFilePath)
	if err != nil {
		log.Println("create pid file error", err.Error())
	}
	pidFile.WriteString(strconv.Itoa(os.Getpid()))
	defer pidFile.Close()
	gracehttp.Serve(&http.Server{Addr: ":" + strconv.Itoa(common.Config.PORT), Handler: loggedRouter})
}
