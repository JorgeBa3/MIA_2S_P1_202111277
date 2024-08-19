package structures

import (
	"bytes"           // Paquete para manipulación de buffers
	"encoding/binary" // Paquete para codificación y decodificación de datos binarios
	"fmt"             // Paquete para formateo de E/S
	"os"              // Paquete para funciones del sistema operativo
	"strings"
	"time" // Paquete para manipulación de tiempo
)

type MBR struct {
	mbr_tamano         int32        // Tamaño del MBR en bytes
	mbr_fecha_creacion float32      // Fecha y hora de creación del MBR
	mbr_dsk_signature  int32        // Firma del disco
	dsk_fit            [1]byte      // Tipo de ajuste
	mbr_partitions     [4]PARTITION // Particiones del MBR
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
	for i := 0; i < len(mbr.mbr_partitions); i++ {
		// Si el start de la partición es -1, entonces está disponible
		if mbr.mbr_partitions[i].part_start == -1 {
			// Devolver la partición, el offset y el índice
			return &mbr.mbr_partitions[i], offset, i
		} else {
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(mbr.mbr_partitions[i].part_s)
		}
	}
	return nil, -1, -1
}

// Método para obtener una partición por nombre
func (mbr *MBR) GetPartitionByName(name string) (*PARTITION, int) {
	// Recorrer las particiones del MBR
	for i, partition := range mbr.mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionName := strings.Trim(string(partition.part_name[:]), "\x00 ")
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
	// Convertir mbr_fecha_creacion a time.Time
	creationTime := time.Unix(int64(mbr.mbr_fecha_creacion), 0)

	// Convertir dsk_fit a char
	diskFit := rune(mbr.dsk_fit[0])

	fmt.Printf("MBR Size: %d\n", mbr.mbr_tamano)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.mbr_dsk_signature)
	fmt.Printf("Disk Fit: %c\n", diskFit)
}

// Método para imprimir las particiones del MBR
func (mbr *MBR) PrintPartitions() {
	for i, partition := range mbr.mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.part_status[0])
		partType := rune(partition.part_type[0])
		partFit := rune(partition.part_fit[0])

		// Convertir Part_name a string
		partName := string(partition.part_name[:])
		// Convertir Part_id a string
		partID := string(partition.part_id[:])

		fmt.Printf("Partition %d:\n", i+1)
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.part_start)
		fmt.Printf("  Size: %d\n", partition.part_s)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.part_correlative)
		fmt.Printf("  ID: %s\n", partID)
	}
}
