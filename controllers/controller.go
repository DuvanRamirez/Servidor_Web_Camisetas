package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/DuvanRamirez/serividor_de_camisetas/models"
	repositorio "github.com/DuvanRamirez/serividor_de_camisetas/repository"
)

var (
	updateQuery = "UPDATE camisetas SET %s WHERE id=:id;"
	deleteQuery = "DELETE FROM camisetas WHERE id=$1;"
	selectQuery = "SELECT id, tipo, color, talla, marca, equipo, foto FROM camisetas WHERE id=$1;"
	listQuery   = "SELECT id, tipo, color, talla, marca, equipo, foto FROM camisetas limit $1 offset $2"
	createQuery = "INSERT INTO camisetas (tipo, color, talla, marca, equipo, foto) VALUES (:tipo, :color, :talla, :marca, :equipo, :foto) returning id;"
)

type Controller struct {
	repo repositorio.Repository[models.Camiseta]
}

func NewController(repo repositorio.Repository[models.Camiseta]) (*Controller, error) {
	if repo == nil {
		return nil, fmt.Errorf("para instanciar un controlador se necesita un repositorio no nulo")
	}
	return &Controller{
		repo: repo,
	}, nil
}

func (c *Controller) ActualizarUnaCamiseta(reqBody []byte, id string) error {
	nuevasValoresCamiseta := make(map[string]any)
	err := json.Unmarshal(reqBody, &nuevasValoresCamiseta)
	if err != nil {
		log.Printf("fallo al actualizar una Camiseta, con error: %s", err.Error())
		return fmt.Errorf("fallo al actualizar una Camiseta, con error: %s", err.Error())
	}

	if len(nuevasValoresCamiseta) == 0 {
		log.Printf("fallo al actualizar una Camiseta, con error: %s", err.Error())
		return fmt.Errorf("fallo al actualizar una Camiseta, con error: %s", err.Error())
	}

	query := construirUpdateQuery(nuevasValoresCamiseta)
	nuevasValoresCamiseta["id"] = id
	err = c.repo.Update(context.TODO(), query, nuevasValoresCamiseta)
	if err != nil {
		log.Printf("fallo al actualizar una Camiseta, con error: %s", err.Error())
		return fmt.Errorf("fallo al actualizar una Camiseta, con error: %s", err.Error())
	}
	return nil
}

func construirUpdateQuery(nuevasValores map[string]any) string {
	columns := []string{}
	for key := range nuevasValores {
		columns = append(columns, fmt.Sprintf("%s=:%s", key, key))
	}
	columnsString := strings.Join(columns, ",")
	return fmt.Sprintf(updateQuery, columnsString)
}

func (c *Controller) EliminarUnaCamiseta(id string) error {
	err := c.repo.Delete(context.TODO(), deleteQuery, id)
	if err != nil {
		log.Printf("fallo al eliminar una Camiseta, con error: %s", err.Error())
		return fmt.Errorf("fallo al eliminar una Camiseta, con error: %s", err.Error())
	}
	return nil
}

func (c *Controller) LeerUnaCamiseta(id string) ([]byte, error) {
	Camiseta, err := c.repo.Read(context.TODO(), selectQuery, id)
	if err != nil {
		log.Printf("fallo al leer una Camiseta, con error: %s", err.Error())
		return nil, fmt.Errorf("fallo al leer una Camiseta, con error: %s", err.Error())
	}

	CamisetaJson, err := json.Marshal(Camiseta)
	if err != nil {
		log.Printf("fallo al leer una Camiseta, con error: %s", err.Error())
		return nil, fmt.Errorf("fallo al leer una Camiseta, con error: %s", err.Error())
	}
	return CamisetaJson, nil
}

func (c *Controller) Leercamisetas(limit, offset int) ([]byte, error) {
	camisetas, _, err := c.repo.List(context.TODO(), listQuery, limit, offset)
	if err != nil {
		log.Printf("fallo al leer camisetas, con error: %s", err.Error())
		return nil, fmt.Errorf("fallo al leer camisetas, con error: %s", err.Error())
	}

	jsoncamisetas, err := json.Marshal(camisetas)
	if err != nil {
		log.Printf("fallo al leer camisetas, con error: %s", err.Error())
		return nil, fmt.Errorf("fallo al leer camisetas, con error: %s", err.Error())
	}
	return jsoncamisetas, nil
}

func (c *Controller) CrearCamiseta(reqBody []byte) (int64, error) {
	nuevaCamiseta := &models.Camiseta{}
	err := json.Unmarshal(reqBody, nuevaCamiseta)
	if err != nil {
		log.Printf("fallo al crear una nueva Camiseta, con error: %s", err.Error())
		return 0, fmt.Errorf("fallo al crear una nueva Camiseta, con error: %s", err.Error())
	}

	valoresColumnasnuevaCamiseta := map[string]any{
		"tipo":     nuevaCamiseta.Tipo,
		"color":    nuevaCamiseta.Color,
		"talla": 	nuevaCamiseta.Talla,
		"marca":  	nuevaCamiseta.Marca,
		"equipo":   nuevaCamiseta.Equipo,
		"foto":		nuevaCamiseta.Foto,
	}

	nuevaId, err := c.repo.Create(context.TODO(), createQuery, valoresColumnasnuevaCamiseta)
	if err != nil {
		log.Printf("fallo al crear un nueva Camiseta, con error: %s", err.Error())
		return 0, fmt.Errorf("fallo al crear un nueva Camiseta, con error: %s", err.Error())
	}
	return nuevaId, nil
}