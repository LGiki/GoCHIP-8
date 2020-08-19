package chip8

type CPU struct {
	// 8-bit general purpose register(from V0 to VF)
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
	// internal stack to store return addresses when calling procedures
	Stack [16]uint16
	// Memory
	Memory Memory
}

func (cpu *CPU) Cycle(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF
	switch opcode & 0xF000 {
	// 0NNN
	//case 0x0000:
	// 1NNN: goto NNN
	case 0x1000:
		cpu.PC = nnn
	// 2NNN: Calls subroutine at NNN
	case 0x2000:
		cpu.SP++
		cpu.Stack[cpu.SP] = cpu.PC
		cpu.PC = nnn
	// 3XNN: Skips the next instruction if VX equals NN
	case 0x3000:
		if cpu.V[x] == nn {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	// 4XNN: Skips the next instruction if VX doesn't equal NN
	case 0x4000:
		if cpu.V[x] != nn {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	// 5XY0: Skips the next instruction if VX equals VY
	case 0x5000:
		if cpu.V[x] == cpu.V[y] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	// 6XNN: Sets VX to NN
	case 0x6000:
		cpu.V[x] = nn
		cpu.PC += 2
	}

}
