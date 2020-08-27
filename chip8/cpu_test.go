package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCPU(t *testing.T) {
	cpu := NewCPU()
	assert.NotNil(t, cpu)
}

func TestCPU_LoadROM(t *testing.T) {
	cpu := NewCPU()
	err := cpu.LoadROM("../roms/PONG")
	assert.Nil(t, err)
	err = cpu.LoadROM("null")
	assert.NotNil(t, err)
}

func TestCPU_Debug(t *testing.T) {
	cpu := NewCPU()
	_ = cpu.LoadROM("../roms/PONG")
	cpu.Run()
	cpu.Debug()
}

func TestCPU_Cycle(t *testing.T) {
	cpu := NewCPU()
	_ = cpu.LoadROM("../roms/PONG")
	cpu.Register.ST = 10
	cpu.Register.DT = 18
	cpu.Run()
	assert.Equal(t, byte(9), cpu.Register.ST)
	assert.Equal(t, byte(17), cpu.Register.DT)
}

func TestCPU_Reset(t *testing.T) {
	cpu := NewCPU()
	_ = cpu.LoadROM("../roms/PONG")
	cpu.Run()
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

func TestCPU_ClearDisplay(t *testing.T) {
	cpu := NewCPU()
	cpu.Display[0][1] = 0x01
	cpu.Display[5][9] = 0x01
	cpu.Display[4][5] = 0x01
	cpu.Display[2][1] = 0x01
	cpu.Display[3][0] = 0x01
	cpu.Display[8][8] = 0x01
	cpu.ClearDisplay()
	newCPU := NewCPU()
	assert.Equal(t, newCPU, cpu)
}

func TestGetOpCode(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory.Memory[0x0A] = 0x10
	cpu.Memory.Memory[0x0B] = 0x18
	cpu.Register.PC = 0x0A
	opCode := cpu.getOpCode()
	assert.Equal(t, uint16(0x1018), opCode)
}
