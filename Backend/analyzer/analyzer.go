package analyzer

import (
	"errors"
	"fmt"
	"main/commands" // Asegúrate de importar el paquete commands
	"strings"
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (interface{}, error) {
	lines := strings.Split(input, "\n")
	var results []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			// Es un comentario, agrega la línea tal cual
			results = append(results, line)
		} else {
			tokens := strings.Fields(line)
			if len(tokens) == 0 {
				continue
			}

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
			default:
				results = append(results, fmt.Sprintf("Comando desconocido: %s", tokens[0]))
			}

		}
	}

	if len(results) == 0 {
		return nil, errors.New("no se proporcionó ningún comando o comentario")
	}

	return results, nil
}
