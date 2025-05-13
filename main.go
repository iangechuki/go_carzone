package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/iangechuki/go_carzone/driver"
	carHandler "github.com/iangechuki/go_carzone/handler/car"
	carService "github.com/iangechuki/go_carzone/service/car"
	carStore "github.com/iangechuki/go_carzone/store/car"

	engineHandler "github.com/iangechuki/go_carzone/handler/engine"
	engineService "github.com/iangechuki/go_carzone/service/engine"
	engineStore "github.com/iangechuki/go_carzone/store/engine"
	"github.com/joho/godotenv"
)


func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ",err)
	}
	driver.InitDB()
	defer driver.CloseDB()
	db := driver.GetDB()
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)
	carHandler := carHandler.NewCarHandler(carService)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	router := mux.NewRouter()

	schemaFile := "store/schema.sql"

	if err := executeSchemaFile(db,schemaFile); err != nil {
		log.Fatal("Error executing schema file: ",err)
	}
	router.HandleFunc("/health",func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_,err := w.Write([]byte("OK"))
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			log.Println("Error writing messages ",err)
			return
		}
	}).Methods("GET")
	router.HandleFunc("/cars/{id}",carHandler.GetCarByID).Methods("GET")
	router.HandleFunc("/cars",carHandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/cars",carHandler.CreateCar).Methods("POST")
	router.HandleFunc("/cars/{id}",carHandler.UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}",carHandler.DeleteCar).Methods("DELETE")

	router.HandleFunc("/engines/{id}",engineHandler.GetEngineByID).Methods("GET")
	router.HandleFunc("/engines",engineHandler.CreateEngine).Methods("POST")
	router.HandleFunc("/engines/{id}",engineHandler.UpdateEngine).Methods("PUT")
	router.HandleFunc("/engines/{id}",engineHandler.DeleteEngine).Methods("DELETE")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s",port)
	log.Printf("Listening on %s",addr)
	log.Fatal(http.ListenAndServe(addr,router))
}
func executeSchemaFile(db *sql.DB,fileName string)error{
	sqlFile,err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error reading schema file: ",err)
	}
	_,err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal("Error executing schema file: ",err)
	}
	return nil
}