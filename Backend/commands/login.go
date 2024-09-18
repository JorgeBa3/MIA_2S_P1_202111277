package commands

import (
	"fmt"
	global "main/utils"
	"os"
	"path/filepath"
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

func ParserLogin(tokens []string) (string, error) {
	if !ValidarDatosLOGIN(tokens) {
		return "", fmt.Errorf("Se necesitan parámetros obligatorios para el comando LOGIN.")
	}

	// Recorrer los parámetros
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Removemos el prefijo '-' si existe
		if strings.HasPrefix(token, "-") {
			token = strings.TrimPrefix(token, "-")
		}

		// Separar en clave=valor
		tk := strings.Split(token, "=")
		if len(tk) == 2 {
			if tk[0] == "id" {
				Logged.Id = tk[1]
			} else if tk[0] == "user" {
				Logged.User = tk[1]
			} else if tk[0] == "pass" {
				Logged.Password = tk[1]
			}
		}
	}

	// Verificar si el usuario ya está logueado
	if Logged.User != "" {
		return Logged.User, nil
	}

	// Verificar si el archivo users.txt existe y obtener el contenido
	usersData, err := ReadUsersFile(Logged.Id)
	if err != nil {
		return "", fmt.Errorf("Error al leer el archivo users.txt: %v", err)
	}

	// Verificar si el usuario y la contraseña coinciden
	if ValidarUsuario(usersData, Logged.User, Logged.Password) {
		Logged.Uid = 1 // Aquí podrías cambiar el UID basado en el archivo
		Logged.Gid = 1 // También puedes ajustar el GID según los datos del usuario
		return fmt.Sprintf("Usuario logueado correctamente %s"), nil
	}

	return "", fmt.Errorf("Usuario o contraseña incorrectos.")
}

func GetPartitionPath(partitionID string) string {
	mounted
	mountedPartition, partitionPath, err := global.GetMountedPartition(partitionID)
	if mountedPartition == true {
		print("la verdad no hace nada, pero igual me pide usarlo :/")
	}
	if err != nil {
		return "Error al obtener id"
	}

	return partitionPath
}

func ReadUsersFile(partitionID string) (string, error) {
	// Obtener el path de la partición montada
	partitionPath := GetPartitionPath(partitionID) // Debes implementar GetPartitionPath para devolver la ruta correcta

	// Ruta del archivo users.txt dentro de la partición
	usersFilePath := filepath.Join(partitionPath, "users.txt")

	// Abrir el archivo users.txt
	file, err := os.Open(usersFilePath)
	if err != nil {
		return "", fmt.Errorf("Error al abrir el archivo users.txt: %v", err)
	}
	defer file.Close()

	// Leer el contenido del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("Error al obtener la información del archivo: %v", err)
	}

	// Leer el contenido en memoria
	usersData := make([]byte, fileInfo.Size())
	_, err = file.Read(usersData)
	if err != nil {
		return "", fmt.Errorf("Error al leer el archivo: %v", err)
	}

	return string(usersData), nil
}

func ValidarUsuario(usersData string, user string, pass string) bool {
	// Supongamos que el archivo tiene líneas con el formato: "usuario,contraseña"
	lines := strings.Split(usersData, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			if parts[0] == user && parts[1] == pass {
				return true
			}
		}
	}
	return false
}
