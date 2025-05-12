package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iangechuki/go_carzone/models"
	"github.com/iangechuki/go_carzone/service"
)

type EngineHandler struct {
	engineService service.EngineServiceInterface
}

func NewEngineHandler(engineService service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		engineService: engineService,
	}
}

func (h *EngineHandler)GetEngineByID(w http.ResponseWriter,r *http.Request){
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	engine,err := h.engineService.GetEngineByID(ctx,id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	body,err := json.Marshal(engine)
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

func (h *EngineHandler)CreateEngine(w http.ResponseWriter,r *http.Request){
	ctx := r.Context()

	body,err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	var engineReq models.EngineRequest
	err = json.Unmarshal(body,&engineReq)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	createdEngine,err := h.engineService.CreateEngine(ctx,&engineReq)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	
	responseBody ,err := json.Marshal(createdEngine)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error while marshallin",err)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err = w.Write(responseBody)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error writing messages ",err)
		return
	}
}

func (h *EngineHandler)DeleteEngine(w http.ResponseWriter,r *http.Request){
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	deletedEngine,err := h.engineService.DeleteEngine(ctx,id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error deleting engine: ",err)
		return
	}
	if err := json.NewEncoder(w).Encode(deletedEngine); err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
}	

func (h *EngineHandler)UpdateEngine(w http.ResponseWriter,r *http.Request){
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var engineReq models.EngineRequest
	if err:= json.NewDecoder(r.Body).Decode(&engineReq); err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		log.Println("Error decoding req: ",err)
		return
	}
	updatedEngine,err := h.engineService.UpdateEngine(ctx,id,&engineReq)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error updating engine: ",err)
		return
	}
	if err := json.NewEncoder(w).Encode(updatedEngine); err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}

}