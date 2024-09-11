package commands

import (
	"bytes"
	"encoding/binary"
	"main/structures"
	"main/utils"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// exec -path=/home/daniel/Escritorio/ArchivosPrueba/ArchivoEjemplo2.script
var DiscMont [99]DiscoMontado

type DiscoMontado struct {
	Path        [150]byte
	Estado      byte
	Particiones [26]ParticionMontada
}

type ParticionMontada struct {
	Letra  byte
	Estado byte
	Nombre [20]byte
}

func ValidarDatosUsers(context []string, action string) {
	usr := ""
	pwd := ""
	grp := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if utils.Comparar(tk[0], "usuario") {
			usr = tk[1]
		} else if Comparar(tk[0], "pwd") {
			pwd = tk[1]
		} else if Comparar(tk[0], "grp") {
			grp = tk[1]
		}
	}
	if Comparar(action, "MK") {
		if usr == "" || pwd == "" || grp == "" {
			Error(action+"USER", "Se necesitan parametros obligatorio para crear un usuario.")
			return
		}
		mkuser(usr, pwd, grp)
	} else if Comparar(action, "RM") {
		if usr == "" {
			Error(action+"USER", "Se necesitan parametros obligatorios para eliminar un usuario.")
			return
		}
		rmuser(usr)
	} else {
		Error(action+"USER", "No se reconoce este comando.")
		return
	}
}

func mkuser(usr string, pwd string, grp string) {
	if !Comparar(Logged.User, "root") {
		Error("MKUSER", "Solo el usuario \"root\" puede acceder a estos comandos.")
		return
	}

	var path string
	partition := GetMountCommand("MKUSER", Logged.Id, &path)
	if string(partition.Part_s) == "0" {
		Error("MKUSER", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKUSER", "No se ha encontrado el disco.")
		return
	}

	super := structures.SuperBlock{}
	file.Seek(int64(partition.Part_start), 0)
	data := leerBytes(file, int(unsafe.Sizeof(structures.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKUSER", "Error al leer el archivo")
		return
	}
	inode := structures.Inode{}
	file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(structures.Inode{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(structures.Inode{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKUSER", "Error al leer el archivo")
		return
	}

	var fb structures.BloqueArchivo
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(structures.BloqueCarpeta{}))+int64(unsafe.Sizeof(structures.BloqueArchivo{}))*int64(bloque-1), 0)

		data = leerBytes(file, int(unsafe.Sizeof(structures.BloqueArchivo{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			Error("MKUSER", "Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	//TODO REVISAR QUE EXISTA EL GRUPO
	existe := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'G' || linea[2] == 'g') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[2] == grp {
				existe = true
				break
			}
		}
	}
	if !existe {
		Error("MKUSER", "No se encontró el grupo \""+grp+"\".")
		return
	}

	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'U' || linea[2] == 'u' {
			c++
			in := strings.Split(linea, ",")
			if in[3] == usr {
				if linea[0] != '0' {
					Error("MKUSER", "EL nombre "+usr+", ya está en uso.")
					return
				}
			}
		}
	}
	txt += strconv.Itoa(c+1) + ",U," + grp + "," + usr + "," + pwd + "\n"
	tam := len(txt)
	var cadenasS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		Error("MKUSER", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKUSER", "No se ha encontrado el disco.")
		return
	}

	for i := 0; i < len(cadenasS); i++ {

		var fbAux structures.BloqueArchivo
		if inode.I_block[i] == -1 {
			file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(structures.BloqueCarpeta{}))+int64(unsafe.Sizeof(structures.BloqueArchivo{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			EscribirBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenasS[i])

		file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(structures.BloqueCarpeta{}))+int64(unsafe.Sizeof(structures.BloqueArchivo{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		EscribirBytes(file, bin6.Bytes())

	}
	for i := 0; i < len(cadenasS); i++ {
		inode.I_block[i] = int32(0)
	}
	file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(structures.Inode{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	EscribirBytes(file, inodos.Bytes())

	Mensaje("MKUSER", "Usuario "+usr+", creado correctamente!")

	file.Close()
}

func rmuser(n string) {
	if !Comparar(Logged.User, "root") {
		Error("RMUSER", "Solo el usuario \"root\" puede acceder a estos comandos.")
		return
	}

	var path string
	partition := GetMountCommand("RMUSER", Logged.Id, &path)
	if string(partition.Part_s) == "0" {
		Error("RMUSER", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMUSER", "No se ha encontrado el disco.")
		return
	}

	super := structures.SuperBlock{}
	file.Seek(int64(partition.Part_start), 0)
	data := leerBytes(file, int(unsafe.Sizeof(structures.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("RMUSER", "Error al leer el archivo")
		return
	}
	inode := structures.Inode{}
	file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(structures.Inode{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(structures.Inode{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("RMUSER", "Error al leer el archivo")
		return
	}

	var fb structures.BloqueArchivo
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(structures.BloqueCarpeta{}))+int64(unsafe.Sizeof(structures.BloqueArchivo{}))*int64(bloque-1), 0)

		data = leerBytes(file, int(unsafe.Sizeof(structures.BloqueArchivo{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			Error("RMUSER", "Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""

	vctr := strings.Split(txt, "\n")
	existe := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'U' || linea[2] == 'u') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[3] == n {
				existe = true
				aux += strconv.Itoa(0) + ",U," + in[2] + "," + in[3] + "," + in[4] + "\n"
				continue
			}
		}
		aux += linea + "\n"
	}
	if !existe {
		Error("RMUSER", "No se encontró el usuario  \""+n+"\".")
		return
	}
	txt = aux

	tam := len(txt)
	var cadenasS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		Error("RMUSER", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMUSER", "No se ha encontrado el disco.")
		return
	}

	for i := 0; i < len(cadenasS); i++ {

		var fbAux structures.BloqueArchivo
		if inode.I_block[i] == -1 {
			file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(structures.BloqueCarpeta{}))+int64(unsafe.Sizeof(structures.BloqueArchivo{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			EscribirBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenasS[i])

		file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(structures.BloqueCarpeta{}))+int64(unsafe.Sizeof(structures.BloqueArchivo{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		EscribirBytes(file, bin6.Bytes())

	}
	for i := 0; i < len(cadenasS); i++ {
		inode.I_block[i] = int32(0)
	}
	file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(structures.Inode{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	EscribirBytes(file, inodos.Bytes())

	Mensaje("RMUSER", "Usuario "+n+", eliminado correctamente!")

	file.Close()
}
