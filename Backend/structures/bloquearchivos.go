package structures

// BloqueArchivo define el bloque de archivo que contiene el contenido del archivo.
type BloqueArchivo struct {
	b_content [64]byte // Array con el contenido del archivo
}
