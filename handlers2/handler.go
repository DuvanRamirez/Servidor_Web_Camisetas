package handlers2

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DuvanRamirez/serividor_de_camisetas/controllers"
	"github.com/gorilla/mux"
)

type Handler struct {
	controller *controllers.Controller
}

func NewHandler(controller *controllers.Controller) (*Handler, error) {
	if controller == nil {
		return nil, fmt.Errorf("para instanciar un handler se necesita un controlador no nulo")
	}
	return &Handler{
		controller: controller,
	}, nil
}

func (h *Handler) ActualizarUnaCamiseta(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("fallo al actualizar un Camiseta, con error: %s", err.Error())
		http.Error(writer, fmt.Sprintf("fallo al actualizar un Camiseta, con error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err = h.controller.ActualizarUnaCamiseta(body, id)
	if err != nil {
		log.Printf("fallo al actualizar un Camiseta, con error: %s", err.Error())
		http.Error(writer, fmt.Sprintf("fallo al actualizar un Camiseta, con error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *Handler) EliminarUnaCamiseta(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	err := h.controller.EliminarUnaCamiseta(id)
	if err != nil {
		log.Printf("fallo al eliminar un Camiseta, con error: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("fallo al eliminar un Camiseta con id %s", id)))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (h *Handler) LeerUnaCamiseta(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	Camiseta, err := h.controller.LeerUnaCamiseta(id)
	if err != nil {
		log.Printf("fallo al leer un Camiseta, con error: %s", err.Error())
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf("el Camiseta con id %s no se pudo encontrar", id)))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(Camiseta)
}

func (h *Handler) Leercamisetas(writer http.ResponseWriter, req *http.Request) {
	camisetas, err := h.controller.Leercamisetas(100, 0)
	if err != nil {
		log.Printf("fallo al leer camisetas, con error: %s", err.Error())
		http.Error(writer, "fallo al leer los camisetas", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(camisetas)
}

func (h *Handler) CrearCamiseta(writer http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("fallo al crear un nueva Camiseta, con error: %s", err.Error())
		http.Error(writer, "fallo al crear un nueva Camiseta", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	nuevaId, err := h.controller.CrearCamiseta(body)
	if err != nil {
		log.Println("fallo al crear un nueva Camiseta, con error:", err.Error())
		http.Error(writer, "fallo al crear un nueva Camiseta", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(fmt.Sprintf("id nueva Camiseta: %d", nuevaId)))
}

