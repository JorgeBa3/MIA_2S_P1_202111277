package structures

import (
	"time" // Para manejo de fechas y horas
)

type SuperBlock struct {
	s_filesystem_type   int32     // Número que identifica el sistema de archivos utilizado
	s_inodes_count      int32     // Número total de inodos
	s_blocks_count      int32     // Número total de bloques
	s_free_blocks_count int32     // Número de bloques libres
	s_free_inodes_count int32     // Número de inodos libres
	s_mtime             time.Time // Última fecha en que el sistema fue montado
	s_umtime            time.Time // Última fecha en que el sistema fue desmontado
	s_mnt_count         int32     // Número de veces que se ha montado el sistema
	s_magic             int32     // Valor que identifica el sistema de archivos (0xEF53)
	s_inode_s           int32     // Tamaño del inodo
	s_block_s           int32     // Tamaño del bloque
	s_first_ino         int32     // Dirección del primer inodo libre
	s_first_blo         int32     // Dirección del primer bloque libre
	s_bm_inode_start    int32     // Inicio del bitmap de inodos
	s_bm_block_start    int32     // Inicio del bitmap de bloques
	s_inode_start       int32     // Inicio de la tabla de inodos
	s_block_start       int32     // Inicio de la tabla de bloques
}
