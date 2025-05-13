package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)
var db *sql.DB
func InitDB(){
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(connStr)
	fmt.Println("Trying to connect to db")
	var err error
	time.Sleep(5 * time.Second)
	db,err = sql.Open("postgres",connStr)
	if err != nil {
		log.Fatal("Error opening db: ",err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging db: ",err)
	}
	fmt.Println("Successfully connected to db")
}
func GetDB()*sql.DB{
	return db
}
func CloseDB(){
	err := db.Close()
	if err != nil {
		log.Fatal("Error closing db: ",err)
	}
}