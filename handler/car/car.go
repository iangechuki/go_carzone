package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/iangechuki/go_carzone/models"
	"github.com/iangechuki/go_carzone/service"
	"go.opentelemetry.io/otel"

	"github.com/gorilla/mux"
)

type CarHandler struct {
	carService service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{
		carService: service,
	}
}

func (h *CarHandler)GetCarByID(w http.ResponseWriter,r *http.Request){
	tracer := otel.Tracer("CarHandler")
	ctx,span := tracer.Start(r.Context(), "GetCarByID-Handler")
	defer span.End()
	
	vars := mux.Vars(r)
	id := vars["id"]

	car,err := h.carService.GetCarByID(ctx,id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	body,err := json.Marshal(car)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err = w.Write(body)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error writing messages ",err)
		return
	}
}
func (h *CarHandler)GetCarByBrand(w http.ResponseWriter,r *http.Request){
	tracer := otel.Tracer("CarHandler")
	ctx,span := tracer.Start(r.Context(), "GetCarByBrand-Handler")
	defer span.End()

	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	cars,err := h.carService.GetCarsByBrand(ctx,brand,isEngine)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	body,err := json.Marshal(cars)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err = w.Write(body)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error writing messages ",err)
		return
	}
}
func (h *CarHandler)CreateCar(w http.ResponseWriter,r *http.Request){
	tracer := otel.Tracer("CarHandler")
	ctx,span := tracer.Start(r.Context(), "CreateCar-Handler")
	defer span.End()

	body,err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	var carReq models.CarRequest
	err = json.Unmarshal(body,&carReq)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	createdCar,err := h.carService.CreateCar(ctx,&carReq)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	
	responseBody ,err := json.Marshal(createdCar)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error while marshallin",err)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	_,err = w.Write(responseBody)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error writing messages ",err)
		return
	}
}
func (h *CarHandler)UpdateCar(w http.ResponseWriter,r *http.Request){
	tracer := otel.Tracer("CarHandler")
	ctx,span := tracer.Start(r.Context(), "UpdateCar-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	var carReq models.CarRequest
	if err:= json.NewDecoder(r.Body).Decode(&carReq); err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		log.Println("Error decoding req: ",err)
		return
	}
	updatedCar,err := h.carService.UpdateCar(ctx,id,&carReq)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error updating car: ",err)
		return
	}
	if err := json.NewEncoder(w).Encode(updatedCar); err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}

}
func (h *CarHandler)DeleteCar(w http.ResponseWriter,r *http.Request){
	tracer := otel.Tracer("CarHandler")
	ctx,span := tracer.Start(r.Context(), "DeleteCar-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]
	deletedCar,err := h.carService.DeleteCar(ctx,id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error deleting car: ",err)
		return
	}
	if err := json.NewEncoder(w).Encode(deletedCar); err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
}