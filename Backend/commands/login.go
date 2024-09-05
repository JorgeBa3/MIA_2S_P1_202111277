package commands

import (
	"strings"
)

type UsuarioActivo struct {
	User     string
	Password string
	Id       string
	Uid      int
	Gid      int
}

var Logged UsuarioActivo

func ValidarDatosLOGIN(context []string) bool {
	id := ""
	user := ""
	pass := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if tk[0] == "id" {
			id = tk[1]
		} else if tk[0] == "usuario" {
			user = tk[1]
		} else if tk[0] == "password" {
			pass = tk[1]
		}
	}
	if id == "" || user == "" || pass == "" {
		//Error("LOGIN", "Se necesitan parámetros obligatorios para el comando LOGIN.")
		return false
	}
	return true
}
