package chip8

import "math/rand"

func (cpu *CPU) exec00E0() {
	cpu.ClearDisplay()
	cpu.Register.PC += 2
}

func (cpu *CPU) exec00EE() {
	cpu.Register.SP--
	cpu.Register.PC = cpu.Stack[cpu.Register.SP] + 2
}

func (cpu *CPU) exec1NNN(nnn uint16) {
	cpu.Register.PC = nnn
}

func (cpu *CPU) exec2NNN(nnn uint16) {
	cpu.Stack[cpu.Register.SP] = cpu.Register.PC
	cpu.Register.SP++
	cpu.Register.PC = nnn
}

func (cpu *CPU) exec3XNN(x uint16, nn byte) {
	if cpu.Register.V[x] == nn {
		cpu.Register.PC += 4
	} else {
		cpu.Register.PC += 2
	}
}

func (cpu *CPU) exec4XNN(x uint16, nn byte) {
	if cpu.Register.V[x] != nn {
		cpu.Register.PC += 4
	} else {
		cpu.Register.PC += 2
	}
}

func (cpu *CPU) exec5XY0(x, y uint16) {
	if cpu.Register.V[x] == cpu.Register.V[y] {
		cpu.Register.PC += 4
	} else {
		cpu.Register.PC += 2
	}
}

func (cpu *CPU) exec6XNN(x uint16, nn byte) {
	cpu.Register.V[x] = nn
	cpu.Register.PC += 2
}

func (cpu *CPU) exec7XNN(x uint16, nn byte) {
	cpu.Register.V[x] += nn
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY0(x, y uint16) {
	cpu.Register.V[x] = cpu.Register.V[y]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY1(x, y uint16) {
	cpu.Register.V[x] |= cpu.Register.V[y]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY2(x, y uint16) {
	cpu.Register.V[x] &= cpu.Register.V[y]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY3(x, y uint16) {
	cpu.Register.V[x] ^= cpu.Register.V[y]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY4(x, y uint16) {
	if cpu.Register.V[y] > (0xFF - cpu.Register.V[x]) {
		cpu.Register.V[0xF] = 1
	} else {
		cpu.Register.V[0xF] = 0
	}
	cpu.Register.V[x] += cpu.Register.V[y]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY5(x, y uint16) {
	if cpu.Register.V[y] > cpu.Register.V[x] {
		cpu.Register.V[0xF] = 0
	} else {
		cpu.Register.V[0xF] = 1
	}
	cpu.Register.V[x] -= cpu.Register.V[y]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY6(x uint16) {
	cpu.Register.V[0xF] = cpu.Register.V[x] & 0x01
	cpu.Register.V[x] >>= 1
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XY7(x, y uint16) {
	if cpu.Register.V[x] > cpu.Register.V[y] {
		cpu.Register.V[0xF] = 0
	} else {
		cpu.Register.V[0xF] = 1
	}
	cpu.Register.V[x] = cpu.Register.V[y] - cpu.Register.V[x]
	cpu.Register.PC += 2
}

func (cpu *CPU) exec8XYE(x uint16) {
	cpu.Register.V[0xF] = cpu.Register.V[x] >> 7
	cpu.Register.V[x] <<= 1
	cpu.Register.PC += 2
}

func (cpu *CPU) exec9XY0(x, y uint16) {
	if cpu.Register.V[x] != cpu.Register.V[y] {
		cpu.Register.PC += 4
	} else {
		cpu.Register.PC += 2
	}
}

func (cpu *CPU) execANNN(nnn uint16) {
	cpu.Register.I = nnn
	cpu.Register.PC += 2
}

func (cpu *CPU) execBNNN(nnn uint16) {
	cpu.Register.PC = uint16(cpu.Register.V[0]) + nnn
}

func (cpu *CPU) execCXNN(x uint16, nn byte) {
	cpu.Register.V[x] = byte(rand.Float32()*255) & nn
	cpu.Register.PC += 2
}

func (cpu *CPU) execDXYN(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	xValue := cpu.Register.V[x]
	yValue := cpu.Register.V[y]
	height := byte(opcode & 0x000F)
	cpu.Register.V[0xF] = 0x00
	for i := yValue; i < yValue+height; i++ {
		for j := xValue; j < xValue+8; j++ {
			bit := (cpu.Memory.Memory[cpu.Register.I+uint16(i-yValue)] >> (7 - j + cpu.Register.V[x])) & 0x01
			xIndex, yIndex := j, i
			if j >= DisplayWidth {
				xIndex = j % DisplayWidth
			}
			if i >= DisplayHeight {
				yIndex = i % DisplayHeight
			}
			if bit == 0x01 && cpu.Display[yIndex][xIndex] == 0x01 {
				cpu.Register.V[0xF] = 0x01
			}
			cpu.Display[yIndex][xIndex] ^= bit
		}
	}
	cpu.NeedDraw = true
	cpu.Register.PC += 2
}

func (cpu *CPU) execEX9E(x uint16) {
	if cpu.KeyState[cpu.Register.V[x]] == 0x01 {
		cpu.Register.PC += 4
	} else {
		cpu.Register.PC += 2
	}
}

func (cpu *CPU) execEXA1(x uint16) {
	if cpu.KeyState[cpu.Register.V[x]] == 0x00 {
		cpu.Register.PC += 4
	} else {
		cpu.Register.PC += 2
	}
}

func (cpu *CPU) execFX07(x uint16) {
	cpu.Register.V[x] = cpu.Register.DT
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX0A(x uint16) {
	cpu.WaitInput = true
	for i, k := range cpu.KeyState {
		if k != 0 {
			cpu.Register.V[x] = byte(i)
			cpu.WaitInput = false
			cpu.Register.PC += 2
			break
		}
	}
}

func (cpu *CPU) execFX15(x uint16) {
	cpu.Register.DT = cpu.Register.V[x]
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX18(x uint16) {
	cpu.Register.ST = cpu.Register.V[x]
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX1E(x uint16) {
	cpu.Register.I += uint16(cpu.Register.V[x])
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX29(x uint16) {
	cpu.Register.I = uint16(cpu.Register.V[x]) * 5
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX33(x uint16) {
	cpu.Memory.Memory[cpu.Register.I] = cpu.Register.V[x] / 100
	cpu.Memory.Memory[cpu.Register.I+1] = (cpu.Register.V[x] / 10) % 10
	cpu.Memory.Memory[cpu.Register.I+2] = (cpu.Register.V[x] % 100) % 10
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX55(x uint16) {
	for i := uint16(0); i <= x; i++ {
		cpu.Memory.Memory[cpu.Register.I+i] = cpu.Register.V[i]
	}
	cpu.Register.PC += 2
}

func (cpu *CPU) execFX65(x uint16) {
	for i := uint16(0); i <= x; i++ {
		cpu.Register.V[i] = cpu.Memory.Memory[cpu.Register.I+i]
	}
	cpu.Register.PC += 2
}
