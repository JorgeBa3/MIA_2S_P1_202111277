package structures

// BloqueApuntadores define el bloque de apuntadores que contiene un array de 16 int.
type BloqueApuntadores struct {
	b_pointers [16]int32 // Array con los apuntadores a bloques (de archivo o carpeta)
}
