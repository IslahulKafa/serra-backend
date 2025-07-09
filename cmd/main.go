package main

import (
	"database/sql"
	"log"
	"serra/api"
	"serra/config"
	"serra/db"

	"github.com/go-sql-driver/mysql"
)

func main() {
	dbConn, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true})
	if err != nil {
		log.Fatalf("Failed to conenct to database: %v", err.Error())
	}

	initStorage(dbConn)

	server := api.NewAPIServer(":8080", dbConn)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
