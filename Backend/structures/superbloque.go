package structures

import (
	"time" // Para manejo de fechas y horas
)

type SuperBlock struct {
	S_filesystem_type   int32     // Número que identifica el sistema de archivos utilizado
	S_inodes_count      int32     // Número total de inodos
	S_blocks_count      int32     // Número total de bloques
	S_free_blocks_count int32     // Número de bloques libres
	S_free_inodes_count int32     // Número de inodos libres
	S_mtime             time.Time // Última fecha en que el sistema fue montado
	S_umtime            time.Time // Última fecha en que el sistema fue desmontado
	S_mnt_count         int32     // Número de veces que se ha montado el sistema
	S_magic             int32     // Valor que identifica el sistema de archivos (0xEF53)
	S_inode_s           int32     // Tamaño del inodo
	S_block_s           int32     // Tamaño del bloque
	S_first_ino         int32     // Dirección del primer inodo libre
	S_first_blo         int32     // Dirección del primer bloque libre
	S_bm_inode_start    int32     // Inicio del bitmap de inodos
	S_bm_block_start    int32     // Inicio del bitmap de bloques
	S_inode_start       int32     // Inicio de la tabla de inodos
	S_block_start       int32     // Inicio de la tabla de bloques
}
