package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ParserMkgrp(tokens []string) (string, error) {
	cmd := &MKDISK{} // Crea una nueva instancia de MKDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkdisk
	re := regexp.MustCompile(`(?i)-name=[^\s]`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Elimina las comillas del valor si están presentes
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-name":
			if value != "BF" && value != "FF" && value != "WF" {
				return "", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value

		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Crear el disco con los parámetros proporcionados
	//err = commandMkgrp(cmd)
	//if err != nil {
	//	return "", fmt.Errorf("Error al crear el disco: %v", err)
	//}

	// Construye un mensaje detallado con las especificaciones del comando ejecutado
	//result := fmt.Sprintf("Comando mkdisk ejecutado con éxito. -Tamaño: %d bytes- Ajuste: %s- Ruta: %s",
	//	sizeBytes, cmd.fit, cmd.path)

	return "", nil // Devuelve el mensaje detallado
}
