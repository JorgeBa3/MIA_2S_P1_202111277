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

		// Removemos el prefijo '-' si existe
		if strings.HasPrefix(token, "-") {
			token = strings.TrimPrefix(token, "-")
		}

		// Separar en clave=valor
		tk := strings.Split(token, "=")
		if len(tk) == 2 {
			if tk[0] == "id" {
				id = tk[1]
			} else if tk[0] == "user" {
				user = tk[1]
			} else if tk[0] == "pass" {
				pass = tk[1]
			}
		}
	}

	// Verificamos si todos los valores requeridos fueron proporcionados
	if id == "" || user == "" || pass == "" {
		//Error("LOGIN", "Se necesitan parámetros obligatorios para el comando LOGIN.")
		return false
	}
	return true
}
