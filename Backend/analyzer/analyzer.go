package analyzer

import (
	"errors"
	"fmt"
	"main/commands"
	"strings"
)

var logued = false

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (interface{}, error) {

	lines := strings.Split(input, "\n")
	var results []string

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
					commands.Error("LOGIN", "Ya hay un usuario en línea.")
					return nil, errors.New("ya hay un usuario en línea")
				} else {
					results = append(results, "Usuario logueado correctamente")
					logued, err := commands.ParserLogin(tokens[1:])
					results = append(results, fmt.Sprintf("Usuario logueado: %s", logued))
					if err != nil {
						results = append(results, fmt.Sprintf("Error en el comando login: %s", err))
					}
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
