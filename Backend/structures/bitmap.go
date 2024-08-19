package structures

// Bitmap define un mapa de bits para inodos o bloques, utilizando un array de bytes.
type Bitmap struct {
	Bits []byte // Cada bit en este array representa el estado de un inodo o bloque.
}

// NewBitmap crea un nuevo bitmap con el tamaño especificado.
func NewBitmap(size int) *Bitmap {
	return &Bitmap{
		Bits: make([]byte, size),
	}
}

// SetBit establece el bit en el índice dado a 1 (ocupado).
func (b *Bitmap) SetBit(index int) {
	b.Bits[index] = 1
}

// ClearBit establece el bit en el índice dado a 0 (libre).
func (b *Bitmap) ClearBit(index int) {
	b.Bits[index] = 0
}

// GetBit devuelve el estado del bit en el índice dado.
func (b *Bitmap) GetBit(index int) byte {
	return b.Bits[index]
}
