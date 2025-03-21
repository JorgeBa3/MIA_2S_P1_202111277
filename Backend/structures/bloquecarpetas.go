package structures

// Content define el contenido de una carpeta, con su nombre y el inodo al que apunta.
type Content struct {
	B_name  [12]byte // Nombre de la carpeta o archivo
	B_inodo int32    // Apuntador hacia un inodo asociado al archivo o carpeta
}

// BloqueCarpeta define el bloque de carpeta que contiene un array de Content.
type BloqueCarpeta struct {
	B_content [4]Content // Array con el contenido de la carpeta
}
