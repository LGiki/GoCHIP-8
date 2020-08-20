package chip8

import (
	"fmt"
	"math/rand"
)

type CPU struct {
	Register Register
	Memory   Memory
	// internal stack to store return addresses when calling procedures
	Stack [16]uint16
	// 2D array representing 64 x 32 grid
	Display [64][32]byte
	// State of the keys
	KeyState [16]byte
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
			cpu.Register.V[0xF] = cpu.Register.V[x] & 0x01
			cpu.Register.V[x] >>= 1
			cpu.Register.PC += 2
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
			cpu.Register.V[0xF] = cpu.Register.V[x] >> 7
			cpu.Register.V[x] <<= 1
			cpu.Register.PC += 2
		}
	// 9XY0: Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
	case 0x9000:
		if cpu.Register.V[x] != cpu.Register.V[y] {
			cpu.Register.PC += 4
		} else {
			cpu.Register.PC += 2
		}
	// ANNN: Sets I to the address NNN
	case 0xA000:
		cpu.Register.I = nnn
		cpu.Register.PC += 2
	// BNNN: Jumps to the address NNN plus V0
	case 0xB000:
		cpu.Register.PC = uint16(cpu.Register.V[0]) + nnn
	// CXNN: Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
	case 0xC000:
		cpu.Register.V[x] = byte(rand.Float32()*255) & nn
		cpu.Register.PC += 2
	// DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory
	//	     location I; I value doesn't change after the execution of this instruction. As described above, VF is set to 1 if any screen pixels are flipped from set to
	//       unset when the sprite is drawn, and to 0 if that doesn’t happen
	case 0xD000:
		//TODO: Implement DXYN
	case 0xE000:
		switch opcode & 0x00FF {
		// EX9E: Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
		case 0x009E:
			if cpu.KeyState[cpu.Register.V[x]] == 0x01 {
				cpu.Register.PC += 4
			} else {
				cpu.Register.PC += 2
			}
		// EXA1: Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)
		case 0x00A1:
			if cpu.KeyState[cpu.Register.V[x]] == 0x00 {
				cpu.Register.PC += 4
			} else {
				cpu.Register.PC += 2
			}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		// FX07: Sets VX to the value of the delay timer
		case 0x0007:
			cpu.Register.V[x] = cpu.Register.DT
			cpu.Register.PC += 2
		// FX0A: A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
		case 0x000A:
			for i, k := range cpu.KeyState {
				if k != 0 {
					cpu.Register.V[x] = byte(i)
					cpu.Register.PC += 2
					break
				}
			}
		// FX15: Sets the delay timer to VX
		case 0x0015:
			cpu.Register.DT = cpu.Register.V[x]
			cpu.Register.PC += 2
		// FX18: Sets the sound timer to VX
		case 0x0018:
			cpu.Register.ST = cpu.Register.V[x]
			cpu.Register.PC += 2
		// FX1E: Adds VX to I, VF is not affected
		case 0x001E:
			cpu.Register.I += uint16(cpu.Register.V[x])
			cpu.Register.PC += 2
		// FX29: Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
		case 0x0029:
			cpu.Register.I = uint16(cpu.Register.V[x]) * 5
			cpu.Register.PC += 2
		// FX33: Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the
		//       least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit
		//       at location I+1, and the ones digit at location I+2.)
		case 0x0033:
			cpu.Memory.Memory[cpu.Register.I] = cpu.Register.V[x] / 100
			cpu.Memory.Memory[cpu.Register.I+1] = (cpu.Register.V[x] / 10) % 10
			cpu.Memory.Memory[cpu.Register.I+2] = (cpu.Register.V[x] % 100) % 10
			cpu.Register.PC += 2
		// FX55: Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified
		case 0x0055:
			for i := uint16(0); i <= x; i++ {
				cpu.Memory.Memory[cpu.Register.I+i] = cpu.Register.V[i]
			}
			cpu.Register.PC += 2
		// FX65: Fills V0 to VX (including VX) with values from memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified
		case 0x0065:
			for i := uint16(0); i <= x; i++ {
				cpu.Register.V[i] = cpu.Memory.Memory[cpu.Register.I+i]
			}
			cpu.Register.PC += 2
		}
	}
}

func (cpu *CPU) Debug() {
	fmt.Printf("PC: %d\n", cpu.Register.PC)
	fmt.Printf("SP: %d\n", cpu.Register.SP)
	fmt.Printf("I: %d\n", cpu.Register.I)
	for i, value := range cpu.Register.V {
		fmt.Printf("V%d: %d\n", i, value)
	}
}
