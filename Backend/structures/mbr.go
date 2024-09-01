package structures

import (
	"bytes"           // Paquete para manipulación de buffers
	"encoding/binary" // Paquete para codificación y decodificación de datos binarios
	"encoding/json"
	"errors"
	"fmt" // Paquete para formateo de E/S
	"os"  // Paquete para funciones del sistema operativo
	"strings"
	"time" // Paquete para manipulación de tiempo
)

type MBR struct {
	Mbr_tamano         int32        // Tamaño del MBR en bytes
	Mbr_fecha_creacion float32      // Fecha y hora de creación del MBR
	Mbr_dsk_signature  int32        // Firma del disco
	Dsk_fit            [1]byte      // Tipo de ajuste
	Mbr_partitions     [4]PARTITION // Particiones del MBR
}

// Convertir MBR a JSON
func (mbr *MBR) ToJSON() (string, error) {
	creationTime := time.Unix(int64(mbr.Mbr_fecha_creacion), 0)
	diskFit := string(mbr.Dsk_fit[0])

	mbrJSON := map[string]interface{}{
		"tamano":            mbr.Mbr_tamano,
		"fit":               diskFit,
		"path":              "", // Se debe agregar la lógica para obtener la ubicación
		"fecha de creacion": creationTime.Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(mbrJSON, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// SerializeMBR escribe la estructura MBR al inicio de un archivo binario
func (mbr *MBR) Serialize(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeMBR lee la estructura MBR desde el inicio de un archivo binario
func (mbr *MBR) Deserialize(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener el tamaño de la estructura MBR
	mbrSize := binary.Size(mbr)
	if mbrSize <= 0 {
		return fmt.Errorf("invalid MBR size: %d", mbrSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura MBR
	buffer := make([]byte, mbrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura MBR
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// Método para obtener la primera partición disponible
func (mbr *MBR) GetFirstAvailablePartition() (*PARTITION, int, int) {
	// Calcular el offset para el start de la partición
	offset := binary.Size(mbr) // Tamaño del MBR en bytes

	// Recorrer las particiones del MBR
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Si el start de la partición es -1, entonces está disponible
		if mbr.Mbr_partitions[i].Part_start == -1 {
			// Devolver la partición, el offset y el índice
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(mbr.Mbr_partitions[i].Part_s)
		}
	}
	return nil, -1, -1
}

// Función para obtener una partición por ID
func (mbr *MBR) GetPartitionByID(id string) (*PARTITION, error) {
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionID := strings.Trim(string(mbr.Mbr_partitions[i].Part_id[:]), "\x00 ")
		// Convertir el id a string y eliminar los caracteres nulos
		inputID := strings.Trim(id, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición
		if strings.EqualFold(partitionID, inputID) {
			return &mbr.Mbr_partitions[i], nil
		}
	}
	return nil, errors.New("partición no encontrada")
}

// Método para obtener una partición por nombre
func (mbr *MBR) GetPartitionByName(name string) (*PARTITION, int) {
	// Recorrer las particiones del MBR
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionName := strings.Trim(string(partition.Part_name[:]), "\x00 ")
		// Convertir el nombre de la partición a string y eliminar los caracteres nulos
		inputName := strings.Trim(name, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionName, inputName) {
			return &partition, i
		}
	}
	return nil, -1
}

// Método para imprimir los valores del MBR
func (mbr *MBR) Print() {
	// Convertir Mbr_fecha_creacion a time.Time
	creationTime := time.Unix(int64(mbr.Mbr_fecha_creacion), 0)

	// Convertir Dsk_fit a char
	diskFit := rune(mbr.Dsk_fit[0])

	fmt.Printf("MBR Size: %d\n", mbr.Mbr_tamano)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_dsk_signature)
	fmt.Printf("Disk Fit: %c\n", diskFit)
}

// Método para imprimir las particiones del MBR
func (mbr *MBR) PrintPartitions() {
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.Part_status[0])
		partType := rune(partition.Part_type[0])
		partFit := rune(partition.Part_fit[0])

		// Convertir Part_name a string
		partName := string(partition.Part_name[:])
		// Convertir Part_id a string
		partID := string(partition.Part_id[:])

		fmt.Printf("Partition %d:\n", i+1)
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.Part_start)
		fmt.Printf("  Size: %d\n", partition.Part_s)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.Part_correlative)
		fmt.Printf("  ID: %s\n", partID)
	}
}
