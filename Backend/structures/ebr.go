package structures

import "fmt"

type EBR struct {
	part_mount [1]byte  // Indica si la partición está montada o no
	part_fit   [1]byte  // Tipo de ajuste de la partición (B, F, W)
	part_start int32    // Byte de inicio de la partición
	part_s     int32    // Tamaño total de la partición en bytes
	part_next  int32    // Byte en el que está el próximo EBR (-1 si no hay siguiente)
	part_name  [16]byte // Nombre de la partición
}

// Crear una partición extendida con los parámetros proporcionados
func (e *EBR) CreateEBR(partStart, partSize, partNext int, partFit, partName string) {
	// Asignar el byte de inicio de la partición
	e.part_start = int32(partStart)

	// Asignar el tamaño de la partición
	e.part_s = int32(partSize)

	// Asignar el byte del siguiente EBR
	e.part_next = int32(partNext)

	// Asignar el tipo de ajuste de la partición
	if len(partFit) > 0 {
		e.part_fit[0] = partFit[0]
	}

	// Asignar el nombre de la partición
	copy(e.part_name[:], partName)

	// Inicialmente la partición no está montada
	e.part_mount[0] = '0'
}

// Montar una partición extendida
func (e *EBR) MountEBR() {
	e.part_mount[0] = '1' // Indica que la partición está montada
}

// Imprimir los valores de la partición extendida
func (e *EBR) Print() {
	fmt.Printf("part_mount: %c\n", e.part_mount[0])
	fmt.Printf("part_fit: %c\n", e.part_fit[0])
	fmt.Printf("part_start: %d\n", e.part_start)
	fmt.Printf("part_s: %d\n", e.part_s)
	fmt.Printf("part_next: %d\n", e.part_next)
	fmt.Printf("part_name: %s\n", string(e.part_name[:]))
}
