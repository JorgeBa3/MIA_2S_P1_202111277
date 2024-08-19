package structures

import (
	"time" // Para manejo de fechas y horas
)

type Inode struct {
	I_uid   int32     // UID del usuario propietario del archivo o carpeta
	I_gid   int32     // GID del grupo al que pertenece el archivo o carpeta
	I_s     int32     // Tamaño del archivo en bytes
	I_atime time.Time // Última fecha en que se leyó el inodo sin modificarlo
	I_ctime time.Time // Fecha en la que se creó el inodo
	I_mtime time.Time // Última fecha en la que se modificó el inodo
	I_block [15]int32 // Array para los bloques directos e indirectos
	I_type  byte      // Indica si es archivo o carpeta (1 = Archivo, 0 = Carpeta)
	I_perm  [3]byte   // Guardará los permisos del archivo o carpeta (UGO en forma octal)
}
