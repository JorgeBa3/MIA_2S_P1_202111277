package structures

import "fmt"

type PARTITION struct {
	part_status      [1]byte  // Estado de la partición
	part_type        [1]byte  // Tipo de partición
	part_fit         [1]byte  // Ajuste de la partición
	part_start       int32    // Byte de inicio de la partición
	part_s           int32    // Tamaño de la partición
	part_name        [16]byte // Nombre de la partición
	part_correlative int32    // Correlativo de la partición
	part_id          [4]byte  // ID de la partición
}

/*
Part Status:
	9: Disponible
	0: Creado
	1: Montado

Esto queda a su criterio.
*/

// Crear una partición con los parámetros proporcionados
func (p *PARTITION) CreatePartition(partStart, partSize int, partType, partFit, partName string) {
	// Asignar status de la partición
	p.part_status[0] = '0' // El valor '0' indica que la partición ha sido creada

	// Asignar el byte de inicio de la partición
	p.part_start = int32(partStart)

	// Asignar el tamaño de la partición
	p.part_s = int32(partSize)

	// Asignar el tipo de partición
	if len(partType) > 0 {
		p.part_type[0] = partType[0]
	}

	// Asignar el ajuste de la partición
	if len(partFit) > 0 {
		p.part_fit[0] = partFit[0]
	}

	// Asignar el nombre de la partición
	copy(p.part_name[:], partName)
}

func (p *PARTITION) MountPartition(correlative int, id string) error {
	// Asignar correlativo a la partición
	p.part_correlative = int32(correlative)

	// Asignar ID a la partición
	copy(p.part_id[:], id)

	return nil
}

// Imprimir los valores de la partición
func (p *PARTITION) Print() {
	fmt.Printf("part_status: %c\n", p.part_status[0])
	fmt.Printf("part_type: %c\n", p.part_type[0])
	fmt.Printf("part_fit: %c\n", p.part_fit[0])
	fmt.Printf("part_start: %d\n", p.part_start)
	fmt.Printf("part_s: %d\n", p.part_s)
	fmt.Printf("part_name: %s\n", string(p.part_name[:]))
	fmt.Printf("part_correlative: %d\n", p.part_correlative)
	fmt.Printf("part_id: %s\n", string(p.part_id[:]))
}
