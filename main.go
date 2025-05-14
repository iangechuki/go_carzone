package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/iangechuki/go_carzone/driver"
	carHandler "github.com/iangechuki/go_carzone/handler/car"
	engineHandler "github.com/iangechuki/go_carzone/handler/engine"
	loginHandler "github.com/iangechuki/go_carzone/handler/login"
	"github.com/iangechuki/go_carzone/middleware"
	carService "github.com/iangechuki/go_carzone/service/car"
	engineService "github.com/iangechuki/go_carzone/service/engine"
	carStore "github.com/iangechuki/go_carzone/store/car"
	engineStore "github.com/iangechuki/go_carzone/store/engine"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelmux "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)


func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ",err)
	}
	traceProvider,err := startTracing()
	if err != nil {
		log.Fatal("Error starting tracing: ",err)
	}
	defer func(){
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatal("Error shutting down tracing: ",err)
		}
	}()

	otel.SetTracerProvider(traceProvider)

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

	router.Use(otelmux.Middleware("CarZone"))
	
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
	router.HandleFunc("/login",loginHandler.LoginHandler).Methods("POST")
	
	protected := router.PathPrefix("/").Subrouter()

	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/cars/{id}",carHandler.GetCarByID).Methods("GET")
	protected.HandleFunc("/cars",carHandler.GetCarByBrand).Methods("GET")
	protected.HandleFunc("/cars",carHandler.CreateCar).Methods("POST")
	protected.HandleFunc("/cars/{id}",carHandler.UpdateCar).Methods("PUT")
	protected.HandleFunc("/cars/{id}",carHandler.DeleteCar).Methods("DELETE")

	protected.HandleFunc("/engines/{id}",engineHandler.GetEngineByID).Methods("GET")
	protected.HandleFunc("/engines",engineHandler.CreateEngine).Methods("POST")
	protected.HandleFunc("/engines/{id}",engineHandler.UpdateEngine).Methods("PUT")
	protected.HandleFunc("/engines/{id}",engineHandler.DeleteEngine).Methods("DELETE")
	
	router.Handle("/metrics",promhttp.Handler())
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

func startTracing()(*trace.TracerProvider,error){
	header := map[string]string{
		"Content-Type":"application/json",

	}
	exporter,err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(header),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil,err
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay * time.Millisecond),
			),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("carzone"),
		),
	),
	)

		
	
	return traceProvider,nil
}