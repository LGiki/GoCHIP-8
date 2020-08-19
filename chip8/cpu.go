package chip8

type CPU struct {
	Register Register
	Memory   Memory
	// internal stack to store return addresses when calling procedures
	Stack [16]uint16
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
		cpu.Register.PC = nnn
	// 2NNN: Calls subroutine at NNN
	case 0x2000:
		cpu.Register.SP++
		cpu.Stack[cpu.Register.SP] = cpu.Register.PC
		cpu.Register.PC = nnn
	// 3XNN: Skips the next instruction if VX equals NN
	case 0x3000:
		if cpu.Register.V[x] == nn {
			cpu.Register.PC += 4
		} else {
			cpu.Register.PC += 2
		}
	// 4XNN: Skips the next instruction if VX doesn't equal NN
	case 0x4000:
		if cpu.Register.V[x] != nn {
			cpu.Register.PC += 4
		} else {
			cpu.Register.PC += 2
		}
	// 5XY0: Skips the next instruction if VX equals VY
	case 0x5000:
		if cpu.Register.V[x] == cpu.Register.V[y] {
			cpu.Register.PC += 4
		} else {
			cpu.Register.PC += 2
		}
	// 6XNN: Sets VX to NN
	case 0x6000:
		cpu.Register.V[x] = nn
		cpu.Register.PC += 2
	// 7XNN: Adds NN to VX (Carry flag is not changed)
	case 0x7000:
		cpu.Register.V[x] += nn
		cpu.Register.PC += 2
	case 0x8000:
		switch opcode & 0x000F {
		// 8XY0: Sets VX to the value of VY
		case 0x0000:
			cpu.Register.V[x] = cpu.Register.V[y]
			cpu.Register.PC += 2
		// 8XY1: Sets VX to VX or VY (Bitwise OR operation)
		case 0x0001:
			cpu.Register.V[x] |= cpu.Register.V[y]
			cpu.Register.PC += 2
		// 8XY2: Sets VX to VX and VY (Bitwise AND operation)
		case 0x0002:
			cpu.Register.V[x] &= cpu.Register.V[y]
			cpu.Register.PC += 2
		// 8XY3: Sets VX to VX xor VY
		case 0x0003:
			cpu.Register.V[x] ^= cpu.Register.V[y]
			cpu.Register.PC += 2
		// 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
		case 0x0004:
			if cpu.Register.V[y] > (0xFF - cpu.Register.V[x]) {
				cpu.Register.V[0xF] = 1
			} else {
				cpu.Register.V[0xF] = 0
			}
			cpu.Register.V[x] += cpu.Register.V[y]
			cpu.Register.PC += 2
		// 8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		case 0x0005:
			if cpu.Register.V[y] > cpu.Register.V[x] {
				cpu.Register.V[0xF] = 0
			} else {
				cpu.Register.V[0xF] = 1
			}
			cpu.Register.V[x] -= cpu.Register.V[y]
			cpu.Register.PC += 2
		// 8XY6: Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
		case 0x0006:
			//TODO: Complete 8XY6
		// 8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		case 0x0007:
			if cpu.Register.V[x] > cpu.Register.V[y] {
				cpu.Register.V[0xF] = 0
			} else {
				cpu.Register.V[0xF] = 1
			}
			cpu.Register.V[x] = cpu.Register.V[y] - cpu.Register.V[x]
			cpu.Register.PC += 2
		// 8XYE: Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
		case 0x000E:
			//TODO: Complete 8XYE
		}
	// 9XY0: Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
	case 0x9000:
		if cpu.Register.V[x] != cpu.Register.V[y] {
			cpu.Register.PC += 4
		} else {
			cpu.Register.PC += 2
		}
	}

}
