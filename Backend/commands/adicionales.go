package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"main/structures"
	"os"
	"strings"
	"unsafe"
)

func Comparar(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

func Error(op string, mensaje string) {
	fmt.Println("\tERROR: " + op + "\n\tTIPO: " + mensaje)
}

func Mensaje(op string, mensaje string) {
	fmt.Println("COMANDO: " + op + ";\nMENSAJE: " + mensaje)
}

func Confirmar(mensaje string) bool {
	fmt.Println(mensaje + " (y/n)")
	var respuesta string
	fmt.Scanln(&respuesta)
	if Comparar(respuesta, "y") {
		return true
	}
	return false
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func leerDisco(path string) *structures.MBR {
	m := structures.MBR{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	data := leerBytes(file, int(unsafe.Sizeof(structures.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &m)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return nil
	}
	var mDir *structures.MBR = &m
	return mDir
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func BuscarParticiones(mbr structures.MBR, name string, path string) *structures.PARTITION {
	var particiones [4]structures.PARTITION
	particiones[0] = mbr.Mbr_partitions[0]
	particiones[1] = mbr.Mbr_partitions[1]
	particiones[2] = mbr.Mbr_partitions[2]
	particiones[3] = mbr.Mbr_partitions[3]

	ext := false
	extended := structures.PARTITION{}
	for i := 0; i < len(particiones); i++ {
		particion := particiones[i]
		if particion.Part_status == [1]byte{'1'} {
			nombre := ""
			for j := 0; j < len(particion.Part_name); j++ {
				if particion.Part_name[j] != 0 {
					nombre += string(particion.Part_name[j])
				}
			}
			if Comparar(nombre, name) {
				return &particion
			} else if particion.Part_type == [1]byte{'E'} || particion.Part_type == [1]byte{'e'} {
				ext = true
				extended = particion
			}
		}
	}

	if ext {
		ebrs := GetLogicas(extended, path)
		for i := 0; i < len(ebrs); i++ {
			ebr := ebrs[i]
			if ebr.Part_s == '1' {
				nombre := ""
				for j := 0; j < len(ebr.Part_name); j++ {
					if ebr.Part_name[j] != 0 {
						nombre += string(ebr.Part_name[j])
					}
				}
				if Comparar(nombre, name) {
					tmp := structures.PARTITION{}
					tmp.Part_status = [1]byte{'1'}
					tmp.Part_type = [1]byte{'L'}
					tmp.Part_fit = ebr.Part_fit
					tmp.Part_start = ebr.Part_start
					tmp.Part_s = ebr.Part_s
					copy(tmp.Part_name[:], ebr.Part_name[:])
					return &tmp
				}
			}
		}
	}
	return nil
}
func GetLogicas(particion structures.PARTITION, path string) []structures.EBR {
	var ebrs []structures.EBR
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	tmp := structures.EBR{}
	file.Seek(int64(particion.Part_start), 0)

	data := leerBytes(file, int(unsafe.Sizeof(structures.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return nil
	}
	for {
		if int(tmp.Part_next) != -1 && int(tmp.Part_s) != 0 {
			ebrs = append(ebrs, tmp)
			file.Seek(int64(tmp.Part_next), 0)

			data = leerBytes(file, int(unsafe.Sizeof(structures.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &tmp)
			if err_ != nil {
				Error("FDSIK", "Error al leer el archivo")
				return nil
			}
		} else {
			file.Close()
			break
		}
	}

	return ebrs
}
