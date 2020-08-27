package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExec00E0(t *testing.T) {
	cpu := NewCPU()
	cpu.Display[0][0] = 0x1
	cpu.Display[0][1] = 0x2
	cpu.Display[5][2] = 0x3
	cpu.Display[10][18] = 0x4
	cpu.Display[15][12] = 0xF
	cpu.exec00E0()
	for i := 0; i < DisplayHeight; i++ {
		for j := 0; j < DisplayWidth; j++ {
			assert.Equal(t, byte(0), cpu.Display[i][j])
		}
	}
}

func TestExec00EE(t *testing.T) {
	cpu := NewCPU()
	cpu.Stack[0x0] = 0x1018
	cpu.Register.SP = 1
	cpu.exec00EE()
	newCPU := NewCPU()
	newCPU.Stack[0x0] = 0x1018
	newCPU.Register.SP = 0
	newCPU.Register.PC = 0x101A
	assert.Equal(t, newCPU, cpu)
}

func TestExec1NNN(t *testing.T) {
	cpu := NewCPU()
	cpu.exec1NNN(0x1018)
	newCPU := NewCPU()
	newCPU.Register.PC = 0x1018
	assert.Equal(t, newCPU, cpu)
}

func TestExec2NNN(t *testing.T) {
	cpu := NewCPU()
	cpu.exec2NNN(0x1018)
	newCPU := NewCPU()
	newCPU.Stack[0x0] = 0x200
	newCPU.Register.SP = 0x1
	newCPU.Register.PC = 0x1018
	assert.Equal(t, newCPU, cpu)
}

func TestExec3XNN(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x18
	cpu.exec3XNN(0xA, 0x18)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x18
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xA] = 0x19
	cpu.exec3XNN(0xA, 0x18)
	newCPU.Register.V[0xA] = 0x19
	newCPU.Register.PC = 0x206
	assert.Equal(t, newCPU, cpu)
}

func TestExec4XNN(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x19
	cpu.exec4XNN(0xA, 0x18)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x19
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xA] = 0x18
	cpu.exec4XNN(0xA, 0x18)
	newCPU.Register.V[0xA] = 0x18
	newCPU.Register.PC = 0x206
	assert.Equal(t, newCPU, cpu)
}

func TestExec5XY0(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x18
	cpu.Register.V[0xB] = 0x18
	cpu.exec5XY0(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x18
	newCPU.Register.V[0xB] = 0x18
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xB] = 0x19
	cpu.exec5XY0(0xA, 0xB)
	newCPU.Register.V[0xB] = 0x19
	newCPU.Register.PC = 0x206
	assert.Equal(t, newCPU, cpu)
}

func TestExec6XNN(t *testing.T) {
	cpu := NewCPU()
	cpu.exec6XNN(0xA, 0x18)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x18
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec7XNN(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x02
	cpu.exec7XNN(0xA, 0x18)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x1A
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY0(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xB] = 0x18
	cpu.exec8XY0(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x18
	newCPU.Register.V[0xB] = 0x18
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY1(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x1C
	cpu.Register.V[0xB] = 0xF8
	cpu.exec8XY1(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x1C | 0xF8
	newCPU.Register.V[0xB] = 0xF8
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY2(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x1C
	cpu.Register.V[0xB] = 0xF8
	cpu.exec8XY2(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x1C & 0xF8
	newCPU.Register.V[0xB] = 0xF8
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY3(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x1C
	cpu.Register.V[0xB] = 0xF8
	cpu.exec8XY3(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x1C ^ 0xF8
	newCPU.Register.V[0xB] = 0xF8
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY4(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0xFF
	cpu.Register.V[0xB] = 0x01
	cpu.exec8XY4(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x00
	newCPU.Register.V[0xB] = 0x01
	newCPU.Register.V[0xF] = 1
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xA] = 0xC3
	cpu.Register.V[0xB] = 0x0F
	cpu.exec8XY4(0xA, 0xB)
	newCPU.Register.V[0xA] = 0xD2
	newCPU.Register.V[0xB] = 0x0F
	newCPU.Register.V[0xF] = 0
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY5(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x0F
	cpu.Register.V[0xB] = 0x06
	cpu.exec8XY5(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x09
	newCPU.Register.V[0xB] = 0x06
	newCPU.Register.V[0xF] = 1
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xA] = 0x06
	cpu.Register.V[0xB] = 0x0F
	cpu.exec8XY5(0xA, 0xB)
	newCPU.Register.V[0xA] = 0xF7
	newCPU.Register.V[0xB] = 0x0F
	newCPU.Register.V[0xF] = 0
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY6(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0b10101010
	cpu.exec8XY6(0xA)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0b01010101
	newCPU.Register.V[0xF] = 0
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xA] = 0b01010101
	cpu.exec8XY6(0xA)
	newCPU.Register.V[0xA] = 0b00101010
	newCPU.Register.PC = 0x204
	newCPU.Register.V[0xF] = 1
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XY7(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0b10101010
	cpu.Register.V[0xB] = 0b00111100
	cpu.exec8XY7(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0b10010010
	newCPU.Register.V[0xB] = 0b00111100
	newCPU.Register.V[0xF] = 0
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xA] = 0b10101010
	cpu.Register.V[0xB] = 0b11111000
	cpu.exec8XY7(0xA, 0xB)
	newCPU.Register.V[0xA] = 0b01001110
	newCPU.Register.V[0xB] = 0b11111000
	newCPU.Register.PC = 0x204
	newCPU.Register.V[0xF] = 1
	assert.Equal(t, newCPU, cpu)
}

func TestExec8XYE(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0b10101010
	cpu.exec8XYE(0xA)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0b01010100
	newCPU.Register.V[0xF] = 0b00000001
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExec9XY0(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xA] = 0x18
	cpu.Register.V[0xB] = 0x19
	cpu.exec9XY0(0xA, 0xB)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x18
	newCPU.Register.V[0xB] = 0x19
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)

	cpu.Register.V[0xB] = 0x18
	cpu.exec9XY0(0xA, 0xB)
	newCPU.Register.V[0xB] = 0x18
	newCPU.Register.PC = 0x206
	assert.Equal(t, newCPU, cpu)
}

func TestExecANNN(t *testing.T) {
	cpu := NewCPU()
	cpu.execANNN(0x1018)
	newCPU := NewCPU()
	newCPU.Register.I = 0x1018
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecBNNN(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0] = 0x12
	cpu.execBNNN(0x1018)
	newCPU := NewCPU()
	newCPU.Register.V[0] = 0x12
	newCPU.Register.PC = 0x102A
	assert.Equal(t, newCPU, cpu)
}

func TestExecCXNN(t *testing.T) {
	cpu := NewCPU()
	cpu.execCXNN(0xA, 0x00)
	newCPU := NewCPU()
	newCPU.Register.V[0xA] = 0x00
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
	cpu.execCXNN(0xA, 0xFF)
	assert.LessOrEqual(t, cpu.Register.V[0xA], byte(0xFF))
	assert.GreaterOrEqual(t, cpu.Register.V[0xA], byte(0x00))
}

func TestExecDXYN(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.I = 0x300
	cpu.Register.V[0x3] = 0
	cpu.Register.V[0xD] = 0
	cpu.Memory.Memory[0x300] = 0x10
	cpu.Memory.Memory[0x301] = 0x18
	cpu.execDXYN(0xD3D2)
	newCPU := NewCPU()
	newCPU.Register.I = 0x300
	newCPU.Register.V[0x3] = 0
	newCPU.Register.V[0xD] = 0
	newCPU.Memory.Memory[0x300] = 0x10
	newCPU.Memory.Memory[0x301] = 0x18
	newCPU.Display[0][3] = 0x01
	newCPU.Display[1][3] = 0x01
	newCPU.Display[1][4] = 0x01
	newCPU.NeedDraw = true
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)

	cpu = NewCPU()
	cpu.Register.I = 0x300
	cpu.Memory.Memory[0x300] = 0x10
	cpu.Memory.Memory[0x301] = 0x18
	cpu.Register.V[0x3] = 63
	cpu.Register.V[0xD] = 31
	cpu.execDXYN(0xD3D2)
	newCPU = NewCPU()
	newCPU.Register.I = 0x300
	newCPU.Register.V[0x3] = 63
	newCPU.Register.V[0xD] = 31
	newCPU.Memory.Memory[0x300] = 0x10
	newCPU.Memory.Memory[0x301] = 0x18
	newCPU.Display[31][2] = 0x01
	newCPU.Display[0][2] = 0x01
	newCPU.Display[0][3] = 0x01
	newCPU.NeedDraw = true
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)

	cpu = NewCPU()
	cpu.Register.I = 0x300
	cpu.Memory.Memory[0x300] = 0x10
	cpu.Memory.Memory[0x301] = 0x18
	cpu.Register.V[0x3] = 63
	cpu.Register.V[0xD] = 31
	cpu.Display[31][2] = 0x01
	cpu.execDXYN(0xD3D2)
	newCPU = NewCPU()
	newCPU.Register.I = 0x300
	newCPU.Register.V[0x3] = 63
	newCPU.Register.V[0xD] = 31
	newCPU.Memory.Memory[0x300] = 0x10
	newCPU.Memory.Memory[0x301] = 0x18
	newCPU.Display[31][2] = 0x00
	newCPU.Display[0][2] = 0x01
	newCPU.Display[0][3] = 0x01
	newCPU.Register.V[0xF] = 0x01
	newCPU.NeedDraw = true
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecEX9E(t *testing.T) {
	cpu := NewCPU()
	cpu.KeyState[0x06] = 0x01
	cpu.Register.V[0x05] = 0x06
	cpu.execEX9E(0x05)
	newCPU := NewCPU()
	newCPU.KeyState[0x06] = 0x01
	newCPU.Register.V[0x05] = 0x06
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)

	cpu.KeyState[0x06] = 0x00
	cpu.execEX9E(0x05)
	newCPU.Register.PC = 0x206
	newCPU.KeyState[0x06] = 0x00
	assert.Equal(t, newCPU, cpu)
}

func TestExecEXA1(t *testing.T) {
	cpu := NewCPU()
	cpu.KeyState[0x06] = 0x00
	cpu.Register.V[0x05] = 0x06
	cpu.execEXA1(0x05)
	newCPU := NewCPU()
	newCPU.KeyState[0x06] = 0x00
	newCPU.Register.V[0x05] = 0x06
	newCPU.Register.PC = 0x204
	assert.Equal(t, newCPU, cpu)

	cpu.KeyState[0x06] = 0x01
	cpu.execEXA1(0x05)
	newCPU.Register.PC = 0x206
	newCPU.KeyState[0x06] = 0x01
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX07(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.DT = 0x6C
	cpu.execFX07(0xB)
	newCPU := NewCPU()
	newCPU.Register.DT = 0x6C
	newCPU.Register.V[0xB] = 0x6C
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX0A(t *testing.T) {
	cpu := NewCPU()
	cpu.execFX0A(0xB)
	newCPU := NewCPU()
	newCPU.WaitInput = true
	assert.Equal(t, newCPU, cpu)

	cpu.execFX0A(0xB)
	assert.Equal(t, newCPU, cpu)

	cpu.execFX0A(0xB)
	assert.Equal(t, newCPU, cpu)

	cpu.KeyState[0xD] = 1
	cpu.execFX0A(0xB)
	newCPU.KeyState[0xD] = 1
	newCPU.WaitInput = false
	newCPU.Register.V[0xB] = 0xD
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX15(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xB] = 0x6C
	cpu.execFX15(0xB)
	newCPU := NewCPU()
	newCPU.Register.DT = 0x6C
	newCPU.Register.V[0xB] = 0x6C
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX18(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xB] = 0x6C
	cpu.execFX18(0xB)
	newCPU := NewCPU()
	newCPU.Register.ST = 0x6C
	newCPU.Register.V[0xB] = 0x6C
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX1E(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0xE] = 0xA6
	cpu.Register.I = 0x1018
	cpu.execFX1E(0xE)
	newCPU := NewCPU()
	newCPU.Register.V[0xE] = 0xA6
	newCPU.Register.I = 0x10BE
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX29(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0x3] = 0xF
	cpu.execFX29(0x3)
	newCPU := NewCPU()
	newCPU.Register.I = 0x004B
	newCPU.Register.V[0x3] = 0xF
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX33(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.V[0x5] = 128
	cpu.Register.I = 0x000A
	cpu.execFX33(0x5)
	newCPU := NewCPU()
	newCPU.Register.V[0x5] = 128
	newCPU.Register.I = 0x000A
	newCPU.Memory.Memory[0x000A] = 1
	newCPU.Memory.Memory[0x000B] = 2
	newCPU.Memory.Memory[0x000C] = 8
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX55(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.I = 0x0100
	cpu.Register.V[0x0] = 0x10
	cpu.Register.V[0x1] = 0x18
	cpu.Register.V[0x2] = 0xAB
	cpu.execFX55(0x0002)
	newCPU := NewCPU()
	newCPU.Register.I = 0x0100
	newCPU.Register.V[0x0] = 0x10
	newCPU.Register.V[0x1] = 0x18
	newCPU.Register.V[0x2] = 0xAB
	newCPU.Memory.Memory[0x0100] = 0x10
	newCPU.Memory.Memory[0x0101] = 0x18
	newCPU.Memory.Memory[0x0102] = 0xAB
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}

func TestExecFX65(t *testing.T) {
	cpu := NewCPU()
	cpu.Register.I = 0x0100
	cpu.Memory.Memory[0x0100] = 0x10
	cpu.Memory.Memory[0x0101] = 0x18
	cpu.Memory.Memory[0x0102] = 0xAB
	cpu.execFX65(0x0002)
	newCPU := NewCPU()
	newCPU.Register.I = 0x0100
	newCPU.Memory.Memory[0x0100] = 0x10
	newCPU.Memory.Memory[0x0101] = 0x18
	newCPU.Memory.Memory[0x0102] = 0xAB
	newCPU.Register.V[0x0] = 0x10
	newCPU.Register.V[0x1] = 0x18
	newCPU.Register.V[0x2] = 0xAB
	newCPU.Register.PC = 0x202
	assert.Equal(t, newCPU, cpu)
}
