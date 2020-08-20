package chip8

type Register struct {
	// 8-bit general purpose register (from V0 to VF)
	V [16]byte
	// Address register
	I uint16
	// Program counter
	PC uint16
	// Stack pointer
	SP byte
	// Delay timer
	DT byte
	// Sound timer
	ST byte
}
