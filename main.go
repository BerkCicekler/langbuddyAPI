package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User: "root",
		Passwd: "",
		Addr: "127.0.0.1:3306",
		DBName: "golang",
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	}

	sqlStorage := NewMySqlStorage(cfg)

	db, err := sqlStorage.Init()

	if err != nil {
		log.Fatal(err)
	}

	err = initFirebase()

	roomMap = make( map[string]*room)

	if err != nil {
		log.Fatal(err)
	}

	store := NewStorage(db)
	api := NewAPIServer(":8080", store)
	api.Serve()
}