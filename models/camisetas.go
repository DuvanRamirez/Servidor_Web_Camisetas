package models

/*
es mejor conservar un estándar entre las etiquetas de json y db para no tener problemas al parsear
de json a db en el método ActualizarUnCamiseta
*/
type Camiseta struct {
	Id          int    `db:"id" json:"id"`
	Tipo      	string `db:"tipo" json:"tipo"`
	Color       string   `db:"color" json:"color"`
	Talla 		string `db:"talla" json:"talla"`
	Marca  		string   `db:"marca" json:"marca"`
	Equipo   	string `db:"equipo" json:"equipo"`
	Foto 		string `db:"foto" json:"foto"`
}