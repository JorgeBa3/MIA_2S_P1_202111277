package analyzer

import (
	"errors"
	"fmt"
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
				results = append(results, "Ejecutando comando mkdisk")
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
