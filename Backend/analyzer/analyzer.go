package analyzer

import (
	"errors"
	"fmt"
	"main/commands"
	"strings"
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (interface{}, error) {

	lines := strings.Split(input, "\n")
	var results []string
	var logued = false
	var LoguedUser string = ""
	//var loguedPassword string = ""
	//var loguedId string = ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// Ignora líneas vacías o comentarios
			continue
		}

		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			// Es un comentario, agrega la línea tal cual
			results = append(results, line)
		} else {
			switch tokens[0] {
			case "mkdisk":
				result, err := commands.ParserMkdisk(tokens[1:])
				if err != nil {
					results = append(results, fmt.Sprintf("Error en el comando mkdisk: %s", err))
				} else {
					results = append(results, fmt.Sprintf("%v", result))
				}
			case "rmdisk":
				result, err := commands.ParserRmdisk(tokens[1:])
				if err != nil {
					results = append(results, fmt.Sprintf("Error en el comando rmdisk: %s", err))
				} else {
					results = append(results, result)
				}
			case "fdisk":
				result, err := commands.ParserFdisk(tokens[1:])
				if err != nil {
					results = append(results, fmt.Sprintf("Error en el comando fdisk: %s", err))
				} else {
					results = append(results, result)
				}
			case "mount":
				result, err := commands.ParserMount(tokens[1:])
				if err != nil {
					results = append(results, fmt.Sprintf("Error en el comando mount: %s", err))
				} else {
					results = append(results, result)

				}
			case "mkfs":
				// Llama a la función CommandMkfs del paquete commands con los argumentos restantes
				result, err := commands.ParserMkfs(tokens[1:])
				if err != nil {
					results = append(results, fmt.Sprintf("Error en el comando mount: %s", err))
				} else {
					results = append(results, result)
				}
			case "login":
				if logued {
					commands.Error("LOGIN", "Ya hay un usuario en línea."+LoguedUser)
					results = append(results, fmt.Sprintf("Sesion no iniciada, usuario ya logueado: %s", LoguedUser))
					println("LoguedUser: ", LoguedUser)
				} else {
					// Solo asigna LoguedUser si no hay error
					parsedUser, err := commands.ParserLogin(tokens[1:])
					if err != nil {
						results = append(results, fmt.Sprintf("Error en el comando login: %s", err))
						println("LoguedUser: ", LoguedUser) // Aquí debería ser el valor anterior
					} else {
						logued = true
						LoguedUser = parsedUser // Solo asigna aquí
						results = append(results, "Usuario logueado correctamente")
						results = append(results, fmt.Sprintf("Usuario logueado: %s", LoguedUser))
					}
				}
			case "logout":
				if logued {
					logued = false
					LoguedUser = ""
					results = append(results, "Sesión cerrada correctamente")
				} else {
					results = append(results, "No hay una sesión activa")
				}

			default:
				results = append(results, fmt.Sprintf("Comando desconocido: %s", tokens[0]))
			}
		}
	}
	if len(results) == 0 {
		return nil, errors.New("no se proporcionó ningún comando")
	}
	results = append(results, (commands.CommandListMounts()))
	return results, nil
}
