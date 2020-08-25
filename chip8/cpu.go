package chip8

import (
	"fmt"
	"log"
)

const (
	DisplayHeight = 32
	DisplayWidth  = 64
)

type CPU struct {
	Register Register
	Memory   Memory
	// internal stack to store return addresses when calling procedures
	Stack [16]uint16
	// 2D array representing 32 x 64 grid
	Display [DisplayHeight][DisplayWidth]byte
	// State of the keys
	KeyState [16]byte
	// Need draw or not
	NeedDraw bool
	// Is wait for input (used by FX0A)
	WaitInput bool
}

func NewCPU() CPU {
	cpu := CPU{}
	cpu.Register.PC = 0x200
	cpu.Memory = Memory{}
	cpu.Memory.LoadFontSet()
	return cpu
}

func (cpu *CPU) Reset() {
	cpu.Register.PC = 0x200
	cpu.Register.I = 0
	cpu.Register.SP = 0
	cpu.Register.DT = 0
	cpu.Register.ST = 0
	cpu.NeedDraw = false
	cpu.WaitInput = false
	for i := 0; i < len(cpu.Register.V); i++ {
		cpu.Register.V[i] = 0
	}
	for i := 0; i < len(cpu.Memory.Memory); i++ {
		cpu.Memory.Memory[i] = 0
	}
	for i := 0; i < len(cpu.Stack); i++ {
		cpu.Stack[i] = 0
	}
	for i := 0; i < len(cpu.KeyState); i++ {
		cpu.KeyState[i] = 0
	}
	cpu.Memory.LoadFontSet()
	cpu.ClearDisplay()
}

func (cpu *CPU) ClearDisplay() {
	for x := 0; x < DisplayHeight; x++ {
		for y := 0; y < DisplayWidth; y++ {
			cpu.Display[x][y] = 0
		}
	}
}

func (cpu *CPU) LoadROM(romPath string) error {
	return cpu.Memory.LoadROM(romPath)
}

func (cpu *CPU) Run() {
	cpu.Cycle()
	if cpu.Register.DT > 0 {
		cpu.Register.DT--
	}
	if cpu.Register.ST > 0 {
		cpu.Register.ST--
	}
}

func (cpu *CPU) getOpCode() uint16 {
	return uint16(cpu.Memory.Memory[cpu.Register.PC])<<8 | uint16(cpu.Memory.Memory[cpu.Register.PC+1])
}

func (cpu *CPU) Cycle() {
	defer func() {
		if err := recover(); err != nil {
			cpu.Debug()
			log.Fatalf("%s\n", err)
		}
	}()
	opcode := cpu.getOpCode()
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00FF {
		// 00E0: Clears the screen
		case 0x00E0:
			cpu.exec00E0()
		// 00EE: Returns from a subroutine
		case 0x00EE:
			cpu.exec00EE()
		}
	// 1NNN: goto NNN
	case 0x1000:
		cpu.exec1NNN(nnn)
	// 2NNN: Calls subroutine at NNN
	case 0x2000:
		cpu.exec2NNN(nnn)
	// 3XNN: Skips the next instruction if VX equals NN
	case 0x3000:
		cpu.exec3XNN(x, nn)
	// 4XNN: Skips the next instruction if VX doesn't equal NN
	case 0x4000:
		cpu.exec4XNN(x, nn)
	// 5XY0: Skips the next instruction if VX equals VY
	case 0x5000:
		cpu.exec5XY0(x, y)
	// 6XNN: Sets VX to NN
	case 0x6000:
		cpu.exec6XNN(x, nn)
	// 7XNN: Adds NN to VX (Carry flag is not changed)
	case 0x7000:
		cpu.exec7XNN(x, nn)
	case 0x8000:
		switch opcode & 0x000F {
		// 8XY0: Sets VX to the value of VY
		case 0x0000:
			cpu.exec8XY0(x, y)
		// 8XY1: Sets VX to VX or VY (Bitwise OR operation)
		case 0x0001:
			cpu.exec8XY1(x, y)
		// 8XY2: Sets VX to VX and VY (Bitwise AND operation)
		case 0x0002:
			cpu.exec8XY2(x, y)
		// 8XY3: Sets VX to VX xor VY
		case 0x0003:
			cpu.exec8XY3(x, y)
		// 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
		case 0x0004:
			cpu.exec8XY4(x, y)
		// 8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		case 0x0005:
			cpu.exec8XY5(x, y)
		// 8XY6: Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
		case 0x0006:
			cpu.exec8XY6(x)
		// 8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		case 0x0007:
			cpu.exec8XY7(x, y)
		// 8XYE: Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
		case 0x000E:
			cpu.exec8XYE(x)
		}
	// 9XY0: Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
	case 0x9000:
		cpu.exec9XY0(x, y)
	// ANNN: Sets I to the address NNN
	case 0xA000:
		cpu.execANNN(nnn)
	// BNNN: Jumps to the address NNN plus V0
	case 0xB000:
		cpu.execBNNN(nnn)
	// CXNN: Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
	case 0xC000:
		cpu.execCXNN(x, nn)
	// DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory
	//	     location I; I value doesn't change after the execution of this instruction. As described above, VF is set to 1 if any screen pixels are flipped from set to
	//       unset when the sprite is drawn, and to 0 if that doesn’t happen
	case 0xD000:
		cpu.execDXYN(opcode)
	case 0xE000:
		switch opcode & 0x00FF {
		// EX9E: Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
		case 0x009E:
			cpu.execEX9E(x)
		// EXA1: Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)
		case 0x00A1:
			cpu.execEXA1(x)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		// FX07: Sets VX to the value of the delay timer
		case 0x0007:
			cpu.execFX07(x)
		// FX0A: A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
		case 0x000A:
			cpu.execFX0A(x)
		// FX15: Sets the delay timer to VX
		case 0x0015:
			cpu.execFX15(x)
		// FX18: Sets the sound timer to VX
		case 0x0018:
			cpu.execFX18(x)
		// FX1E: Adds VX to I, VF is not affected
		case 0x001E:
			cpu.execFX1E(x)
		// FX29: Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
		case 0x0029:
			cpu.execFX29(x)
		// FX33: Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the
		//       least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit
		//       at location I+1, and the ones digit at location I+2.)
		case 0x0033:
			cpu.execFX33(x)
		// FX55: Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified
		case 0x0055:
			cpu.execFX55(x)
		// FX65: Fills V0 to VX (including VX) with values from memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified
		case 0x0065:
			cpu.execFX65(x)
		}
	}
}

func (cpu *CPU) Debug() {
	fmt.Printf("OpCode: %X\n", cpu.getOpCode())
	fmt.Printf("PC: %d\n", cpu.Register.PC)
	fmt.Printf("SP: %d\n", cpu.Register.SP)
	fmt.Printf("I: %d\n", cpu.Register.I)
	for i, value := range cpu.Register.V {
		fmt.Printf("V%d: %d\n", i, value)
	}
}
