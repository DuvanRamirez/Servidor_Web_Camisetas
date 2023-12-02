package main

import (
	"log"
	"net/http"

	"github.com/DuvanRamirez/serividor_de_camisetas/controllers"
	"github.com/DuvanRamirez/serividor_de_camisetas/handlers2"
	"github.com/DuvanRamirez/serividor_de_camisetas/models"
	repositorio "github.com/DuvanRamirez/serividor_de_camisetas/repository" 
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/gorilla/handlers"
)

func ConectarDB(url, driver string) (*sqlx.DB, error) {
	pgUrl, _ := pq.ParseURL(url)
	db, err := sqlx.Connect(driver, pgUrl) // driver: postgres
	if err != nil {
		log.Printf("fallo la conexion a PostgreSQL, error: %s", err.Error())
		return nil, err
	}

	log.Printf("Nos conectamos bien a la base de datos db: %#v", db)
	return db, nil
}

func main() {

	
	db, err := ConectarDB("postgres://qsmomarz:W6VhelAgEODDlCxLeXSTD9QTetktqNF3@batyr.db.elephantsql.com/qsmomarz", "postgres")
	if err != nil {
		log.Fatalln("error conectando a la base de datos", err.Error())
		return
	}

	repo, err := repositorio.NewRepository[models.Camiseta](db)
	if err != nil {
		log.Fatalln("fallo al crear una instancia de repositorio", err.Error())
		return
	}

	controller, err := controllers.NewController(repo)
	if err != nil {
		log.Fatalln("fallo al crear una instancia de controller", err.Error())
		return
	}

	handler, err := handlers2.NewHandler(controller)
	if err != nil {
		log.Fatalln("fallo al crear una instancia de handler", err.Error())
		return
	}

	
	router := mux.NewRouter()

	router.Handle("/camisetas", http.HandlerFunc(handler.Leercamisetas)).Methods(http.MethodGet)
	router.Handle("/camisetas", http.HandlerFunc(handler.CrearCamiseta)).Methods(http.MethodPost)
	router.Handle("/camisetas/{id}", http.HandlerFunc(handler.LeerUnaCamiseta)).Methods(http.MethodGet)
	router.Handle("/camisetas/{id}", http.HandlerFunc(handler.ActualizarUnaCamiseta)).Methods(http.MethodPatch)
	router.Handle("/camisetas/{id}", http.HandlerFunc(handler.EliminarUnaCamiseta)).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET","POST","PATCH","DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"}),
	)(router))

}