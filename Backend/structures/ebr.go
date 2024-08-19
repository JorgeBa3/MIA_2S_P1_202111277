package structures

import "fmt"

type EBR struct {
	Part_mount [1]byte  // Indica si la partición está montada o no
	Part_fit   [1]byte  // Tipo de ajuste de la partición (B, F, W)
	Part_start int32    // Byte de inicio de la partición
	Part_s     int32    // Tamaño total de la partición en bytes
	Part_next  int32    // Byte en el que está el próximo EBR (-1 si no hay siguiente)
	Part_name  [16]byte // Nombre de la partición
}

// Crear una partición extendida con los parámetros proporcionados
func (e *EBR) CreateEBR(partStart, partSize, partNext int, partFit, partName string) {
	// Asignar el byte de inicio de la partición
	e.Part_start = int32(partStart)

	// Asignar el tamaño de la partición
	e.Part_s = int32(partSize)

	// Asignar el byte del siguiente EBR
	e.Part_next = int32(partNext)

	// Asignar el tipo de ajuste de la partición
	if len(partFit) > 0 {
		e.Part_fit[0] = partFit[0]
	}

	// Asignar el nombre de la partición
	copy(e.Part_name[:], partName)

	// Inicialmente la partición no está montada
	e.Part_mount[0] = '0'
}

// Montar una partición extendida
func (e *EBR) MountEBR() {
	e.Part_mount[0] = '1' // Indica que la partición está montada
}

// Imprimir los valores de la partición extendida
func (e *EBR) Print() {
	fmt.Printf("part_mount: %c\n", e.Part_mount[0])
	fmt.Printf("part_fit: %c\n", e.Part_fit[0])
	fmt.Printf("part_start: %d\n", e.Part_start)
	fmt.Printf("part_s: %d\n", e.Part_s)
	fmt.Printf("part_next: %d\n", e.Part_next)
	fmt.Printf("part_name: %s\n", string(e.Part_name[:]))
}
