package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	// Para manejo de fechas y horas
)

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_inodes_count int32
	S_free_blocks_count int32
	S_mtime             float32
	S_umtime            float32
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
	// Total: 68 bytes
}

// Serialize escribe la estructura SuperBlock en un archivo binario en la posición especificada
func (sb *SuperBlock) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura SuperBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize lee la estructura SuperBlock desde un archivo binario en la posición especificada
func (sb *SuperBlock) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Obtener el tamaño de la estructura SuperBlock
	sbSize := binary.Size(sb)
	if sbSize <= 0 {
		return fmt.Errorf("invalid SuperBlock size: %d", sbSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura SuperBlock
	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura SuperBlock
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}
