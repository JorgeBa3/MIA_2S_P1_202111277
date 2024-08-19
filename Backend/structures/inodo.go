package structures

import (
	"time" // Para manejo de fechas y horas
)

type Inode struct {
	i_uid   int32     // UID del usuario propietario del archivo o carpeta
	i_gid   int32     // GID del grupo al que pertenece el archivo o carpeta
	i_s     int32     // Tamaño del archivo en bytes
	i_atime time.Time // Última fecha en que se leyó el inodo sin modificarlo
	i_ctime time.Time // Fecha en la que se creó el inodo
	i_mtime time.Time // Última fecha en la que se modificó el inodo
	i_block [15]int32 // Array para los bloques directos e indirectos
	i_type  byte      // Indica si es archivo o carpeta (1 = Archivo, 0 = Carpeta)
	i_perm  [3]byte   // Guardará los permisos del archivo o carpeta (UGO en forma octal)
}
