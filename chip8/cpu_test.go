package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCPU(t *testing.T) {
	cpu := NewCPU()
	assert.NotNil(t, cpu)
}

func TestCPU_Reset(t *testing.T) {
	cpu := NewCPU()
	_ = cpu.LoadROM("roms/PONG")
	cpu.Register.V[0xA] = 0x10
	cpu.Register.PC = 0x1018
	cpu.Stack[0x0] = 0x1018
	cpu.Memory.Memory[0x0] = 0x12
	cpu.Display[10][18] = 0x4
	cpu.NeedDraw = true
	cpu.WaitInput = true
	cpu.KeyState[1] = 0x34
	cpu.Reset()
	newCPU := NewCPU()
	assert.Equal(t, newCPU, cpu)
}
