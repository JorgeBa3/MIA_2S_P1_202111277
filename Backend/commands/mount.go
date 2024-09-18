package commands

import (
	"bytes"
	"encoding/binary"
	structures "main/structures"
	global "main/utils"
	utils "main/utils"
	"os"
	"strconv"
	"unsafe"

	"errors" // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"    // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"regexp" // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas

	// Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MOUNT struct {
	path string // Ruta del archivo del disco
	name string // Nombre de la partición
}

/*
	mount -path=/home/Disco1.mia -name=Part1 #id=341a
	mount -path=/home/Disco2.mia -name=Part1 #id=342a
	mount -path=/home/Disco3.mia -name=Part2 #id=343a
*/

// CommandMount parsea el comando mount y devuelve una instancia de MOUNT
func ParserMount(tokens []string) (string, error) {
	cmd := &MOUNT{} // Crea una nueva instancia de MOUNT

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`(?i)-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)
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

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -path y -name hayan sido proporcionados
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	err := commandMount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Generar el ID de partición para el resultado final
	idPartition, err := GenerateIdPartition(cmd, 0) // En lugar de 0, debes pasar el índice real de la partición montada si lo tienes disponible
	if err != nil {
		return "", err
	}

	// Devuelve el nombre de la partición junto con su ID
	result := fmt.Sprintf("Partición: %s montada correctamente, con ID: %s", cmd.name, idPartition)
	return result, nil
}

func commandMount(mount *MOUNT) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(mount.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Buscar la partición con el nombre especificado
	partition, indexPartition := mbr.GetPartitionByName(mount.name)
	if partition == nil {
		fmt.Println("Error: la partición no existe")
		return errors.New("la partición no existe")
	}

	// Generar un id único para la partición
	idPartition, err := GenerateIdPartition(mount, indexPartition)
	if err != nil {
		fmt.Println("Error generando el id de partición:", err)
		return err
	}

	//  Guardar la partición montada en la lista de montajes globales
	utils.MountedPartitions[idPartition] = mount.path

	// Modificamos la partición para indicar que está montada
	partition.MountPartition(indexPartition, idPartition)

	// Guardar la partición modificada en el MBR
	mbr.Mbr_partitions[indexPartition] = *partition

	// Serializar la estructura MBR en el archivo binario
	err = mbr.Serialize(mount.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return err
	}

	return nil
}

func GenerateIdPartition(mount *MOUNT, indexPartition int) (string, error) {
	// Asignar una letra a la partición
	letter, err := utils.GetLetter(mount.path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", err
	}

	// Crear id de partición
	idPartition := fmt.Sprintf("%s%d%s", global.Carnet, indexPartition+1, letter)

	return idPartition, nil
}

// CommandListMounts devuelve una lista de las particiones montadas
func CommandListMounts() string {
	// Verifica si hay particiones montadas
	if len(utils.MountedPartitions) == 0 {
		return "No hay particiones montadas."
	}

	// Variable para construir el resultado
	var result strings.Builder
	result.WriteString("Particiones montadas:\n")

	// Recorre el mapa de particiones montadas
	for id, path := range utils.MountedPartitions {
		result.WriteString(fmt.Sprintf("ID: %s, Ruta: %s\n", id, path))
	}

	// Devuelve el resultado como una cadena
	return result.String()
}

func GetMountCommand(comando string, id string, p *string) structures.PARTITION {
	if !(id[0] == '7' && id[1] == '7') {
		Error(comando, "El primer identificador no es válido.")
		return structures.PARTITION{}
	}
	letra := id[len(id)-1]
	id = strings.ReplaceAll(id, "77", "")
	i, _ := strconv.Atoi(string(id[0] - 1))
	if i < 0 {
		Error(comando, "El primer identificador no es válido.")
		return structures.PARTITION{}
	}
	for j := 0; j < 26; j++ {
		if DiscMont[i].Particiones[j].Estado == 1 {
			if DiscMont[i].Particiones[j].Letra == letra {

				path := ""
				for k := 0; k < len(DiscMont[i].Path); k++ {
					if DiscMont[i].Path[k] != 0 {
						path += string(DiscMont[i].Path[k])
					}
				}

				file, error := os.Open(strings.ReplaceAll(path, "\"", ""))
				if error != nil {
					Error(comando, "No se ha encontrado el disco")
					return structures.PARTITION{}
				}
				disk := structures.MBR{}
				file.Seek(0, 0)

				data := leerBytes(file, int(unsafe.Sizeof(structures.MBR{})))
				buffer := bytes.NewBuffer(data)
				err_ := binary.Read(buffer, binary.BigEndian, &disk)

				if err_ != nil {
					Error("FDSIK", "Error al leer el archivo")
					return structures.PARTITION{}
				}
				file.Close()

				nombreParticion := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombreParticion += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				*p = path
				return *BuscarParticiones(disk, nombreParticion, path)
			}
		}
	}
	return structures.PARTITION{}
}
